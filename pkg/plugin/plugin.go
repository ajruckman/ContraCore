package plugin

import (
    "context"

    "github.com/caddyserver/caddy"
    "github.com/coredns/coredns/core/dnsserver"
    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal"
    "github.com/ajruckman/ContraCore/internal/process"
)

func init() {
    caddy.RegisterPlugin("contracore", caddy.Plugin{
        ServerType: "dns",
        Action:     setup,
    })
}

func setup(c *caddy.Controller) error {
    internal.Setup(c)

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
    return process.DNS(e.Name(), e.Next, ctx, w, r)
}
