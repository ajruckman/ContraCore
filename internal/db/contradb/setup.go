package contradb

import (
    "net/url"

    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/stdlib"
    "github.com/jmoiron/sqlx"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/system"
)

var (
    pdb *pgx.Conn
    xdb *sqlx.DB
)

func Setup() {
    readConfig()

    var err error
    var dbURL *url.URL

    dbURL, err = dbURL.Parse(config.ContraDBURL)
    if err != nil {
        system.Console.Errorf("invalid ContraCore database URL")
        panic(err)
    }
    system.Console.Info("ContraDB address:  ", dbURL.Host+":"+dbURL.Port())

    xdb, err = sqlx.Connect("pgx", config.ContraDBURL)
    if err != nil {
        system.Console.Errorf("could not connect to PostgreSQL database server")
        panic(err)
    } else {
        system.PostgresOnline.Store(true)
    }

    if system.PostgresOnline.Load() {
        pdb, err = stdlib.AcquireConn(xdb.DB)
        if err != nil {
            system.Console.Error("could not acquire SQLX connection")
            panic(err)
        }
    }
}
