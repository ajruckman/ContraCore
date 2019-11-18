package db

//go:generate go run generate.go

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
            Host:     "10.3.0.16",
            User:     "contradbmgr",
            Password: "contradbmgr",
            Database: "contradb",
        },
        MaxConnections: 32,
    }

    PDB, err = pgx.NewConnPool(conf)
    Err(err)

    XDB = sqlx.NewDb(stdlib.OpenDBFromPool(PDB), "pgx")
}
