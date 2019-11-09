package db

import (
    "github.com/jackc/pgx"
    "github.com/jackc/pgx/stdlib"
    "github.com/jmoiron/sqlx"

    . "github.com/ajruckman/xlib"
)

var (
    PDB *pgx.ConnPool
    XDB *sqlx.DB
)

func init() {
    var err error

    conf := pgx.ConnPoolConfig{
        ConnConfig: pgx.ConnConfig{
            Host:     "localhost",
            User:     "contradbmgr",
            Password: "contradbmgr",
            Database: "contradb",
        },
        MaxConnections: 32,
    }

    PDB, err = pgx.NewConnPool(conf)
    Err(err)

    XDB = sqlx.NewDb(stdlib.OpenDBFromPool(PDB), "pgx")

    PDB.Exec(`

CREATE TABLE IF NOT EXISTS "Log"
(
    "ID"           BIGSERIAL NOT NULL
        CONSTRAINT "Log_pk" PRIMARY KEY,
    "Time"         TIMESTAMP DEFAULT now(),
    "Client"       TEXT,
    "Question"     TEXT,
    "QuestionType" TEXT,
    "Answers"      TEXT[]
);

`)
}
