package dbschema

type Blacklist struct {
	ID      int    `db:"id"`
	Pattern string `db:"pattern"`
	Class   int    `db:"class"`
	Domain  string `db:"domain"`
	TLD     string `db:"tld"`
	SLD     string `db:"sld"`
}