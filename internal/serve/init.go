package serve

func init() {
    go logWorker()
    cacheDHCP()
    cacheRules()
}
