package db

import (
    "fmt"

    "github.com/coredns/coredns/plugin/pkg/log"
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/stdlib"
    "github.com/jmoiron/sqlx"
    "go.uber.org/atomic"

    _ "github.com/ClickHouse/clickhouse-go"

    "github.com/ajruckman/ContraCore/cmd/vars"
)

var (
    PDB *pgx.Conn
    XDB *sqlx.DB
    CDB *sqlx.DB

    PostgresOnline   atomic.Bool
    ClickHouseOnline atomic.Bool

    clog = log.NewWithPlugin("contradomain")
)

func Setup() {
    var err error

    fmt.Println(vars.ContraDBURL)
    fmt.Println(vars.ContraLogURL)

    //contraDBConf := `postgres://contradbmgr:contradbmgr@10.2.0.104/contradb`
    //contraDBConf := `postgres://contradbmgr:contradbmgr@10.3.0.16/contradb`
    //contraDBConf := `postgres://contradbmgr:contradbmgr@127.0.0.1/contradb`

    XDB, err = sqlx.Connect("pgx", vars.ContraDBURL)
    if err != nil {
        clog.Errorf("could not connect to PostgreSQL database server")
        panic(err)
    } else {
        PostgresOnline.Store(true)
    }

    if PostgresOnline.Load() {
        PDB, err = stdlib.AcquireConn(XDB.DB)
        if err != nil {
            clog.Error("could not acquire SQLX connection")
            panic(err)
        }
    }

    //contraLogConf := `tcp://10.3.0.16:9000?username=contralogmgr&password=contralogmgr&database=contralog`
    CDB, err = sqlx.Connect("clickhouse", vars.ContraLogURL)
    if err != nil {
        clog.Error("could not connect to ClickHouse database server")
        panic(err)
    } else {
        ClickHouseOnline.Store(true)
    }
}
