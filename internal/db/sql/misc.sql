DROP TABLE "Log";
TRUNCATE TABLE "Log";
ALTER SEQUENCE "Log_ID_seq" RESTART;

SELECT count(*)
FROM "Log"
WHERE "Question" = 'nobody.invalid';

SELECT DISTINCT "Log"."QuestionType"
FROM "Log";

SELECT *
FROM "Log"
WHERE "QuestionType" = 'SOA';

SELECT DISTINCT unnest("Log"."Answers") AS answer
FROM "Log"
WHERE "QuestionType" = 'A'
  AND "answer"::INET IS NOT NULL;


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

CREATE OR REPLACE VIEW "Reverse" AS
WITH answers AS (
    SELECT DISTINCT unnest(l."Answers") AS "Answer", l."Question"
    FROM "Log" l
    ORDER BY "Question"
)
SELECT *
FROM answers a
WHERE is_inet(a."Answer");


SELECT * FROM "Reverse";