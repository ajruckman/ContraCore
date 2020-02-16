package functions

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

// Reads a CoreDNS resource record and returns its string representation.
// From: coredns/plugin/test/helpers.go
func ReadRR(val dns.RR) string {
	var res string

	switch x := val.(type) {
	case *dns.SRV:
		res = fmt.Sprintf("%d|%d|%d|%s", x.Priority, x.Weight, x.Port, x.Target)

	case *dns.RRSIG:
		res = fmt.Sprintf("%d|%d|%s", x.TypeCovered, x.Labels, x.SignerName)

	case *dns.NSEC:
		res = x.NextDomain

	case *dns.A:
		res = RT(x.A.String())

	case *dns.AAAA:
		res = RT(x.AAAA.String())

	case *dns.TXT:
		res = strings.Join(x.Txt, "|")

	case *dns.HINFO:
		res = fmt.Sprintf("%s|%s", x.Cpu, x.Os)

	case *dns.SOA:
		res = x.Ns

	case *dns.PTR:
		res = RT(x.Ptr)

	case *dns.CNAME:
		res = RT(x.Target)

	case *dns.MX:
		res = fmt.Sprintf("%s|%d", x.Mx, x.Preference)

	case *dns.NS:
		res = x.Ns

	case *dns.OPT:
		res = fmt.Sprintf("%d|%t", x.UDPSize(), x.Do())
	}

	return res
}
