\c contradb contracore_mgr

CREATE OR REPLACE VIEW distinct_lease_clients AS
SELECT l.time, l.mac, l.hostname, o.vendor
FROM lease l
         LEFT OUTER JOIN oui o ON o.mac = left(l.mac::TEXT, 8)
WHERE l.id IN (
    SELECT max(id) AS id
    FROM lease l
    GROUP BY mac
)
ORDER BY time;

CREATE OR REPLACE VIEW lease_vendor_counts AS
SELECT vendor, count(vendor) AS c
FROM lease l
WHERE id IN (
    SELECT max(id) AS id
    FROM lease
    WHERE time > now() - INTERVAL '7 days'
    GROUP BY mac
)
GROUP BY vendor
ORDER BY c DESC;