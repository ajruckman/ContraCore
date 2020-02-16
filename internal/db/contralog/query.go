package contralog

import (
    "database/sql"
    "time"

    "github.com/ajruckman/ContraCore/internal/db/contralog/dbschema"
    "github.com/ajruckman/ContraCore/internal/system"
)

func BeginBatch() (tx *sql.Tx, stmt *sql.Stmt, err error) {
    if !system.ContraLogOnline.Load() {
        return nil, nil, &ErrContraLogOffline{}
    }

    tx, err = cdb.Begin()
    if err != nil {
        return nil, nil, errOfflineOrOriginal(err)
    }

    stmt, err = tx.Prepare(`
        INSERT INTO contralog.log(time, client, question, question_type, action, answers, client_mac, client_hostname, client_vendor, query_id)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)

    return
}

func SaveLog(stmt *sql.Stmt, log dbschema.Log) error {
    //var mac *string
    //if log.ClientMAC != nil {
    //    r := log.ClientMAC.String()
    //    mac = &r
    //}

    _, err := stmt.Exec(log.Time, log.Client, log.Question, log.QuestionType, log.Action, log.Answers, log.ClientMAC, log.ClientHostname, log.ClientVendor, log.QueryID)
    return errOfflineOrOriginal(err)
}

func CommitBatch(tx *sql.Tx) (err error) {
    err = tx.Commit()
    return errOfflineOrOriginal(err)
}

func GetLastNLogs(limit int) (res []dbschema.Log, err error) {
    err = cdb.Select(&res, `SELECT time, client, question, question_type, action, answers, client_mac, client_hostname, client_vendor, query_id FROM contralog.log ORDER BY time DESC LIMIT ?`, limit)
    //if err != nil {
    //    return
    //}

    //for rows.Next() {
    //    var l = _log{}
    //    err = rows.Scan(&l.Time, &l.Client, &l.Question, &l.QuestionType, &l.Action, &l.Answers, &l.ClientMAC, &l.ClientHostname, &l.ClientVendor, &l.QueryID)
    //    if err != nil {
    //        return
    //    }
    //
    //    var clientMAC net.HardwareAddr
    //    if l.ClientMAC != nil {
    //        clientMAC, err = net.ParseMAC(*l.ClientMAC)
    //        if err != nil {
    //            return
    //        }
    //    }
    //
    //    res = append(res, dbschema.Log{
    //        Time:           l.Time,
    //        Client:         l.Client,
    //        Question:       l.Question,
    //        QuestionType:   l.QuestionType,
    //        Action:         l.Action,
    //        Answers:        l.Answers,
    //        ClientMAC:      &clientMAC,
    //        ClientHostname: l.ClientHostname,
    //        ClientVendor:   l.ClientVendor,
    //        QueryID:        l.QueryID,
    //    })
    //}

    return
}

type _log struct {
    ID             int       `db:"id"`
    Time           time.Time `db:"time"`
    Client         string    `db:"client"`
    Question       string    `db:"question"`
    QuestionType   string    `db:"question_type"`
    Action         string    `db:"action"`
    Answers        []string  `db:"answers"`
    ClientMAC      *string   `db:"client_mac"`
    ClientHostname *string   `db:"client_hostname"`
    ClientVendor   *string   `db:"client_vendor"`
    QueryID        uint16    `db:"query_id"`
}
