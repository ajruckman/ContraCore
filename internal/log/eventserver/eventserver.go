package eventserver

import (
    "encoding/json"
    "net"
    "time"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/schema"
    "github.com/ajruckman/ContraCore/internal/system"
)

var (
    transmitQueue = make(chan schema.Log)
    clients       = map[string]net.Conn{}
)

func Setup() {
    go transmitWorker()
}

func Transmit(log schema.Log) {
    transmitQueue <- log
}

func Onboard(conn net.Conn, buffer []schema.Log) {
    clients[conn.RemoteAddr().String()] = conn

    for _, v := range buffer {
        content := marshal(v)

        _, err := conn.Write(content)
        if err != nil {
            if _, ok := err.(*net.OpError); ok {
                system.Console.Info("Deleting disconnected client:", conn.RemoteAddr().String())
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

func transmitWorker() {
    for queuedLog := range transmitQueue {
        content := marshal(queuedLog)

        for _, conn := range clients {
            _, err := conn.Write(content)

            if err != nil {
                if _, ok := err.(*net.OpError); ok {
                    system.Console.Info("Deleting disconnected client:", conn.RemoteAddr().String())
                    delete(clients, conn.RemoteAddr().String())
                } else {
                    Err(err)
                }
            }
        }
    }
}
