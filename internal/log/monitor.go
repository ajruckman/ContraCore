package log

import (
    "time"

    "go.uber.org/atomic"

    "github.com/ajruckman/ContraCore/internal/state"
)

var (
    monInterval = time.Second * 15
    monCount    atomic.Uint32
    monAnyNew   atomic.Bool
)

func statWorker() {
    for range time.Tick(monInterval) {
        // No new events have been logged
        if !monAnyNew.Swap(false) {
            continue
        }

        c := monCount.Swap(0)

        queryBufferLock.Lock()
        queryLogBufferLen := len(queryBuffer)
        queryBufferLock.Unlock()

        var (
            avgDurAns  = float64(AnsweredTotDuration.Swap(0).Milliseconds()) / float64(AnsweredTotCount.Swap(0))
            avgDurPass = float64(PassedTotDuration.Swap(0).Milliseconds()) / float64(PassedTotCount.Swap(0))
        )

        state.Console.Infof("Log buffer size: %d | New log rows: %d | Rows/second: %.3f | Avg. ms answered reqs.: %.2f | Avg. ms passed reqs.: %.2f",
            queryLogBufferLen,
            c,
            float64(c)/monInterval.Seconds(),
            avgDurAns,
            avgDurPass,
        )
    }
}
