package config

import (
    "errors"

    "github.com/caddyserver/caddy"
)

func ParseCorefile(c *caddy.Controller) {
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
                    ContraDBURL = c.Val()

                case "ContraLogURL":
                    c.Next()
                    ContraLogURL = c.Val()

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

var fields = [...]string{"ContraDBURL", "ContraLogURL"}
