package netmgr

// Package setup function.
func Setup() {
	loadCache()

	go listen()
	go transmitWorker()
}
