\c contradb contracore_mgr

DROP VIEW IF EXISTS contracore.lease_details_by_mac;
CREATE OR REPLACE VIEW contracore.lease_details_by_mac AS
SELECT l.time, l.mac, l.ip, l.hostname, o.vendor
FROM contracore.lease l
         LEFT OUTER JOIN oui o ON o.mac = left(l.mac::TEXT, 8)
WHERE time > now() - INTERVAL '3 days'
  AND l.id IN (
    SELECT max(id) AS id
    FROM contracore.lease
    GROUP BY mac
)
  AND op != 'del'
ORDER BY id DESC;

DROP VIEW IF EXISTS contracore.lease_details_by_ip;
CREATE OR REPLACE VIEW contracore.lease_details_by_ip AS
SELECT l.time, l.ip, l.mac, l.hostname, l.vendor
FROM contracore.lease l
WHERE time > now() - INTERVAL '3 days'
  AND (id) IN (
    SELECT max(id)
    FROM contracore.lease
    GROUP BY ip
)
  AND op != 'del'
ORDER BY id DESC;

DROP VIEW IF EXISTS contracore.lease_details_by_ip_hostname;
CREATE OR REPLACE VIEW contracore.lease_details_by_ip_hostname AS
SELECT l.time, l.mac, l.ip, l.hostname, o.vendor
FROM contracore.lease l
         LEFT OUTER JOIN oui o ON o.mac = left(l.mac::TEXT, 8)
WHERE time > now() - INTERVAL '7 days'
  AND l.id IN (
    SELECT max(id) AS id
    FROM contracore.lease
    GROUP BY hostname, ip
)
  AND op != 'del'
  AND hostname IS NOT NULL
ORDER BY id DESC;

CREATE OR REPLACE VIEW contracore.lease_vendor_counts AS
SELECT vendor, count(vendor) AS c
FROM contracore.lease l
WHERE id IN (
    SELECT max(id) AS id
    FROM contracore.lease
    WHERE time > now() - INTERVAL '7 days'
    GROUP BY mac
)
GROUP BY vendor
ORDER BY c DESC;
