package dbschema

type LogCountPerHour struct {
    Hour  string `db:"hour" json:"aggHour"`
    Count uint64 `db:"c" json:"queryCount"`
}
