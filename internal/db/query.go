package db

import (
    "time"

    . "github.com/ajruckman/xlib"
    "github.com/jackc/pgx/pgtype"

    "github.com/ajruckman/ContraCore/internal/schema"
)

func GetLeaseDetails() (res []schema.LeaseDetails) {
    rows, err := XDB.Queryx(`SELECT id, time, op, mac, ip, coalesce(hostname, '') AS hostname, coalesce(vendor, '') AS vendor FROM lease_details;`)
    Err(err)

    defer rows.Close()

    for rows.Next() {
        var n = internalLeaseDetails{}
        err = rows.StructScan(&n)
        Err(err)

        res = append(res, schema.LeaseDetails{
            ID:       n.ID,
            Time:     n.Time,
            Op:       n.Op,
            MAC:      n.MAC,
            IP:       n.IP.IPNet.IP,
            Hostname: n.Hostname,
            Vendor:   n.Vendor,
        })
    }

    return
}

type internalLeaseDetails struct {
    ID       uint64      `db:"id"`
    Time     time.Time   `db:"time"`
    Op       string      `db:"op"`
    MAC      string      `db:"mac"`
    IP       pgtype.Inet `db:"ip"`
    Hostname string      `db:"hostname"`
    Vendor   string      `db:"vendor"`
}

func GetConfig() (res schema.Config, err error) {
    err = XDB.Get(&res, `SELECT * FROM config ORDER BY id DESC LIMIT 1`)

    return
}
