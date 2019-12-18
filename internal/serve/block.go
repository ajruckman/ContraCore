package serve

import (
    "github.com/miekg/dns"
)

func respondWithBlock(q *queryContext) (ret bool, rcode int, err error) {
    if ruleCache.check(q._domain) {
        q.action = "block"
        var m *dns.Msg
        var v string

        switch q._qu.Qtype {
        case dns.TypeA:
            v = "0.0.0.0"
            m = GenResponse(q.r, q._qu.Qtype, "0.0.0.0")

        case dns.TypeAAAA:
            v = "::"
            m = GenResponse(q.r, q._qu.Qtype, "::")

        case dns.TypeCNAME:
            v = ""
            m = GenResponse(q.r, q._qu.Qtype, "")

        default:
            v = "-"
        }

        m = GenResponse(q.r, q._qu.Qtype, v)

        clog.Info("Blocking query '", q._domain, "' with value '", v, "'")
        err = q.Respond(m)

        return true, dns.RcodeSuccess, err
    }

    return
}
