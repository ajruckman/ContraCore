package serve

import (
    "context"
    "errors"
    "strings"
    "sync"

    . "github.com/ajruckman/xlib"
    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema"
)

func DNS(name string, next plugin.Handler, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
    ip := strings.Split(w.RemoteAddr().String(), ":")[0]

    // https://stackoverflow.com/a/4083071/9911189
    if len(r.Question) != 1 {
        return 0, errors.New("this should never happen")
    }

    store := schema.Log{
        Client:       ip,
        Question:     rt(r.Question[0].Name),
        QuestionType: dns.TypeToString[r.Question[0].Qtype],
    }

    storedQueries.Store(r.Id, store)

    clog.Info(ip, " -> ", r.Id, " ", dns.TypeToString[r.Question[0].Qtype])

    iw := NewResponseInterceptor(w)

    return plugin.NextOrFailure(name, next, ctx, iw, r)
}

func init() {
    go logWorker()
}

var logChannel = make(chan schema.Log)

func logWorker() {
    for v := range logChannel {
        err := db.Log(v)
        Err(err)
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
        goto done
    }

    stored = loaded.(schema.Log)

    clog.Info(stored.Client, " <- ", r.Id, " ", stored.QuestionType)

    for _, v := range r.Answer {
        stored.Answers = append(stored.Answers, read(v))
    }

    logChannel <- stored

    storedQueries.Delete(r.Id)

done:
    return ri.ResponseWriter.WriteMsg(r)
}
