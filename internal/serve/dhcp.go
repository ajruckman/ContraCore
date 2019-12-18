package serve

import (
    "strings"

    . "github.com/ajruckman/xlib"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema"
)

var DHCPCache = map[string][]schema.LeaseDetails{}

func cacheDHCP() {
    leases, err := db.GetLeaseDetails()
    Err(err)

    for _, lease := range leases {
        if lease.Hostname == "" {
            continue
        }

        hostname := strings.ToLower(lease.Hostname)
        if _, exists := DHCPCache[hostname]; !exists {
            DHCPCache[hostname] = []schema.LeaseDetails{}
        }
        DHCPCache[hostname] = append(DHCPCache[hostname], lease)
    }
}

func respondWithDynamicDNS(q *queryContext) (ret bool, rcode int, err error) {
    if v, ok := DHCPCache[strings.ToLower(q._domain)]; ok {
        q.action = "answer" // TODO: Special action for RcodeServerFailure?

        var m *dns.Msg

        for _, lease := range v {
            if lease.IP.To4() != nil {
                if q._qu.Qtype == dns.TypeA {
                    m = GenResponse(q.r, q._qu.Qtype, lease.IP.To4().String())
                    err = q.Respond(m)
                    clog.Debug("lease IP is IPv4, question is A")
                    return true, dns.RcodeSuccess, err
                } else {
                    clog.Debug("lease IP is IPv4, question is AAAA")
                    continue
                }

            } else if lease.IP.To16() != nil {
                if q._qu.Qtype == dns.TypeAAAA {
                    m = GenResponse(q.r, q._qu.Qtype, lease.IP.To16().String())
                    err = q.Respond(m)
                    clog.Debug("lease IP is IPv6, question is AAAA")
                    return true, dns.RcodeSuccess, err
                } else {
                    clog.Debug("lease IP is IPv6, question is A")
                    continue
                }
            }
        }

        clog.Error("lease with hostname '", q._domain, "' exists but query type is not A or AAAA")
        m = RespondWithCode(q.r, dns.RcodeServerFailure)
        err = q.Respond(m)
        return true, dns.RcodeServerFailure, err
    }

    return
}
