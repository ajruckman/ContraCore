package internal

import (
	"fmt"
	"time"

	"github.com/caddyserver/caddy"

	"github.com/ajruckman/ContraCore/internal/cache"
	"github.com/ajruckman/ContraCore/internal/db/contradb"
	"github.com/ajruckman/ContraCore/internal/db/contralog"
	"github.com/ajruckman/ContraCore/internal/log"
	"github.com/ajruckman/ContraCore/internal/netmgr"
	"github.com/ajruckman/ContraCore/internal/process"
	"github.com/ajruckman/ContraCore/internal/system"
)

// The master ContraCore setup function. Calls other package setup functions.
func Setup(c *caddy.Controller) {
	fmt.Println("<><><>")
	if !system.HasInitialized.Load() {
		system.HasInitialized.Store(true)

		system.Console.Info("Setting up ContraDomain plugin")

		began := time.Now()

		system.Console.Info("ContraDomain setup: Parsing Corefile directive")
		system.ParseCorefile(c)

		system.Console.Info("ContraDomain setup: Setting up ContraDB")
		contradb.Setup()
		system.Console.Info("ContraDomain setup: Setting up ContraLog")
		contralog.Setup()

		system.Console.Info("ContraDomain setup: Setting up log system")
		log.Setup()

		system.Console.Info("ContraDomain setup: Setting up netmgr")
		go netmgr.Setup()

		//system.Console.Info("ContraDomain setup: running provisioner")
		//provision.Setup()

		//system.Console.Info("ContraDomain setup: caching whitelist rules")
		//cache.ReadWhitelist()
		system.Console.Info("ContraDomain setup: Caching blacklist rules")
		cache.ReadBlacklist(func(s string) bool {
			system.Console.Info("ContraDomain setup: " + s)
			return false
		})
		system.Console.Info("ContraDomain setup: Caching whitelist rules")
		cache.ReadWhitelist(func(s string, err error) bool {
			if err == nil {
				system.Console.Info("ContraDomain setup: " + s)
			} else {
				system.Console.Warningf("ContraDomain setup: " + s)
				system.Console.Warning(err.Error())
			}
			return false
		})

		system.Console.Info("ContraDomain setup: Setting up query processor")
		process.Setup()

		system.Console.Infof("ContraDomain setup: All systems set up in %v", time.Since(began))

	} else {
		system.Console.Info("ContraDomain plugin has already been set up in another zone; skipping setup")
	}
}
