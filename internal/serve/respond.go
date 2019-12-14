package serve

import (
    "net"
    "strings"

    "github.com/coredns/coredns/request"
    "github.com/miekg/dns"
)

func respondWithDynamicDNS(q *queryContext) (ret bool, rcode int, err error) {
    if v, ok := DHCPCache[strings.ToLower(q.domain)]; ok {
        q.action = "answer" // TODO: Special action for RcodeServerFailure?

        var m *dns.Msg

        for _, lease := range v {
            if lease.IP.To4() != nil {
                if q.qu.Qtype == dns.TypeA {
                    m = GenResponse(q.r, q.qu.Qtype, lease.IP.To4().String())
                    err = q.Respond(m)
                    clog.Info("lease IP is IPv4, question is A")
                    return true, dns.RcodeSuccess, err
                } else {
                    clog.Info("lease IP is IPv4, question is AAAA")
                    continue
                }

            } else if lease.IP.To16() != nil {
                if q.qu.Qtype == dns.TypeAAAA {
                    m = GenResponse(q.r, q.qu.Qtype, lease.IP.To16().String())
                    err = q.Respond(m)
                    clog.Info("lease IP is IPv6, question is AAAA")
                    return true, dns.RcodeSuccess, err
                } else {
                    clog.Info("lease IP is IPv6, question is A")
                    continue
                }
            }
        }

        clog.Error("lease with hostname '", q.domain, "' exists but query type is not A or AAAA")
        m = RespondWithCode(q.r, dns.RcodeServerFailure)
        err = q.Respond(m)
        return true, dns.RcodeServerFailure, err
    }

    //clog.Info("lease does not exist: ", nm)

    return
}

func respondWithBlock(q *queryContext) (ret bool, rcode int, err error) {
    if ruleCache.Check(q.domain) {
        q.action = "block"
        var m *dns.Msg
        var v string

        switch q.qu.Qtype {
        case dns.TypeA:
            v = "0.0.0.0"
            m = GenResponse(q.r, q.qu.Qtype, "0.0.0.0")

        case dns.TypeAAAA:
            v = "::"
            m = GenResponse(q.r, q.qu.Qtype, "::")

        case dns.TypeCNAME:
            v = ""
            m = GenResponse(q.r, q.qu.Qtype, "")

        default:
            v = "-"
        }

        m = GenResponse(q.r, q.qu.Qtype, v)

        clog.Info("Blocking query '", q.domain, "' with value '", v, "'")
        err = q.Respond(m)

        return true, dns.RcodeSuccess, err
    }

    return
}

//func GenResponse(w dns.ResponseWriter, r *dns.Msg, qtype uint16, value string) *dns.Msg {
func GenResponse(r *dns.Msg, qtype uint16, value string) *dns.Msg {
    m := new(dns.Msg)
    m.SetReply(r)
    m.Authoritative = true
    m.RecursionAvailable = true
    m.Rcode = dns.RcodeSuccess

    var rr dns.RR
    state := request.Request{Req: r}

    switch qtype {
    case dns.TypeA:
        rr = &dns.A{
            Hdr: dns.RR_Header{
                Name:   state.QName(),
                Rrtype: dns.TypeA,
                Class:  state.QClass(),
            },
            A: net.ParseIP(value).To4(),
        }
    case dns.TypeAAAA:
        rr = &dns.AAAA{
            Hdr: dns.RR_Header{
                Name:   state.QName(),
                Rrtype: dns.TypeAAAA,
                Class:  state.QClass(),
            },
            AAAA: net.ParseIP(value),
        }
    case dns.TypeCNAME:
        rr = &dns.CNAME{
            Hdr: dns.RR_Header{
                Name:   state.QName(),
                Rrtype: dns.TypeCNAME,
                Class:  state.QClass(),
            },
            Target: dns.Fqdn(value),
        }
    default:
        rr = &dns.ANY{
            Hdr: dns.RR_Header{
                Name:   state.QName(),
                Rrtype: dns.TypeANY,
                Class:  state.QClass(),
            },
        }
    }
    m.Answer = append(m.Answer, rr)

    return m
}

//func RespondWithCode(w dns.ResponseWriter, r *dns.Msg, code int) error {
func RespondWithCode(r *dns.Msg, code int) *dns.Msg {
    m := new(dns.Msg)
    m.SetReply(r)
    m.Authoritative = true
    //m.Compress = true
    m.RecursionAvailable = true
    m.Rcode = code

    return m
}
