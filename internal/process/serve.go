package process

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"

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
		_client:        net.ParseIP(strings.Split(w.RemoteAddr().String(), ":")[0]),

		received: time.Now().UTC(),
	}

	dcs := strings.Split(q._domain, ".")
	if len(dcs) > 1 {
		suffix := strings.Join(dcs[1:], ".")
		q._suffix = &suffix
	}

	if q._domain == "!runprobe" {
		return dns.RcodeSuccess, w.WriteMsg(responseWithCode(r, 15)) // 15 = max valid unassigned RCODE
	}

	lease, found := getLeaseByIP(q._client)
	if found {
		m := lease.MAC.String()
		q.mac = &m
		q.hostname = lease.Hostname
		q.vendor = lease.Vendor
	}

	system.Console.Info("incoming: ", q.String())

	if whitelisted := whitelist(&q); whitelisted {
		q.action = ActionWhitelisted // Might be overwritten by other returners
		goto skip
	}

	if ret, rcode, err := blacklist(&q); ret {
		return rcode, err
	}

skip:

	// Always check this; queries with search domains will contain periods
	if ret, rcode, err := respondByHostname(&q); ret {
		return rcode, err
	}

	if ret, rcode, err := respondByPTR(&q); ret {
		return rcode, err
	}

	// TODO: strip search domain to check DomainNeeded safely
	if system.DomainNeeded && (strings.Count(q._domain, ".") == 0 || q._matchedSearchDomain != nil) {
		if q._question.Qtype == dns.TypeNS && q._domain == "" {
			// Permit looking up root servers
			goto next
		}

		if q._matchedSearchDomain == nil {
			system.Console.Infof("DomainNeeded is true and question '%s' does not contain any periods; returning NXDomain", q._domain)
		} else {
			system.Console.Infof("DomainNeeded is true and question '%s' has a search domain that was not found; returning NXDomain", q._domain)
		}
		q.action = ActionDomainNeeded
		m := responseWithCode(r, dns.RcodeNameError)
		err := q.respond(m)
		return dns.RcodeRefused, err
	}
next:

	q.action = ActionNotBlacklisted
	return plugin.NextOrFailure(name, next, ctx, q, r)
}
