package contradb

import (
    "database/sql"
    "fmt"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/config"
)

func ReadConfig() {
    conf, err := GetConfig()

    if err == sql.ErrNoRows {
        fmt.Println("Generating default config")

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
    _, err := XDB.Exec(`INSERT INTO config (search_domains) VALUES(default);`)
    Err(err)
}
