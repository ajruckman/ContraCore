SELECT client, client_hostname, client_vendor, count(*) AS c
FROM contralog.log
WHERE time > now() - toIntervalDay(7)
GROUP BY client, client_hostname, client_vendor
ORDER BY c DESC;

SELECT question,
       question_type,
       action,
       count(question) AS q,
       min(time)       AS first
FROM contralog.log
WHERE action LIKE 'respond.ddnsptr'
  AND time > toIntervalDay(2)
GROUP BY question, question_type, action
ORDER BY q DESC;


SELECT formatDateTime(toStartOfHour(time) + (now() - toStartOfHour(now())), '%F %H:%M') AS hour,
       action,
       count(action)                                                                    AS c
FROM log
GROUP BY hour, action
ORDER BY hour DESC, count(action) DESC;



SELECT formatDateTime(time, '%F %H:%M') AS hour
FROM log
ORDER BY time DESC;


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



SELECT toStartOfHour(now() - toIntervalDay(7)) + (number * 60 * 60) AS t
FROM numbers(7 * 24);

SELECT *
FROM (
         SELECT formatDateTime(arrayJoin(arrayMap(x -> now() - (x * 60 * 60), range(7 * 24))), '%F %H:%M') AS hour
         )
         CROSS JOIN
     (
         SELECT DISTINCT action FROM log
         )
    SETTINGS joined_subquery_requires_alias = 0;



SELECT formatDateTime(arrayJoin(arrayMap(x -> now() - (x * 60 * 60), range(7 * 24))), '%F %H:%M'),
       (SELECT DISTINCT action FROM log),
       0
UNION ALL
SELECT formatDateTime(toStartOfHour(time) + (now() - toStartOfHour(now())), '%F %H:%M') AS hour,
       action,
       count(action)                                                                    AS c
FROM log
WHERE time > now() - toIntervalDay(7)
GROUP BY hour, action
ORDER BY hour DESC, count(action) DESC;

-- SELECT formatDateTime(s1.t, '%F %H') AS hour, s1.action, s2.c
-- FROM (
--          SELECT toStartOfHour(now() - toIntervalDay(7)) + (number * 60 * 60) AS t, action
--          FROM (
--                   SELECT DISTINCT action
--                   FROM contralog.log
--                   ) AS actions,
--          numbers(7 * 24)
--          ORDER BY t, action
--              SETTINGS joined_subquery_requires_alias=0
--          ) AS s1
--          LEFT OUTER JOIN
--      (
--          SELECT toStartOfHour(time) AS t,
--                 action,
--                 count(*)            AS c
--          FROM contralog.log
--          GROUP BY toStartOfHour(time), action
--          ) AS s2
--      ON s1.t = s2.t AND s1.action = s2.action;
