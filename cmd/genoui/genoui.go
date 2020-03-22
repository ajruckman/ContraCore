package main

import (
	"github.com/ajruckman/ContraCore/internal/db/contradb"
	"github.com/ajruckman/ContraCore/internal/db/contradb/ouigen"
	"github.com/ajruckman/ContraCore/internal/system"
)

func main() {
	system.ContraDBURL = "postgres://contra_usr:EvPvkro59Jb7RK3o@10.3.0.16/contradb"
	contradb.Setup()

	ouigen.GenOUI()
}
