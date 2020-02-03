package internal

import (
    "github.com/caddyserver/caddy"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/db/contradb"
    "github.com/ajruckman/ContraCore/internal/db/contralog"
    "github.com/ajruckman/ContraCore/internal/log"
    "github.com/ajruckman/ContraCore/internal/process"
    "github.com/ajruckman/ContraCore/internal/provision"
)

func Setup(c *caddy.Controller) {
    config.ParseCorefile(c)

    log.Setup()

    contradb.Setup()
    contralog.Setup()

    provision.Setup()

    process.Setup()
}
