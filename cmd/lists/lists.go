package main

import (
    "container/list"
    "fmt"
    "time"

    "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/db/contralog"
    "github.com/ajruckman/ContraCore/internal/db/contralog/dbschema"
)

func main() {
    config.ContraLogURL = "tcp://10.3.0.16:9000?username=contralogmgr&password=contralogmgr&database=contralog"
    contralog.Setup()

    _, _ = contralog.GetLastNLogs(1)

    began := time.Now()
    logs, err := contralog.GetLastNLogs(10000)
    xlib.Err(err)
    fmt.Println(time.Since(began))

    linkedLogs := list.New()
    for _, v := range logs {
       linkedLogs.PushFront(v)
    }

    _ = logs

    return

    for true {
        n := dbschema.Log{
            Time:           time.Now(),
            Client:         "1.2.3.4",
            Question:       "google.com",
            QuestionType:   "A",
            Action:         "pass.notblacklisted",
            Answers:        []string{"1.1.1.1", "2.2.2.2"},
            ClientMAC:      nil,
            ClientHostname: nil,
            ClientVendor:   nil,
        }

        //began = time.Now()
        //logs = append(logs, n)
        //logs = logs[1:]
        //fmt.Println(time.Since(began))

        //began = time.Now()
        //logs = append([]dbschema.Log{n}, logs...)
        //logs = logs[:10000]
        //fmt.Println(time.Since(began))

        began = time.Now()
        linkedLogs.PushFront(n)
        linkedLogs.Remove(linkedLogs.Back())
        fmt.Println(time.Since(began))
    }

    //fmt.Println(logs)
}
