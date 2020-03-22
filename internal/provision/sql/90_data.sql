INSERT INTO config (sources, search_domains, domain_needed, spoofed_a, spoofed_aaaa, spoofed_cname, spoofed_default)
VALUES (('{https://hosts-file.net/grm.txt, https://hosts-file.net/exp.txt, https://hosts-file.net/emd.txt, https://hosts-file.net/psh.txt, https://hosts-file.net/ad_servers.txt, https://v.firebog.net/hosts/Easyprivacy.txt}'),
        ('{hb.ruckman.dev}'),
        TRUE,
        '0.0.0.0',
        '::',
        '',
        '-');

INSERT INTO whitelist
    (pattern, expires, ips, subnets, macs, vendors, hostnames)
VALUES ('.*\.google\.com',
        NULL,
        ARRAY ['10.2.0.10'::INET, '10.2.0.11'::INET],
        NULL,
        NULL,
        NULL,
        NULL);

INSERT INTO whitelist
    (pattern, expires, ips, subnets, macs, vendors, hostnames)
VALUES ('.*\.google\.com',
        current_timestamp + (30 * INTERVAL '1 day'),
        ARRAY ['10.2.0.10'::INET, '10.2.0.11'::INET],
        NULL,
        NULL,
        NULL,
        NULL);

INSERT INTO whitelist
    (pattern, expires, ips, subnets, macs, vendors, hostnames)
VALUES ('.*\.(google\.com',
        current_timestamp + (30 * INTERVAL '1 day'),
        ARRAY ['10.2.0.10'::INET, '10.2.0.11'::INET],
        NULL,
        NULL,
        NULL,
        NULL);

INSERT INTO whitelist
    (pattern, expires, ips, subnets, macs, vendors, hostnames)
VALUES ('.*\.google\.com',
        NULL,
        NULL,
        ARRAY ['10.2.0.0/24'::CIDR, '10.2.1.0/24'::CIDR, '10.0.0.0/8', '10.0.0.0/30', '2001:0db8::/64', '2001:0db8::1/128'],
        NULL,
        NULL,
        NULL);

INSERT INTO whitelist
    (pattern, expires, ips, subnets, macs, vendors, hostnames)
VALUES ('.*\.google\.com',
        NULL,
        NULL,
        NULL,
        ARRAY ['00:11:22:33:44:55'::MACADDR, '11:22:33:44:55:66'::MACADDR],
        NULL,
        NULL);

INSERT INTO whitelist
    (pattern, expires, ips, subnets, macs, vendors, hostnames)
VALUES ('.*\.google\.com',
        NULL,
        NULL,
        NULL,
        NULL,
        ARRAY ['Microsoft Corporation', 'EDUP INTERNATIONAL (HK) CO., LTD'],
        NULL);

INSERT INTO whitelist
    (pattern, expires, ips, subnets, macs, vendors, hostnames)
VALUES ('.*\.google\.com',
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        ARRAY ['ajr-desktop', 'ajr-laptop']);
