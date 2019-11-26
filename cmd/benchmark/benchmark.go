package main

import (
    "testing"

    "github.com/ajruckman/ContraCore/internal/rulegen"
)

func BenchmarkRuleGenV2(b *testing.B) {
    benchmarkRuleGen(rulegen.BlockV2, b)
}

func BenchmarkRuleGenV3(b *testing.B) {
    benchmarkRuleGen(rulegen.BlockV3, b)
}

func BenchmarkRuleGenV4(b *testing.B) {
    benchmarkRuleGen(rulegen.BlockV4, b)
}

func benchmarkRuleGen(evaluator func(*rulegen.Node, string, []string), b *testing.B) {
    b.ReportAllocs()

    b.ResetTimer()

    for n := 0; n < b.N; n++ {
        //began := time.Now()
        rulegen.ReadDomainScanners(evaluator, contents...)
        //fmt.Println("Time: ", time.Since(began))
        //fmt.Println("Total:", total)
        //fmt.Println("Kept: ", len(res))
        //fmt.Println("Ratio:", float64(len(res))/float64(total))
    }

    //fmt.Println("Done")
}
