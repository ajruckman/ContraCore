package serve

import (
    "net"

    "github.com/coredns/coredns/request"
    "github.com/miekg/dns"
)

func respondWithDynamicDNS(w dns.ResponseWriter, r *dns.Msg, qu dns.Question, nm string, iw *ResponseInterceptor) (ret bool, rcode int, err error) {
    if v, ok := dhcpCacheMap[nm]; ok {
        var m *dns.Msg

        for _, lease := range v {
            if lease.IP.To4() != nil {
                if qu.Qtype == dns.TypeA {
                    m = GenResponse(w, r, qu.Qtype, lease.IP.To4().String())
                    err = iw.WriteMsg(m)
                    clog.Info("lease IP is IPv4, question is A")
                    return true, dns.RcodeSuccess, err
                } else {
                    clog.Info("lease IP is IPv4, question is AAAA")
                    continue
                }

            } else if lease.IP.To16() != nil {
                if qu.Qtype == dns.TypeAAAA {
                    m = GenResponse(w, r, qu.Qtype, lease.IP.To16().String())
                    err = iw.WriteMsg(m)
                    clog.Info("lease IP is IPv6, question is AAAA")
                    return true, dns.RcodeSuccess, err
                } else {
                    clog.Info("lease IP is IPv6, question is A")
                    continue
                }
            }
        }

        err = RespondWithCode(w, r, dns.RcodeServerFailure)
        return true, dns.RcodeServerFailure, err
    }

    //clog.Info("lease does not exist: ", nm)

    return
}

func GenResponse(w dns.ResponseWriter, r *dns.Msg, qtype uint16, value string) *dns.Msg {
    m := new(dns.Msg)
    m.SetReply(r)
    m.Authoritative = true
    m.RecursionAvailable = true
    m.Rcode = dns.RcodeSuccess

    var rr dns.RR
    state := request.Request{W: w, Req: r}

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
                Rrtype: dns.TypeCNAME,
                Class:  state.QClass(),
            },
        }
    }
    m.Answer = append(m.Answer, rr)

    return m
}

func RespondWithCode(w dns.ResponseWriter, r *dns.Msg, code int) error {
    m := new(dns.Msg)
    m.SetReply(r)
    m.Authoritative = true
    //m.Compress = true
    m.RecursionAvailable = true
    m.Rcode = code

    return w.WriteMsg(m)
}
