package contralog

import (
    "database/sql"

    "github.com/pkg/errors"

    "github.com/ajruckman/ContraCore/internal/db/contralog/dbschema"
)

func BeginBatch() (tx *sql.Tx, stmt *sql.Stmt, err error) {
    tx, err = CDB.Begin()
    if err != nil {
        return nil, nil, err
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
    return err
}

func CommitBatch(tx *sql.Tx) (err error) {
    err = tx.Commit()
    if err != nil {
        return errors.Wrap(err, "could not commit ContraLog transaction")
    }
    return
}
