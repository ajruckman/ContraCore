package contralog

import (
    "database/sql"

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
    var mac *string
    if log.ClientMAC != nil {
        r := log.ClientMAC.String()
        mac = &r
    }

    _, err := stmt.Exec(log.Time, log.Client, log.Question, log.QuestionType, log.Action, log.Answers, mac, log.ClientHostname, log.ClientVendor, log.QueryID)
    return errOfflineOrOriginal(err)
}

func CommitBatch(tx *sql.Tx) (err error) {
    err = tx.Commit()
    return errOfflineOrOriginal(err)
}
