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
    query_id        INT
)
    ENGINE = MergeTree(event_date, (client, question, question_type), 8192);

-- CREATE VIEW log_top_blocked AS
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
