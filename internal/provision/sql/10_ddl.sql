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
    client_vendor   TEXT,

    CONSTRAINT log_pk PRIMARY KEY (id),
    CONSTRAINT log_action_chk CHECK (action IN ('ddns-hostname', 'ddns-ptr', 'restrict', 'block', 'pass'))
);

BEGIN TRANSACTION;
ALTER TABLE log
    DROP CONSTRAINT log_action_chk;
UPDATE log
SET action = 'ddns-hostname'
WHERE action = 'answer';
ALTER TABLE log
    ADD CONSTRAINT log_action_chk CHECK (action IN ('ddns-hostname', 'ddns-ptr', 'restrict', 'block', 'pass'));
END TRANSACTION;

----- Rule
CREATE TABLE IF NOT EXISTS rule
(
    id      SERIAL NOT NULL,
    pattern TEXT   NOT NULL,
    class   INT    NOT NULL,
    domain  TEXT,
    tld     TEXT,
    sld     TEXT,

    CONSTRAINT rule_pk PRIMARY KEY (id),
    CONSTRAINT rule_class_chk CHECK (0 <= class AND class <= 3),

    CONSTRAINT rule_nonnull_chk CHECK
        (
            (class = 0 AND domain IS NULL AND tld IS NULL AND sld IS NULL)
            OR
            (class = 1 AND domain IS NOT NULL AND tld IS NOT NULL AND sld IS NULL)
            OR
            (class = 2 AND domain IS NOT NULL AND tld IS NOT NULL AND sld IS NOT NULL)
        )
);

----- Lease
CREATE TABLE IF NOT EXISTS lease
(
    id       BIGSERIAL NOT NULL,
    time     TIMESTAMP NOT NULL DEFAULT now(),
    source   TEXT      NOT NULL,
    op       CHAR(3)   NOT NULL CHECK (op IN ('add', 'old', 'del')),
    mac      TEXT,
    ip       INET,
    hostname TEXT,
    vendor   TEXT,

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
SELECT lease.time, lease.op, lease.mac, lease.ip, lease.hostname, lease.vendor
FROM lease
WHERE (id, ip) IN (
    SELECT max(id), ip
    FROM lease
    GROUP BY ip)
ORDER BY id DESC;

CREATE OR REPLACE VIEW log_details_recent AS
SELECT l.id,
       l.time,
       l.client,
       l.question,
       l.question_type,
       l.action,
       l.answers,
       l.client_hostname,
       l.client_mac,
       l.client_vendor
FROM log l
ORDER BY l.id DESC
LIMIT 1000;

----- Log blocks by client, question, hostname, vendor, and count
CREATE OR REPLACE VIEW log_block_details AS
SELECT client, client_hostname AS hostname, client_vendor AS vendor, question, count(question) AS c
FROM log
WHERE action = 'block'
GROUP BY client, question, client_hostname, client_vendor
HAVING count(question) > 3
ORDER BY c DESC;