package serve

import (
    "context"
    "errors"
    "strings"
    "time"

    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/log"
)

func DNS(name string, next plugin.Handler, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

    // https://stackoverflow.com/a/4083071/9911189
    // https://groups.google.com/forum/#!topic/comp.protocols.dns.bind/uOWxNkm7AVg
    if len(r.Question) != 1 {
        return 0, errors.New("this should never happen")
    }

    q := log.QueryInfo{
        ResponseWriter: w,
        R:              r,

        QU_:     r.Question[0],
        Domain_: rt(r.Question[0].Name),
        Client_: strings.Split(w.RemoteAddr().String(), ":")[0],

        Received: time.Now(),
    }

    if q.Domain_ == "!runprobe" {
        return dns.RcodeSuccess, w.WriteMsg(responseWithCode(r, 15)) // 15 = max valid unassigned RCODE
    }

    var began time.Time

    // These will be needed when we implement whitelisting anyway, so I don't mind looking them up for all requests
    began = time.Now()
    lease, found := getLeaseByIP(q.Client_)
    if found {
        q.MAC = &lease.MAC
        q.Hostname = lease.Hostname
        q.Vendor = lease.Vendor

        //vendor, ok := ouiMACPrefixToVendor.Load(trimmedMAC)
        //if ok {
        //    s := vendor.(string)
        //    q.vendor = &s
        //}
    }
    q.Durations.TimeLookupLease = time.Since(began)

    log.CLOG.Infof("%s -> %d %s", q.Client_, r.Id, dns.TypeToString[q.QU_.Qtype])

    began = time.Now()
    if strings.Count(q.Domain_, ".") == 0 {
        if ret, rcode, err := respondByHostname(&q); ret {
            return rcode, err
        }
    }
    q.Durations.TimeCheckRespondByHostname = time.Since(began)

    began = time.Now()
    if ret, rcode, err := respondByPTR(&q); ret {
        return rcode, err
    }
    q.Durations.TimeCheckRespondByPTR = time.Since(began)

    if q.Hostname != nil && strings.ToLower(*q.Hostname) == "syd-laptop" {
        log.CLOG.Infof("This is Syd's laptop; skipping respondWithBlock")
        goto skip
    }
    began = time.Now()
    if ret, rcode, err := respondWithBlock(&q); ret {
        return rcode, err
    }
    q.Durations.TimeCheckRespondWithBlock = time.Since(began)
skip:

    if config.Config.DomainNeeded && strings.Count(q.Domain_, ".") == 0 {
        if q.QU_.Qtype == dns.TypeNS && q.Domain_ == "" {
            // Permit looking up root servers
            goto next
        }

        log.CLOG.Infof("DomainNeeded is true and question '%s' does not contain any periods; returning RcodeServerFailure", q.Domain_)
        q.Action = "restrict"
        m := responseWithCode(q.R, dns.RcodeServerFailure)
        err := q.Respond(m)
        return dns.RcodeServerFailure, err
    }
next:

    q.Action = "pass"
    return plugin.NextOrFailure(name, next, ctx, q, r)
}
