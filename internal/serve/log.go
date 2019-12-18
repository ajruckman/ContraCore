package serve

import (
    "github.com/coredns/coredns/plugin/pkg/log"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema"
)

var (
    logChannel = make(chan schema.Log)
    clog       = log.NewWithPlugin("contradomain")
)

func logWorker() {
    for v := range logChannel {
        err := db.Log(v)
        if err != nil {
            clog.Warning("could not insert log for query '" + v.Question + "'")
            clog.Warning(err.Error())
        }

        clog.Info(v.Client, " <- ", v.QueryID, " ", v.QuestionType, " ", v.Duration)
    }
}
