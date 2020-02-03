package contradb

import (
    "time"

    "github.com/jackc/pgx/pgtype"
    "github.com/jmoiron/sqlx"

    contradbschema "github.com/ajruckman/ContraCore/internal/db/contradb/dbschema"
    contralogschema "github.com/ajruckman/ContraCore/internal/db/contralog/dbschema"
    "github.com/ajruckman/ContraCore/internal/state"
)

func GetLeaseDetails() (res []contradbschema.LeaseDetails, err error) {
    if state.PostgresOnline.Load() {
        var rows *sqlx.Rows

        rows, err = XDB.Queryx(`SELECT time, op, ip, mac, hostname, vendor FROM lease_details;`)
        if err != nil {
            return
        }

        defer rows.Close()

        for rows.Next() {
            var n = _leaseDetails{}
            err = rows.StructScan(&n)
            if err != nil {
                return
            }

            res = append(res, contradbschema.LeaseDetails{
                Time:     n.Time,
                Op:       n.Op,
                IP:       n.IP.IPNet.IP,
                MAC:      n.MAC.Addr,
                Hostname: n.Hostname,
                Vendor:   n.Vendor,
            })
        }
    }

    return
}

type _leaseDetails struct {
    Time     time.Time      `db:"time"`
    Op       string         `db:"op"`
    IP       pgtype.Inet    `db:"ip"`
    MAC      pgtype.Macaddr `db:"mac"`
    Hostname *string        `db:"hostname"`
    Vendor   *string        `db:"vendor"`
}

func GetOUI() (res []contradbschema.OUI, err error) {
    if state.PostgresOnline.Load() {
        err = XDB.Select(&res, `SELECT * FROM oui;`)
    }
    return
}

func GetConfig() (res contradbschema.Config, err error) {
    if state.PostgresOnline.Load() {
        err = XDB.Get(&res, `SELECT * FROM config ORDER BY id DESC LIMIT 1`)
    }
    return
}

func GetRules() (res []contradbschema.Rule, err error) {
    if state.PostgresOnline.Load() {
        err = XDB.Select(&res, `SELECT id, pattern, class, COALESCE(domain, '') AS domain, COALESCE(tld, '') AS tld, COALESCE(sld, '') AS sld FROM rule;`)
    }
    return
}

func GetHourly() (res []contralogschema.LogCountPerHour, err error) {
    //err = CDB.Select(&res, `SELECT * FROM log_count_per_hour;`)
    return
}
