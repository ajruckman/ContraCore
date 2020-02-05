SELECT *
FROM contradb.public.log_details_recent;

SELECT array_to_json(answers)
FROM log_details_recent;

SELECT array_to_string(answers, ',')
FROM log_details_recent;

SELECT time::DATE                                                     AS eventdate,
       date_trunc('second', time)::TEXT,
       client,
       question,
       question_type,
       action,
       '"' || replace(array_to_json(answers)::TEXT, '"', '''') || '"' AS answers,
       client_hostname,
       client_mac,
       client_vendor
FROM log;
