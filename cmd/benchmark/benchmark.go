package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    _ "net/http/pprof"
    "testing"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/rulegen"
)

var (
    urls = []string{
        "https://raw.githubusercontent.com/EnergizedProtection/block/master/unified/formats/domains.txt",
        "https://someonewhocares.org/hosts/hosts",
        "https://gist.githubusercontent.com/angristan/20a398983c5b1daa9c13a1cbadb78fd6/raw/58d54b172b664ee5a0b53bb2e25c391433f2cc7a/hosts",
        //"https://www.encrypt-the-planet.com/downloads/hosts",
    }
    contents [][]byte

    contentTotals int64
)

func init() {
    for _, url := range urls {
        fmt.Print("Reading ", url, "... ")
        resp, err := http.Get(url)
        Err(err)
        fmt.Print(resp.StatusCode, " ")
        read, err := ioutil.ReadAll(resp.Body)
        Err(err)
        contents = append(contents, read)
        fmt.Println("done")

        contentTotals += int64(len(read))
    }
    fmt.Println()
}

func main() {
    testing.Benchmark(benchmarkRuleGenV2) // Seed

    rV2 := testing.Benchmark(benchmarkRuleGenV2)
    fmt.Println("BlockV2:", rV2.T.Milliseconds(), rV2.String()+" -> "+rV2.MemString())

    rV3 := testing.Benchmark(benchmarkRuleGenV3)
    fmt.Println("BlockV3:", rV3.T.Milliseconds(), rV3.String()+" -> "+rV3.MemString())

    rV4 := testing.Benchmark(benchmarkRuleGenV4)
    fmt.Println("BlockV4:", rV4.T.Milliseconds(), rV4.String()+" -> "+rV4.MemString())

    rV5 := testing.Benchmark(benchmarkRuleGenV5)
    fmt.Println("BlockV5:", rV5.T.Milliseconds(), rV5.String()+" -> "+rV5.MemString())

    rV6 := testing.Benchmark(benchmarkRuleGenV6)
    fmt.Println("BlockV6:", rV6.T.Milliseconds(), rV6.String()+" -> "+rV6.MemString())

    //err := http.ListenAndServe(":8080", nil)
    //Err(err)
}

func benchmarkRuleGenV2(b *testing.B) {
    benchmarkRuleGen(rulegen.BlockV2, b)
}

func benchmarkRuleGenV3(b *testing.B) {
    benchmarkRuleGen(rulegen.BlockV3, b)
}

func benchmarkRuleGenV4(b *testing.B) {
    benchmarkRuleGen(rulegen.BlockV4, b)
}

func benchmarkRuleGenV5(b *testing.B) {
    benchmarkRuleGen(rulegen.BlockV5, b)
}

func benchmarkRuleGenV6(b *testing.B) {
    benchmarkRuleGen(rulegen.BlockV6, b)
}

const numBench = 15

func benchmarkRuleGen(evaluator func(*rulegen.Node, string, []string), b *testing.B) {
    b.ReportAllocs()
    b.SetBytes(contentTotals * numBench)

    //for n := 0; n < b.N; n++ {
    for n := 0; n < numBench; n++ {
        rulegen.ReadDomainScanners(evaluator, contents...)
    }
}
