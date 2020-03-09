package log

import (
	"github.com/ajruckman/ContraCore/internal/netmgr"
)

// Package setup function.
func Setup() {
	go queryBufferFlushScheduled()
	//go listen()
	go statWorker()
	go inputMonitor()

	go netmgr.Setup()
}
