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
    _, err := db.XDB.Exec(`INSERT INTO lease (op, mac, ip, hostname) VALUES ($1, $2, $3, nullif($4, ''))`, op, mac, ip, hostname)
    Err(err)
}

func coalesce(index int) string {
    if len(os.Args) > index {
        return os.Args[index]
    } else {
        return ""
    }
}
