package serve

import (
    "fmt"
    "strings"

    "github.com/miekg/dns"
)

// coredns/plugin/test/helpers.go
func read(val dns.RR) string {
    var res string

    switch x := val.(type) {
    case *dns.SRV:
        res = fmt.Sprintf("%d|%d|%d|%s", x.Priority, x.Weight, x.Port, x.Target)

    case *dns.RRSIG:
        res = fmt.Sprintf("%d|%d|%s", x.TypeCovered, x.Labels, x.SignerName)

    case *dns.NSEC:
        res = x.NextDomain

    case *dns.A:
        res = rt(x.A.String())

    case *dns.AAAA:
        res = rt(x.AAAA.String())

    case *dns.TXT:
        res = strings.Join(x.Txt, "|")

    case *dns.HINFO:
        res = fmt.Sprintf("%s|%s", x.Cpu, x.Os)

    case *dns.SOA:
        res = x.Ns

    case *dns.PTR:
        res = rt(x.Ptr)

    case *dns.CNAME:
        res = rt(x.Target)

    case *dns.MX:
        res = fmt.Sprintf("%s|%d", x.Mx, x.Preference)

    case *dns.NS:
        res = x.Ns

    case *dns.OPT:
        res = fmt.Sprintf("%d|%t", x.UDPSize(), x.Do())
    }

    return res
}
