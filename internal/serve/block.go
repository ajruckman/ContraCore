package serve

import (
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/log"
)

func respondWithBlock(q *log.QueryInfo) (ret bool, rcode int, err error) {
    if ruleCache.check(q.Domain_) {
        q.Action = "block"
        var m *dns.Msg
        var v string

        switch q.QU_.Qtype {
        case dns.TypeA:
            v = config.Config.SpoofedA
            m = genResponse(q.R, q.QU_.Qtype, v)

        case dns.TypeAAAA:
            v = config.Config.SpoofedAAAA
            m = genResponse(q.R, q.QU_.Qtype, v)

        case dns.TypeCNAME:
            v = config.Config.SpoofedCNAME
            m = genResponse(q.R, q.QU_.Qtype, v)

        default:
            v = config.Config.SpoofedDefault
        }

        log.CLOG.Infof("Blocking query '%s' with value '%s'", q.Domain_, v)

        m = genResponse(q.R, q.QU_.Qtype, v)
        err = q.Respond(m)

        return true, dns.RcodeSuccess, err
    }

    return
}
