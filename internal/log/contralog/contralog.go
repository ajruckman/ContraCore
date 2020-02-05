package contralog

import (
    "database/sql"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/db/contralog"
    "github.com/ajruckman/ContraCore/internal/schema"
    "github.com/ajruckman/ContraCore/internal/system"
)

//func QueryLogWorker() {
//    var (
//        query log.QueryInfo
//        timer = time.NewTimer(log.queryLogBufferFlushInterval)
//    )
//
//    // Debounce code from: https://drailing.net/2018/01/debounce-function-for-golang/
//    for {
//        select {
//        case query = <-log.inputChannel:
//            log.queryLogBufferLock.Lock()
//            v := dbschema.Log{
//                Time:         query.Received,
//                Client:       query.Client_,
//                Question:     query.Domain_,
//                QuestionType: dns.TypeToString[query.QU_.Qtype],
//                Action:       query.Action,
//                Answers:      query.Answers,
//
//                //QueryID: query.r.Id, // This isn't unique
//                QueryID: uint16(rand.Intn(65536)),
//
//                Duration: time.Now().Sub(query.Received),
//
//                ClientMAC:      query.MAC,
//                ClientHostname: query.Hostname,
//                ClientVendor:   query.Vendor,
//            }
//            log.Console.Infof("%s <- %d %s %s %v", v.Client, v.QueryID, v.Question, v.QuestionType, v.Duration)
//
//            log.queryLogBuffer = append(log.queryLogBuffer, v)
//            if len(log.queryLogBuffer) > log.queryLogBufferSaveThreshold {
//                log.Console.Infof("log buffer contains %d queries (more than threshold %d); saving to database", len(log.queryLogBuffer), log.queryLogBufferSaveThreshold)
//
//                // Copy buffer
//                buffer := make([]dbschema.Log, len(log.queryLogBuffer))
//                copy(buffer, log.queryLogBuffer)
//                log.queryLogBuffer = []dbschema.Log{}
//
//                // Save buffer
//                saveQueryLogBuffer(buffer)
//            }
//
//            log.queryLogBufferLock.Unlock()
//
//            timer.Reset(log.queryLogBufferFlushInterval)
//
//        case <-timer.C:
//            buffer := log.ReadAndResetBuffer()
//
//            // Check buffer
//            if len(buffer) == 0 {
//                log.Console.Info("nothing new to log; continuing")
//                timer.Reset(log.queryLogBufferFlushInterval)
//                continue
//            } else {
//                log.Console.Infof("timer expired and log buffer contains %d queries; saving to database", len(log.queryLogBuffer))
//            }
//
//            // Save buffer
//            saveQueryLogBuffer(buffer)
//
//            timer.Reset(log.queryLogBufferFlushInterval)
//        }
//    }
//}

func SaveQueryLogBuffer(buffer []schema.Log) {
    // Create transactions
    var (
        cdbTX   *sql.Tx
        cdbSTMT *sql.Stmt
        err     error
    )

    if system.ClickHouseOnline.Load() {
        cdbTX, cdbSTMT, err = contralog.BeginBatch()
        if err != nil {
            system.Console.Warningf("could not begin cdb transaction")
            system.Console.Warning(err.Error())
        }
    }

    // Save queries
    for _, v := range buffer {
        if cdbTX != nil {
            err = contralog.SaveLog(cdbSTMT, v.ToContraLog())
            if err != nil {
                system.Console.Warningf("could not insert secondary log for query '%s'", v.Question)
                system.Console.Warning(err.Error())
            }
        }
    }

    if cdbTX != nil {
        err = contralog.CommitBatch(cdbTX)
        Err(err)
    }
}
