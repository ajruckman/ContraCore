package schema

type Rule struct {
    ID      int    `db:"id"`
    Pattern string `db:"pattern"`
    Domain  string `db:"domain"`
    TLD     string `db:"tld"`
    SLD     string `db:"sld"`
}
