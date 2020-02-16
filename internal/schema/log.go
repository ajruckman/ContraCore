package schema

import (
	"time"

	"github.com/ajruckman/ContraCore/internal/db/contralog/dbschema"
)

// A query log record with additional internal fields.
type Log struct {
	dbschema.Log

	Duration time.Duration `json:"-"`
}

// Converts a slice of ContraLog log structs to a slice of wrapped log structs.
func LogsFromContraLogs(contraLogs []dbschema.Log) (res []Log) {
	for _, l := range contraLogs {
		res = append(res, Log{
			Log: dbschema.Log{
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
			},
		})
	}
	return
}
