SELECT formatDateTime(arrayJoin(arrayMap(x -> now() - (x * 60 * 60), range(7 * 24))), '%F %H:%M') AS hour;

CREATE VIEW IF NOT EXISTS log_actions_per_sliding_hour AS
SELECT toStartOfHour(time) + (now() - toStartOfHour(now())) AS h,
       formatDateTime(h, '%F %H:%M')                        AS h_f,
       action,
       count(*)                                             AS c
FROM log
GROUP BY h, action
ORDER BY h DESC WITH fill step -3600;


drop table log_actions_per_sliding_hour;