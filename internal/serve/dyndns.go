package serve

import (
    "regexp"
    "strings"
    "sync"
    "time"

    . "github.com/ajruckman/xlib"
    "github.com/miekg/dns"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema/contradb"
)

var (
    //dhcpHostnameToLeases     map[string][]schema.LeaseDetails
    dhcpHostnameToLeases     sync.Map
    dhcpHostnameToLeasesLock sync.Mutex
    //dhcpIPToLease            map[string]schema.LeaseDetails
    dhcpIPToLease     sync.Map
    dhcpIPToLeaseLock sync.Mutex
)

func cacheDHCP() (ipsSeen, hostnamesSeen int) {
    leases, err := db.GetLeaseDetails()
    Err(err)

    dhcpHostnameToLeasesLock.Lock()
    dhcpIPToLeaseLock.Lock()

    dhcpHostnameToLeases = sync.Map{}
    dhcpIPToLease = sync.Map{}

    dhcpHostnameToLeasesTemp := map[string][]contradb.LeaseDetails{}

    for _, lease := range leases {
        dhcpIPToLease.Store(lease.IP.String(), lease)
        ipsSeen++
        //dhcpIPToLease[lease.IP.String()] = lease

        if lease.Hostname == nil {
            continue
        }

        hostname := strings.ToLower(*lease.Hostname)
        if _, exists := dhcpHostnameToLeasesTemp[hostname]; !exists {
            hostnamesSeen++
            dhcpHostnameToLeasesTemp[hostname] = []contradb.LeaseDetails{}
        }
        dhcpHostnameToLeasesTemp[hostname] = append(dhcpHostnameToLeasesTemp[hostname], lease)
    }

    for hostname, leases := range dhcpHostnameToLeasesTemp {
        dhcpHostnameToLeases.Store(hostname, leases)
    }

    dhcpHostnameToLeasesLock.Unlock()
    dhcpIPToLeaseLock.Unlock()

    return
}

func dhcpRefreshWorker() {
    for range time.Tick(time.Duration(dhcpRefreshInterval) * time.Second) {
        clog.Info("Refreshing DHCP lease cache")
        began := time.Now()
        ipsSeen, hostnamesSeen := cacheDHCP()
        clog.Infof("DHCP lease cache refreshed in %v. %d distinct IPs and %d distinct hostnames found.", time.Since(began), ipsSeen, hostnamesSeen)
    }
}

func getLeasesByHostname(hostname string) ([]contradb.LeaseDetails, bool) {
    dhcpHostnameToLeasesLock.Lock()
    v, ok := dhcpHostnameToLeases.Load(hostname)
    dhcpHostnameToLeasesLock.Unlock()

    if ok {
        return v.([]contradb.LeaseDetails), ok
    } else {
        return []contradb.LeaseDetails{}, ok
    }
}

func getLeaseByIP(ip string) (contradb.LeaseDetails, bool) {
    dhcpIPToLeaseLock.Lock()
    v, ok := dhcpIPToLease.Load(ip)
    dhcpIPToLeaseLock.Unlock()

    if ok {
        return v.(contradb.LeaseDetails), ok
    } else {
        return contradb.LeaseDetails{}, ok
    }
}

func respondByHostname(q *queryContext) (ret bool, rcode int, err error) {
    if v, ok := getLeasesByHostname(q._domain); ok {
        var m *dns.Msg

        for _, lease := range v {
            if lease.IP.To4() != nil {
                if q._qu.Qtype == dns.TypeA {
                    clog.Debug("lease IP is IPv4, question is A")
                    q.action = "ddns-hostname" // TODO: Special action for RcodeServerFailure?
                    m = genResponse(q.r, q._qu.Qtype, lease.IP.To4().String())
                    err = q.Respond(m)
                    return true, dns.RcodeSuccess, err
                } else {
                    clog.Debug("lease IP is IPv4, question is AAAA")
                    continue
                }

            } else if lease.IP.To16() != nil {
                if q._qu.Qtype == dns.TypeAAAA {
                    clog.Debug("lease IP is IPv6, question is AAAA")
                    q.action = "ddns-hostname" // TODO: Special action for RcodeServerFailure?
                    m = genResponse(q.r, q._qu.Qtype, lease.IP.To16().String())
                    err = q.Respond(m)
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

        ip := bits[4] + "." + bits[3] + "." + bits[2] + "." + bits[1]

        if v, ok := getLeaseByIP(ip); ok {
            if v.Hostname == nil {
                return
            }

            q.action = "ddns-ptr"

            m := genResponse(q.r, q._qu.Qtype, *v.Hostname)
            err := q.Respond(m)
            return true, dns.RcodeSuccess, err
        }
    }

    return
}
