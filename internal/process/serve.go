package process

import (
    "context"
    "errors"
    "strings"
    "time"

    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/functions"
    "github.com/ajruckman/ContraCore/internal/system"
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
        _question:      r.Question[0],
        _domain:        functions.RT(r.Question[0].Name),
        _client:        strings.Split(w.RemoteAddr().String(), ":")[0],

        received: time.Now().UTC(),
    }

    if q._domain == "!runprobe" {
        return dns.RcodeSuccess, w.WriteMsg(responseWithCode(r, 15)) // 15 = max valid unassigned RCODE
    }

    // These will be needed when we implement whitelisting anyway, so I don't mind looking them up for all requests
    lease, found := getLeaseByIP(q._client)
    if found {
        q.mac = &lease.MAC
        q.hostname = lease.Hostname
        q.vendor = lease.Vendor
    }

    system.Console.Infof("%s -> %d %s", q._client, r.Id, dns.TypeToString[q._question.Qtype])

    //if strings.Count(q._domain, ".") == 0 {
        // Always check this; queries with search domains will contain periods
        if ret, rcode, err := respondByHostname(&q); ret {
            return rcode, err
        }
    //}

    if ret, rcode, err := respondByPTR(&q); ret {
        return rcode, err
    }

    if q.hostname != nil && strings.ToLower(*q.hostname) == "syd-laptop" {
        system.Console.Infof("This is Syd's laptop; skipping respondWithBlock")
        goto skip
    }
    if ret, rcode, err := respondWithBlock(&q); ret {
        return rcode, err
    }
skip:

    // TODO: strip search domain to check DomainNeeded safely
    if config.DomainNeeded && strings.Count(q._domain, ".") == 0 {
        if q._question.Qtype == dns.TypeNS && q._domain == "" {
            // Permit looking up root servers
            goto next
        }

        system.Console.Infof("DomainNeeded is true and question '%s' does not contain any periods; returning RcodeServerFailure", q._domain)
        q.action = ActionDomainNeeded
        m := responseWithCode(r, dns.RcodeServerFailure)
        err := q.respond(m)
        return dns.RcodeServerFailure, err
    }
next:

    q.action = ActionNotBlacklisted
    return plugin.NextOrFailure(name, next, ctx, q, r)
}
