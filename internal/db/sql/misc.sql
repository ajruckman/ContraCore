SELECT count(*)
FROM log;

DROP TABLE log CASCADE;
TRUNCATE TABLE log;
ALTER SEQUENCE log_id_seq RESTART;

----- Remove trailing periods

UPDATE log
SET answers =
    array(SELECT trim(TRAILING '.' FROM elem) FROM unnest(answers) elem);

-- select * from log where id = 475755;
---- {wildcard.weather.microsoft.com.edgekey.net.,e15275.g.akamaiedge.net.,104.98.58.140}
---- {wildcard.weather.microsoft.com.edgekey.net ,e15275.g.akamaiedge.net ,104.98.58.140}

-----

SELECT DISTINCT log.question_type
FROM log;

SELECT *
FROM log
WHERE question_type = 'TKEY';

SELECT *
FROM log
WHERE question_type = 'SOA'
  AND answers IS NOT NULL
ORDER BY id DESC;

-----

CREATE OR REPLACE VIEW question_counts AS
SELECT l.question, count(l.question) AS count
FROM log l
GROUP BY l.question
ORDER BY count DESC;

-----

CREATE OR REPLACE VIEW question_counts_by_client AS
SELECT l.question, l.client, count(l.question) AS count
FROM log l
GROUP BY l.question, l.client
ORDER BY count DESC;

-----

CREATE OR REPLACE VIEW client_counts AS
SELECT l.client, count(l.client) AS count
FROM log l
GROUP BY l.client
ORDER BY count DESC;

-----

-- https://www.depesz.com/2010/10/22/grouping-data-into-time-ranges/
CREATE FUNCTION ts_round(TIMESTAMPTZ, INT4) RETURNS TIMESTAMPTZ AS
$$
SELECT 'epoch'::TIMESTAMPTZ + '1 second'::INTERVAL * ($2 * (EXTRACT(EPOCH FROM $1)::INT4 / $2));
$$ LANGUAGE SQL;

CREATE OR REPLACE VIEW question_counts_per_hour AS
SELECT ts_round(time, 3600)
    AS hour,
    count(l.id)
FROM log l
GROUP BY 1
ORDER BY 1;

-----

DROP FUNCTION IF EXISTS is_inet(S TEXT);
CREATE OR REPLACE FUNCTION is_inet(s TEXT) RETURNS BOOLEAN AS
$$
BEGIN
    PERFORM s::INET;
    RETURN TRUE;
EXCEPTION
    WHEN OTHERS THEN RETURN FALSE;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE VIEW reverse AS
WITH answers AS (
    SELECT DISTINCT unnest(l.answers) AS answer, l.question
    FROM log l
    ORDER BY question
)
SELECT *
FROM answers a
WHERE is_inet(a.answer);

-----

SELECT l.question, count(l.question) AS count
FROM log l
WHERE l.question LIKE '%roku%'
GROUP BY l.question
ORDER BY count DESC;

-----

SELECT DISTINCT client
FROM log
WHERE question LIKE '%spotify%'
ORDER BY client::INET;