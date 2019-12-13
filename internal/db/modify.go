package db

import (
    . "github.com/ajruckman/xlib"
)

func InsertDefaultConfig() {
    _, err := XDB.Exec(`INSERT INTO config (search_domains) VALUES(default);`)
    Err(err)
}
