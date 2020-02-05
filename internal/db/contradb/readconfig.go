package contradb

import (
    "database/sql"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/config"
    "github.com/ajruckman/ContraCore/internal/system"
)

func readConfig() {
    conf, err := GetConfig()

    if err == sql.ErrNoRows {
        system.Console.Info("Generating default ContraDB config")

        InsertDefaultConfig()
        conf, err = GetConfig()
        Err(err)
    } else if err != nil {
        Err(err)
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

    config.RuleSources = ruleSources
    config.SearchDomains = searchDomains
    config.DomainNeeded = conf.DomainNeeded
    config.SpoofedA = conf.SpoofedA
    config.SpoofedAAAA = conf.SpoofedAAAA
    config.SpoofedCNAME = conf.SpoofedCNAME
    config.SpoofedDefault = conf.SpoofedDefault
}

func InsertDefaultConfig() {
    _, err := xdb.Exec(`INSERT INTO config (search_domains) VALUES(default);`)
    Err(err)
}
