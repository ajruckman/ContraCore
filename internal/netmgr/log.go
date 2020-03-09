package netmgr

import (
	"encoding/json"
	"net"

	"github.com/ajruckman/ContraCore/internal/schema"
	"github.com/ajruckman/ContraCore/internal/system"
)

var transmitQueue = make(chan schema.Log)

func ProcessQuery(log schema.Log) {
	transmitQueue <- log

	cache = append(cache, log)
	if len(cache) > cacheSize {
		over := len(cache) - cacheSize
		cache = cache[over:]
	}
}

// Sends logs in the transmit queue to every connected client, and removes
// disconnected clients from the client list.
func transmitWorker() {
	for queuedLog := range transmitQueue {
		data, err := json.Marshal(queuedLog)
		if err != nil {
			system.Console.Warning("netmgr: error serializing query:")
			system.Console.Warning(err.Error())
			continue
		}

		data = append([]byte("query "), data...)

		for _, c := range clients {
			err = sendBytes(c, data)
			if err != nil {
				if _, ok := err.(*net.OpError); ok {
					system.Console.Info("netmgr: deleting disconnected client:", c.address)
					delete(clients, c.address)
				} else {
					system.Console.Warningf("netmgr: unhandled error encountered sending query to client %s:", c.address)
					system.Console.Warning(err.Error())
				}
			}
		}
	}
}
