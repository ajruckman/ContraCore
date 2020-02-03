package contralog

import (
    "fmt"

    _ "github.com/ClickHouse/clickhouse-go"
    "github.com/jmoiron/sqlx"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/state"
)

var (
    CDB *sqlx.DB
)

func Setup() {
    var err error

    fmt.Println("ContraLog URL:", config.ContraLogURL)

    CDB, err = sqlx.Connect("clickhouse", config.ContraLogURL)
    if err != nil {
        state.Console.Error("could not connect to ClickHouse database server")
        panic(err)
    } else {
        state.ClickHouseOnline.Store(true)
    }
}
