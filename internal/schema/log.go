package schema

import (
    "time"
)

type Log struct {
    ID           int       `db:"ID"`
    Time         time.Time `db:"Time"`
    Client       string    `db:"Client"`
    Question     string    `db:"Question"`
    QuestionType string    `db:"QuestionType"`
    Answers      []string  `db:"Answers"`
}
