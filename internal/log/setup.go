package log

// Package setup function.
func Setup() {
	go queryBufferFlushScheduled()
	//go listen()
	go statWorker()
	go inputMonitor()
}
