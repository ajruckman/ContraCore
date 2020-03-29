package dbschema

import (
	"time"
)

type Blacklist struct {
	ID      int        `db:"id"`
	Pattern string     `db:"pattern"`
	Expires *time.Time `db:"expires"`
	Class   int        `db:"class"`
	Domain  string     `db:"domain"`
	TLD     string     `db:"tld"`
	SLD     string     `db:"sld"`
}
