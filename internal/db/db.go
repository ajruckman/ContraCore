package db

import (
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/stdlib"
    "github.com/jmoiron/sqlx"

    . "github.com/ajruckman/xlib"
)

var (
    PDB *pgx.Conn
    XDB *sqlx.DB
)

func init() {
    conf := `postgres://contradbmgr:contradbmgr@10.3.0.16/contradb`
    var err error

    XDB, err = sqlx.Connect("pgx", conf)
    Err(err)

    PDB, err = stdlib.AcquireConn(XDB.DB)
    Err(err)
}
