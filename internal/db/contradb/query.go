package contradb

import (
    "context"
    "time"

    "github.com/jackc/pgconn"
    "github.com/jackc/pgx/pgtype"
    "github.com/jackc/pgx/v4"
    "github.com/jmoiron/sqlx"

    contradbschema "github.com/ajruckman/ContraCore/internal/db/contradb/dbschema"
    contralogschema "github.com/ajruckman/ContraCore/internal/db/contralog/dbschema"
    "github.com/ajruckman/ContraCore/internal/system"
)

func GetLeaseDetails() (res []contradbschema.LeaseDetails, err error) {
    if system.PostgresOnline.Load() {
        var rows *sqlx.Rows

        rows, err = xdb.Queryx(`SELECT time, op, ip, mac, hostname, vendor FROM lease_details;`)
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
    if system.PostgresOnline.Load() {
        err = xdb.Select(&res, `SELECT * FROM oui;`)
    }
    return
}

func GetConfig() (res contradbschema.Config, err error) {
    if system.PostgresOnline.Load() {
        err = xdb.Get(&res, `SELECT * FROM config ORDER BY id DESC LIMIT 1`)
    }
    return
}

func GetRules() (res []contradbschema.Rule, err error) {
    if system.PostgresOnline.Load() {
        err = xdb.Select(&res, `SELECT id, pattern, class, COALESCE(domain, '') AS domain, COALESCE(tld, '') AS tld, COALESCE(sld, '') AS sld FROM rule;`)
    }
    return
}

func GetHourly() (res []contralogschema.LogCountPerHour, err error) {
    //err = CDB.Select(&res, `SELECT * FROM log_count_per_hour;`)
    return
}

func Exec(query string, args ...interface{}) (pgconn.CommandTag, error) {
    return pdb.Exec(context.Background(), query, args...)
}

func CopyFrom(tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
    return pdb.CopyFrom(context.Background(), tableName, columnNames, rowSrc)
}

func Select(dest interface{}, query string, args ...interface{}) error {
    return xdb.Select(dest, query, args...)
}
