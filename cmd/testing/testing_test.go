package main

import (
    "bufio"
    "bytes"
    "fmt"
    "testing"

    "github.com/ajruckman/ContraCore/internal/rulegen"
)

func BenchmarkRuleGenV1(b *testing.B) {
    benchmarkRuleGen(rulegen.BlockV1, b)
    fmt.Println(rulegen.IncV1)
    return
}

func BenchmarkRuleGenV2(b *testing.B) {
    benchmarkRuleGen(rulegen.BlockV2, b)
    fmt.Println(rulegen.IncV2)
    return
}

func BenchmarkRuleGenV3(b *testing.B) {
    benchmarkRuleGen(rulegen.BlockV3, b)
    fmt.Println(rulegen.IncV3)
    return
}

func benchmarkRuleGen(evaluator func(*rulegen.Node, string, []string), b *testing.B) {
    b.ReportAllocs()

    for n := 0; n < b.N; n++ {
        _, total, _, _, _ := rulegen.ReadDomainList(evaluator, bufio.NewScanner(bytes.NewReader(body)))
        fmt.Println(total)
    }
}
