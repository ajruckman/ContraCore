package serve

import (
    "fmt"
    "regexp"
    "strings"
    "sync"
    "time"

    . "github.com/ajruckman/xlib"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema"
)

var (
    dhcpHostnameToLeases     map[string][]schema.LeaseDetails
    dhcpHostnameToLeasesLock sync.Mutex
    dhcpIPToLease            map[string]schema.LeaseDetails
    dhcpIPToLeaseLock        sync.Mutex
)

func cacheDHCP() {
    leases, err := db.GetLeaseDetails()
    Err(err)

    dhcpHostnameToLeasesLock.Lock()
    dhcpIPToLeaseLock.Lock()

    dhcpHostnameToLeases = map[string][]schema.LeaseDetails{}
    dhcpIPToLease = map[string]schema.LeaseDetails{}

    for _, lease := range leases {
        if lease.Hostname == "" {
            continue
        }

        hostname := strings.ToLower(lease.Hostname)
        if _, exists := dhcpHostnameToLeases[hostname]; !exists {
            dhcpHostnameToLeases[hostname] = []schema.LeaseDetails{}
        }
        dhcpHostnameToLeases[hostname] = append(dhcpHostnameToLeases[hostname], lease)

        dhcpIPToLease[lease.IP.String()] = lease
    }

    dhcpHostnameToLeasesLock.Unlock()
    dhcpIPToLeaseLock.Unlock()
}

func dhcpRefreshWorker() {
    for range time.Tick(time.Duration(dhcpRefreshInterval) * time.Second) {
        clog.Info("Refreshing DHCP lease cache")
        began := time.Now()
        cacheDHCP()
        clog.Infof("DHCP lease cache refreshed in %v. %d distinct hostnames and %d distinct IPs found.", time.Since(began), len(dhcpHostnameToLeases), len(dhcpIPToLease))
    }
}

func getLeasesByHostname(hostname string) ([]schema.LeaseDetails, bool) {
    dhcpHostnameToLeasesLock.Lock()
    v, ok := dhcpHostnameToLeases[hostname]
    dhcpHostnameToLeasesLock.Unlock()
    return v, ok
}

func getLeaseByIP(ip string) (schema.LeaseDetails, bool) {
    dhcpIPToLeaseLock.Lock()
    v, ok := dhcpIPToLease[ip]
    dhcpIPToLeaseLock.Unlock()
    return v, ok
}

func respondByHostname(q *queryContext) (ret bool, rcode int, err error) {
    if v, ok := getLeasesByHostname(q._domain); ok {
        var m *dns.Msg

        for _, lease := range v {
            if lease.IP.To4() != nil {
                if q._qu.Qtype == dns.TypeA {
                    m = genResponse(q.r, q._qu.Qtype, lease.IP.To4().String())
                    err = q.Respond(m)
                    clog.Debug("lease IP is IPv4, question is A")
                    q.action = "ddns-hostname" // TODO: Special action for RcodeServerFailure?
                    return true, dns.RcodeSuccess, err
                } else {
                    clog.Debug("lease IP is IPv4, question is AAAA")
                    continue
                }

            } else if lease.IP.To16() != nil {
                if q._qu.Qtype == dns.TypeAAAA {
                    m = genResponse(q.r, q._qu.Qtype, lease.IP.To16().String())
                    err = q.Respond(m)
                    clog.Debug("lease IP is IPv6, question is AAAA")
                    q.action = "ddns-hostname" // TODO: Special action for RcodeServerFailure?
                    return true, dns.RcodeSuccess, err
                } else {
                    clog.Debug("lease IP is IPv6, question is A")
                    continue
                }
            }
        }

        //clog.Error("lease with hostname '", q._domain, "' exists but query type is not A or AAAA")
        //m = responseWithCode(q.r, dns.RcodeSuccess)
        //err = q.Respond(m)
        //return true, dns.RcodeSuccess, err
    }

    return
}

var (
    matchReverse = regexp.MustCompile(`^(?:(?:1?[0-9]{1,2}|2[0-4][0-9]|25[0-5])\.){4}in-addr\.arpa$`)
    getReverse   = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)\.(\d+)\.in-addr.arpa$`)
)

func respondByPTR(q *queryContext) (ret bool, rcode int, err error) {
    if q._qu.Qtype != dns.TypePTR {
        return
    }

    if strings.HasSuffix(q._domain, ".in-addr.arpa") && matchReverse.MatchString(q._domain) {
        bits := getReverse.FindStringSubmatch(q._domain)
        fmt.Println(bits)

        ip := bits[4] + "." + bits[3] + "." + bits[2] + "." + bits[1]

        if v, ok := getLeaseByIP(ip); ok {
            q.action = "ddns-ptr"
            m := genResponse(q.r, q._qu.Qtype, v.IP.String())
            err := q.Respond(m)
            return true, dns.RcodeSuccess, err
        }
    }

    return
}
