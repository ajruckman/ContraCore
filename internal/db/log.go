package db

import (
    "context"

    "github.com/ajruckman/ContraCore/internal/schema"
)

func Log(log schema.Log) error {
    _, err := PDB.Exec(context.Background(), `

INSERT
INTO log (client, question, question_type, action, answers, client_hostname, client_mac, client_vendor)
SELECT values.*, l.hostname, l.mac, l.vendor
FROM (
    SELECT $1::INET    AS client,
           $2          AS question,
           $3          AS question_type,
           $4          AS action,
           $5::TEXT[]  AS answers
) values
     LEFT OUTER JOIN lease_details l ON l.ip = values.client
GROUP BY values.client, values.question, values.question_type, values.action, values.answers, l.hostname, l.mac, l.vendor;

`, log.Client, log.Question, log.QuestionType, log.Action, log.Answers)

    return err
}
