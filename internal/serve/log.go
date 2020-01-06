package serve

import (
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

    logMonInterval = 30
    logCount       atomic.Uint32

    logAnsweredTotDuration atomic.Duration
    logAnsweredTotCount    atomic.Uint32
    logPassedTotDuration   atomic.Duration
    logPassedTotCount      atomic.Uint32

    dhcpRefreshInterval = 15

    logDurations = true
)

func logWorker() {
    for q := range logChannel {

        began := time.Now()
        v := schema.Log{
            Time:         q.received,
            Client:       q._client,
            Question:     q._domain,
            QuestionType: dns.TypeToString[q._qu.Qtype],
            Action:       q.action,
            Answers:      q.answers,

            QueryID:  q.r.Id,
            Duration: time.Now().Sub(q.received),

            ClientMAC:      q.mac,
            ClientHostname: q.hostname,
            ClientVendor:   q.vendor,
        }
        q.durations.timeGenLogStruct = time.Since(began)

        began = time.Now()
        err := db.Log(v)
        if err != nil {
            spew.Dump(v)

            clog.Warningf("could not insert log for query '%s'", v.Question)
            clog.Warning(err.Error())
        }
        q.durations.timeSaveLogToPG = time.Since(began)

        began = time.Now()
        err = db.LogC(v)
        if err != nil {
            clog.Warningf("could not insert secondary log for query '%s'", v.Question)
            clog.Warning(err.Error())
        }
        q.durations.timeSaveLogToCH = time.Since(began)

        began = time.Now()
        eventserver.Tick(v)
        q.durations.timeSendLogToEventClients = time.Since(began)

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

        logCount.Inc()

        if v.Action == "pass" {
            logPassedTotCount.Inc()
            logPassedTotDuration.Add(v.Duration)
        } else {
            logAnsweredTotCount.Inc()
            logAnsweredTotDuration.Add(v.Duration)
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
