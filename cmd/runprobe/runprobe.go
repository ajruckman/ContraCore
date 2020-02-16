package main

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
)

func main() {
	c := new(dns.Client)
	c.Timeout = 1 * time.Second

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("!runprobe"), dns.TypeA)

	t := time.NewTicker(time.Second * 5)

	for range t.C {
		fmt.Println("TICK")

		r, _, err := c.Exchange(m, "127.0.0.1:5300")

		if err == nil && r.Rcode == 15 {
			if !lastUp {
				onServerUp()
			}

		} else {
			if lastUp {
				onServerDown()
			}
		}
	}
}

var lastUp = false

func onServerUp() {
	fmt.Println("onServerUp()")

	lastUp = true
}

func onServerDown() {
	fmt.Println("onServerDown()")

	lastUp = false
}
