package process

import (
    "time"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/db/contradb"
    "github.com/ajruckman/ContraCore/internal/system"
)

func Setup() {
    began := time.Now()
    ipsSeen, hostnamesSeen, err := cacheDHCP()
    if _, ok := err.(*contradb.ErrContraDBOffline); ok {
        system.Console.Warning("failed to load lease cache because ContraDB is not connected")
    } else if err != nil {
        Err(err)
    } else {
        system.Console.Infof("DHCP lease cache loaded in %v; %d distinct IPs and %d distinct hostnames found", time.Since(began), ipsSeen, hostnamesSeen)
    }

    readWhitelistRules()
    readBlacklistRules()
    go dhcpRefreshWorker()
}
