package main

import (
	"log"
	"log/syslog"
	"os"
	"strings"

	. "github.com/ajruckman/xlib"

	"github.com/ajruckman/ContraCore/internal/db/contradb"
	"github.com/ajruckman/ContraCore/internal/system"
)

func main() {
	logwriter, err := syslog.New(syslog.LOG_NOTICE, "loglease")
	Err(err)
	log.SetOutput(logwriter)
	log.Print(strings.Join(os.Args, " | "))

	system.ContraDBURL = "postgres://contra_usr:EvPvkro59Jb7RK3o@127.0.0.1/contradb"
	contradb.Setup()

	op, mac, ip, hostname := coalesce(1), coalesce(2), coalesce(3), coalesce(4)
	logNewEntry(op, mac, ip, hostname)
}

func logNewEntry(op, mac, ip, hostname string) {
	_, err := contradb.Exec(`

INSERT INTO lease (source, op, ip, mac, hostname, vendor)
SELECT values.*, o.vendor
FROM (
     SELECT $1             AS source,
            $2             AS op,
            $3::INET       AS ip,
            $4::MACADDR    AS mac,
            NULLIF($5, '') AS hostname
) values
LEFT OUTER JOIN oui o ON o.mac = left($4::TEXT, 8)
GROUP BY values.source, values.op, values.ip, values.mac, values.hostname, o.vendor;

`,
		"dnsmasq", op, ip, mac, hostname)

	Err(err)
}

func coalesce(index int) string {
	if len(os.Args) > index {
		return os.Args[index]
	} else {
		return ""
	}
}
