package log

import (
    "sync"
    "time"

    "github.com/ajruckman/ContraCore/internal/log/contralog"
    "github.com/ajruckman/ContraCore/internal/schema"
    "github.com/ajruckman/ContraCore/internal/state"
)

var (
    queryBuffer     []schema.Log
    queryBufferLock sync.Mutex

    queryBufferSaveThreshold = 25               // Save all in buffer if the buffer contains this many queries
    queryBufferSaveInterval  = time.Second * 30 // Save all in buffer if no new logs have been added after this time

    queryBufferMonitorTimer = time.NewTimer(queryBufferSaveInterval)
)

func enqueue(log schema.Log) {
    queryBufferLock.Lock()

    queryBufferMonitorTimer.Reset(queryBufferSaveInterval)
    queryBuffer = append(queryBuffer, log)

    if len(queryBuffer) >= queryBufferSaveThreshold {
        state.Console.Infof("log buffer contains %d queries (more than threshold %d); flushing to database immediately", len(queryBuffer), queryBufferSaveThreshold)
        contralog.SaveQueryLogBuffer(queryBuffer)
        queryBuffer = []schema.Log{}
    }

    queryBufferLock.Unlock()
}

func queryBufferDebouncer() {
    for range queryBufferMonitorTimer.C {
        queryBufferLock.Lock()

        if len(queryBuffer) == 0 {
            queryBufferLock.Unlock()
            return
        }

        state.Console.Infof("timer expired and log buffer contains %d queries; flushing to database", len(queryBuffer))

        contralog.SaveQueryLogBuffer(queryBuffer)
        queryBuffer = []schema.Log{}
        queryBufferLock.Unlock()
    }
}
