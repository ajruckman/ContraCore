package log

import (
	"time"

	"go.uber.org/atomic"

	"github.com/ajruckman/ContraCore/internal/system"
)

var (
	monInterval = time.Second * 15
	monCount    atomic.Uint32
	monAnyNew   atomic.Bool
)

// Prints query statistics to STDOUT on a schedule.
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
			avgDurPassed    = float64(PassedTotDuration.Swap(0).Milliseconds()) / float64(PassedTotCount.Swap(0))
			avgDurBlocked   = float64(BlockedTotDuration.Swap(0).Milliseconds()) / float64(BlockedTotCount.Swap(0))
			avgDurResponded = float64(RespondedTotDuration.Swap(0).Milliseconds()) / float64(RespondedTotCount.Swap(0))
		)

		system.Console.Infof(""+
			"Log buffer size: %d | "+
			"New log rows: %d | "+
			"Rows/second: %.3f | "+
			"Avg. ms passed reqs.: %.2f | "+
			"Avg. ms blocked reqs.: %.2f | "+
			"Avg. ms responded reqs.: %.2f | "+
			"ContraDB: %t | "+
			"ContraLog: %t",

			queryLogBufferLen,
			c,
			float64(c)/monInterval.Seconds(),
			avgDurPassed,
			avgDurBlocked,
			avgDurResponded,
			system.ContraDBOnline.Load(),
			system.ContraLogOnline.Load(),
		)
	}
}
