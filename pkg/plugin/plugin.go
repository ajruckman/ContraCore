package plugin

import (
    "context"
    "fmt"
    "strings"

    "github.com/caddyserver/caddy"
    "github.com/coredns/coredns/core/dnsserver"
    "github.com/coredns/coredns/plugin"
    clog "github.com/coredns/coredns/plugin/pkg/log"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/dnsvar"
)

var (
    Log = clog.NewWithPlugin("contradomain")
)

func init() {
    Log.Info("--------------------------------------------------- v11")

    caddy.RegisterPlugin("contracore", caddy.Plugin{
        ServerType: "dns",
        Action:     setup,
    })
}

func setup(c *caddy.Controller) error {
    dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
        return ContraCore{Next: next}
    })

    return nil
}

type ContraCore struct {
    Next plugin.Handler
}

func (e ContraCore) Name() string {
    return "contracore"
}

func (e ContraCore) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
    ip := strings.Split(w.RemoteAddr().String(), ":")[0]

    var queries []string

    for _, v := range r.Question {
        qtype := dnsvar.TypeMap[v.Qtype]
        qname := strings.TrimSuffix(v.Name, ".")

        queries = append(queries, fmt.Sprintf("[%s] %s", qtype, qname))
    }

    Log.Info(ip, " -> ", strings.Join(queries, ", "))

    iw := NewResponseInterceptor(w)

    return plugin.NextOrFailure(e.Name(), e.Next, ctx, iw, r)
}

type ResponseInterceptor struct {
    dns.ResponseWriter
}

func NewResponseInterceptor(w dns.ResponseWriter) *ResponseInterceptor {
    return &ResponseInterceptor{ResponseWriter: w}
}

func (r *ResponseInterceptor) WriteMsg(res *dns.Msg) error {
    Log.Info("--- ", res.Id)
    return r.ResponseWriter.WriteMsg(res)
}
