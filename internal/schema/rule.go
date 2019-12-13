package schema

type Rule struct {
    ID      int    `db:"id"`
    Domain  string `db:"domain"`
    Pattern string `db:"pattern"`
}
