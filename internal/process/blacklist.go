package process

import (
	"github.com/miekg/dns"

	"github.com/ajruckman/ContraCore/internal/cache"
	"github.com/ajruckman/ContraCore/internal/system"
)

func blacklist(q *queryContext) (ret bool, rcode int, err error) {
	if cache.BlacklistCache.Check(q._domain) {
		q.action = ActionBlock
		var m *dns.Msg
		var v string

		switch q._question.Qtype {
		case dns.TypeA:
			v = system.SpoofedA
			m = genResponse(q.r, q._question.Qtype, v)

		case dns.TypeAAAA:
			v = system.SpoofedAAAA
			m = genResponse(q.r, q._question.Qtype, v)

		case dns.TypeCNAME:
			v = system.SpoofedCNAME
			m = genResponse(q.r, q._question.Qtype, v)

		default:
			v = system.SpoofedDefault
		}

		system.Console.Infof("blocking query %d with value '%s'", q.r.Id, v)

		m = genResponse(q.r, q._question.Qtype, v)
		err = q.respond(m)

		return true, dns.RcodeSuccess, err
	}

	return
}
