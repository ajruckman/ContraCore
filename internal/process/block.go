package process

import (
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/state"
)

func respondWithBlock(q *queryContext) (ret bool, rcode int, err error) {
    if ruleCache.check(q._domain) {
        //q.Action = "block"
        var m *dns.Msg
        var v string

        switch q._question.Qtype {
        case dns.TypeA:
            v = config.SpoofedA
            m = genResponse(q.r, q._question.Qtype, v)

        case dns.TypeAAAA:
            v = config.SpoofedAAAA
            m = genResponse(q.r, q._question.Qtype, v)

        case dns.TypeCNAME:
            v = config.SpoofedCNAME
            m = genResponse(q.r, q._question.Qtype, v)

        default:
            v = config.SpoofedDefault
        }

        state.Console.Infof("Blocking query '%s' with value '%s'", q._domain, v)

        m = genResponse(q.r, q._question.Qtype, v)
        err = q.respond(m)

        return true, dns.RcodeSuccess, err
    }

    return
}
