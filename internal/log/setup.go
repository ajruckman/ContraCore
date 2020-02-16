package log

import (
	"github.com/ajruckman/ContraCore/internal/log/eventserver"
)

// Package setup function.
func Setup() {
	go queryBufferFlushScheduled()
	go listen()
	go statWorker()
	go inputMonitor()

	go eventserver.Setup()
}
