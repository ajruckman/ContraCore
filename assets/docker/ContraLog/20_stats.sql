-- CREATE VIEW IF NOT EXISTS contralog.stats_top_100_blocked AS
SELECT question,
       count(question) AS q,
       min(time)       AS first
FROM contralog.log
WHERE action LIKE 'block.%'
  AND time > now() - toIntervalDay(7)
GROUP BY question
ORDER BY q DESC
LIMIT 100;

SELECT DISTINCT action
FROM contralog.log
WHERE time > now() - toIntervalDay(1);

SELECT toStartOfHour(time) AS h, count(*) as c
FROM contralog.log
GROUP BY toStartOfHour(time)
ORDER BY toStartOfHour(time) DESC;

select toStartOfFifteenMinutes()