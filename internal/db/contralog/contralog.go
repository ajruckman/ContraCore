// Package contralog contains code for interacting with ContraLog.
package contralog

import (
	"database/sql/driver"
	"net"
	"net/url"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"go.uber.org/atomic"

	"github.com/ajruckman/ContraCore/internal/system"
)

var (
	cdb        *sqlx.DB
	failedOnce atomic.Bool
)

// Package setup function.
func Setup() {
	var err error
	var dbURL *url.URL

	dbURL, err = url.Parse(system.ContraLogURL)
	if err != nil {
		system.Console.Error("invalid ContraLog database URL")
		panic(err)
	}
	system.Console.Info("ContraLog address: ", dbURL.Host)

	connect()

	ping()

	go monitor()
}

// Attempts to connect to ContraLog.
func connect() {
	var err error

	cdb, err = sqlx.Connect("clickhouse", system.ContraLogURL)
	if err != nil {
		if !failedOnce.Load() {
			system.Console.Error("failed to connect to ClickHouse database server with error:")
			system.Console.Error(err.Error())
			system.ContraLogOnline.Store(false)
			failedOnce.Store(true)
		}
	} else {
		system.Console.Info("connected to ClickHouse database server")
		system.ContraLogOnline.Store(true)
		failedOnce.Store(false)
	}
}

func monitor() {
	for range time.Tick(time.Second * 10) {
		ping()
	}
}

type ErrContraLogOffline struct {
}

func (e *ErrContraLogOffline) Error() string {
	return "ContraLog is disconnected"
}

func checkOfflineError(err error) bool {
	_, isOpErr := err.(*net.OpError)
	return isOpErr || err == driver.ErrBadConn
}

func errOfflineOrOriginal(err error) error {
	if checkOfflineError(err) {
		offline(err)
		return &ErrContraLogOffline{}
	} else {
		return err
	}
}

// Pings ContraLog to trigger online/offline code.
func ping() {
	var err error

	if cdb == nil {
		connect()

		return
	}

	err = cdb.Ping()
	if err != nil {
		offline(err)
	} else {
		online()
	}
}

func offline(err error) {
	if system.ContraLogOnline.Load() {
		if !checkOfflineError(err) {
			system.Console.Error("failed to ping ContraLog with unanticipated error:")
			system.Console.Error(err.Error())
		} else {
			system.Console.Error("failed to ping ContraLog because it is offline")
		}
		system.ContraLogOnline.Store(false)
	}

	if !checkOfflineError(err) {
		panic(err)
	}
}

func online() {
	if !system.ContraLogOnline.Load() {
		system.Console.Info("ClickHouse health check succeeded")
		system.ContraLogOnline.Store(true)
	}
}
