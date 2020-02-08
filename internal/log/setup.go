package log

import (
    "github.com/ajruckman/ContraCore/internal/log/eventserver"
)

func Setup() {
    go listen()
    go queryBufferFlushScheduled()
    go statWorker()
    go inputMonitor()

    go eventserver.Setup()
}
