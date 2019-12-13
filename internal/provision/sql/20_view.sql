----- Lease details
CREATE OR REPLACE VIEW lease_details AS
SELECT lease.time, lease.op, lease.mac, lease.ip, lease.hostname, o.vendor
FROM lease
     LEFT OUTER JOIN oui o ON trunc(o.mac)::TEXT ILIKE (left(lease.mac, 9) || '%')
WHERE (id, ip) IN (
    SELECT max(id), ip
    FROM lease
    GROUP BY ip)
ORDER BY id;
