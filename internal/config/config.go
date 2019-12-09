package config

import (
    "database/sql"
    "fmt"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema"
)

var Config schema.Config

func init() {
    var err error

    Config, err = db.GetConfig()

    if err == sql.ErrNoRows || Config.ID == 0 {
        fmt.Println("Generating default config")

        db.InsertDefaultConfig()
        Config, err = db.GetConfig()
        Err(err)
    }
}
