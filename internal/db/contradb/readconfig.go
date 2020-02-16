package contradb

import (
	"database/sql"

	"go.uber.org/atomic"

	"github.com/ajruckman/ContraCore/internal/system"
)

var configLoaded = atomic.Bool{}

func readConfig() {
	conf, err := GetConfig()

	if err == sql.ErrNoRows {
		system.Console.Info("Generating default ContraDB config")

		err = insertDefaultConfig()
		if err != nil {
			system.Console.Error("failed to save default config to ContraDB; loading hardcoded config")
			system.Console.Error(err.Error())
			loadOfflineConfig()
			return
		}

		conf, err = GetConfig()
		if err != nil {
			system.Console.Error("failed to load saved default config from ContraDB; loading hardcoded config")
			system.Console.Error(err.Error())
			loadOfflineConfig()
			return
		}

	} else if err != nil {
		system.Console.Warning("failed to load default config from ContraDB; loading hardcoded config")
		system.Console.Warning(err.Error())
		loadOfflineConfig()
		return
	}

	//

	var ruleSources []string
	for _, v := range conf.Sources.Elements {
		ruleSources = append(ruleSources, v.String)
	}

	var searchDomains []string
	for _, v := range conf.SearchDomains.Elements {
		searchDomains = append(searchDomains, v.String)
	}

	//

	system.RuleSources = ruleSources
	system.SearchDomains = searchDomains
	system.DomainNeeded = conf.DomainNeeded
	system.SpoofedA = conf.SpoofedA
	system.SpoofedAAAA = conf.SpoofedAAAA
	system.SpoofedCNAME = conf.SpoofedCNAME
	system.SpoofedDefault = conf.SpoofedDefault

	configLoaded.Store(true)

	system.Console.Info("loaded config from ContraDB")
}

func loadOfflineConfig() {
	system.RuleSources = []string{}
	system.SearchDomains = []string{}
	system.DomainNeeded = true
	system.SpoofedA = "0.0.0.0"
	system.SpoofedAAAA = "::0"
	system.SpoofedCNAME = ""
	system.SpoofedDefault = "-"
}
