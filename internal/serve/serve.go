package serve

import (
    "context"
    "errors"
    "net"
    "strings"
    "time"

    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/config"
)

func DNS(name string, next plugin.Handler, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

    // https://stackoverflow.com/a/4083071/9911189
    // https://groups.google.com/forum/#!topic/comp.protocols.dns.bind/uOWxNkm7AVg
    if len(r.Question) != 1 {
        return 0, errors.New("this should never happen")
    }

    q := queryContext{
        ResponseWriter: w,
        r:              r,

        _qu:     r.Question[0],
        _domain: rt(r.Question[0].Name),
        _client: strings.Split(w.RemoteAddr().String(), ":")[0],

        received: time.Now(),
    }

    if q._domain == "!runprobe" {
        return dns.RcodeSuccess, w.WriteMsg(responseWithCode(r, dns.RcodeSuccess))
    }

    var began time.Time

    // These will be needed when we implement whitelisting anyway, so I don't mind looking them up for all requests
    began = time.Now()
    lease, found := getLeaseByIP(q._client)
    if found {
        q.mac = &lease.MAC
        q.hostname = lease.Hostname
        q.vendor = lease.Vendor

        //vendor, ok := ouiMACPrefixToVendor.Load(trimmedMAC)
        //if ok {
        //    s := vendor.(string)
        //    q.vendor = &s
        //}
    }
    q.durations.timeLookupLease = time.Since(began)

    clog.Infof("%s -> %d %s", q._client, r.Id, dns.TypeToString[q._qu.Qtype])

    began = time.Now()
    if strings.Count(q._domain, ".") == 0 {
        if ret, rcode, err := respondByHostname(&q); ret {
            return rcode, err
        }
    }
    q.durations.timeCheckRespondByHostname = time.Since(began)

    began = time.Now()
    if ret, rcode, err := respondByPTR(&q); ret {
        return rcode, err
    }
    q.durations.timeCheckRespondByPTR = time.Since(began)

    if q.hostname != nil && strings.ToLower(*q.hostname) == "syd-laptop" {
        clog.Infof("This is Syd's laptop; skipping respondWithBlock")
        goto skip
    }
    began = time.Now()
    if ret, rcode, err := respondWithBlock(&q); ret {
        return rcode, err
    }
    q.durations.timeCheckRespondWithBlock = time.Since(began)
skip:

    if config.Config.DomainNeeded && strings.Count(q._domain, ".") == 0 {
        if q._qu.Qtype == dns.TypeNS && q._domain == "" {
            // Permit looking up root servers
            goto next
        }

        clog.Infof("DomainNeeded is true and question '%s' does not contain any periods; returning RcodeServerFailure", q._domain)
        q.action = "restrict"
        m := responseWithCode(q.r, dns.RcodeServerFailure)
        err := q.Respond(m)
        return dns.RcodeServerFailure, err
    }
next:

    q.action = "pass"
    return plugin.NextOrFailure(name, next, ctx, q, r)
}

type queryContext struct {
    dns.ResponseWriter
    r *dns.Msg

    _qu     dns.Question
    _domain string
    _client string

    received time.Time
    action   string

    mac      *net.HardwareAddr
    hostname *string
    vendor   *string

    answers []string

    durations durations
}

type durations struct {
    timeLookupLease            time.Duration
    timeCheckRespondByHostname time.Duration
    timeCheckRespondByPTR      time.Duration
    timeCheckRespondWithBlock  time.Duration
    timeGenLogStruct           time.Duration
    timeSaveLogToPG            time.Duration
    timeSaveLogToCH            time.Duration
    timeSendLogToEventClients  time.Duration
}

func (q *queryContext) Respond(res *dns.Msg) (err error) {
    var answers []string
    for _, v := range res.Answer {
        answers = append(answers, rrToString(v))
    }
    q.answers = answers

    logChannel <- *q

    err = q.ResponseWriter.WriteMsg(res)
    return
}

func (q queryContext) WriteMsg(r *dns.Msg) error {
    return q.Respond(r)
}
