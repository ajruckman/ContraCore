package log

import (
	"strings"

	"go.uber.org/atomic"

	"github.com/ajruckman/ContraCore/internal/netmgr"
	"github.com/ajruckman/ContraCore/internal/schema"
	"github.com/ajruckman/ContraCore/internal/system"
)

var (
	PassedTotDuration    atomic.Duration // The total time spent processing and forwarding requests that were not answered by ContraCore.
	PassedTotCount       atomic.Uint32   // The count of queries that were not answered by ContraCore.
	BlockedTotDuration   atomic.Duration // The total time spent processing and responding to requests that were blocked by ContraCore.
	BlockedTotCount      atomic.Uint32   // The count of queries that were blocked by ContraCore.
	RespondedTotDuration atomic.Duration // The total time spent processing requests that were answered by ContraCore.
	RespondedTotCount    atomic.Uint32   // The count of queries that were answered by ContraCore.

	queryChannel = make(chan schema.Log) // Channel holding unprocessed query logs.
)

// Pushes a new query log into the query log channel for processing and updates query statistics.
func Query(log schema.Log) {
	monAnyNew.Store(true)
	monCount.Add(1)

	if strings.HasPrefix(log.Action, "pass.") {
		PassedTotCount.Inc()
		PassedTotDuration.Add(log.Duration)
	} else if strings.HasPrefix(log.Action, "block.") {
		BlockedTotCount.Inc()
		BlockedTotDuration.Add(log.Duration)
	} else if strings.HasPrefix(log.Action, "respond.") {
		RespondedTotCount.Inc()
		RespondedTotDuration.Add(log.Duration)
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
