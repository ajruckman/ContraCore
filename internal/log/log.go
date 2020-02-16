package log

import (
    "strings"

    "go.uber.org/atomic"

    "github.com/ajruckman/ContraCore/internal/log/eventserver"
    "github.com/ajruckman/ContraCore/internal/schema"
    "github.com/ajruckman/ContraCore/internal/system"
)

var (
    AnsweredTotDuration atomic.Duration
    AnsweredTotCount    atomic.Uint32
    PassedTotDuration   atomic.Duration
    PassedTotCount      atomic.Uint32

    LogRuleLookupDurations = true

    queryChannel = make(chan schema.Log)
)

func Query(log schema.Log) {
    monAnyNew.Store(true)
    monCount.Add(1)

    if strings.HasPrefix(log.Action, "pass.") {
        PassedTotCount.Inc()
        PassedTotDuration.Add(log.Duration)
    } else {
        AnsweredTotCount.Inc()
        AnsweredTotDuration.Add(log.Duration)
    }

    select {
    case queryChannel <- log:
        break
    default:
        system.Console.Warningf("couldn't immediately push to queryChannel: %s", log.Question)
        queryChannel <- log
    }
}

func inputMonitor() {
    for log := range queryChannel {
        eventserver.Transmit(log)
        enqueue(log)
    }
}
