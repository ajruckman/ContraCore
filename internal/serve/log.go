package serve

import (
    "context"
    "database/sql"
    "math/rand"
    "sync"
    "time"

    "github.com/coredns/coredns/plugin/pkg/log"
    "github.com/davecgh/go-spew/spew"
    "github.com/jackc/pgx/v4"
    "github.com/miekg/dns"
    "go.uber.org/atomic"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/eventserver"
    "github.com/ajruckman/ContraCore/internal/schema"
)

var (
    ctx = context.Background()

    logChannel = make(chan queryInfo)
    clog       = log.NewWithPlugin("contradomain")

    queryLogBuffer              []queryInfo
    queryLogBufferLock          sync.Mutex
    queryLogBufferSaveThreshold = 30               // Save all in buffer if the buffer contains this many queries
    queryLogBufferFlushInterval = time.Second * 15 // Save all in buffer if no new logs have been added after this time

    logMonInterval = 15
    logCount       atomic.Uint32

    logAnsweredTotDuration atomic.Duration
    logAnsweredTotCount    atomic.Uint32
    logPassedTotDuration   atomic.Duration
    logPassedTotCount      atomic.Uint32

    dhcpRefreshInterval = 15

    logDurations = true
)

func queryLogWorker() {
    var (
        query queryInfo
        timer = time.NewTimer(queryLogBufferFlushInterval)
    )

    // Debounce code from: https://drailing.net/2018/01/debounce-function-for-golang/
    for {
        select {
        case query = <-logChannel:
            queryLogBufferLock.Lock()
            queryLogBuffer = append(queryLogBuffer, query)

            if len(queryLogBuffer) > queryLogBufferSaveThreshold {
                clog.Infof("log buffer contains %d queries (more than threshold %d); saving to database", len(queryLogBuffer), queryLogBufferSaveThreshold)

                // Save buffer
                saveQueryLogBuffer()
            }

            queryLogBufferLock.Unlock()

            timer.Reset(queryLogBufferFlushInterval)

        case <-timer.C:
            queryLogBufferLock.Lock()

            // Check buffer
            if len(queryLogBuffer) == 0 {
                queryLogBufferLock.Unlock()
                clog.Info("nothing new to log; continuing")
                timer.Reset(queryLogBufferFlushInterval)
                continue
            } else {
                clog.Infof("timer expired and log buffer contains %d queries; saving to database", len(queryLogBuffer))
            }

            // Save buffer
            saveQueryLogBuffer()

            queryLogBufferLock.Unlock()

            timer.Reset(queryLogBufferFlushInterval)
        }
    }
}

func saveQueryLogBuffer() {
    // Create transactions
    var (
        pdbTX   pgx.Tx
        cdbTX   *sql.Tx
        cdbSTMT *sql.Stmt
        err     error
    )

    if db.PostgresOnline.Load() {
        pdbTX, err = db.PDB.Begin(context.Background())
        if err != nil {
            clog.Warningf("could not begin PDB transaction")
            clog.Warning(err.Error())
        }
    }

    if db.ClickHouseOnline.Load() {
        cdbTX, cdbSTMT, err = db.LogCBeginBatch()
        if err != nil {
            clog.Warningf("could not begin CDB transaction")
            clog.Warning(err.Error())
        }
    }

    // Copy buffer
    buffer := make([]queryInfo, len(queryLogBuffer))
    copy(buffer, queryLogBuffer)
    queryLogBuffer = []queryInfo{}

    // Save queries
    for _, q := range buffer {

        began := time.Now()
        v := schema.Log{
            Time:         q.received,
            Client:       q._client,
            Question:     "!internal!" + q._domain,
            QuestionType: dns.TypeToString[q._qu.Qtype],
            Action:       q.action,
            Answers:      q.answers,

            //QueryID: q.r.Id, // This isn't unique
            QueryID: uint16(rand.Intn(65536)),

            Duration: time.Now().Sub(q.received),

            ClientMAC:      q.mac,
            ClientHostname: q.hostname,
            ClientVendor:   q.vendor,
        }
        q.durations.timeGenLogStruct = time.Since(began)

        began = time.Now()
        eventserver.Transmit(v)
        q.durations.timeSendLogToEventClients = time.Since(began)

        if pdbTX != nil {
            began = time.Now()
            err := db.Log(pdbTX, v)
            if err != nil {
                spew.Dump(v)

                clog.Warningf("could not insert log for query '%s'", v.Question)
                clog.Warning(err.Error())
            }
            q.durations.timeSaveLogToPG = time.Since(began)
        } else {
            clog.Warningf("not logging query '%s' because pdbTX is nil", v.Question)
        }

        if cdbTX != nil {
            began = time.Now()
            err = db.LogC(cdbSTMT, v)
            if err != nil {
                clog.Warningf("could not insert secondary log for query '%s'", v.Question)
                clog.Warning(err.Error())
            }
            q.durations.timeSaveLogToCH = time.Since(began)
        }

        if !logDurations {
            clog.Infof("%s <- %d %s %s %v", v.Client, v.QueryID, v.Question, v.QuestionType, v.Duration)
        } else {
            clog.Infof("%s <- %d %s %s %v\n\tLookup lease:        %v\n\tRespond by hostname: %v\n\tRespond by PTR:      %v\n\tRespond with block:  %v\n\tGen log struct:      %v\n\tSave to PG:          %v\n\tSave to CH:          %v\n\tSend to clients:     %v",
                v.Client,
                v.QueryID,
                v.Question,
                v.QuestionType,
                v.Duration,
                q.durations.timeLookupLease,
                q.durations.timeCheckRespondByHostname,
                q.durations.timeCheckRespondByPTR,
                q.durations.timeCheckRespondWithBlock,
                q.durations.timeGenLogStruct,
                q.durations.timeSaveLogToPG,
                q.durations.timeSaveLogToCH,
                q.durations.timeSendLogToEventClients,
            )
        }

        if v.Action == "pass" {
            logPassedTotCount.Inc()
            logPassedTotDuration.Add(v.Duration)
        } else {
            logAnsweredTotCount.Inc()
            logAnsweredTotDuration.Add(v.Duration)
        }
    }

    if pdbTX != nil {
        err = pdbTX.Commit(ctx)
        if err != nil {
            clog.Warningf("could not commit pdbTX: %s", err.Error())
        }
    }

    if cdbTX != nil {
        err = cdbTX.Commit()
        if err != nil {
            clog.Warningf("could not commit cdbTX: %s", err.Error())
        }
    }

    logCount.Add(uint32(len(buffer)))
}

func logMonitor() {
    for range time.Tick(time.Duration(logMonInterval) * time.Second) {
        c := logCount.Swap(0)

        var (
            avgDurAns  = float64(logAnsweredTotDuration.Swap(0).Milliseconds()) / float64(logAnsweredTotCount.Swap(0))
            avgDurPass = float64(logPassedTotDuration.Swap(0).Milliseconds()) / float64(logPassedTotCount.Swap(0))
        )

        clog.Infof("Log channel backlog: %d | New log rows: %d | Rows/second: %.3f | Avg. ms answered reqs.: %.2f | Avg. ms passed reqs.: %.2f",
            len(logChannel),
            c,
            float64(c)/float64(logMonInterval),
            avgDurAns,
            avgDurPass,
        )
    }
}
