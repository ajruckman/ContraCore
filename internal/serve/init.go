package serve

import (
    "github.com/ajruckman/ContraCore/internal/eventserver"
)

func init() {
    go logWorker()
    go logMonitor()

    cacheDHCP()
    go dhcpRefreshWorker()

    cacheOUI()

    readRules()

    go eventserver.Serve()
}
