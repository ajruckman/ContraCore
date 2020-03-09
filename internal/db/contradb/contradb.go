// Package contradb contains functions for interacting with ContraDB.
package contradb

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/atomic"

	"github.com/ajruckman/ContraCore/internal/system"
)

var (
	xdb        *sqlx.DB    // The sqlx connection to ContraDB.
	pdb        *pgx.Conn   // The PGX connection to ContraDB.
	failedOnce atomic.Bool // Whether the last connection failed.
)

// Package setup function.
func Setup() {
	dbURL, err := url.Parse(system.ContraDBURL)
	if err != nil {
		system.Console.Error("invalid ContraCore database URL")
		panic(err)
	}
	system.Console.Info("ContraDB address: ", dbURL.Host)

	connect()

	readConfig()

	ping()

	go monitor()
}

// Attempts to connect to ContraDB.
func connect() {
	var err error

	xdb, err = sqlx.Connect("pgx", system.ContraDBURL)
	if err != nil {
		if !failedOnce.Load() {
			system.Console.Error("failed to connect to PostgreSQL database server with error:")
			system.Console.Error(errors.Unwrap(err)) // Don't print username + password
			system.ContraDBOnline.Store(false)
			failedOnce.Store(true)
		}
	} else {
		pdb, err = stdlib.AcquireConn(xdb.DB)
		if err != nil {
			system.Console.Error("failed to acquire PGX connection with error:")
			system.Console.Error(err)
			system.ContraDBOnline.Store(false)
		} else {
			system.Console.Info("connected to ContraDB")
			system.ContraDBOnline.Store(true)
			failedOnce.Store(false)

			if !configLoaded.Load() {
				readConfig()
			}
		}
	}
}

func monitor() {
	for range time.Tick(time.Second * 10) {
		ping()
	}
}

type ErrContraDBOffline struct {
}

func (e *ErrContraDBOffline) Error() string {
	return "ContraDB is disconnected"
}

func checkOfflineError(err error) bool {
	if err == nil {
		return false
	}
	_, isOpErr := errors.Unwrap(err).(*net.OpError)
	return isOpErr || reflect.TypeOf(err).String() == "*pgconn.connLockError"
}

func errOfflineOrOriginal(err error) error {
	if checkOfflineError(err) {
		offline(err)
		return &ErrContraDBOffline{}
	} else {
		return err
	}
}

// Pings ContraDB to trigger online/offline code.
func ping() {
	//var err error

	if xdb == nil || pdb == nil {
		connect()

		return
	}

	//err = pdb.Ping(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//	offline(err)
	//} else {
	//	online()
	//}
}

func offline(err error) {
	if system.ContraDBOnline.Load() {
		if !checkOfflineError(err) {
			system.Console.Error("failed to ping ContraDB with unanticipated error:")
			system.Console.Error(err.Error())
		} else {
			fmt.Println(string(debug.Stack()))
			system.Console.Error("failed to ping ContraDB because it is offline")
		}
		system.ContraDBOnline.Store(false)
	}
}

func online() {
	if !system.ContraDBOnline.Load() {
		system.Console.Info("PostgreSQL health check succeeded")
		system.ContraDBOnline.Store(true)
	}
}
