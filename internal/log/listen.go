package log

import (
    "net"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/log/eventserver"
    "github.com/ajruckman/ContraCore/internal/schema"
    "github.com/ajruckman/ContraCore/internal/system"
)

func listen() {
    loadCache()

    ln, err := net.Listen("tcp", "0.0.0.0:64417")
    Err(err)

    for {
        conn, err := ln.Accept()
        Err(err)

        system.Console.Info("New client: ", conn.RemoteAddr().String())

        queryBufferLock.Lock()
        buffer := make([]schema.Log, len(queryBuffer))
        copy(buffer, queryBuffer)
        queryBufferLock.Unlock()

        go eventserver.Onboard(conn, cache)
    }
}
