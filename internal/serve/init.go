package serve

func init() {
    go logWorker()
    go logMonitor()

    cacheDHCP()
    go dhcpRefreshWorker()

    readRules()
}
