package db

import (
    "github.com/coredns/coredns/plugin/pkg/log"
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/stdlib"
    "github.com/jmoiron/sqlx"
    "go.uber.org/atomic"

    _ "github.com/ClickHouse/clickhouse-go"
)

var (
    PDB *pgx.Conn
    XDB *sqlx.DB
    CDB *sqlx.DB

    PostgresOnline   atomic.Bool
    ClickHouseOnline atomic.Bool

    clog = log.NewWithPlugin("contradomain")
)

func init() {
    var err error

    //contraDBConf := `postgres://contradbmgr:contradbmgr@10.2.0.104/contradb`
    contraDBConf := `postgres://contradbmgr:contradbmgr@10.3.0.16/contradb`
    //contraDBConf := `postgres://contradbmgr:contradbmgr@127.0.0.1/contradb`

    XDB, err = sqlx.Connect("pgx", contraDBConf)
    if err != nil {
        clog.Errorf("could not connect to PostgreSQL database server")
    } else {
        PostgresOnline.Store(true)
    }

    if PostgresOnline.Load() {
        PDB, err = stdlib.AcquireConn(XDB.DB)
        if err != nil {
            clog.Error("could not acquire SQLX connection")
        }
    }

    contraLogConf := `tcp://10.3.0.16:9000?username=contralogmgr&password=contralogmgr&database=contralog`
    CDB, err = sqlx.Connect("clickhouse", contraLogConf)
    if err != nil {
        clog.Error("could not connect to ClickHouse database server")
    } else {
        ClickHouseOnline.Store(true)
    }
}
