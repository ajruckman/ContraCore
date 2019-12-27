package serve

import (
    "time"

    "github.com/coredns/coredns/plugin/pkg/log"
    "go.uber.org/atomic"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema"
)

var (
    logChannel = make(chan schema.Log)
    clog       = log.NewWithPlugin("contradomain")

    logMonInterval = 30
    logCount       atomic.Uint32

    logAnsweredTotDuration atomic.Duration
    logAnsweredTotCount    atomic.Uint32
    logPassedTotDuration   atomic.Duration
    logPassedTotCount      atomic.Uint32

    dhcpRefreshInterval = 15
)

func logWorker() {
    for v := range logChannel {
        err := db.Log(v)
        if err != nil {
            clog.Warningf("could not insert log for query '%s'", v.Question)
            clog.Warning(err.Error())
        }

        clog.Infof("%s <- %d %s %v", v.Client, v.QueryID, v.QuestionType, v.Duration)

        logCount.Inc()

        if v.Action == "pass" {
            logPassedTotCount.Inc()
            logPassedTotDuration.Add(v.Duration)
        } else {
            logAnsweredTotCount.Inc()
            logAnsweredTotDuration.Add(v.Duration)
        }
    }
}

func logMonitor() {
    for range time.Tick(time.Duration(logMonInterval) * time.Second) {
        c := logCount.Swap(0)

        var (
            avgDurAns    = float64(logAnsweredTotDuration.Swap(0).Milliseconds()) / float64(logAnsweredTotCount.Swap(0))
            avgDurPassed = float64(logPassedTotDuration.Swap(0).Milliseconds()) / float64(logPassedTotCount.Swap(0))
        )

        clog.Infof("Log channel backlog: %d | New log rows: %d | Rows/second: %.3f | Avg. ms answered reqs.: %.2f | Avg. ms passed reqs.: %.2f",
            len(logChannel),
            c,
            float64(c)/float64(logMonInterval),
            avgDurAns,
            avgDurPassed,
        )
    }
}
