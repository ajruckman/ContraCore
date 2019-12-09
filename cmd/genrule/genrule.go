package main

import (
    "fmt"
    "testing"

    "github.com/ajruckman/ContraCore/internal/rulegen"
)

var (
    urls = []string{
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/spark/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/blu/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/basic/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/ultimate/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/unified/formats/domains.txt",
        //"https://someonewhocares.org/hosts/hosts",
        //"https://gist.githubusercontent.com/angristan/20a398983c5b1daa9c13a1cbadb78fd6/raw/58d54b172b664ee5a0b53bb2e25c391433f2cc7a/hosts",
        //"https://www.encrypt-the-planet.com/downloads/hosts",

        "http://localhost/contradomain/unified",
        "http://localhost/contradomain/someonewhocares",
        "http://localhost/contradomain/win",
    }
    contents []string
)

func main() {
    //began := time.Now()
    //res, _ := rulegen.ProcessFromURLs(urls)
    //fmt.Println(len(res), time.Since(began))

    var b testing.BenchmarkResult

    minP := 1
    maxP := 10
    minC := 1000
    maxC := 30000
    checks := int64(5)

    fmt.Printf("%-10v", "")
    for c := minC; c < maxC; c += 1000 {
        fmt.Print(fmt.Sprintf("%-10d", c))
    }
    fmt.Println()

    for p := minP; p < maxP; p++ {
        rulegen.MaxPar = p

        fmt.Printf("%-10d", p)

        for c := minC; c < maxC; c += 1000 {
            rulegen.ChunkSize = c

            var total int64
            for check := int64(0); check < checks; check++ {
                b = testing.Benchmark(BenchmarkProcessFromURLsWithPointers)
                total += b.T.Milliseconds()
            }

            fmt.Printf("%-10d", total/checks)
        }

        fmt.Println()
    }
}

func difference(a, b []string) []string {
    mb := make(map[string]struct{}, len(b))
    for _, x := range b {
        mb[x] = struct{}{}
    }
    var diff []string
    for _, x := range a {
        if _, found := mb[x]; !found {
            diff = append(diff, x)
        }
    }
    return diff
}
