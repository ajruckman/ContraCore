package provision

import `fmt`
import `context`
import `github.com/ajruckman/ContraCore/internal/db`

func init() {
    fmt.Println(`Provisioning database`)
    _, err := db.PDB.Exec(context.Background(), `
----- Log
CREATE TABLE IF NOT EXISTS log
(
    id              BIGSERIAL NOT NULL,
    time            TIMESTAMP NOT NULL DEFAULT now(),
    client          INET      NOT NULL,
    question        TEXT      NOT NULL,
    question_type   TEXT      NOT NULL,
    action          TEXT      NOT NULL,
    answers         TEXT[],
    client_hostname TEXT,
    client_mac      TEXT,

    CONSTRAINT log_pk PRIMARY KEY (id),
    CONSTRAINT log_action_chk CHECK (action IN ('answer', 'restrict', 'block', 'pass'))
);

CREATE INDEX IF NOT EXISTS "log_question_answers_idx" ON log (question, answers);

----- Rule
CREATE TABLE IF NOT EXISTS rule
(
    id      SERIAL NOT NULL,
    pattern TEXT   NOT NULL,
    domain  TEXT   NOT NULL,
    tld     TEXT   NOT NULL,
    sld     TEXT   NOT NULL,

    CONSTRAINT rules_pk PRIMARY KEY (id)
);

----- Lease
CREATE TABLE IF NOT EXISTS lease
(
    id       BIGSERIAL NOT NULL,
    time     TIMESTAMP NOT NULL DEFAULT now(),
    op       CHAR(3) CHECK (op IN ('add', 'old', 'del')),
    mac      TEXT,
    ip       INET,
    hostname TEXT,

    CONSTRAINT lease_pk PRIMARY KEY (id)
);

----- OUI
CREATE TABLE IF NOT EXISTS oui
(
    mac    TEXT,
    vendor TEXT
);

----- Config
CREATE TABLE IF NOT EXISTS config
(
    id             SERIAL  NOT NULL,
    sources        TEXT[]  NOT NULL DEFAULT ARRAY [] ::TEXT[],
    search_domains TEXT[]  NOT NULL DEFAULT ARRAY [] ::TEXT[],
    domain_needed  BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT config_pk PRIMARY KEY (id)
);

----- Lease details
CREATE OR REPLACE VIEW lease_details AS
SELECT lease.time, lease.op, lease.mac, lease.ip, lease.hostname, o.vendor
FROM lease
     LEFT OUTER JOIN oui o ON o.mac::TEXT ILIKE (left(lease.mac, 9) || '%')
WHERE (id, ip) IN (
    SELECT max(id), ip
    FROM lease
    GROUP BY ip)
ORDER BY id;

----- Recent log with details
CREATE OR REPLACE VIEW log_details_recent AS
SELECT l.id,
    l.time,
    l.client,
    l.question,
    l.question_type,
    l.answers,
    l.client_hostname,
    l.client_mac,
    o.vendor AS client_vendor
FROM log l
     LEFT OUTER JOIN oui o ON o.mac::TEXT LIKE (left(l.client_mac, 9) || '%')
ORDER BY l.id DESC
LIMIT 1000;
    `)
    if err != nil { panic(err) }
}