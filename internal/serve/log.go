package serve

import (
    "context"
    "math/rand"
    "sync"
    "time"

    "github.com/coredns/coredns/plugin/pkg/log"
    "github.com/davecgh/go-spew/spew"
    "github.com/miekg/dns"
    "go.uber.org/atomic"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/eventserver"
    "github.com/ajruckman/ContraCore/internal/schema"
)

var (
    logChannel = make(chan queryContext)
    clog       = log.NewWithPlugin("contradomain")

    logBuffer     []queryContext
    logBufferLock sync.Mutex

    logMonInterval = 30
    logCount       atomic.Uint32

    logAnsweredTotDuration atomic.Duration
    logAnsweredTotCount    atomic.Uint32
    logPassedTotDuration   atomic.Duration
    logPassedTotCount      atomic.Uint32

    dhcpRefreshInterval = 15

    logDurations = true
)

func saveQueries(buffer []queryContext) {
    // Create transactions
    pdbTX, err := db.PDB.Begin(context.Background())
    if err != nil {
        clog.Warningf("could not begin PDB transaction")
        clog.Warning(err.Error())
    }

    cdbTX, err := db.CDB.Begin()
    if err != nil {
        clog.Warningf("could not begin CDB transaction")
        clog.Warning(err.Error())
    }

    // Save queries
    for _, q := range buffer {

        began := time.Now()
        v := schema.Log{
            Time:         q.received,
            Client:       q._client,
            Question:     q._domain,
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
            err = db.LogC(cdbTX, v)
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

    logCount.Add(uint32(len(buffer)))
}

func logWorker() {
    var (
        query queryContext
        timer = time.NewTimer(time.Second * 10)
    )

    // Debounce code from: https://drailing.net/2018/01/debounce-function-for-golang/
    for {
        select {
        case query = <-logChannel:
            logBufferLock.Lock()
            logBuffer = append(logBuffer, query)

            if len(logBuffer) > 10 {
                clog.Infof("log buffer contains %d queries; saving to database", len(logBuffer))
            }

            logBufferLock.Unlock()

            timer.Reset(time.Second * 10)

        case <-timer.C:
            logBufferLock.Lock()

            // Check buffer
            if len(logBuffer) == 0 {
                clog.Info("nothing new to log; continuing")
                continue
            } else {
                clog.Infof("timer expired and log buffer contains %d queries; saving to database", len(logBuffer))
            }

            // Copy buffer
            buffer := make([]queryContext, len(logBuffer))
            copy(buffer, logBuffer)
            logBuffer = []queryContext{}

            logBufferLock.Unlock()

            // Save buffer
            saveQueries(buffer)
        }
    }
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
