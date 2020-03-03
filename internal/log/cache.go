package log

import (
	. "github.com/ajruckman/xlib"

	"github.com/ajruckman/ContraCore/internal/db/contralog"
	"github.com/ajruckman/ContraCore/internal/schema"
	"github.com/ajruckman/ContraCore/internal/system"
)

// The maximum number of logs in the query log cache.
const cacheSize = 5000

// A slice containing the latest query logs so that they may be sent to new
// ContraWeb clients.
var cache []schema.Log

// Reads at most cacheSize logs into the query log cache.
func loadCache() {
	if !system.ContraLogOnline.Load() {
		system.Console.Warning("ContraLog is disconnected; not loading recent query log cache yet")
		return
	}
	logs, err := contralog.GetLastNLogs(cacheSize)
	Err(err)

	cache = schema.LogsFromContraLogs(logs)
	system.Console.Infof("loaded %d recent logs from ContraLog", len(cache))
}
