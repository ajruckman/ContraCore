package log

import (
	"strings"

	"go.uber.org/atomic"

	"github.com/ajruckman/ContraCore/internal/netmgr"
	"github.com/ajruckman/ContraCore/internal/schema"
	"github.com/ajruckman/ContraCore/internal/system"
)

var (
	AnsweredTotDuration atomic.Duration // The total time spent processing requests that were answered by ContraCore.
	AnsweredTotCount    atomic.Uint32   // The count of queries that were answered by ContraCore.
	PassedTotDuration   atomic.Duration // The total time spent processing and forwarding requests that were not answered by ContraCore.
	PassedTotCount      atomic.Uint32   // The count of queries that were not answered by ContraCore.

	LogRuleLookupDurations = true // Log time spent looking up whitelist and blacklist rules if true.

	queryChannel = make(chan schema.Log) // Channel holding unprocessed query logs.
)

// Pushes a new query log into the query log channel for processing and updates query statistics.
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

// Processes records on the query log channel.
func inputMonitor() {
	for log := range queryChannel {
		netmgr.ProcessQuery(log)
		enqueue(log)
	}
}
