\c contradb contracore_mgr

DROP VIEW IF EXISTS contracore.lease_details_by_mac;
CREATE OR REPLACE VIEW contracore.lease_details_by_mac AS
SELECT l.time, l.mac, l.ip, l.hostname, o.vendor
FROM contracore.lease l
         LEFT OUTER JOIN oui o ON o.mac = left(l.mac::TEXT, 8)
WHERE time > now() - INTERVAL '5 days'
  AND l.id IN (
    SELECT max(id) AS id
    FROM contracore.lease
    GROUP BY mac
)
  AND op != 'del'
ORDER BY id DESC;

--

DROP VIEW IF EXISTS contracore.lease_details_by_ip;
CREATE OR REPLACE VIEW contracore.lease_details_by_ip AS
SELECT l.time, l.ip, l.mac, l.hostname, l.vendor
FROM contracore.lease l
WHERE time > now() - INTERVAL '5 days'
  AND (id) IN (
    SELECT max(id)
    FROM contracore.lease
    GROUP BY ip
)
  AND op != 'del'
ORDER BY id DESC;

--

DROP VIEW IF EXISTS contracore.lease_details_by_ip_hostname;
CREATE OR REPLACE VIEW contracore.lease_details_by_ip_hostname AS
SELECT l.time, l.mac, l.ip, l.hostname, o.vendor
FROM contracore.lease l
         LEFT OUTER JOIN oui o ON o.mac = left(l.mac::TEXT, 8)
WHERE time > now() - INTERVAL '5 days'
  AND l.id IN (
    SELECT max(id) AS id
    FROM contracore.lease
    GROUP BY ip, hostname
)
  AND op != 'del'
  AND hostname IS NOT NULL
ORDER BY id DESC;

--

DROP VIEW IF EXISTS contracore.lease_vendor_counts;
CREATE OR REPLACE VIEW contracore.lease_vendor_counts AS
WITH maxes AS (
    SELECT max(id) AS id
    FROM contracore.lease
    WHERE time > now() - INTERVAL '7 days'
    GROUP BY mac
),
     counts AS (
         SELECT vendor, count(vendor) AS c
         FROM lease
                  INNER JOIN maxes m ON lease.id = m.id
         GROUP BY vendor
     ),
     total AS (
         SELECT sum(counts.c) AS total
         FROM counts
     )
SELECT counts.*,
       CASE WHEN counts.c = 0 THEN 0::FLOAT ELSE counts.c::FLOAT / sum(total.total)::FLOAT END AS ratio
FROM counts,
     total
GROUP BY counts.vendor, counts.c
ORDER BY ratio DESC;
