package db

import (
    "net"
    "strconv"
    "time"

    "github.com/jackc/pgx/pgtype"
    "github.com/jmoiron/sqlx"

    "github.com/ajruckman/ContraCore/internal/schema"
    "github.com/ajruckman/ContraCore/internal/schema/contradb"
    "github.com/ajruckman/ContraCore/internal/schema/contralog"
)

func GetLeaseDetails() (res []contradb.LeaseDetails, err error) {
    var rows *sqlx.Rows

    rows, err = XDB.Queryx(`SELECT time, op, ip, mac, hostname, vendor FROM lease_details;`)
    if err != nil {
        return
    }

    //err = XDB.Select(&res, `SELECT * FROM lease_details;`)
    //Err(err)

    defer rows.Close()

    for rows.Next() {
        var n = internalLeaseDetails{}
        err = rows.StructScan(&n)
        if err != nil {
            return
        }

        res = append(res, contradb.LeaseDetails{
            Time:     n.Time,
            Op:       n.Op,
            IP:       n.IP.IPNet.IP,
            MAC:      n.MAC.Addr,
            Hostname: n.Hostname,
            Vendor:   n.Vendor,
        })
    }

    return
}

type internalLeaseDetails struct {
    Time     time.Time      `db:"time"`
    Op       string         `db:"op"`
    IP       pgtype.Inet    `db:"ip"`
    MAC      pgtype.Macaddr `db:"mac"`
    Hostname *string        `db:"hostname"`
    Vendor   *string        `db:"vendor"`
}

func GetOUI() (res []contradb.OUI, err error) {
    err = XDB.Select(&res, `SELECT * FROM oui;`)
    return
}

func GetConfig() (res contradb.Config, err error) {
    err = XDB.Get(&res, `SELECT * FROM config ORDER BY id DESC LIMIT 1`)
    return
}

func GetRules() (res []contradb.Rule, err error) {
    err = XDB.Select(&res, `SELECT id, pattern, class, COALESCE(domain, '') AS domain, COALESCE(tld, '') AS tld, COALESCE(sld, '') AS sld FROM rule;`)
    return
}

func GetHourly() (res []contralog.LogCountPerHour, err error) {
    //err = CDB.Select(&res, `SELECT * FROM log_count_per_hour;`)
    return
}

type InternalLog struct {
    Time           time.Time `db:"time"`
    Client         string    `db:"client"`
    Question       string    `db:"question"`
    QuestionType   string    `db:"question_type"`
    Action         string    `db:"action"`
    Answers        []string  `db:"answers"`
    ClientMAC      *string   `db:"client_mac"`
    ClientHostname *string   `db:"client_hostname"`
    ClientVendor   *string   `db:"client_vendor"`
}

func GetLastNLogs(n int) (res []schema.Log, err error) {
    rows, err := CDB.Queryx(`SELECT time, client, question, question_type, action, answers, client_hostname, client_mac, client_vendor FROM log ORDER BY time DESC LIMIT ` + strconv.Itoa(n) + `;`)
    if err != nil {
        return
    }

    for rows.Next() {
        var n = InternalLog{}
        err = rows.StructScan(&n)
        if err != nil {
            return
        }

        l := schema.Log{
            Time:           n.Time,
            Client:         n.Client,
            Question:       n.Question,
            QuestionType:   n.QuestionType,
            Action:         n.Action,
            Answers:        n.Answers,
            ClientHostname: n.ClientHostname,
            ClientVendor:   n.ClientVendor,
        }

        if n.ClientMAC != nil {
            m, err := net.ParseMAC(*n.ClientMAC)
            if err != nil {
                return nil, err
            }
            l.ClientMAC = &m
        }

        res = append(res, l)
    }

    return
}
