package db

import (
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/stdlib"
    "github.com/jmoiron/sqlx"

    _ "github.com/ClickHouse/clickhouse-go"

    . "github.com/ajruckman/xlib"
)

var (
    PDB *pgx.Conn
    XDB *sqlx.DB
    CDB *sqlx.DB
)

func init() {
    var err error

    //contraDBConf := `postgres://contradbmgr:contradbmgr@10.2.0.104/contradb`
    //contraDBConf := `postgres://contradbmgr:contradbmgr@10.3.0.16/contradb`
    contraDBConf := `postgres://contradbmgr:contradbmgr@127.0.0.1/contradb`

    XDB, err = sqlx.Connect("pgx", contraDBConf)
    Err(err)

    PDB, err = stdlib.AcquireConn(XDB.DB)
    Err(err)

    contraLogConf := `tcp://10.3.0.16:9000?username=contralogmgr&password=contralogmgr&database=contralog`
    CDB, err = sqlx.Connect("clickhouse", contraLogConf)
    Err(err)
}
