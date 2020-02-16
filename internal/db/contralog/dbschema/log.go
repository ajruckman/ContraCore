package dbschema

import (
	"time"
)

type Log struct {
	Time           time.Time `db:"time"`
	Client         string    `db:"client"`
	Question       string    `db:"question"`
	QuestionType   string    `db:"question_type"`
	Action         string    `db:"action"`
	Answers        []string  `db:"answers"`
	ClientMAC      *string   `db:"client_mac"`
	ClientHostname *string   `db:"client_hostname"`
	ClientVendor   *string   `db:"client_vendor"`
	QueryID        uint16    `db:"query_id"`
}
