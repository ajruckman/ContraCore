CREATE TABLE IF NOT EXISTS log
(
    id            BIGSERIAL NOT NULL,
    time          TIMESTAMP NOT NULL DEFAULT now(),
    client        INET,
    question      TEXT,
    question_type TEXT,
    answers       TEXT[],

    CONSTRAINT log_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS rule
(
    id      SERIAL NOT NULL,
--     domain  TEXT   NOT NULL,
--     tld     TEXT   NOT NULL,
--     sld     TEXT,
    pattern TEXT   NOT NULL,

    CONSTRAINT rules_pk PRIMARY KEY (id)
);
