package netmgr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	. "github.com/ajruckman/xlib"

	"github.com/ajruckman/ContraCore/internal/cache"
	"github.com/ajruckman/ContraCore/internal/db/contradb"
	"github.com/ajruckman/ContraCore/internal/rule"
	"github.com/ajruckman/ContraCore/internal/system"
)

var clients = map[string]client{}

type client struct {
	net.Conn
	address string
}

func listen() {
	system.Console.Info("netmgr: Listening")
	ln, err := net.Listen("tcp", "0.0.0.0:64417")
	Err(err)

	for {
		conn, err := ln.Accept()
		Err(err)

		c := client{
			conn,
			conn.RemoteAddr().String(),
		}

		system.Console.Infof("netmgr: New client: %s", c.address)

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
				system.Console.Warningf("unhandled error encountered reading from client %s:", c.address)
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
	handleErr := func(err error, msg string, status string) {
		if err != nil {
			system.Console.Error(msg + ":")
			system.Console.Error(err)

			if status != "" {
				err = sendString(c, status)
				if err != nil {
					system.Console.Errorf("netmgr: Failed to send error status code to client %s:", c.address)
					system.Console.Error(err)
				}
			}
		}
	}

	cmd := strings.Split(line, " ")[0]
	cmd = cmd[:len(cmd)-1]

	switch cmd {
	case "ping":
		err := sendString(c, "ping.pong")
		handleErr(err, "netmgr: failed to send pong to client "+c.address, "")

	// TODO: pong

	case "get_cache":
		data, err := json.Marshal(logCache)
		handleErr(err, "netmgr: Error serializing cache", "")
		if err != nil {
			return
		}
		data = append([]byte("get_cache.cache "), data...)

		err = sendBytes(c, data)
		handleErr(err, "netmgr: Failed to send cache to client "+c.address, "")
		if err != nil {
			return
		}

	case "gen_rules":
		if online := system.ContraDBOnline.Load(); !online {
			handleErr(&contradb.ErrContraDBOffline{}, "netmgr: gen_rules command received but ContraDB is offline; doing nothing", "gen_rules.contradb_offline")
			return
		}

		// Get sources
		sources := system.RuleSources
		msg := fmt.Sprintf("Generating rules from %d sources: %s", len(sources), strings.Join(sources, ", "))
		system.Console.Infof("netmgr: " + msg)

		data := append([]byte("gen_rules.sources "), []byte(msg)...)
		err := sendBytes(c, data)
		handleErr(err, "netmgr: Failed to send sources to client", "")
		if err != nil {
			return
		}

		// Generate rules
		rules, _ := rule.GenFromURLs(system.RuleSources, func(progress string) bool {
			system.Console.Infof("netmgr: gen_rules: %s", progress)
			err = sendString(c, fmt.Sprintf("gen_rules.gen_progress %s", progress))
			return false
		})

		// Save rules
		begin := time.Now()
		rule.SaveRules(rules, func(progress string) bool {
			system.Console.Infof("netmgr: gen_rules: %s", progress)
			err = sendString(c, fmt.Sprintf("gen_rules.save_progress %s", progress))
			return false
		})
		end := time.Now()

		msg = "Blacklist rules saved in " + end.Sub(begin).String()
		system.Console.Infof("netmgr: " + msg)
		err = sendString(c, "gen_rules.saved_in " + msg)
		handleErr(err, "netmgr: Failed to send rule save time to client "+c.address, "")
		if err != nil {
			return
		}

		// Re-cache blacklist rules
		begin = time.Now()
		cache.ReadBlacklist(func(progress string) bool {
			system.Console.Info(progress)
			err = sendString(c, fmt.Sprintf("gen_rules.recache_progress %s", progress))
			return err != nil
		})
		end = time.Now()

		msg = "Blacklist rules re-cached in " + end.Sub(begin).String()
		system.Console.Infof("netmgr: " + msg)
		err = sendString(c, fmt.Sprintf("gen_rules.recached_in "+msg))
		handleErr(err, "netmgr: Failed to send blacklist rule re-cache time to client "+c.address, "")
		if err != nil {
			return
		}

		// Complete
		err = sendString(c, fmt.Sprintf("gen_rules.complete"))
		handleErr(err, "netmgr: Failed to send gen_rules.complete message to client "+c.address, "")

	case "reload_config":
		if online := system.ContraDBOnline.Load(); !online {
			handleErr(&contradb.ErrContraDBOffline{}, "netmgr: reload_config command received but ContraDB is offline; doing nothing", "reload_config.contradb_offline")
			return
		}

		system.Console.Info("netmgr: Received reload_config; reloading config from ContraDB")
		contradb.ReadConfig()
		err := sendString(c, fmt.Sprintf("reload_config.complete"))
		handleErr(err, "netmgr: Failed to send reload_config.complete message to client "+c.address, "")

	default:
		system.Console.Warningf("netmgr: Unknown command received from c %s: '%s'", c.address, cmd)
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
