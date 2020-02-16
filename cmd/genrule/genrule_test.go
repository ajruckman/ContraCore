package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ajruckman/ContraCore/internal/rule"
)

func BenchmarkGenRule(b *testing.B) {
	b.ReportAllocs()

	fmt.Println()
	fmt.Println()
	fmt.Println(strings.Repeat("-", 10))

	var totMilliseconds int64

	for i := 0; i < b.N; i++ {
		begin := time.Now()
		res, total := rule.GenFromURLs(urls)
		end := time.Now()

		kept := len(res)
		ratio := float64(kept) / float64(total)

		totMilliseconds += end.Sub(begin).Milliseconds()

		fmt.Println(ratio, kept, total, end.Sub(begin))
	}
	fmt.Println(strings.Repeat("-", 10))
	fmt.Println(float64(totMilliseconds)/float64(b.N), "milliseconds per op")
	fmt.Println()
}
