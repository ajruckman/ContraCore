CREATE TABLE contralog.log
(
    event_date      Date     DEFAULT now(),
--     uuid            UUID     DEFAULT generateUUIDv4(),
    time            DateTime DEFAULT now(),
    client          String,
    question        String,
    question_type   String,
    action          String,
    answers         Array(String),
    client_hostname Nullable(String),
    client_mac      Nullable(String),
    client_vendor   Nullable(String),
    query_id        INT
)
    ENGINE = MergeTree(event_date, (client, question, question_type), 8192);

DROP TABLE contralog.log_buffer;
-- num_layers, min_time, max_time, min_rows, max_rows, min_bytes, max_bytes
CREATE TABLE contralog.log_buffer AS contralog.log ENGINE = Buffer(contralog, log, 1, 10, 120, 10, 300, 0, 16000);

-- CREATE TABLE contralog.log2
-- (
--     event_date      Date     DEFAULT now(),
--     time            DateTime DEFAULT now(),
--     client          String,
--     question        String,
--     question_type   String,
--     action          String,
--     answers         Array(String),
--     client_hostname Nullable(String),
--     client_mac      Nullable(String),
--     client_vendor   Nullable(String)
-- )
--     ENGINE = MergeTree(event_date, (client, question, question_type), 8192);

CREATE VIEW log_top_blocked AS
SELECT client, client_hostname AS hostname, client_vendor AS vendor, question, count(question) AS c
FROM contralog.log
WHERE action = 'block'
GROUP BY client, hostname, vendor, question
ORDER BY c DESC;

-- CREATE VIEW log_top_blocked_per_day AS
SELECT event_date, client, client_hostname AS hostname, client_vendor AS vendor, question, count(question) AS c
FROM contralog.log
WHERE action = 'block'
GROUP BY event_date, client, hostname, vendor, question
HAVING c > 10
ORDER BY event_date, c DESC;

DROP TABLE log_count_per_hour;
-- CREATE VIEW log_count_per_hour AS
SELECT *
FROM (
      SELECT formatDateTime(toStartOfHour(time), '%F %H:%M') AS hour, count(*) AS c
      FROM log
      GROUP BY toStartOfHour(time)
      ORDER BY hour DESC
      LIMIT 168
         )
ORDER BY hour;

SELECT count(*)
FROM log;

select * from log_buffer where question = 'xmpp010.hpeprint.com' and query_id != 0

/*

 cat log_for_clickhouse_1.csv | clickhouse-client --query "INSERT INTO contralog.log FORMAT CSV" --user contralogmgr --password contralogmgr

 */

