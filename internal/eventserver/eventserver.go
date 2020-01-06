package eventserver

import (
    "encoding/json"
    "fmt"
    "net"
    "time"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/schema"
)

var clients = map[string]net.Conn{}

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

    //recent, err := db.GetLastNLogs(1000)
    //Err(err)
    //
    //for i := len(recent) - 1; i >= 0; i-- {
    //    content := marshal(recent[i])
    //
    //    _, err = conn.Write(content)
    //    if err != nil {
    //        if _, ok := err.(*net.OpError); ok {
    //            fmt.Println("Deleting broken client:", conn.RemoteAddr().String())
    //            delete(clients, conn.RemoteAddr().String())
    //        } else {
    //            fmt.Println(reflect.TypeOf(err))
    //            Err(err)
    //        }
    //    }
    //}
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
    }

    content, err := json.Marshal(l)
    Err(err)

    return append(content, '\n')
}

func Tick(log schema.Log) {
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
