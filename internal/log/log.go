package log

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "math/rand"
    "net"
    "strings"
    "sync"
    "time"

    "github.com/coredns/coredns/plugin/pkg/log"
    "github.com/davecgh/go-spew/spew"
    "github.com/jackc/pgx/v4"
    "github.com/miekg/dns"
    "go.uber.org/atomic"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema"

    . "github.com/ajruckman/xlib"
)

type QueryInfo struct {
    dns.ResponseWriter
    R *dns.Msg

    QU_     dns.Question
    Domain_ string
    Client_ string

    Received time.Time
    Action   string

    MAC      *net.HardwareAddr
    Hostname *string
    Vendor   *string

    Answers []string

    Durations Durations
}

type Durations struct {
    TimeLookupLease            time.Duration
    TimeCheckRespondByHostname time.Duration
    TimeCheckRespondByPTR      time.Duration
    TimeCheckRespondWithBlock  time.Duration
    TimeGenLogStruct           time.Duration
    TimeSaveLogToPG            time.Duration
    TimeSaveLogToCH            time.Duration
    TimeSendLogToEventClients  time.Duration
}

func (q *QueryInfo) Respond(res *dns.Msg) (err error) {
    var answers []string
    for _, v := range res.Answer {
        answers = append(answers, rrToString(v))
    }
    q.Answers = answers

    LogChannel <- *q

    err = q.ResponseWriter.WriteMsg(res)
    return
}

func (q QueryInfo) WriteMsg(r *dns.Msg) error {
    return q.Respond(r)
}

func rt(in string) string {
    return strings.TrimSuffix(in, ".")
}

// coredns/plugin/test/helpers.go
func rrToString(val dns.RR) string {
    var res string

    switch x := val.(type) {
    case *dns.SRV:
        res = fmt.Sprintf("%d|%d|%d|%s", x.Priority, x.Weight, x.Port, x.Target)

    case *dns.RRSIG:
        res = fmt.Sprintf("%d|%d|%s", x.TypeCovered, x.Labels, x.SignerName)

    case *dns.NSEC:
        res = x.NextDomain

    case *dns.A:
        res = rt(x.A.String())

    case *dns.AAAA:
        res = rt(x.AAAA.String())

    case *dns.TXT:
        res = strings.Join(x.Txt, "|")

    case *dns.HINFO:
        res = fmt.Sprintf("%s|%s", x.Cpu, x.Os)

    case *dns.SOA:
        res = x.Ns

    case *dns.PTR:
        res = rt(x.Ptr)

    case *dns.CNAME:
        res = rt(x.Target)

    case *dns.MX:
        res = fmt.Sprintf("%s|%d", x.Mx, x.Preference)

    case *dns.NS:
        res = x.Ns

    case *dns.OPT:
        res = fmt.Sprintf("%d|%t", x.UDPSize(), x.Do())
    }

    return res
}

var (
    ctx = context.Background()

    LogChannel = make(chan QueryInfo)
    CLOG       = log.NewWithPlugin("contradomain")

    QueryLogBuffer              []schema.Log
    QueryLogBufferLock          sync.Mutex
    queryLogBufferSaveThreshold = 15               // Save all in buffer if the buffer contains this many queries
    queryLogBufferFlushInterval = time.Second * 30 // Save all in buffer if no new logs have been added after this time

    logMonInterval = time.Second * 15
    logMonCount    atomic.Uint32
    logMonAnyNew   atomic.Bool

    logAnsweredTotDuration atomic.Duration
    logAnsweredTotCount    atomic.Uint32
    logPassedTotDuration   atomic.Duration
    logPassedTotCount      atomic.Uint32

    LogDurations = true
)

func QueryLogWorker() {
    var (
        query QueryInfo
        timer = time.NewTimer(queryLogBufferFlushInterval)
    )

    // Debounce code from: https://drailing.net/2018/01/debounce-function-for-golang/
    for {
        select {
        case query = <-LogChannel:
            QueryLogBufferLock.Lock()
            v := schema.Log{
                Time:         query.Received,
                Client:       query.Client_,
                Question:     query.Domain_,
                QuestionType: dns.TypeToString[query.QU_.Qtype],
                Action:       query.Action,
                Answers:      query.Answers,

                //QueryID: query.r.Id, // This isn't unique
                QueryID: uint16(rand.Intn(65536)),

                Duration: time.Now().Sub(query.Received),

                ClientMAC:      query.MAC,
                ClientHostname: query.Hostname,
                ClientVendor:   query.Vendor,
            }
            CLOG.Infof("%s <- %d %s %s %v", v.Client, v.QueryID, v.Question, v.QuestionType, v.Duration)

            QueryLogBuffer = append(QueryLogBuffer, v)
            if len(QueryLogBuffer) > queryLogBufferSaveThreshold {
                CLOG.Infof("log buffer contains %d queries (more than threshold %d); saving to database", len(QueryLogBuffer), queryLogBufferSaveThreshold)

                // Copy buffer
                buffer := make([]schema.Log, len(QueryLogBuffer))
                copy(buffer, QueryLogBuffer)
                QueryLogBuffer = []schema.Log{}

                // Save buffer
                saveQueryLogBuffer(buffer)
            }

            QueryLogBufferLock.Unlock()

            timer.Reset(queryLogBufferFlushInterval)

        case <-timer.C:
            QueryLogBufferLock.Lock()

            // Check buffer
            if len(QueryLogBuffer) == 0 {
                QueryLogBufferLock.Unlock()
                CLOG.Info("nothing new to log; continuing")
                timer.Reset(queryLogBufferFlushInterval)
                continue
            } else {
                CLOG.Infof("timer expired and log buffer contains %d queries; saving to database", len(QueryLogBuffer))
            }

            // Copy buffer
            buffer := make([]schema.Log, len(QueryLogBuffer))
            copy(buffer, QueryLogBuffer)
            QueryLogBuffer = []schema.Log{}

            // Save buffer
            saveQueryLogBuffer(buffer)

            QueryLogBufferLock.Unlock()

            timer.Reset(queryLogBufferFlushInterval)
        }
    }
}

func saveQueryLogBuffer(buffer []schema.Log) {
    logMonAnyNew.Store(true)

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
            CLOG.Warningf("could not begin PDB transaction")
            CLOG.Warning(err.Error())
        }
    }

    if db.ClickHouseOnline.Load() {
        cdbTX, cdbSTMT, err = db.LogCBeginBatch()
        if err != nil {
            CLOG.Warningf("could not begin CDB transaction")
            CLOG.Warning(err.Error())
        }
    }

    // Save queries
    for _, v := range buffer {

        //began := time.Now()

        //q.durations.timeGenLogStruct = time.Since(began)

        //began = time.Now()
        Transmit(v)
        //q.durations.timeSendLogToEventClients = time.Since(began)

        if pdbTX != nil {
            //began = time.Now()
            err := db.Log(pdbTX, v)
            if err != nil {
                spew.Dump(v)

                CLOG.Warningf("could not insert log for query '%s'", v.Question)
                CLOG.Warning(err.Error())
            }
            //q.durations.timeSaveLogToPG = time.Since(began)
        } else {
            CLOG.Warningf("not logging query '%s' because pdbTX is nil", v.Question)
        }

        if cdbTX != nil {
            //began = time.Now()
            err = db.LogC(cdbSTMT, v)
            if err != nil {
                CLOG.Warningf("could not insert secondary log for query '%s'", v.Question)
                CLOG.Warning(err.Error())
            }
            //q.durations.timeSaveLogToCH = time.Since(began)
        }

        //if !LogDurations {
        //    CLOG.Infof("%s <- %d %s %s %v", v.Client, v.QueryID, v.Question, v.QuestionType, v.Duration)
        //}
        //else {
        //    CLOG.Infof("%s <- %d %s %s %v\n\tLookup lease:        %v\n\tRespond by hostname: %v\n\tRespond by PTR:      %v\n\tRespond with block:  %v\n\tGen log struct:      %v\n\tSave to PG:          %v\n\tSave to CH:          %v\n\tSend to clients:     %v",
        //        v.Client,
        //        v.QueryID,
        //        v.Question,
        //        v.QuestionType,
        //        v.Duration,
        //        q.durations.timeLookupLease,
        //        q.durations.timeCheckRespondByHostname,
        //        q.durations.timeCheckRespondByPTR,
        //        q.durations.timeCheckRespondWithBlock,
        //        q.durations.timeGenLogStruct,
        //        q.durations.timeSaveLogToPG,
        //        q.durations.timeSaveLogToCH,
        //        q.durations.timeSendLogToEventClients,
        //    )
        //}

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
            CLOG.Warningf("could not commit pdbTX: %s", err.Error())
        }
    }

    if cdbTX != nil {
        err = cdbTX.Commit()
        if err != nil {
            CLOG.Warningf("could not commit cdbTX: %s", err.Error())
        }
    }

    logMonCount.Add(uint32(len(buffer)))
}

func LogMonitor() {
    for range time.Tick(logMonInterval) {
        if !logMonAnyNew.Swap(false) {
            continue
        }

        c := logMonCount.Swap(0)

        QueryLogBufferLock.Lock()
        queryLogBufferLen := len(QueryLogBuffer)
        QueryLogBufferLock.Unlock()

        var (
            avgDurAns  = float64(logAnsweredTotDuration.Swap(0).Milliseconds()) / float64(logAnsweredTotCount.Swap(0))
            avgDurPass = float64(logPassedTotDuration.Swap(0).Milliseconds()) / float64(logPassedTotCount.Swap(0))
        )

        CLOG.Infof("Log channel backlog: %d | New log rows: %d | Rows/second: %.3f | Avg. ms answered reqs.: %.2f | Avg. ms passed reqs.: %.2f",
            queryLogBufferLen,
            c,
            float64(c)/logMonInterval.Seconds(),
            avgDurAns,
            avgDurPass,
        )
    }
}

var (
    clients = map[string]net.Conn{}
    queue   = make(chan schema.Log)
)

func init() {
    go transmitWorker()
}

func Serve() {
    listen()
}

func listen() {
    ln, err := net.Listen("tcp", "0.0.0.0:64417")
    Err(err)

    for {
        conn, err := ln.Accept()
        Err(err)

        fmt.Println("New client:", conn.RemoteAddr().String())
        go setup(conn)
    }
}

func setup(conn net.Conn) {
    clients[conn.RemoteAddr().String()] = conn

    QueryLogBufferLock.Lock()

    // Copy buffer
    buffer := make([]schema.Log, len(QueryLogBuffer))
    copy(buffer, QueryLogBuffer)
    QueryLogBuffer = []schema.Log{}

    QueryLogBufferLock.Unlock()

    for _, v := range buffer {
        content := marshal(v)

        _, err := conn.Write(content)
        if err != nil {
            if _, ok := err.(*net.OpError); ok {
                fmt.Println("Deleting disconnected client:", conn.RemoteAddr().String())
                delete(clients, conn.RemoteAddr().String())
            } else {
                Err(err)
            }
        }
    }
}

func marshal(log schema.Log) []byte {
    var m *string
    if log.ClientMAC != nil {
        s := log.ClientMAC.String()
        m = &s
    }

    l := struct {
        Time           time.Time
        Client         string
        Question       string
        QuestionType   string
        Action         string
        Answers        []string
        ClientMAC      *string
        ClientHostname *string
        ClientVendor   *string
        QueryID        uint16
    }{
        Time:           log.Time,
        Client:         log.Client,
        Question:       log.Question,
        QuestionType:   log.QuestionType,
        Action:         log.Action,
        Answers:        log.Answers,
        ClientMAC:      m,
        ClientHostname: log.ClientHostname,
        ClientVendor:   log.ClientVendor,
        QueryID:        log.QueryID,
    }

    content, err := json.Marshal(l)
    Err(err)

    return append(content, '\n')
}

func Transmit(log schema.Log) {
    queue <- log
}

func transmitWorker() {
    for log := range queue {
        content := marshal(log)

        for _, conn := range clients {
            _, err := conn.Write(content)

            if err != nil {
                if _, ok := err.(*net.OpError); ok {
                    fmt.Println("Deleting disconnected client:", conn.RemoteAddr().String())
                    delete(clients, conn.RemoteAddr().String())
                } else {
                    Err(err)
                }
            }
        }
    }
}
