CREATE TABLE contralog.log
(
    event_date      Date     DEFAULT now(),
    time            DateTime DEFAULT now(),
    client          String,
    question        String,
    question_type   String,
    action          String,
    answers         Array(String),
    client_hostname Nullable(String),
    client_mac      Nullable(String),
    client_vendor   Nullable(String),
    query_id        UInt16
)
    ENGINE = MergeTree(event_date, (client, question, question_type), 8192);

ALTER TABLE log UPDATE action = 'pass.notblacklisted'  WHERE action = 'pass';
ALTER TABLE log UPDATE action = 'respond.block'        WHERE action = 'block';
ALTER TABLE log UPDATE action = 'respond.ddnshostname' WHERE action = 'ddns-hostname';
ALTER TABLE log UPDATE action = 'respond.ddnsptr'      WHERE action = 'restrict';
ALTER TABLE log UPDATE action = 'respond.domainneeded' WHERE action = 'ddns-ptr';

SELECT DISTINCT action FROM log;

DROP TABLE IF EXISTS log_top_blocked;
CREATE VIEW log_top_blocked AS
SELECT client, client_hostname AS hostname, client_vendor AS vendor, question, count(question) AS c
FROM contralog.log
WHERE action = 'respond.block'
GROUP BY client, hostname, vendor, question
ORDER BY c DESC;

DROP TABLE IF EXISTS log_top_blocked_per_day;
CREATE VIEW log_top_blocked_per_day AS
SELECT event_date, client, client_hostname AS hostname, client_vendor AS vendor, question, count(question) AS c
FROM contralog.log
WHERE action = 'respond.block'
GROUP BY event_date, client, hostname, vendor, question
HAVING c > 10
ORDER BY event_date, c DESC;

CREATE VIEW log_count_per_hour AS
SELECT *
FROM (
      SELECT formatDateTime(toStartOfHour(time), '%F %H:%M') AS hour, count(*) AS c
      FROM log
      GROUP BY toStartOfHour(time)
      ORDER BY hour DESC
      LIMIT 168
         )
ORDER BY hour;

SELECT concat(database, '.', table)                         AS table,
       formatReadableSize(sum(bytes))                       AS size,
       sum(rows)                                            AS rows,
       max(modification_time)                               AS latest_modification,
       sum(bytes)                                           AS bytes_size,
       any(engine)                                          AS engine,
       formatReadableSize(sum(primary_key_bytes_in_memory)) AS primary_keys_size
FROM system.parts
WHERE active
GROUP BY database, table
ORDER BY bytes_size DESC;
