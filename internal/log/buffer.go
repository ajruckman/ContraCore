package log

import (
	"database/sql"
	"sync"
	"time"

	. "github.com/ajruckman/xlib"

	"github.com/ajruckman/ContraCore/internal/db/contralog"
	"github.com/ajruckman/ContraCore/internal/schema"
	"github.com/ajruckman/ContraCore/internal/system"
)

var (
	queryBuffer     []schema.Log
	queryBufferLock sync.Mutex

	queryBufferSaveThreshold = 100              // Save all in buffer if the buffer contains this many queries.
	queryBufferSaveInterval  = time.Second * 30 // Save all in buffer if no new logs have been added after this time.
)

// Adds a query log to the query log buffer. Flushes the buffer to ContraLog if
// the cache then contains equal to or greater than queryBufferSaveThreshold
// logs.
func enqueue(log schema.Log) {
	queryBufferLock.Lock()

	queryBuffer = append(queryBuffer, log)

	if len(queryBuffer) >= queryBufferSaveThreshold {
		if system.ContraLogOnline.Load() {
			system.Console.Infof("log buffer contains %d queries (more than threshold %d); flushing to database immediately", len(queryBuffer), queryBufferSaveThreshold)
			saveQueryLogBuffer(queryBuffer)
		} else {
			system.Console.Infof("log buffer contains %d queries (more than threshold %d), but ContraLog is not connected; clearing buffer", len(queryBuffer), queryBufferSaveThreshold)
		}
		queryBuffer = []schema.Log{}
	}

	queryBufferLock.Unlock()
}

// Flushes the query log buffer, if it contains any logs, to the database on a
// timer.
func queryBufferFlushScheduled() {
	for range time.Tick(queryBufferSaveInterval) {
		queryBufferLock.Lock()

		if len(queryBuffer) == 0 {
			queryBufferLock.Unlock()
			continue
		}

		if system.ContraLogOnline.Load() {
			system.Console.Infof("log buffer timer expired and log buffer contains %d queries; flushing to database", len(queryBuffer))
			saveQueryLogBuffer(queryBuffer)
		} else {
			system.Console.Infof("log buffer timer expired and log buffer contains %d queries, but ContraLog is not connected; clearing buffer", len(queryBuffer))
		}

		queryBuffer = []schema.Log{}
		queryBufferLock.Unlock()
	}
}

// Immediately saves a log slice to ContraLog.
func saveQueryLogBuffer(buffer []schema.Log) {
	// Create transactions
	var (
		cdbTX   *sql.Tx
		cdbSTMT *sql.Stmt
		err     error
	)

	if system.ContraLogOnline.Load() {
		cdbTX, cdbSTMT, err = contralog.BeginLogBatch()
		if err != nil {
			system.Console.Error("failed to begin cdb transaction with error:")
			system.Console.Error(err.Error())
			return
		}
	}

	// Save queries
	for _, v := range buffer {
		if cdbTX != nil {
			err = contralog.SaveLog(cdbSTMT, v.Log)
			if err != nil {
				system.Console.Errorf("failed to insert log for query '%s' with error:", v.Question)
				system.Console.Error(err.Error())
			}
		}
	}

	if cdbTX != nil {
		err = contralog.CommitLogBatch(cdbTX)
		Err(err)
	}
}
