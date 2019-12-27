package main

import (
    "log"
    "log/syslog"
    "os"
    "strings"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/db"
)

func main() {
    logwriter, err := syslog.New(syslog.LOG_NOTICE, "loglease")
    Err(err)
    log.SetOutput(logwriter)
    log.Print(strings.Join(os.Args, " | "))

    op, mac, ip, hostname := coalesce(1), coalesce(2), coalesce(3), coalesce(4)
    logNewEntry(op, mac, ip, hostname)
}

func logNewEntry(op, mac, ip, hostname string) {
    //_, err := db.XDB.Exec(`INSERT INTO lease (source, op, mac, ip, hostname) VALUES ($1, $2, $3, $4, nullif($5, ''))`, "dnsmasq", op, mac, ip, hostname)
    //Err(err)

    _, err := db.XDB.Exec(`

INSERT INTO lease (source, op, mac, ip, hostname, vendor)
SELECT values.*, o.vendor
FROM (
     SELECT $1             AS source,
            $2             AS op,
            $3             AS mac,
            $4::INET       AS ip,
            NULLIF($5, '') AS hostname
) values
LEFT OUTER JOIN oui o ON left(o.mac::TEXT, 9) = left($3, 9)
GROUP BY values.source, values.op, values.mac, values.ip, values.hostname, o.vendor;

`,
        "dnsmasq", op, mac, ip, hostname)

    Err(err)
}

func coalesce(index int) string {
    if len(os.Args) > index {
        return os.Args[index]
    } else {
        return ""
    }
}
