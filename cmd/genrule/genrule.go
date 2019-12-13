package main

import (
    "fmt"
    "time"

    "github.com/ajruckman/ContraCore/internal/rulegen"
)

var (
    urls = []string{
        "http://localhost/contradomain/spark",
        "http://localhost/contradomain/bluGo",
        "http://localhost/contradomain/blu",
        "http://localhost/contradomain/basic",
        "http://localhost/contradomain/ultimate",
        "http://localhost/contradomain/unified",

        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/spark/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/bluGo/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/blu/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/basic/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/ultimate/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/unified/formats/domains.txt",
        //"https://someonewhocares.org/hosts/hosts",
        //"https://gist.githubusercontent.com/angristan/20a398983c5b1daa9c13a1cbadb78fd6/raw/58d54b172b664ee5a0b53bb2e25c391433f2cc7a/hosts",
        //"https://www.encrypt-the-planet.com/downloads/hosts",

        //"http://localhost/contradomain/unified",
        //"http://localhost/contradomain/someonewhocares",
        //"http://localhost/contradomain/win",
    }
)

func bench() {
    checks := 10
    //for _, url := range urls {

        var totMilliseconds int64

        for i := 0; i < checks+1; i++ {
            begin := time.Now()
            res, total := rulegen.ProcessFromURLsPointers(urls)
            end := time.Now()

            if i != 0 {
                totMilliseconds += end.Sub(begin).Milliseconds()

                kept := len(res)
                ratio := float64(kept) / float64(total)

                fmt.Println(ratio, kept, total)
            }
        }

        fmt.Println(float64(totMilliseconds)/float64(checks), urls)
    //}
}

func load() {
    begin := time.Now()
    rules, total := rulegen.ProcessFromURLsPointers(urls)
    end := time.Now()

    kept := len(rules)
    ratio := float64(kept) / float64(total)

    fmt.Println(ratio, kept, total, end.Sub(begin))
}

func main() {
    //bench()
    load()
    //began := time.Now()
    //res, _ := rulegen.ProcessFromURLs(urls)
    //fmt.Println(len(res), time.Since(began))

    //var b testing.BenchmarkResult
    //
    //minP := 4
    //maxP := 4
    //minC := 10000
    //maxC := 10000
    //checks := int64(5)
    //
    ////fmt.Printf("%-10v", "")
    //for c := minC; c < maxC; c += 1000 {
    //    fmt.Print(fmt.Sprintf("%-10d", c))
    //}
    ////fmt.Println()
    //
    //for p := minP; p <= maxP; p++ {
    //    rulegen.MaxPar = p
    //
    //    //fmt.Printf("%-10d", p)
    //
    //    for c := minC; c <= maxC; c += 1000 {
    //        rulegen.ChunkSize = c
    //
    //        var total int64
    //        for check := int64(0); check <= checks+1; check++ {
    //            b = testing.Benchmark(BenchmarkProcessFromURLsWithPointers)
    //            if check != 0 {
    //                total += b.T.Milliseconds()
    //            }
    //
    //        }
    //
    //        fmt.Println(total/checks, checks)
    //
    //        //fmt.Printf("%-10d", total/checks)
    //    }
    //
    //    //fmt.Println()
    //}
}
