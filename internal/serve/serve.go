package serve

import (
    "context"
    "errors"
    "strings"
    "time"

    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/schema"
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

    // These will be needed when we implement whitelisting anyway, so I don't mind looking them up for all requests
    lease, found := getLeaseByIP(q._client)
    if found {
        q.mac = lease.MAC
        q.hostname = lease.Hostname
    }

    clog.Infof("%s -> %d %s", q._client, r.Id, dns.TypeToString[q._qu.Qtype])

    if strings.Count(q._domain, ".") == 0 {
        if ret, rcode, err := respondByHostname(&q); ret {
            return rcode, err
        }
    }

    if ret, rcode, err := respondByPTR(&q); ret {
        return rcode, err
    }

    if ret, rcode, err := respondWithBlock(&q); ret {
        return rcode, err
    }

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

    mac      string
    hostname string
}

func (q *queryContext) Respond(res *dns.Msg) (err error) {
    var answers []string
    for _, v := range res.Answer {
        answers = append(answers, rrToString(v))
    }

    logChannel <- schema.Log{
        Client:       q._client,
        Question:     q._domain,
        QuestionType: dns.TypeToString[q._qu.Qtype],
        Action:       q.action,
        Answers:      answers,

        QueryID:  q.r.Id,
        Duration: time.Now().Sub(q.received),
    }

    err = q.ResponseWriter.WriteMsg(res)
    return
}

func (q queryContext) WriteMsg(r *dns.Msg) error {
    return q.Respond(r)
}
