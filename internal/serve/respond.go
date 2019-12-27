package serve

import (
    "net"

    "github.com/coredns/coredns/request"
    "github.com/miekg/dns"
)

func genResponse(r *dns.Msg, qtype uint16, value string) *dns.Msg {
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
    case dns.TypePTR:
        rr = &dns.PTR{
            Hdr: dns.RR_Header{
                Name: state.QName(),
                Rrtype: dns.TypePTR,
                Class: state.QClass(),
            },
            Ptr: dns.Fqdn(value),
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

func responseWithCode(r *dns.Msg, code int) *dns.Msg {
    m := new(dns.Msg)
    m.SetReply(r)
    m.Authoritative = true
    //m.Compress = true
    m.RecursionAvailable = true
    m.Rcode = code

    return m
}
