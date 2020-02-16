package main

import (
    "fmt"
    "time"

    "github.com/ajruckman/xlib"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/db/contradb"
    "github.com/ajruckman/ContraCore/internal/db/contradb/dbschema"
)

/*

Classes 0-2: 713500   / 938306   | Average request time (ms): 0.449        | diversificando.org

*/

func main() {
    c := new(dns.Client)

    config.ContraDBURL = "postgres://contradbmgr:contradbmgr@127.0.0.1/contradb"
    contradb.Setup()

    var rules []dbschema.Blacklist
    err := contradb.Select(&rules, `SELECT id, pattern, class, COALESCE(domain, '') AS domain, COALESCE(tld, '') AS tld, COALESCE(sld, '') AS sld FROM rule ORDER BY random()`)
    xlib.Err(err)

    var total int64 = 0
    var count = 0

    last := "sub."

    for i, rule := range rules {
        m := new(dns.Msg)
        m.SetQuestion(dns.Fqdn(last+rule.Domain), dns.TypeA)
        count++

        began := time.Now()
        r, _, err := c.Exchange(m, "127.0.0.1:5300")
        xlib.Err(err)
        dur := time.Since(began)

        total += dur.Milliseconds()

        _ = r

        last = rule.Domain + "."

        //if r.Rcode == dns.RcodeSuccess {
        //    panic(errors.New(fmt.Sprintf("%d", r.Rcode)))
        //}

        if i%1 == 0 {
            fmt.Printf("%-8d / %-8d | Average request time (ms): %-8s     | %s\n", i, len(rules), fmt.Sprintf("%.3f", float64(total)/float64(count)), rule.Domain)
        }

        //for _, v := range r.Answer {
        //    if !strings.HasSuffix(v.String(), "	0.0.0.0") {
        //        panic(errors.New("Query was unblocked: " + v.Header().Name))
        //    }
        //}
    }

    fmt.Printf("\n---")
    fmt.Printf("Average request time (ms): %.3f\n", float64(total)/float64(count))
}
