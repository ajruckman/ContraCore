CREATE TABLE contralog.oui
(
    mac    FixedString(8),
    vendor String
)
    ENGINE = Log();

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
    client_vendor   Nullable(String)
)
    ENGINE = AggregatingMergeTree(event_date, (client, question, question_type), 8192);

CREATE VIEW log_top_blocked AS
SELECT client, client_hostname AS hostname, client_vendor AS vendor, question, count(question) AS c
FROM contralog.log
WHERE action = 'block'
GROUP BY client, hostname, vendor, question
HAVING count(question) > 3
ORDER BY c DESC;

/*

 cat log_for_clickhouse_1.csv | clickhouse-client --query "INSERT INTO contralog.log FORMAT CSV" --user contralogmgr --password contralogmgr

 */

INSERT INTO log(client, question, question_type, action, answers, client_hostname, client_mac, client_vendor)
SELECT '10.2.0.10',
       'google.com',
       'A',
       'pass',
       ['8.8.8.8', '8.8.4.4'],
       'ajr-desktop',
       '34:97:F6:87:52:71',
       vendor
FROM oui
WHERE startsWith(mac, '34:97:f6');

SELECT *
FROM contralog.oui
WHERE startsWith(mac, '34:97:f6');
