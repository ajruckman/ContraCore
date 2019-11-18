package schema

import (
    "time"
)

type QuestionCountsPerHour struct {
    Hour  time.Time `db:"hour"`
    Count int       `db:"count"`
}
