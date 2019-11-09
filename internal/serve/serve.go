package serve

import (
    "context"
    "errors"
    "fmt"
    "strings"
    "sync"

    . "github.com/ajruckman/xlib"
    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema"
    "github.com/ajruckman/ContraCore/internal/shared"
)

func init() {
    go logWorker()
}

func DNS(name string, next plugin.Handler, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
    ip := strings.Split(w.RemoteAddr().String(), ":")[0]

    // https://stackoverflow.com/a/4083071/9911189
    store := schema.Log{
        Client: ip,
    }

    if len(r.Question) > 1 {
        return 0, errors.New("this should never happen")
    }

    for _, v := range r.Question {
        store.Question = strings.TrimSuffix(v.Name, ".")
        store.QuestionType = dns.TypeToString[v.Qtype]
    }

    storedQueryLock.Lock()
    storedQueries[r.Id] = store
    storedQueryLock.Unlock()

    shared.Log.Info(ip, " -> ", r.Id, " ", dns.TypeToString[r.Question[0].Qtype])

    iw := NewResponseInterceptor(w)

    return plugin.NextOrFailure(name, next, ctx, iw, r)
}

var logChannel = make(chan schema.Log)

func logWorker() {
    for v := range logChannel {
        err := db.Log(v)
        Err(err)
    }
}

var storedQueryLock = sync.Mutex{}
var storedQueries = map[uint16]schema.Log{}

type ResponseInterceptor struct {
    dns.ResponseWriter
}

func NewResponseInterceptor(w dns.ResponseWriter) *ResponseInterceptor {
    return &ResponseInterceptor{ResponseWriter: w}
}

func (ri *ResponseInterceptor) WriteMsg(r *dns.Msg) error {
    var stored schema.Log
    var ok bool

    storedQueryLock.Lock()
    if stored, ok = storedQueries[r.Id]; !ok {
        storedQueryLock.Unlock()
        shared.Log.Error("Unmatched query ID ", r.Id)
        goto done
    }
    storedQueryLock.Unlock()
    shared.Log.Info(stored.Client, " <- ", r.Id, " ", stored.QuestionType)

    for _, v := range r.Answer {
        switch v.(type) {
        case *dns.A:
            stored.Answers = append(stored.Answers, v.(*dns.A).A.String())

        case *dns.AAAA:
            stored.Answers = append(stored.Answers, v.(*dns.AAAA).AAAA.String())

        case *dns.CNAME:
            stored.Answers = append(stored.Answers, v.(*dns.CNAME).Target)

        case *dns.SRV:
            stored.Answers = append(stored.Answers, v.(*dns.SRV).Target)

        case *dns.PTR:
            stored.Answers = append(stored.Answers, v.(*dns.PTR).Ptr)

        case *dns.SOA:
            m := v.(*dns.SOA)
            s := fmt.Sprintf("%s|%s|%d|%d|%d|%d|%d", m.Ns, m.Mbox, m.Serial, m.Refresh, m.Retry, m.Expire, m.Minttl)

            stored.Answers = append(stored.Answers, s)
        }
    }

    logChannel <- stored

    storedQueryLock.Lock()
    delete(storedQueries, r.Id)
    storedQueryLock.Unlock()

done:
    return ri.ResponseWriter.WriteMsg(r)
}
