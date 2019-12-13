package serve

import (
    "context"
    "errors"
    "strings"
    "sync"
    "time"

    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema"
)

func DNS(name string, next plugin.Handler, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

    // https://stackoverflow.com/a/4083071/9911189
    if len(r.Question) != 1 {
        return 0, errors.New("this should never happen")
    }

    ip := strings.Split(w.RemoteAddr().String(), ":")[0]
    qu := r.Question[0]
    nm := rt(qu.Name)

    store := schema.Log{
        Client:       ip,
        Question:     rt(qu.Name),
        QuestionType: dns.TypeToString[qu.Qtype],
        Stored:       time.Now(),
    }

    storedQueries.Store(r.Id, store)

    clog.Info(ip, " -> ", r.Id, " ", dns.TypeToString[qu.Qtype])

    iw := NewResponseInterceptor(w)

    if ret, rcode, err := respondWithDynamicDNS(w, r, qu, nm, iw); ret {
        return rcode, err
    }

    //clog.Infof("NM: %s | dn: %v", nm, config.Config.DomainNeeded)

    if config.Config.DomainNeeded && strings.Count(nm, ".") == 0 {
        clog.Infof("DomainNeeded is true and question '%s' does not contain any periods; returning RcodeServerFailure", nm)
        err := RespondWithCode(w, r, dns.RcodeServerFailure)
        return dns.RcodeServerFailure, err
    }

    return plugin.NextOrFailure(name, next, ctx, iw, r)
}

func init() {
    go logWorker()
    dhcpCache()
}

var logChannel = make(chan schema.Log)

func logWorker() {
    for v := range logChannel {
        err := db.Log(v)
        if err != nil {
            clog.Warning("could not insert log for query '" + v.Question + "'")
            clog.Warning(err.Error())
        }
    }
}

var storedQueries sync.Map

type ResponseInterceptor struct {
    dns.ResponseWriter
}

func NewResponseInterceptor(w dns.ResponseWriter) *ResponseInterceptor {
    return &ResponseInterceptor{ResponseWriter: w}
}

func (ri *ResponseInterceptor) WriteMsg(r *dns.Msg) error {
    var (
        loaded interface{}
        ok     bool
        stored schema.Log
    )

    if loaded, ok = storedQueries.Load(r.Id); !ok {
        clog.Error("Unmatched query ID ", r.Id)
        storedQueries.Range(func(key interface{}, value interface{}) bool {
            clog.Debug("    ", key, " -> ", value)

            return true
        })
        goto done
    }

    stored = loaded.(schema.Log)

    if time.Now().Sub(stored.Stored) > (time.Second * 3) {
        clog.Error("Stale query ID ", r.Id)
        storedQueries.Delete(r.Id)
        goto done
    }

    clog.Info(stored.Client, " <- ", r.Id, " ", stored.QuestionType)

    for _, v := range r.Answer {
        stored.Answers = append(stored.Answers, rrToString(v))
    }

    logChannel <- stored

done:
    return ri.ResponseWriter.WriteMsg(r)
}
