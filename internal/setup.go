package internal

import (
    "time"

    "github.com/caddyserver/caddy"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/db/contradb"
    "github.com/ajruckman/ContraCore/internal/db/contralog"
    "github.com/ajruckman/ContraCore/internal/log"
    "github.com/ajruckman/ContraCore/internal/process"
    "github.com/ajruckman/ContraCore/internal/provision"
    "github.com/ajruckman/ContraCore/internal/system"
)

func Setup(c *caddy.Controller) {
    system.HasInitializedLock.Lock()

    if !system.HasInitialized {
        system.HasInitialized = true
        system.HasInitializedLock.Unlock()

        system.Console.Info("setting up ContraDomain plugin")

        began := time.Now()

        system.Console.Info("ContraDomain setup: parsing Corefile directive")
        config.ParseCorefile(c)

        system.Console.Info("ContraDomain setup: setting up ContraDB")
        contradb.Setup()
        system.Console.Info("ContraDomain setup: setting up ContraLog")
        contralog.Setup()

        system.Console.Info("ContraDomain setup: setting up log system")
        log.Setup()

        system.Console.Info("ContraDomain setup: running provisioner")
        provision.Setup()

        system.Console.Info("ContraDomain setup: setting up query processor")
        process.Setup()

        system.Console.Infof("ContraDomain setup: all systems set up in %v", time.Since(began))

    } else {
        system.HasInitializedLock.Unlock()
        system.Console.Info("ContraDomain plugin has already been set up in another zone; skipping setup")
    }
}
