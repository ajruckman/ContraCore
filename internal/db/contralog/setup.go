package contralog

import (
    "net/url"

    _ "github.com/ClickHouse/clickhouse-go"
    "github.com/jmoiron/sqlx"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/system"
)

var (
    cdb *sqlx.DB
)

func Setup() {
    var err error
    var dbURL *url.URL

    dbURL, err = dbURL.Parse(config.ContraLogURL)
    if err != nil {
        system.Console.Errorf("invalid ContraLog database URL")
        panic(err)
    }
    system.Console.Info("ContraLog address: ", dbURL.Host+":"+dbURL.Port())

    cdb, err = sqlx.Connect("clickhouse", config.ContraLogURL)
    if err != nil {
        system.Console.Error("could not connect to ClickHouse database server")
        panic(err)
    } else {
        system.ClickHouseOnline.Store(true)
    }
}
