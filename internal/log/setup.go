package log

import (
    "github.com/ajruckman/ContraCore/internal/log/eventserver"
)

func Setup() {
    go queryBufferFlushScheduled()
    go listen()
    go statWorker()
    go inputMonitor()

    go eventserver.Setup()
}
