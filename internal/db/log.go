package db

import (
    "github.com/ajruckman/ContraCore/internal/schema"
)

func Log(log schema.Log) error {
    _, err := PDB.Exec(`
        INSERT INTO "Log" ("Client", "Question", "QuestionType", "Answers") VALUES ($1, $2, $3, $4)
    `, log.Client, log.Question, log.QuestionType, log.Answers)

    return err
}
