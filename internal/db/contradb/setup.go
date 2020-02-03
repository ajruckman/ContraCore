package contradb

import (
    "fmt"

    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/stdlib"
    "github.com/jmoiron/sqlx"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/state"
)

var (
    PDB *pgx.Conn
    XDB *sqlx.DB
)

func Setup() {
    var err error

    fmt.Println("ContraDB URL: ", config.ContraDBURL)

    XDB, err = sqlx.Connect("pgx", config.ContraDBURL)
    if err != nil {
        state.Console.Errorf("could not connect to PostgreSQL database server")
        panic(err)
    } else {
        state.PostgresOnline.Store(true)
    }

    if state.PostgresOnline.Load() {
        PDB, err = stdlib.AcquireConn(XDB.DB)
        if err != nil {
            state.Console.Error("could not acquire SQLX connection")
            panic(err)
        }
    }
}
