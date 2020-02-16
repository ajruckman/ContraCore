package eventserver

import (
	"encoding/json"
	"net"

	. "github.com/ajruckman/xlib"

	"github.com/ajruckman/ContraCore/internal/schema"
	"github.com/ajruckman/ContraCore/internal/system"
)

var (
	transmitQueue = make(chan schema.Log)
	clients       = map[string]net.Conn{}
)

// Package setup function.
func Setup() {
	go transmitWorker()
}

// Pushes a log onto the transmit channel.
func Transmit(log schema.Log) {
	transmitQueue <- log
}

// Sets up a new client by sending the client every log in the cache and adding
// them to the client list.
func Onboard(conn net.Conn, initial []schema.Log) {
	clients[conn.RemoteAddr().String()] = conn

	system.Console.Infof("Sending client %d rows", len(initial))

	e := func(err error) {
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				system.Console.Info("Deleting disconnected client:", conn.RemoteAddr().String())
				delete(clients, conn.RemoteAddr().String())
			} else {
				Err(err)
			}
		}
	}

	_, err := conn.Write([]byte("initial\n"))
	e(err)

	for _, v := range initial {
		content := marshal(v)

		_, err := conn.Write(content)
		e(err)
	}

	_, err = conn.Write([]byte("!initial\n"))
	e(err)
}

// Marshals a log to a JSON byte slice.
func marshal(log schema.Log) []byte {
	//var m *string
	//if log.ClientMAC != nil {
	//    s := log.ClientMAC.String()
	//    m = &s
	//}

	//l := struct {
	//    Time           time.Time
	//    Client         string
	//    Question       string
	//    QuestionType   string
	//    Action         string
	//    Answers        []string
	//    ClientMAC      *string
	//    ClientHostname *string
	//    ClientVendor   *string
	//    QueryID        uint16
	//}{
	//    Time:           log.Time,
	//    Client:         log.Client,
	//    Question:       log.Question,
	//    QuestionType:   log.QuestionType,
	//    Action:         log.Action,
	//    Answers:        log.Answers,
	//    ClientMAC:      log.ClientMAC,
	//    ClientHostname: log.ClientHostname,
	//    ClientVendor:   log.ClientVendor,
	//    QueryID:        log.QueryID,
	//}

	content, err := json.Marshal(log)
	Err(err)

	return append(content, '\n')
}

// Sends logs in the transmit queue to every connected client, and removes
// disconnected clients from the client list.
func transmitWorker() {
	for queuedLog := range transmitQueue {
		content := marshal(queuedLog)

		for _, conn := range clients {
			_, err := conn.Write(content)

			if err != nil {
				if _, ok := err.(*net.OpError); ok {
					system.Console.Info("Deleting disconnected client:", conn.RemoteAddr().String())
					delete(clients, conn.RemoteAddr().String())
				} else {
					Err(err)
				}
			}
		}
	}
}
