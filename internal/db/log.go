package db

import (
    "context"
    "database/sql"

    "github.com/jackc/pgx/v4"

    "github.com/ajruckman/ContraCore/internal/schema"
)

func Log(tx pgx.Tx, log schema.Log) error {
    _, err := tx.Exec(context.Background(), `
        INSERT INTO log (time, client, question, question_type, action, answers, client_mac, client_hostname, client_vendor) 
        VALUES ($1, $2::INET, $3, $4, $5, $6::TEXT[], $7, $8, $9);
    `, log.Time, log.Client, log.Question, log.QuestionType, log.Action, log.Answers, log.ClientMAC, log.ClientHostname, log.ClientVendor)

    return err
}

func LogCBeginBatch() (tx *sql.Tx, stmt *sql.Stmt, err error) {
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

func LogC(stmt *sql.Stmt, log schema.Log) error {
    var mac *string
    if log.ClientMAC != nil {
        r := log.ClientMAC.String()
        mac = &r
    }

    _, err := stmt.Exec(log.Time, log.Client, log.Question, log.QuestionType, log.Action, log.Answers, mac, log.ClientHostname, log.ClientVendor, log.QueryID)
    return err
}
