package main

import (
    "testing"
    "time"

    "github.com/ajruckman/ContraCore/internal/rulegen"
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
//        res, _ := rulegen.ProcessDomainSources(rulegen.BlockV4, contents)
//        _ = began
//        _ = res
//        fmt.Println(len(res), time.Since(began))
//    }
//}

func BenchmarkProcessFromURLsNoPointers(b *testing.B) {
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        began := time.Now()
        res, _ := rulegen.ProcessFromURLs(urls)
        _ = began
        _ = res
        //fmt.Println(len(res), time.Since(began))
    }
}

func BenchmarkProcessFromURLsWithPointers(b *testing.B) {
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        began := time.Now()
        res, _ := rulegen.ProcessFromURLsPointers(urls)
        _ = began
        _ = res
        //fmt.Println(len(res), time.Since(began))
    }
}
