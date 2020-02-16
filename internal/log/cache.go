package log

import (
    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/db/contralog"
    "github.com/ajruckman/ContraCore/internal/schema"
    "github.com/ajruckman/ContraCore/internal/system"
)

const cacheSize = 5000

var cache []schema.Log

func loadCache() {
    logs, err := contralog.GetLastNLogs(cacheSize)
    Err(err)

    cache = schema.LogsFromContraLogs(logs)
    system.Console.Infof("loaded %d recent logs from ContraLog", len(cache))
}
