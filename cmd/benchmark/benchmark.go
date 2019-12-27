package main

import (
    "bufio"
    "errors"
    "fmt"
    "net/http"
    _ "net/http/pprof"
    "runtime"
    "strings"
    "testing"
    "time"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/rulegen"
)

var (
    urls = []string{
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/spark/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/blu/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/basic/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/ultimate/formats/domains.txt",
        "https://raw.githubusercontent.com/EnergizedProtection/block/master/unified/formats/domains.txt",
        //"https://someonewhocares.org/hosts/hosts",
        //"https://gist.githubusercontent.com/angristan/20a398983c5b1daa9c13a1cbadb78fd6/raw/58d54b172b664ee5a0b53bb2e25c391433f2cc7a/hosts",
        //"https://www.encrypt-the-planet.com/downloads/hosts",
    }
    contents []string

    contentTotals int64
)

func init() {
    for _, url := range urls {
        fmt.Print("Reading ", url, "... ")
        resp, err := http.Get(url)
        Err(err)
        fmt.Print(resp.StatusCode, " ")

        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            contentTotals += int64(len(scanner.Bytes()))
            contents = append(contents, scanner.Text())
        }
        fmt.Println("done")

        contents = append(contents, strings.Repeat("-", 25))
    }

    //f, err := os.Create("all.txt")
    //Err(err)
    //
    //_, err = f.WriteString(strings.Join(contents, "\n"))
    //Err(err)

    fmt.Println()
}

func main() {

    began := time.Now()
    res, totalValid := rulegen.GenFromURLs(urls)

    fmt.Println()
    fmt.Println("Time: ", time.Since(began))
    fmt.Println("Total:", totalValid)
    fmt.Println("Kept: ", len(res))
    fmt.Println("Ratio:", float64(len(res))/float64(totalValid))

    return

    //testing.Benchmark(benchmarkRuleGenV2) // Seed

    //rV2 := testing.Benchmark(benchmarkRuleGenV2)
    //fmt.Println("BlockV2:", rV2.T.Milliseconds(), rV2.String()+" -> "+rV2.MemString())
    //
    //rV3 := testing.Benchmark(benchmarkRuleGenV3)
    //fmt.Println("BlockV3:", rV3.T.Milliseconds(), rV3.String()+" -> "+rV3.MemString())

    for i := 0; i < 50; i++ {

        {
            rV4 := testing.Benchmark(benchmarkRuleGenV4)
            fmt.Println("BlockV4:", rV4.T.Milliseconds(), rV4.String()+" -> "+rV4.MemString())
        }

        //runtime.GC()
        //time.Sleep(time.Second)
        //
        //{
        //    rV5 := testing.Benchmark(benchmarkRuleGenV5)
        //    fmt.Println("BlockV5:", rV5.T.Milliseconds(), rV5.String()+" -> "+rV5.MemString())
        //}

        runtime.GC()
        time.Sleep(time.Second)

        {
            rV4 := testing.Benchmark(benchmarkRuleGenV4)
            fmt.Println("BlockV4:", rV4.T.Milliseconds(), rV4.String()+" -> "+rV4.MemString())
        }

        //runtime.GC()
        //time.Sleep(time.Second)
        //
        //{
        //    rV5 := testing.Benchmark(benchmarkRuleGenV5)
        //    fmt.Println("BlockV5:", rV5.T.Milliseconds(), rV5.String()+" -> "+rV5.MemString())
        //}
        //
        //fmt.Println()
    }

    //err := http.ListenAndServe(":8080", nil)
    //Err(err)
}

//func benchmarkRuleGenV2(b *testing.B) {
//    benchmarkRuleGen(rulegen.BlockV2, b, "v2")
//}
//
//func benchmarkRuleGenV3(b *testing.B) {
//    benchmarkRuleGen(rulegen.BlockV3, b, "v3")
//}

func benchmarkRuleGenV4(b *testing.B) {
    benchmarkRuleGen(rulegen.BlockV4, b, "v4")
}

//func benchmarkRuleGenV5(b *testing.B) {
//    benchmarkRuleGen(rulegen.BlockV5, b, "v5")
//}

const numBench = 50
const expect = 817906

var good []string

func benchmarkRuleGen(evaluator func(*rulegen.Node, string, []string), b *testing.B, name string) {
    b.ReportAllocs()
    b.SetBytes(contentTotals * numBench)

    //for n := 0; n < b.N; n++ {
    for n := 0; n < numBench; n++ {
        res, _ := rulegen.GenFromURLs(urls)
        if len(res) != expect {
            fmt.Println(len(res), name)

            if len(good) > 0 {
                fmt.Println(len(good), len(res), difference(good, res))
            }

            panic(errors.New("length mismatch"))
            //b.FailNow()
        } else {
            good = res
        }

        //res = nil
        //time.Sleep(time.Second)
        //runtime.GC()
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
