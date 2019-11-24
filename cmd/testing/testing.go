package main

import (
    "bufio"
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/rulegen"
)

var body []byte

func init() {
    fmt.Println("Init")
    resp, err := http.Get("https://raw.githubusercontent.com/EnergizedProtection/block/master/unified/formats/domains.txt")
    Err(err)

    body, err = ioutil.ReadAll(resp.Body)
    Err(err)
    fmt.Println("Rule byte length:", len(body))
}

func main() {
    //benchmark(25)
    //var hosts []string

    rulegen.ReadDomainList(rulegen.BlockV1, bufio.NewScanner(bytes.NewReader(body)))

    _, total, kept, ratio, dur := rulegen.ReadDomainList(rulegen.BlockV1,
        bufio.NewScanner(bytes.NewReader(body)),
        //    //"https://raw.githubusercontent.com/EnergizedProtection/block/master/blu/formats/hosts",
        //    //"https://someonewhocares.org/hosts/hosts",
    )

    fmt.Println(total)
    fmt.Println(kept)
    fmt.Println(ratio)
    fmt.Println(dur)
}

var packs = []string{"spark", "bluGo", "blu", "basic", "ultimate", "unified"}

type Benchmark struct {
    Pack     string
    Total    int
    Kept     int
    Ratio    float64
    Duration time.Duration
}

var results = map[string][]Benchmark{}

func benchmark(count int) {
    time.Sleep(time.Second * 5)

    for i := 0; i < count; i++ {
        // Primer
        rulegen.ReadDomainList(rulegen.BlockV1, bufio.NewScanner(bytes.NewReader(body)))

        for _, pack := range packs {
            fmt.Println(fmt.Sprintf("Bench #%-4d of %s", i+1, pack))

            _, total, kept, ratio, dur := rulegen.ReadDomainList(rulegen.BlockV1, bufio.NewScanner(bytes.NewReader(body)))

            if _, ok := results[pack]; !ok {
                results[pack] = []Benchmark{}
            }
            results[pack] = append(results[pack], Benchmark{
                Pack:     pack,
                Total:    total,
                Kept:     kept,
                Ratio:    ratio,
                Duration: dur,
            })

            time.Sleep(time.Second * 3)
        } //1269353, 1.375

        time.Sleep(time.Second * 15)
    }

    fmt.Println()

    for _, v := range results {
        var sum time.Duration
        for _, bench := range v {
            sum += bench.Duration
        }

        c := v[0]
        f := fmt.Sprintf("%s,%d,%d,%f,%s", c.Pack, c.Total, c.Kept, c.Ratio, (sum / time.Duration(count)).String())

        fmt.Println(f)
    }
}
