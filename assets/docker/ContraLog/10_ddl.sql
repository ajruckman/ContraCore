CREATE TABLE IF NOT EXISTS contralog.log
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

--

DROP TABLE IF EXISTS contralog.log_top_blocked;
CREATE VIEW contralog.log_top_blocked AS
SELECT client, client_hostname AS hostname, client_vendor AS vendor, question, count(question) AS c
FROM contralog.log
WHERE action LIKE 'block.%'
GROUP BY client, hostname, vendor, question
ORDER BY c DESC;

--

DROP TABLE IF EXISTS contralog.log_top_blocked_per_day;
CREATE VIEW contralog.log_top_blocked_per_day AS
SELECT event_date, client, client_hostname AS hostname, client_vendor AS vendor, question, count(question) AS c
FROM contralog.log
WHERE action LIKE 'block.%'
GROUP BY event_date, client, hostname, vendor, question
HAVING c > 10
ORDER BY event_date, c DESC;

--

DROP TABLE IF EXISTS contralog.log_count_per_hour;
CREATE VIEW contralog.log_count_per_hour AS
SELECT formatDateTime(toStartOfHour(time), '%F %H') AS hour, count(*) AS c, action
FROM contralog.log
WHERE action LIKE 'pass.%'
GROUP BY toStartOfHour(time), action
ORDER BY hour DESC, action DESC
LIMIT 168;

--

-- DROP TABLE IF EXISTS contralog.log_actions_per_hour;
-- CREATE VIEW contralog.log_actions_per_hour AS
-- SELECT formatDateTime(toStartOfHour(time) + (now() - toStartOfHour(now())), '%F %H:%M') AS hour,
--        action,
--        count(action)                                                                    AS c
-- FROM log
-- WHERE time > now() - toIntervalDay(7)
-- GROUP BY hour, action
-- ORDER BY hour DESC, count(action) DESC;

-- Accurate version
DROP TABLE IF EXISTS contralog.log_actions_per_hour;
CREATE VIEW contralog.log_actions_per_hour AS
SELECT formatDateTime(s1.t, '%F %H:00') AS hour, s1.action, s2.c
FROM (
         SELECT toStartOfHour(now() - toIntervalDay(7)) + (number * 60 * 60) AS t, action
         FROM (
                  SELECT DISTINCT action
                  FROM contralog.log
                  ) AS actions,
         numbers(7 * 24)
         ORDER BY t, action
             SETTINGS joined_subquery_requires_alias=0
         ) AS s1
         LEFT OUTER JOIN
     (
         SELECT toStartOfHour(time) AS t,
                action,
                count(*)            AS c
         FROM contralog.log
         GROUP BY toStartOfHour(time), action
         ) AS s2
     ON s1.t = s2.t AND s1.action = s2.action;


SELECT EventDate, count()
FROM table
GROUP BY EventDate
ORDER BY EventDate WITH FILL;

SELECT toStartOfHour(time) + (now() - toStartOfHour(now())) AS h,
       action,
       count(action)                                        AS c
FROM log
GROUP BY h, action
ORDER BY h DESC WITH FILL STEP -3600;

SELECT toStartOfHour(time) + (now() - toStartOfHour(now())) AS h,
       action,
       count(action)                                             AS c
FROM log
GROUP BY h, action
ORDER BY h DESC WITH fill step -3600;

SELECT DISTINCT action, s.h, s.action, s.c
FROM log
 INNER JOIN
     (SELECT toStartOfHour(time) + (now() - toStartOfHour(now())) AS h,
             action,
             count(action)                                             AS c
      FROM log
     WHERE time > now() - toIntervalDay(7)
      GROUP BY h, action
      ORDER BY h
          DESC WITH fill step -3600
         ) s
ON log.action = s.action
WHERE log.time > now() - toIntervalDay(7);

select distinct l1.action from log l1 where l1.time > now() - toIntervalDay(7)
inner join log l2 on l1.action = l2.action;



-- -- Detailed version
-- DROP TABLE IF EXISTS contralog.log_actions_per_hour;
-- CREATE VIEW contralog.log_actions_per_hour AS
-- SELECT hour, action, max(c)
-- FROM (
--       (
--           WITH now() as n
--           SELECT formatDateTime(
--                     toStartOfHour(time) + (n - toStartOfHour(n)),
--                     '%F %H:%M'
--                  ) AS hour,
--                  action,
--                  count(action) AS c
--           FROM log
--           WHERE time > n - toIntervalDay(7)
--           GROUP BY hour, action
-- --           ORDER BY hour DESC, count(action) DESC
--       )
--       UNION ALL
--       (
--           SELECT *
--           FROM (
--               SELECT formatDateTime(
--                   arrayJoin(arrayMap(x -> now() - (x * 60 * 60), range(7 * 24))),
--                   '%F %H:%M'
--               ) AS hour
--           )
--           CROSS JOIN
--           (
--               SELECT DISTINCT action, 0
--                   FROM log
--                   WHERE time > now() - toIntervalDay(7)
--               )
--               SETTINGS joined_subquery_requires_alias = 0
--       )
--     )
-- GROUP BY hour, action
-- ORDER BY hour DESC, max(c) DESC;

--

DROP TABLE IF EXISTS contralog.log_action_counts;
CREATE VIEW contralog.log_action_counts AS
SELECT action, count(action)
FROM contralog.log
WHERE time > now() - toIntervalDay(7)
GROUP BY action;
