package plugin

import (
    "context"
    "errors"

    "github.com/caddyserver/caddy"
    "github.com/coredns/coredns/core/dnsserver"
    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/cmd/vars"
    "github.com/ajruckman/ContraCore/internal/db"
    _ "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/provision"
    "github.com/ajruckman/ContraCore/internal/serve"
)

func init() {
    caddy.RegisterPlugin("contracore", caddy.Plugin{
        ServerType: "dns",
        Action:     setup,
    })
}

func setup(c *caddy.Controller) error {
    parseConfig(c)

    db.Setup()
    provision.Setup()

    dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
        return ContraCore{Next: next}
    })

    return nil
}

var (
    fields = []string{"ContraDBURL", "ContraLogURL"}
)

func parseConfig(c *caddy.Controller) {
    c.Next()

    if c.Val() != "contracore" {
        panic(errors.New("unexpected plugin name '" + c.Val() + "'"))
    }

    c.Next()

    if c.Val() != "{" {
        panic(errors.New("expected opening brace"))
    }

    for c.Next() {
        if c.Val() == "}" {
            break
        }

        for _, field := range fields {
            if field == c.Val() {

                switch c.Val() {
                case "ContraDBURL":
                    c.Next()
                    vars.ContraDBURL = c.Val()

                case "ContraLogURL":
                    c.Next()
                    vars.ContraLogURL = c.Val()

                default:
                    panic(errors.New("unhandled field '" + c.Val() + "'"))
                }

                goto next
            }
        }

        panic(errors.New("unexpected token '" + c.Val() + "'"))

    next:
    }
}

type ContraCore struct {
    Next plugin.Handler
}

func (e ContraCore) Name() string {
    return "contracore"
}

func (e ContraCore) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
    return serve.DNS(e.Name(), e.Next, ctx, w, r)
}
