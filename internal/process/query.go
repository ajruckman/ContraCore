package process

import (
    "fmt"
    "net"
    "strings"
    "time"

    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/functions"
    "github.com/ajruckman/ContraCore/internal/log"
    "github.com/ajruckman/ContraCore/internal/schema"
)

type action uint8

const (
    ActionAllow action = iota
    ActionRestrict
    ActionPass
    ActionBlock
)

type queryContext struct {
    dns.ResponseWriter

    r         *dns.Msg
    _question dns.Question

    _domain string
    _client string

    received time.Time
    mac      *net.HardwareAddr
    hostname *string
    vendor   *string

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
        Time:           q.received,
        Client:         q._client,
        Question:       q._domain,
        QuestionType:   dns.TypeToString[q._question.Qtype],
        Action:         "-test-",
        Answers:        q.answers,
        ClientMAC:      q.mac,
        ClientHostname: q.hostname,
        ClientVendor:   q.vendor,
        QueryID:        uint16(3),
        Duration:       time.Now().Sub(q.received),
    })

    err = q.ResponseWriter.WriteMsg(res)
    return
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