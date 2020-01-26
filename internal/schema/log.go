package schema

import (
    "net"
    "time"
)

type Log struct {
    //UUID           string            `db:"uuid"`
    ID             int               `db:"id" json:"-"`
    Time           time.Time         `db:"time"`
    Client         string            `db:"client"`
    Question       string            `db:"question"`
    QuestionType   string            `db:"question_type"`
    Action         string            `db:"action"`
    Answers        []string          `db:"answers"`
    ClientMAC      *net.HardwareAddr `db:"client_mac"`
    ClientHostname *string           `db:"client_hostname"`
    ClientVendor   *string           `db:"client_vendor"`

    QueryID  uint16        `db:"-"`
    Duration time.Duration `db:"-" json:"-"`
}
