----- Log
CREATE TABLE IF NOT EXISTS contra.log
(
    id              BIGSERIAL NOT NULL,
    time            TIMESTAMP NOT NULL DEFAULT now(),
    client          INET      NOT NULL,
    question        TEXT      NOT NULL,
    question_type   TEXT      NOT NULL,
    action          TEXT      NOT NULL,
    answers         TEXT[],
    client_mac      MACADDR,
    client_hostname TEXT,
    client_vendor   TEXT,

    CONSTRAINT log_pk PRIMARY KEY (id),
    CONSTRAINT log_action_chk CHECK (action IN ('ddns-hostname', 'ddns-ptr', 'restrict', 'block', 'pass'))
);

----- Whitelist
CREATE TABLE IF NOT EXISTS contra.whitelist
(
    id        SERIAL NOT NULL,
    pattern   TEXT   NOT NULL,
    expires   TIMESTAMP,
    ips       INET[],
    subnets   CIDR[],
    hostnames TEXT[],
    macs      MACADDR[],
    vendors   TEXT[],

    CONSTRAINT whitelist_pk PRIMARY KEY (id),

    CONSTRAINT whitelist_nonnull_chk CHECK
        (
            (ips IS NOT NULL)
            OR
            (subnets IS NOT NULL)
            OR
            (hostnames IS NOT NULL)
            OR
            (macs IS NOT NULL)
            OR
            (vendors IS NOT NULL)
        )
);

----- Blacklist
CREATE TABLE IF NOT EXISTS contra.blacklist
(
    id      SERIAL NOT NULL,
    pattern TEXT   NOT NULL,
    expires TIMESTAMP,
    class   INT    NOT NULL,
    domain  TEXT,
    tld     TEXT,
    sld     TEXT,

    CONSTRAINT blacklist_pk PRIMARY KEY (id),
    CONSTRAINT blacklist_class_chk CHECK (0 <= class AND class <= 3),

    CONSTRAINT blacklist_nonnull_chk CHECK
        (
            (class = 0 AND domain IS NULL AND tld IS NULL AND sld IS NULL)
            OR
            (class = 1 AND domain IS NOT NULL AND tld IS NOT NULL AND sld IS NULL)
            OR
            (class = 2 AND domain IS NOT NULL AND tld IS NOT NULL AND sld IS NOT NULL)
        )
);

----- Lease
CREATE TABLE IF NOT EXISTS contra.lease
(
    id       BIGSERIAL NOT NULL,
    time     TIMESTAMP NOT NULL DEFAULT now(),
    source   TEXT      NOT NULL,
    op       CHAR(3)   NOT NULL CHECK (op IN ('add', 'old', 'del')),
    mac      MACADDR   NOT NULL,
    ip       INET      NOT NULL,
    hostname TEXT,
    vendor   TEXT,

    CONSTRAINT lease_pk PRIMARY KEY (id)
);

----- Reservation
CREATE TABLE IF NOT EXISTS contra.reservation
(
    id      SERIAL    NOT NULL,
    time    TIMESTAMP NOT NULL DEFAULT now(),
    active  BOOLEAN   NOT NULL DEFAULT TRUE,
    mac     MACADDR   NOT NULL,
    ip      INET      NOT NULL,
    creator TEXT,
    comment TEXT,

    CONSTRAINT reservation_pk PRIMARY KEY (id)
);

----- OUI
CREATE TABLE IF NOT EXISTS contra.oui
(
    mac    CHAR(8),
    vendor TEXT
);

----- Config
CREATE TABLE IF NOT EXISTS contra.config
(
    id              SERIAL  NOT NULL,
    sources         TEXT[]  NOT NULL DEFAULT ARRAY [] ::TEXT[],
    search_domains  TEXT[]  NOT NULL DEFAULT ARRAY [] ::TEXT[],
    domain_needed   BOOLEAN NOT NULL DEFAULT TRUE,
    spoofed_a       TEXT    NOT NULL DEFAULT '0.0.0.0',
    spoofed_aaaa    TEXT    NOT NULL DEFAULT '::',
    spoofed_cname   TEXT    NOT NULL DEFAULT '',
    spoofed_default TEXT    NOT NULL DEFAULT '-',

    CONSTRAINT config_pk PRIMARY KEY (id)
);

----- Lease details
CREATE OR REPLACE VIEW lease_details AS
SELECT lease.time, lease.op, lease.mac, lease.ip, lease.hostname, lease.vendor
FROM lease
WHERE (id) IN (
    SELECT max(id)
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
       l.client_mac,
       l.client_hostname,
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

CREATE OR REPLACE VIEW oui_vendors AS
SELECT DISTINCT vendor
FROM oui
ORDER BY vendor;
