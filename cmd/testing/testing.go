package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "time"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/rulegen"
)

var (
    urls = []string{
        "https://raw.githubusercontent.com/EnergizedProtection/block/master/unified/formats/domains.txt",
        "https://someonewhocares.org/hosts/hosts",
        "https://gist.githubusercontent.com/angristan/20a398983c5b1daa9c13a1cbadb78fd6/raw/58d54b172b664ee5a0b53bb2e25c391433f2cc7a/hosts",
        "https://www.encrypt-the-planet.com/downloads/hosts",
    }
    contents [][]byte
)

func init() {
    for _, url := range urls {
        fmt.Print("Reading ", url, "... ")
        resp, err := http.Get(url)
        Err(err)
        read, err := ioutil.ReadAll(resp.Body)
        Err(err)
        contents = append(contents, read)
        fmt.Println("done")
    }
}

func main() {
    began := time.Now()
    res, total := rulegen.ReadDomainScanners(rulegen.BlockV2, contents...)

    fmt.Println("Time: ", time.Since(began))
    fmt.Println("Total:", total)
    fmt.Println("Kept: ", len(res))
    fmt.Println("Ratio:", float64(len(res))/float64(total))
}

//urls     = []string{"spark", "bluGo", "blu", "basic", "ultimate", "unified"}

//type Benchmark struct {
//    Pack     string
//    Total    int
//    Kept     int
//    Ratio    float64
//    Duration time.Duration
//}
//
//var results = map[string][]Benchmark{}
//
//func benchmark(count int) {
//    time.Sleep(time.Second * 5)
//
//    for i := 0; i < count; i++ {
//        // Primer
//        rulegen.ReadDomainScanners(rulegen.BlockV2, bufio.NewScanner(bytes.NewReader(body)))
//
//        for _, pack := range packs {
//            fmt.Println(fmt.Sprintf("Bench #%-4d of %s", i+1, pack))
//
//            _, total, kept, ratio, dur := rulegen.ReadDomainScanners(rulegen.BlockV2, bufio.NewScanner(bytes.NewReader(body)))
//
//            if _, ok := results[pack]; !ok {
//                results[pack] = []Benchmark{}
//            }
//            results[pack] = append(results[pack], Benchmark{
//                Pack:     pack,
//                Total:    total,
//                Kept:     kept,
//                Ratio:    ratio,
//                Duration: dur,
//            })
//
//            time.Sleep(time.Second * 3)
//        } //1269353, 1.375
//
//        time.Sleep(time.Second * 15)
//    }
//
//    fmt.Println()
//
//    for _, v := range results {
//        var sum time.Duration
//        for _, bench := range v {
//            sum += bench.Duration
//        }
//
//        c := v[0]
//        f := fmt.Sprintf("%s,%d,%d,%f,%s", c.Pack, c.Total, c.Kept, c.Ratio, (sum / time.Duration(count)).String())
//
//        fmt.Println(f)
//    }
//}