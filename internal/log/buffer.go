package log

import (
    "sync"
    "time"

    "github.com/ajruckman/ContraCore/internal/log/contralog"
    "github.com/ajruckman/ContraCore/internal/schema"
    "github.com/ajruckman/ContraCore/internal/system"
)

var (
    queryBuffer     []schema.Log
    queryBufferLock sync.Mutex

    queryBufferSaveThreshold = 100              // Save all in buffer if the buffer contains this many queries
    queryBufferSaveInterval  = time.Second * 3 // Save all in buffer if no new logs have been added after this time

    queryBufferFlushTicker = time.NewTicker(queryBufferSaveInterval)
)

func enqueue(log schema.Log) {
    queryBufferLock.Lock()

    queryBuffer = append(queryBuffer, log)

    if len(queryBuffer) >= queryBufferSaveThreshold {
        if system.ContraLogOnline.Load() {
            system.Console.Infof("log buffer contains %d queries (more than threshold %d); flushing to database immediately", len(queryBuffer), queryBufferSaveThreshold)
            contralog.SaveQueryLogBuffer(queryBuffer)
        } else {
            system.Console.Infof("log buffer contains %d queries (more than threshold %d), but ContraLog is not connected; clearing buffer", len(queryBuffer), queryBufferSaveThreshold)
        }
        queryBuffer = []schema.Log{}
    }

    queryBufferLock.Unlock()
}

func queryBufferFlushScheduled() {
    for range queryBufferFlushTicker.C {
        queryBufferLock.Lock()

        if len(queryBuffer) == 0 {
            queryBufferLock.Unlock()
            continue
        }

        if system.ContraLogOnline.Load() {
            system.Console.Infof("log buffer timer expired and log buffer contains %d queries; flushing to database", len(queryBuffer))
            contralog.SaveQueryLogBuffer(queryBuffer)
        } else {
            system.Console.Infof("log buffer timer expired and log buffer contains %d queries, but ContraLog is not connected; clearing buffer", len(queryBuffer))
        }

        queryBuffer = []schema.Log{}
        queryBufferLock.Unlock()
    }
}
