package main

import (
    "testing"

    "github.com/ajruckman/ContraCore/internal/rulegen"
)

// https://www.speedscope.app/

func BenchmarkRuleGen(b *testing.B) {
    b.ReportAllocs()
    b.SetBytes(contentTotals * int64(b.N))

    for n := 0; n < b.N; n++ {
        rulegen.ReadDomainScanners(rulegen.BlockV6, contents...)
    }
}
