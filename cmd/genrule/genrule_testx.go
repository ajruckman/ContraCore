package main

import (
    "fmt"
    "testing"
    "time"

    "github.com/ajruckman/ContraCore/internal/rule"
)

func BenchmarkDummyStart(b *testing.B) {
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        time.Sleep(time.Millisecond)
    }
}

//func BenchmarkPubSub(b *testing.B) {
//    b.ReportAllocs()
//
//    for i := 0; i < b.N; i++ {
//        began := time.Now()
//
//        for _, url := range urls {
//            //fmt.Print("Reading ", url, "... ")
//            resp, err := http.Get(url)
//            Err(err)
//            //fmt.Print(resp.StatusCode, " ")
//
//            scanner := bufio.NewScanner(resp.Body)
//            for scanner.Scan() {
//                contents = append(contents, scanner.Text())
//            }
//            //fmt.Println("done")
//
//            contents = append(contents, strings.Repeat("-", 25))
//        }
//
//        res, _ := rulegen.ProcessDomainSources(rulegen.block, contents)
//        _ = began
//        _ = res
//        fmt.Println(len(res), time.Since(began))
//    }
//}

//func BenchmarkProcessFromURLsNoPointers(b *testing.B) {
//    b.ReportAllocs()
//
//    for i := 0; i < b.N; i++ {
//        began := time.Now()
//        res, _ := rulegen.ProcessFromURLs(urls)
//        _ = began
//        _ = res
//        //fmt.Println(len(res), time.Since(began))
//    }
//}

func benchmarkProcessFromURLsWithPointers(b *testing.B) {
    b.ReportAllocs()
    b.N =1

    //for i := 0; i < b.N; i++ {
        began := time.Now()
        res, total := rule.GenFromURLs(urls)
        _ = began
        _ = res
        fmt.Println("=>", total)
        //fmt.Println(len(res), time.Since(began))
    //}

}
