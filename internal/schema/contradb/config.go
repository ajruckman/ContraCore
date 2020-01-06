package contradb

import (
    "github.com/jackc/pgx/pgtype"
)

type Config struct {
    ID            int              `db:"id"`
    Sources       pgtype.TextArray `db:"sources"`
    SearchDomains pgtype.TextArray `db:"search_domains"`
    DomainNeeded  bool             `db:"domain_needed"`
}
