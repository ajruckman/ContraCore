package process

func Setup() {
    cacheDHCP()
    readRules()
    go dhcpRefreshWorker()
}
