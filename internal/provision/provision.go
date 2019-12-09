package provision

import (
    "fmt"

    "github.com/ajruckman/ContraCore/internal/db"
)

func init() {
    fmt.Println(`Provisioning database`)
    _, err := db.PDB.Exec(`
CREATE TABLE IF NOT EXISTS log
(
    id            BIGSERIAL NOT NULL,
    time          TIMESTAMP DEFAULT now(),
    client        INET,
    question      TEXT,
    question_type TEXT,
    answers       TEXT[],

    CONSTRAINT log_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS rule
(
    id       SERIAL NOT NULL,
    domains4 TEXT   NOT NULL,
    tld      TEXT   NOT NULL,
    sld      TEXT,
    pattern  TEXT   NOT NULL,

    CONSTRAINT rules_pk PRIMARY KEY (id)
);
    `)
    if err != nil { panic(err) }
}