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

with (select splitByString('.', question) from log
    where time > now() - toIntervalDay(7))
 as arr
select arr;

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
         ORDER BY TIME, action
             SETTINGS joined_subquery_requires_alias=0
         ) AS s1
         LEFT OUTER JOIN
     (
         SELECT toStartOfHour(time) AS hour,
                action,
                count(*)            AS c
         FROM log
         GROUP BY toStartOfHour(time), action
         ) AS s2
     ON s1.hour = s2.hour AND s1.action = s2.action;

-- DROP TABLE IF EXISTS log_actions_per_hour;
-- CREATE VIEW log_actions_per_hour AS

SELECT;

SELECT *
FROM (
         SELECT now() - toIntervalDay(7) + (number * 60 * 60) AS h
         FROM numbers(7 * 24)
         ) AS s1
         ANY
         JOIN
     (
         SELECT * FROM log
         ) AS s2
     ON s2.time BETWEEN s1.h - toIntervalHour(1) AND s1.h;


SELECT now() - toIntervalDay(7) + (number * 60 * 60) AS h
FROM numbers(7 * 24) AS n
         ANY
         LEFT JOIN log
                   ON log.time = h AND log.time > h
    SETTINGS joined_subquery_requires_alias = 0
;
--      (
--          SELECT * FROM log
--          ) AS s2
--      ON s2.time BETWEEN s1.h - toIntervalHour(1) AND s1.h;


WITH toStartOfHour(now()) AS s,
    now() - s AS off

SELECT s - toIntervalDay(7) + (number * 60 * 60) + off AS h
FROM numbers(7 * 24);


WITH toStartOfHour(now()) AS s,
    now() - s AS off
SELECT toStartOfHour(log.time) + off AS h
FROM log
ORDER BY time DESC;


WITH toStartOfHour(now()) AS s,
    now() - s AS off

SELECT s - toIntervalDay(7) + (number * 60 * 60) + off AS h
FROM numbers(7 * 24)
         LEFT OUTER JOIN (
    SELECT toStartOfHour(log.time) + off AS h
    FROM log
    ) AS s2
                         ON h = s2.h;




SELECT now();
SELECT arrayJoin(arrayMap(x -> now() - (x * 60 * 60), range(7 * 24)));

SELECT arrayJoin(arrayMap(x -> now() - (x * 60 * 60), range(7 * 24))) as h; asof join log on log.time = h and log.time > h;














WITH toStartOfHour(now()) AS s,
    now() - s AS off
SELECT toStartOfHour(time) + off
FROM log
ORDER BY time


-- SELECT now() - toIntervalDay(7) + (number * 60 * 60) + off as h
-- FROM numbers(7 * 24)
;



SELECT s1.t AS hour, s1.action, s2.c
FROM (
         SELECT now() - toIntervalDay(7) + (number * 60 * 60) AS t, action
         FROM (
                  SELECT DISTINCT action
                  FROM log
                  ) AS actions,
         numbers(7 * 24)
         ORDER BY t, action
             SETTINGS joined_subquery_requires_alias=0
         ) AS s1
         LEFT OUTER JOIN
     (
         SELECT time     AS hour,
                action,
                count(*) AS c
         FROM log
         GROUP BY time, action
         ) AS s2
     ON s1.hour = s2.hour AND s1.action = s2.action;


DROP TABLE IF EXISTS log_action_counts;
CREATE VIEW log_action_counts AS
SELECT action, count(action)
FROM log
WHERE time > now() - toIntervalDay(7)
GROUP BY action;


------------------

SELECT *
FROM (
         SELECT toStartOfHour(now() - toIntervalDay(7)) + (number * 60 * 60) AS hour
         FROM numbers(7 * 24)
         ) AS s1;


SELECT *
FROM (SELECT DISTINCT action FROM log) AS s1
         LEFT OUTER JOIN
     (
         SELECT toStartOfHour(time), action, count(action)
         FROM log
         GROUP BY toStartOfHour(time), action
         ) AS s2
     ON s1.action = s2.action;


SELECT s1.hour, pass_notblacklisted.c AS pass_notblacklisted, block_domainneeded.c AS block_domainneeded
FROM (
         SELECT toStartOfHour(now() - toIntervalDay(7)) + (number * 60 * 60) AS hour
         FROM numbers(7 * 24)
         ) AS s1
         LEFT OUTER JOIN (
    SELECT toStartOfHour(time) AS hour, count(*) AS c
    FROM log
    WHERE action = 'pass.notblacklisted'
    GROUP BY toStartOfHour(time)
    ) AS pass_notblacklisted
                         ON s1.hour = pass_notblacklisted.hour
         LEFT OUTER JOIN (
    SELECT toStartOfHour(time) AS hour, count(*) AS c
    FROM log
    WHERE action = 'block.blacklisted'
    GROUP BY toStartOfHour(time)
    ) AS block_domainneeded
                         ON s1.hour = block_domainneeded.hour ) s1;



--          LEFT OUTER JOIN
--      (
--          SELECT action,
--                 toStartOfHour(time) AS hour,
--                 count(*)            AS c
--          FROM log
--          GROUP BY toStartOfHour(time), action
--          ) s2
--      ON s1.hour = s2.hour;


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
