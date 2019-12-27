package main

import (
    "fmt"
    "time"

    "github.com/ajruckman/ContraCore/internal/rulegen"
)

var (
    urls = []string{
        "http://localhost/contradomain/spark",
        "http://localhost/contradomain/bluGo",
        "http://localhost/contradomain/blu",
        "http://localhost/contradomain/basic",
        //"http://localhost/contradomain/ultimate",
        //"http://localhost/contradomain/unified",

        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/spark/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/bluGo/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/blu/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/basic/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/ultimate/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/unified/formats/domains.txt",
        //"https://someonewhocares.org/hosts/hosts",
        //"https://gist.githubusercontent.com/angristan/20a398983c5b1daa9c13a1cbadb78fd6/raw/58d54b172b664ee5a0b53bb2e25c391433f2cc7a/hosts",
        //"https://www.encrypt-the-planet.com/downloads/hosts",

        //"http://localhost/contradomain/unified",
        //"http://localhost/contradomain/someonewhocares",
        //"http://localhost/contradomain/win",
    }

    // https://v.firebog.net/hosts/lists.php
    // Around ~255,000 distinct domains
    ticked = []string{
        "https://raw.githubusercontent.com/EnergizedProtection/block/master/spark/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/bluGo/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/blu/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/basic/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/ultimate/formats/domains.txt",
        //"https://raw.githubusercontent.com/EnergizedProtection/block/master/unified/formats/domains.txt",
        //
        //"https://hosts-file.net/grm.txt",
        //"https://hosts-file.net/exp.txt",
        //"https://hosts-file.net/emd.txt",
        //"https://hosts-file.net/psh.txt",
        //"https://hosts-file.net/ad_servers.txt",
        //"https://reddestdream.github.io/Projects/MinimalHosts/etc/MinimalHostsBlocker/minimalhosts",
        //"https://raw.githubusercontent.com/StevenBlack/hosts/master/data/KADhosts/hosts",
        //"https://raw.githubusercontent.com/StevenBlack/hosts/master/data/add.Spam/hosts",
        //"https://v.firebog.net/hosts/static/w3kbl.txt",
        //"https://adaway.org/hosts.txt",
        //"https://v.firebog.net/hosts/AdguardDNS.txt",
        //"https://raw.githubusercontent.com/anudeepND/blacklist/master/adservers.txt",
        //"https://s3.amazonaws.com/lists.disconnect.me/simple_ad.txt",
        //"https://v.firebog.net/hosts/Easylist.txt",
        //"https://pgl.yoyo.org/adservers/serverlist.php?hostformat=hosts&showintro=0&mimetype=plaintext",
        //"https://raw.githubusercontent.com/StevenBlack/hosts/master/data/UncheckyAds/hosts",
        //"https://www.squidblacklist.org/downloads/dg-ads.acl",
        //"https://raw.githubusercontent.com/bigdargon/hostsVN/master/hosts",
        //"https://v.firebog.net/hosts/Easyprivacy.txt",
        //"https://v.firebog.net/hosts/Prigent-Ads.txt",
        //"https://gitlab.com/quidsup/notrack-blocklists/raw/master/notrack-blocklist.txt",
        //"https://raw.githubusercontent.com/StevenBlack/hosts/master/data/add.2o7Net/hosts",
        //"https://raw.githubusercontent.com/crazy-max/WindowsSpyBlocker/master/data/hosts/spy.txt",
        //"https://raw.githubusercontent.com/Kees1958/WS3_annual_most_used_survey_blocklist/master/w3tech_hostfile.txt",
        //"https://www.github.developerdan.com/hosts/lists/ads-and-tracking-extended.txt",
        //"https://hostfiles.frogeye.fr/firstparty-trackers-hosts.txt",
        //"https://s3.amazonaws.com/lists.disconnect.me/simple_malvertising.txt",
        //"https://mirror1.malwaredomains.com/files/justdomains",
        //"https://mirror.cedia.org.ec/malwaredomains/immortal_domains.txt",
        //"https://www.malwaredomainlist.com/hostslist/hosts.txt",
        //"https://bitbucket.org/ethanr/dns-blacklists/raw/8575c9f96e5b4a1308f2f12394abd86d0927a4a0/bad_lists/Mandiant_APT1_Report_Appendix_D.txt",
        //"https://v.firebog.net/hosts/Prigent-Malware.txt",
        //"https://v.firebog.net/hosts/Prigent-Phishing.txt",
        //"https://phishing.army/download/phishing_army_blocklist_extended.txt",
        //"https://gitlab.com/quidsup/notrack-blocklists/raw/master/notrack-malware.txt",
        //"https://v.firebog.net/hosts/Shalla-mal.txt",
        //"https://raw.githubusercontent.com/StevenBlack/hosts/master/data/add.Risk/hosts",
        //"https://www.squidblacklist.org/downloads/dg-malicious.acl",
        //"https://gitlab.com/curben/urlhaus-filter/raw/master/urlhaus-filter-hosts.txt",
        //"https://raw.githubusercontent.com/DandelionSprout/adfilt/master/Alternate%20versions%20Anti-Malware%20List/AntiMalwareHosts.txt",
        //"https://zerodot1.gitlab.io/CoinBlockerLists/hosts_browser",
    }
)

func bench() {
    checks := 10

    var totMilliseconds int64

    for i := 0; i < checks+1; i++ {
        begin := time.Now()
        res, total := rulegen.GenFromURLs(urls)
        end := time.Now()

        if i != 0 {
            totMilliseconds += end.Sub(begin).Milliseconds()

            kept := len(res)
            ratio := float64(kept) / float64(total)

            fmt.Println(ratio, kept, total)
        }
    }

    fmt.Println(float64(totMilliseconds)/float64(checks), urls)
}

func load() {
    // Generate rules
    begin := time.Now()
    rules, distinct := rulegen.GenFromURLs(ticked)
    end := time.Now()

    kept := len(rules)
    ratio := float64(kept) / float64(distinct)

    fmt.Printf("%d rules generated from %d distinct domains in %v; ratio = %.3f\n", kept, distinct, end.Sub(begin), ratio)

    // Save to DB
    begin = time.Now()
    rulegen.SaveRules(rules)
    end = time.Now()

    fmt.Println("Rules saved in", end.Sub(begin))
}

func main() {
    load()
}
