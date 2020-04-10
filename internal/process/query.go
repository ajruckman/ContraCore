package process

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"

	"github.com/ajruckman/ContraCore/internal/db/contralog/dbschema"
	"github.com/ajruckman/ContraCore/internal/functions"
	"github.com/ajruckman/ContraCore/internal/log"
	"github.com/ajruckman/ContraCore/internal/schema"
)

type action string

const (
	ActionNotBlacklisted = "pass.notblacklisted"
	ActionWhitelisted    = "pass.whitelisted"
	ActionBlock          = "block.blacklisted"
	ActionDomainNeeded   = "block.domainneeded"
	ActionDDNSHostname   = "respond.ddnshostname"
	ActionDDNSPTR        = "respond.ddnsptr"
)

type queryContext struct {
	dns.ResponseWriter

	r         *dns.Msg
	_question dns.Question

	_domain              string
	_suffix              *string
	_matchedSearchDomain *string
	_client              net.IP

	received time.Time
	mac      *string
	hostname *string
	vendor   *string

	action  action
	answers []string
}

func (q *queryContext) respond(res *dns.Msg) (err error) {
	var answers []string
	for _, v := range res.Answer {
		answers = append(answers, rrToString(v))
	}
	q.answers = answers

	//LogChannel <- *q

	log.Query(schema.Log{
		Log: dbschema.Log{
			Time:           q.received,
			Client:         q._client.String(),
			Question:       q._domain,
			QuestionType:   dns.TypeToString[q._question.Qtype],
			Action:         string(q.action),
			Answers:        q.answers,
			ClientMAC:      q.mac,
			ClientHostname: q.hostname,
			ClientVendor:   q.vendor,
			QueryID:        uint16(rand.Intn(65536)),
		},
		Duration: time.Now().Sub(q.received),
	})

	err = q.ResponseWriter.WriteMsg(res)
	return
}

func (q queryContext) String() string {
	return fmt.Sprintf("{%s | %d -> [%s] %s}", q._client.String(), q.r.Id, dns.TypeToString[q._question.Qtype], q._domain)
}

func (q queryContext) WriteMsg(r *dns.Msg) error {
	return q.respond(r)
}

// coredns/plugin/test/helpers.go
func rrToString(val dns.RR) string {
	var res string

	switch x := val.(type) {
	case *dns.SRV:
		res = fmt.Sprintf("%d|%d|%d|%s", x.Priority, x.Weight, x.Port, x.Target)

	case *dns.RRSIG:
		res = fmt.Sprintf("%d|%d|%s", x.TypeCovered, x.Labels, x.SignerName)

	case *dns.NSEC:
		res = x.NextDomain

	case *dns.A:
		res = functions.RT(x.A.String())

	case *dns.AAAA:
		res = functions.RT(x.AAAA.String())

	case *dns.TXT:
		res = strings.Join(x.Txt, "|")

	case *dns.HINFO:
		res = fmt.Sprintf("%s|%s", x.Cpu, x.Os)

	case *dns.SOA:
		res = x.Ns

	case *dns.PTR:
		res = functions.RT(x.Ptr)

	case *dns.CNAME:
		res = functions.RT(x.Target)

	case *dns.MX:
		res = fmt.Sprintf("%s|%d", x.Mx, x.Preference)

	case *dns.NS:
		res = x.Ns

	case *dns.OPT:
		res = fmt.Sprintf("%d|%t", x.UDPSize(), x.Do())
	}

	return res
}
