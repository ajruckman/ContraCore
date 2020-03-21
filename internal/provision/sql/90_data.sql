INSERT INTO config (sources, search_domains, domain_needed, spoofed_a, spoofed_aaaa, spoofed_cname, spoofed_default)
VALUES (('{https://hosts-file.net/grm.txt, https://hosts-file.net/exp.txt, https://hosts-file.net/emd.txt, https://hosts-file.net/psh.txt, https://hosts-file.net/ad_servers.txt, https://v.firebog.net/hosts/Easyprivacy.txt}'),
        ('{hb.ruckman.dev}'),
        TRUE,
        '0.0.0.0',
        '::',
        '',
        '-');
