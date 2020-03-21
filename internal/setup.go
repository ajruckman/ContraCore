package internal

import (
	"time"

	"github.com/caddyserver/caddy"

	"github.com/ajruckman/ContraCore/internal/cache"
	"github.com/ajruckman/ContraCore/internal/db/contradb"
	"github.com/ajruckman/ContraCore/internal/db/contralog"
	"github.com/ajruckman/ContraCore/internal/log"
	"github.com/ajruckman/ContraCore/internal/netmgr"
	"github.com/ajruckman/ContraCore/internal/process"
	"github.com/ajruckman/ContraCore/internal/provision"
	"github.com/ajruckman/ContraCore/internal/system"
)

// The master ContraCore setup function. Calls other package setup functions.
func Setup(c *caddy.Controller) {
	if !system.HasInitialized.Load() {
		system.HasInitialized.Store(true)

		system.Console.Info("setting up ContraDomain plugin")

		began := time.Now()

		system.Console.Info("ContraDomain setup: parsing Corefile directive")
		system.ParseCorefile(c)

		system.Console.Info("ContraDomain setup: setting up ContraDB")
		contradb.Setup()
		system.Console.Info("ContraDomain setup: setting up ContraLog")
		contralog.Setup()

		system.Console.Info("ContraDomain setup: setting up log system")
		log.Setup()

		system.Console.Info("ContraDomain setup: setting up netmgr")
		go netmgr.Setup()

		system.Console.Info("ContraDomain setup: running provisioner")
		provision.Setup()

		//system.Console.Info("ContraDomain setup: caching whitelist rules")
		//cache.ReadWhitelist()
		system.Console.Info("ContraDomain setup: caching blacklist rules")
		cache.ReadBlacklist(func(s string) bool {
			system.Console.Info("ContraDomain setup: " + s)
			return false
		})

		system.Console.Info("ContraDomain setup: setting up query processor")
		process.Setup()

		system.Console.Infof("ContraDomain setup: all systems set up in %v", time.Since(began))

	} else {
		system.Console.Info("ContraDomain plugin has already been set up in another zone; skipping setup")
	}
}
