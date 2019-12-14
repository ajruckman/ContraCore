package serve

import (
    "context"
    "errors"
    "strings"
    "time"

    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema"
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

        qu:       r.Question[0],
        domain:   rt(r.Question[0].Name),
        client:   strings.Split(w.RemoteAddr().String(), ":")[0],
        received: time.Now(),
    }

    //store := schema.Log{
    //    Client:       ip,
    //    Question:     nm,
    //    QuestionType: dns.TypeToString[qu.Qtype],
    //    Stored:       time.Now(),
    //}

    //storedQueries.Store(r.Id, store)

    clog.Info(q.client, " -> ", r.Id, " ", dns.TypeToString[q.qu.Qtype])

    if strings.Count(q.domain, ".") == 0 {
        if ret, rcode, err := respondWithDynamicDNS(&q); ret {
            return rcode, err
        }
    }

    if ret, rcode, err := respondWithBlock(&q); ret {
        return rcode, err
    }

    //clog.Infof("NM: %s | dn: %v", nm, config.Config.DomainNeeded)

    if config.Config.DomainNeeded && strings.Count(q.domain, ".") == 0 {
        q.action = "restrict"
        clog.Infof("DomainNeeded is true and question '%s' does not contain any periods; returning RcodeServerFailure", q.domain)
        m := RespondWithCode(q.r, dns.RcodeServerFailure)
        err := q.Respond(m)
        return dns.RcodeServerFailure, err
    }

    q.action = "pass"
    return plugin.NextOrFailure(name, next, ctx, q, r)
}

type queryContext struct {
    dns.ResponseWriter
    r *dns.Msg

    qu       dns.Question
    domain   string
    client   string
    received time.Time

    action string
    //answers []string
}

func (q *queryContext) Respond(res *dns.Msg) (err error) {
    var answers []string
    for _, v := range res.Answer {
        answers = append(answers, rrToString(v))
    }

    logChannel <- schema.Log{
        Client:       q.client,
        Question:     q.domain,
        QuestionType: dns.TypeToString[q.qu.Qtype],
        Action:       q.action,
        Answers:      answers,

        QueryID:  q.r.Id,
        Duration: time.Now().Sub(q.received),
    }

    err = q.ResponseWriter.WriteMsg(res)
    return
}

func (q queryContext) WriteMsg(r *dns.Msg) error {
    return q.Respond(r)
}

func init() {
    go logWorker()
    cacheDHCP()
    cacheRules()
}

var logChannel = make(chan schema.Log)

func logWorker() {
    for v := range logChannel {
        err := db.Log(v)
        if err != nil {
            clog.Warning("could not insert log for query '" + v.Question + "'")
            clog.Warning(err.Error())
        }

        clog.Info(v.Client, " <- ", v.QueryID, " ", v.QuestionType, " ", v.Duration)
    }
}

//var storedQueries sync.Map
//
//type ResponseInterceptor struct {
//    dns.ResponseWriter
//}
//
//func NewResponseInterceptor(w dns.ResponseWriter) *ResponseInterceptor {
//    return &ResponseInterceptor{ResponseWriter: w}
//}
//
//func (ri *ResponseInterceptor) WriteMsg(res *dns.Msg) error {
//    var (
//        loaded interface{}
//        ok     bool
//        stored schema.Log
//    )
//
//    if loaded, ok = storedQueries.Load(res.Id); !ok {
//        clog.Error("Unmatched query ID ", res.Id)
//        storedQueries.Range(func(key interface{}, value interface{}) bool {
//            clog.Debug("    ", key, " -> ", value)
//
//            return true
//        })
//        goto done
//    }
//
//    stored = loaded.(schema.Log)
//
//    if time.Now().Sub(stored.Stored) > (time.Second * 3) {
//        clog.Error("Stale query ID ", res.Id)
//        storedQueries.Delete(res.Id)
//        goto done
//    }
//
//    clog.Info(stored.Client, " <- ", res.Id, " ", stored.QuestionType)
//
//    for _, v := range res.Answer {
//        stored.Answers = append(stored.Answers, rrToString(v))
//    }
//
//    logChannel <- stored
//
//done:
//    return ri.ResponseWriter.WriteMsg(res)
//}
