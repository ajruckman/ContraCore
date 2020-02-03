package schema

import (
    "net"
    "time"

    "github.com/ajruckman/ContraCore/internal/db/contralog/dbschema"
)

type Log struct {
    ID             int `json:"-"`
    Time           time.Time
    Client         string
    Question       string
    QuestionType   string
    Action         string
    Answers        []string
    ClientMAC      *net.HardwareAddr
    ClientHostname *string
    ClientVendor   *string

    QueryID  uint16
    Duration time.Duration `json:"-"`
}

func LogsFromContraLogs(contraLogs []dbschema.Log) (res []Log) {
    for _, l := range contraLogs {
        res = append(res, Log{
            ID:             l.ID,
            Time:           l.Time,
            Client:         l.Client,
            Question:       l.Question,
            QuestionType:   l.QuestionType,
            Action:         l.Action,
            Answers:        l.Answers,
            ClientMAC:      l.ClientMAC,
            ClientHostname: l.ClientHostname,
            ClientVendor:   l.ClientVendor,
            QueryID:        l.QueryID,
        })
    }
    return
}

func (log Log) ToContraLog() dbschema.Log {
    return dbschema.Log{
        ID:             log.ID,
        Time:           log.Time,
        Client:         log.Client,
        Question:       log.Question,
        QuestionType:   log.QuestionType,
        Action:         log.Action,
        Answers:        log.Answers,
        ClientMAC:      log.ClientMAC,
        ClientHostname: log.ClientHostname,
        ClientVendor:   log.ClientVendor,
        QueryID:        log.QueryID,
    }
}
