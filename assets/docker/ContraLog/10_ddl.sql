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

DROP TABLE IF EXISTS contralog.log_top_blocked;
CREATE VIEW contralog.log_top_blocked AS
SELECT client, client_hostname AS hostname, client_vendor AS vendor, question, count(question) AS c
FROM contralog.log
WHERE action LIKE 'block.%'
GROUP BY client, hostname, vendor, question
ORDER BY c DESC;

DROP TABLE IF EXISTS contralog.log_top_blocked_per_day;
CREATE VIEW contralog.log_top_blocked_per_day AS
SELECT event_date, client, client_hostname AS hostname, client_vendor AS vendor, question, count(question) AS c
FROM contralog.log
WHERE action LIKE 'block.%'
GROUP BY event_date, client, hostname, vendor, question
HAVING c > 10
ORDER BY event_date, c DESC;

DROP TABLE IF EXISTS contralog.log_count_per_hour;
CREATE VIEW contralog.log_count_per_hour AS
SELECT formatDateTime(toStartOfHour(time), '%F %H') AS hour, count(*) AS c, action
FROM contralog.log
WHERE action LIKE 'pass.%'
GROUP BY toStartOfHour(time), action
ORDER BY hour DESC, action DESC
LIMIT 168;

DROP TABLE IF EXISTS contralog.log_actions_per_hour;
CREATE VIEW contralog.log_actions_per_hour AS
SELECT formatDateTime(s1.t, '%F %H') AS hour, s1.action, s2.c
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

DROP TABLE IF EXISTS contralog.log_action_counts;
CREATE VIEW contralog.log_action_counts AS
SELECT action, count(action)
FROM contralog.log
WHERE time > now() - toIntervalDay(7)
GROUP BY action;
