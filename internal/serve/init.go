package serve

import (
    "github.com/ajruckman/ContraCore/internal/log"
)

func init() {
    go log.QueryLogWorker()
    go log.LogMonitor()

    cacheDHCP()
    go dhcpRefreshWorker()

    cacheOUI()

    readRules()

    go log.Serve()
}
