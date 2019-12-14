package db

import (
    "time"

    "github.com/jackc/pgx/pgtype"
    "github.com/jmoiron/sqlx"

    "github.com/ajruckman/ContraCore/internal/schema"
)

func GetLeaseDetails() (res []schema.LeaseDetails, err error) {
    var rows *sqlx.Rows

    rows, err = XDB.Queryx(`SELECT time, op, mac, ip, coalesce(hostname, '') AS hostname, coalesce(vendor, '') AS vendor FROM lease_details;`)
    if err != nil {
        return
    }

    defer rows.Close()

    for rows.Next() {
        var n = internalLeaseDetails{}
        err = rows.StructScan(&n)
        if err != nil {
            return
        }

        res = append(res, schema.LeaseDetails{
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

func GetRules() (res []schema.Rule, err error) {
    err = XDB.Select(&res, `SELECT * FROM rule;`)
    return
}
