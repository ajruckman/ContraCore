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

ALTER
TABLE
log
UPDATE action = 'pass.notblacklisted'
WHERE action = 'pass';
ALTER
TABLE
log
UPDATE action = 'respond.block'
WHERE action = 'block';
ALTER
TABLE
log
UPDATE action = 'respond.ddnshostname'
WHERE action = 'ddns-hostname';
ALTER
TABLE
log
UPDATE action = 'respond.ddnsptr'
WHERE action = 'restrict';
ALTER
TABLE
log
UPDATE action = 'respond.domainneeded'
WHERE action = 'ddns-ptr';

SELECT DISTINCT action
FROM log;

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

DROP TABLE IF EXISTS log_count_per_hour;
CREATE VIEW log_count_per_hour AS
SELECT formatDateTime(toStartOfHour(time), '%F %H') AS hour, count(*) AS c, action
FROM log
WHERE action LIKE 'pass.%'
GROUP BY toStartOfHour(time), action
ORDER BY hour DESC, action DESC
LIMIT 168;

-- DROP TABLE IF EXISTS log_actions_per_hour;
-- CREATE VIEW log_actions_per_hour AS
--     -- WITH (
-- --     SELECT subtractDays(toStartOfHour(max(time)), 7)
-- --     FROM log
-- -- )
-- --     AS mintime
-- SELECT action, formatDateTime(toStartOfHour(time), '%F %H') AS hour, count(*) AS c
-- FROM log
-- -- WHERE time > mintime
-- GROUP BY toStartOfHour(time), action
-- ORDER BY action, hour;


-- SELECT slot, c
-- FROM (
--          SELECT toStartOfHour(time) AS slot,
--                 count(*)            AS c
--          FROM log
--          WHERE action = 'asdf'
--          GROUP BY slot
--          ) s1
--          ANY
--          RIGHT JOIN
--      (
--          SELECT arrayJoin(timeSlots(now() - toIntervalDay(7), toUInt32(7 * 24 * 60),)) AS slot
--          ) s2 USING (slot)
-- ;
--
-- SELECT metric, time
-- FROM (
--          SELECT toStartOfHour(now() - toIntervalDay(7)) AS time,
--                 toUInt16(0)                             AS metric
--          FROM numbers(7 * 24)
--          ) s1
-- ;


-- select * from actions;
-- SELECT s1.hour
-- FROM (

DROP TABLE IF EXISTS log_actions_per_hour;
CREATE VIEW log_actions_per_hour AS
SELECT formatDateTime(s1.hour, '%F %H') AS hour, s1.action, s2.c
FROM (
         SELECT toStartOfHour(now() - toIntervalDay(7)) + (number * 60 * 60) AS hour, action
                FROM (
                  SELECT DISTINCT action
                  FROM log
         ) AS actions,
         numbers(7 * 24)
         ORDER BY hour, action
    SETTINGS joined_subquery_requires_alias=0
) AS s1
LEFT OUTER JOIN
(
SELECT toStartOfHour(time) AS hour,
    action,
    count(*)               AS c
    FROM log
    GROUP BY toStartOfHour(time), action
) AS s2
ON s1.hour = s2.hour AND s1.action = s2.action;

--          ) s1

--          LEFT OUTER JOIN
--      (
--          SELECT action,
--                 toStartOfHour(time) AS hour,
--                 count(*)            AS c
--          FROM log
--          GROUP BY toStartOfHour(time), action
--          ) s2
--      ON s1.hour = s2.hour
;


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
