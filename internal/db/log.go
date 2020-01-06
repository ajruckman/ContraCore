package db

import (
    "context"

    "github.com/ajruckman/ContraCore/internal/schema"
)

func Log(log schema.Log) error {
    _, err := PDB.Exec(context.Background(), `

INSERT INTO log (client, question, question_type, action, answers, client_mac, client_hostname, client_vendor) 
VALUES ($1::INET, $2, $3, $4, $5::TEXT[], $6, $7, $8);

`, log.Client, log.Question, log.QuestionType, log.Action, log.Answers, log.ClientMAC, log.ClientHostname, log.ClientVendor)

    return err
}

func LogC(log schema.Log) error {
    return nil

    tx, err := CDB.Begin()
    if err != nil {
        return err
    }

    stmt, err := tx.Prepare(`

INSERT INTO contralog.log(client, question, question_type, action, answers, client_mac, client_hostname, client_vendor)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)

`)
    if err != nil {
        return err
    }

    var mac *string
    if log.ClientMAC != nil {
        r := log.ClientMAC.String()
        mac = &r
    }

    _, err = stmt.Exec(log.Client, log.Question, log.QuestionType, log.Action, log.Answers, mac, log.ClientHostname, log.ClientVendor)
    if err != nil {
        return err
    }

    err = tx.Commit()
    return err
}
