package serve

import (
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/config"
)

func respondWithBlock(q *queryInfo) (ret bool, rcode int, err error) {
    if ruleCache.check(q._domain) {
        q.action = "block"
        var m *dns.Msg
        var v string

        switch q._qu.Qtype {
        case dns.TypeA:
            v = config.Config.SpoofedA
            m = genResponse(q.r, q._qu.Qtype, v)

        case dns.TypeAAAA:
            v = config.Config.SpoofedAAAA
            m = genResponse(q.r, q._qu.Qtype, v)

        case dns.TypeCNAME:
            v = config.Config.SpoofedCNAME
            m = genResponse(q.r, q._qu.Qtype, v)

        default:
            v = config.Config.SpoofedDefault
        }

        clog.Infof("Blocking query '%s' with value '%s'", q._domain, v)

        m = genResponse(q.r, q._qu.Qtype, v)
        err = q.Respond(m)

        return true, dns.RcodeSuccess, err
    }

    return
}
