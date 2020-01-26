package eventserver

import (
    "encoding/json"
    "fmt"
    "net"
    "time"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/schema"
)

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
