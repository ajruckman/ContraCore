package netmgr

import (
	"bufio"
	"encoding/json"
	"net"
	"strings"

	. "github.com/ajruckman/xlib"

	"github.com/ajruckman/ContraCore/internal/system"
)

var clients = map[string]client{}

type client struct {
	net.Conn
	address string
}

func listen() {
	system.Console.Info("netmgr: listening")
	ln, err := net.Listen("tcp", "0.0.0.0:64417")
	Err(err)

	for {
		conn, err := ln.Accept()
		Err(err)

		c := client{
			conn,
			conn.RemoteAddr().String(),
		}

		system.Console.Infof("netmgr: new client: %s", c.address)

		clients[c.address] = c
		go read(c)
	}
}

func read(c client) {
	r := bufio.NewReader(c)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if _, ok := err.(*net.OpError); !ok {
				system.Console.Warningf("unhandled error encountered reading from c %s:", c.address)
				system.Console.Warning(err.Error())
				Err(err)
			} else {
				delete(clients, c.address)
				_ = c.Close()
				break
			}
		}

		system.Console.Debugf("netmgr: %s <- %s", c.address, line)
		interpret(c, line)
	}
}

func interpret(c client, line string) {
	cmd := strings.Split(line, " ")[0]
	cmd = cmd[:len(cmd)-1]

	switch cmd {
	case "ping":
		err := sendString(c, "pong")
		if err != nil {
			system.Console.Warningf("netmgr: failed to send pong to c %s:", c.address)
			system.Console.Warning(err)
		}

	// TODO: pong

	case "get_cache":
		data, err := json.Marshal(cache)
		if err != nil {
			system.Console.Warningf("netmgr: error serializing cache:")
			system.Console.Warning(err)
			return
		}

		data = append([]byte("cache "), data...)

		err = sendBytes(c, data)
		if err != nil {
			system.Console.Warningf("netmgr: failed to send cache to c %s:", c.address)
			system.Console.Warning(err)
			return
		}

	default:
		system.Console.Warningf("netmgr: unknown command received from c %s: '%s'", c.address, cmd)
	}
}

func sendString(c client, data string) (err error) {
	_, err = c.Write(append([]byte(data), '\n'))
	return
}

func sendBytes(c client, data []byte) (err error) {
	_, err = c.Write(append(data, '\n'))
	return
}
