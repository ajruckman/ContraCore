package serve

//import (
//    "context"
//    "database/sql"
//    "math/rand"
//    "sync"
//    "time"
//
//    "github.com/coredns/coredns/plugin/pkg/log"
//    "github.com/davecgh/go-spew/spew"
//    "github.com/jackc/pgx/v4"
//    "github.com/miekg/dns"
//    "go.uber.org/atomic"
//
//    "github.com/ajruckman/ContraCore/internal/db"
//    "github.com/ajruckman/ContraCore/internal/eventserver"
//    "github.com/ajruckman/ContraCore/internal/schema"
//)
//
//var (
//    ctx = context.Background()
//
//    clog       = log.NewWithPlugin("contradomain")
//
//    QueryLogBuffer              []schema.Log
//    QueryLogBufferLock          sync.Mutex
//    queryLogBufferSaveThreshold = 15               // Save all in buffer if the buffer contains this many queries
//    queryLogBufferFlushInterval = time.Second * 30 // Save all in buffer if no new logs have been added after this time
//
//    logMonInterval = time.Second * 15
//    logMonCount    atomic.Uint32
//    logMonAnyNew   atomic.Bool
//
//    logAnsweredTotDuration atomic.Duration
//    logAnsweredTotCount    atomic.Uint32
//    logPassedTotDuration   atomic.Duration
//    logPassedTotCount      atomic.Uint32
//
//    logDurations = true
//)
//
//func queryLogWorker() {
//    var (
//        query schema.QueryInfo
//        timer = time.NewTimer(queryLogBufferFlushInterval)
//    )
//
//    // Debounce code from: https://drailing.net/2018/01/debounce-function-for-golang/
//    for {
//        select {
//        case query = <-logChannel:
//            QueryLogBufferLock.Lock()
//            v := schema.Log{
//                Time:         query.received,
//                Client:       query._client,
//                Question:     query._domain,
//                QuestionType: dns.TypeToString[query._qu.Qtype],
//                Action:       query.action,
//                Answers:      query.answers,
//
//                //QueryID: query.r.Id, // This isn't unique
//                QueryID: uint16(rand.Intn(65536)),
//
//                Duration: time.Now().Sub(query.received),
//
//                ClientMAC:      query.mac,
//                ClientHostname: query.hostname,
//                ClientVendor:   query.vendor,
//            }
//            clog.Infof("%s <- %d %s %s %v", v.Client, v.QueryID, v.Question, v.QuestionType, v.Duration)
//
//            QueryLogBuffer = append(QueryLogBuffer, v)
//            if len(QueryLogBuffer) > queryLogBufferSaveThreshold {
//                clog.Infof("log buffer contains %d queries (more than threshold %d); saving to database", len(QueryLogBuffer), queryLogBufferSaveThreshold)
//
//                // Copy buffer
//                buffer := make([]schema.Log, len(QueryLogBuffer))
//                copy(buffer, QueryLogBuffer)
//                QueryLogBuffer = []schema.Log{}
//
//                // Save buffer
//                saveQueryLogBuffer(buffer)
//            }
//
//            QueryLogBufferLock.Unlock()
//
//            timer.Reset(queryLogBufferFlushInterval)
//
//        case <-timer.C:
//            QueryLogBufferLock.Lock()
//
//            // Check buffer
//            if len(QueryLogBuffer) == 0 {
//                QueryLogBufferLock.Unlock()
//                clog.Info("nothing new to log; continuing")
//                timer.Reset(queryLogBufferFlushInterval)
//                continue
//            } else {
//                clog.Infof("timer expired and log buffer contains %d queries; saving to database", len(QueryLogBuffer))
//            }
//
//            // Copy buffer
//            buffer := make([]schema.Log, len(QueryLogBuffer))
//            copy(buffer, QueryLogBuffer)
//            QueryLogBuffer = []schema.Log{}
//
//            // Save buffer
//            saveQueryLogBuffer(buffer)
//
//            QueryLogBufferLock.Unlock()
//
//            timer.Reset(queryLogBufferFlushInterval)
//        }
//    }
//}
//
//func saveQueryLogBuffer(buffer []schema.Log) {
//    logMonAnyNew.Store(true)
//
//    // Create transactions
//    var (
//        pdbTX   pgx.Tx
//        cdbTX   *sql.Tx
//        cdbSTMT *sql.Stmt
//        err     error
//    )
//
//    if db.PostgresOnline.Load() {
//        pdbTX, err = db.PDB.Begin(context.Background())
//        if err != nil {
//            clog.Warningf("could not begin PDB transaction")
//            clog.Warning(err.Error())
//        }
//    }
//
//    if db.ClickHouseOnline.Load() {
//        cdbTX, cdbSTMT, err = db.LogCBeginBatch()
//        if err != nil {
//            clog.Warningf("could not begin CDB transaction")
//            clog.Warning(err.Error())
//        }
//    }
//
//    // Save queries
//    for _, v := range buffer {
//
//        //began := time.Now()
//
//        //q.durations.timeGenLogStruct = time.Since(began)
//
//        //began = time.Now()
//        eventserver.Transmit(v)
//        //q.durations.timeSendLogToEventClients = time.Since(began)
//
//        if pdbTX != nil {
//            //began = time.Now()
//            err := db.Log(pdbTX, v)
//            if err != nil {
//                spew.Dump(v)
//
//                clog.Warningf("could not insert log for query '%s'", v.Question)
//                clog.Warning(err.Error())
//            }
//            //q.durations.timeSaveLogToPG = time.Since(began)
//        } else {
//            clog.Warningf("not logging query '%s' because pdbTX is nil", v.Question)
//        }
//
//        if cdbTX != nil {
//            //began = time.Now()
//            err = db.LogC(cdbSTMT, v)
//            if err != nil {
//                clog.Warningf("could not insert secondary log for query '%s'", v.Question)
//                clog.Warning(err.Error())
//            }
//            //q.durations.timeSaveLogToCH = time.Since(began)
//        }
//
//        //if !logDurations {
//        //    clog.Infof("%s <- %d %s %s %v", v.Client, v.QueryID, v.Question, v.QuestionType, v.Duration)
//        //}
//        //else {
//        //    clog.Infof("%s <- %d %s %s %v\n\tLookup lease:        %v\n\tRespond by hostname: %v\n\tRespond by PTR:      %v\n\tRespond with block:  %v\n\tGen log struct:      %v\n\tSave to PG:          %v\n\tSave to CH:          %v\n\tSend to clients:     %v",
//        //        v.Client,
//        //        v.QueryID,
//        //        v.Question,
//        //        v.QuestionType,
//        //        v.Duration,
//        //        q.durations.timeLookupLease,
//        //        q.durations.timeCheckRespondByHostname,
//        //        q.durations.timeCheckRespondByPTR,
//        //        q.durations.timeCheckRespondWithBlock,
//        //        q.durations.timeGenLogStruct,
//        //        q.durations.timeSaveLogToPG,
//        //        q.durations.timeSaveLogToCH,
//        //        q.durations.timeSendLogToEventClients,
//        //    )
//        //}
//
//        if v.Action == "pass" {
//            logPassedTotCount.Inc()
//            logPassedTotDuration.Add(v.Duration)
//        } else {
//            logAnsweredTotCount.Inc()
//            logAnsweredTotDuration.Add(v.Duration)
//        }
//    }
//
//    if pdbTX != nil {
//        err = pdbTX.Commit(ctx)
//        if err != nil {
//            clog.Warningf("could not commit pdbTX: %s", err.Error())
//        }
//    }
//
//    if cdbTX != nil {
//        err = cdbTX.Commit()
//        if err != nil {
//            clog.Warningf("could not commit cdbTX: %s", err.Error())
//        }
//    }
//
//    logMonCount.Add(uint32(len(buffer)))
//}
//
//func logMonitor() {
//    for range time.Tick(logMonInterval) {
//        if !logMonAnyNew.Swap(false) {
//            continue
//        }
//
//        c := logMonCount.Swap(0)
//
//        QueryLogBufferLock.Lock()
//        queryLogBufferLen := len(QueryLogBuffer)
//        QueryLogBufferLock.Unlock()
//
//        var (
//            avgDurAns  = float64(logAnsweredTotDuration.Swap(0).Milliseconds()) / float64(logAnsweredTotCount.Swap(0))
//            avgDurPass = float64(logPassedTotDuration.Swap(0).Milliseconds()) / float64(logPassedTotCount.Swap(0))
//        )
//
//        clog.Infof("Log channel backlog: %d | New log rows: %d | Rows/second: %.3f | Avg. ms answered reqs.: %.2f | Avg. ms passed reqs.: %.2f",
//            queryLogBufferLen,
//            c,
//            float64(c)/logMonInterval.Seconds(),
//            avgDurAns,
//            avgDurPass,
//        )
//    }
//}
