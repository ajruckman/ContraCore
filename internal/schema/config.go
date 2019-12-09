package schema

import (
    "github.com/jackc/pgx/pgtype"
)

type Config struct {
    ID            int              `db:"id"`
    SearchDomains pgtype.TextArray `db:"search_domains"`
    DomainNeeded  bool             `db:"domain_needed"`
}
