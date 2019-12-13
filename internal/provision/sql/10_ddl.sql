----- Log
CREATE TABLE IF NOT EXISTS log
(
    id              BIGSERIAL NOT NULL,
    time            TIMESTAMP NOT NULL DEFAULT now(),
    client          INET,
    question        TEXT,
    question_type   TEXT,
    answers         TEXT[],
    client_hostname TEXT,
    client_mac      TEXT,

    CONSTRAINT log_pk PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS "log_question_answers_idx" ON log (question, answers);

----- Rule
CREATE TABLE IF NOT EXISTS rule
(
    id      SERIAL NOT NULL,
    domain  TEXT   NOT NULL,
    pattern TEXT   NOT NULL,

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
    search_domains TEXT[]  NOT NULL DEFAULT ARRAY [''],
    domain_needed  BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT config_pk PRIMARY KEY (id)
);
