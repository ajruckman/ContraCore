package schema

import (
    "time"
)

type Log struct {
    ID           int       `db:"id"`
    Time         time.Time `db:"time"`
    Client       string    `db:"client"`
    Question     string    `db:"question"`
    QuestionType string    `db:"question_type"`
    Answers      []string  `db:"answers"`
}
