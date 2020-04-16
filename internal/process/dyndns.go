package process

import (
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	. "github.com/ajruckman/xlib"
	"github.com/miekg/dns"
	"go.uber.org/atomic"

	"github.com/ajruckman/ContraCore/internal/db/contradb"
	"github.com/ajruckman/ContraCore/internal/db/contradb/dbschema"
	"github.com/ajruckman/ContraCore/internal/system"
)

var (
	dhcpHostnameToLeases     sync.Map
	dhcpHostnameToLeasesLock sync.Mutex
	dhcpIPToLeases           sync.Map
	dhcpIPToLeaseLock        sync.Mutex

	dhcpRefreshInterval = time.Second * 15
	dhcpRefreshFailed   atomic.Bool
)

func cacheDHCP() (ipsSeen, hostnamesSeen int, err error) {
	leases, err := contradb.GetLeaseDetails()
	if err != nil {
		return
	}

	dhcpHostnameToLeasesLock.Lock()
	dhcpIPToLeaseLock.Lock()

	dhcpHostnameToLeases = sync.Map{}
	dhcpIPToLeases = sync.Map{}

	var (
		dhcpIPToLeasesTemp       = map[string][]dbschema.LeaseDetailsByIP{}
		dhcpHostnameToLeasesTemp = map[string][]dbschema.LeaseDetailsByIP{}
	)

	for _, lease := range leases {
		//dhcpIPToLeases.Store(lease.IP.String(), lease)

		if _, exists := dhcpIPToLeasesTemp[lease.IP.String()]; !exists {
			ipsSeen++
			dhcpIPToLeasesTemp[lease.IP.String()] = []dbschema.LeaseDetailsByIP{}
		}
		dhcpIPToLeasesTemp[lease.IP.String()] = append(dhcpIPToLeasesTemp[lease.IP.String()], lease)

		if lease.Hostname == nil {
			continue
		}

		hostname := strings.ToLower(*lease.Hostname)
		if _, exists := dhcpHostnameToLeasesTemp[hostname]; !exists {
			hostnamesSeen++
			dhcpHostnameToLeasesTemp[hostname] = []dbschema.LeaseDetailsByIP{}
		}
		dhcpHostnameToLeasesTemp[hostname] = append(dhcpHostnameToLeasesTemp[hostname], lease)
	}

	for ip, leases := range dhcpIPToLeasesTemp {
		dhcpIPToLeases.Store(ip, leases)
	}

	for hostname, leases := range dhcpHostnameToLeasesTemp {
		dhcpHostnameToLeases.Store(hostname, leases)
	}

	dhcpHostnameToLeasesLock.Unlock()
	dhcpIPToLeaseLock.Unlock()

	return
}

func dhcpRefreshWorker() {
	for range time.Tick(dhcpRefreshInterval) {
		began := time.Now()
		ipsSeen, hostnamesSeen, err := cacheDHCP()
		if _, ok := err.(*contradb.ErrContraDBOffline); ok {
			if !dhcpRefreshFailed.Load() {
				system.Console.Warning("failed to refresh lease cache because ContraDB is not connected")
				dhcpRefreshFailed.Store(true)
			}
		} else if err != nil {
			Err(err)
		} else {
			_ = began
			_ = ipsSeen
			_ = hostnamesSeen
			//system.Console.Infof("DHCP lease cache refreshed in %v. %d distinct IPs and %d distinct hostnames found.", time.Since(began), ipsSeen, hostnamesSeen)
			dhcpRefreshFailed.Store(false)
		}
	}
}

func getLeasesByHostname(hostname string) ([]dbschema.LeaseDetailsByIP, bool) {
	dhcpHostnameToLeasesLock.Lock()
	v, ok := dhcpHostnameToLeases.Load(hostname)
	dhcpHostnameToLeasesLock.Unlock()

	if ok {
		return v.([]dbschema.LeaseDetailsByIP), ok
	} else {
		return []dbschema.LeaseDetailsByIP{}, ok
	}
}

func getLeasesByIP(ip net.IP) ([]dbschema.LeaseDetailsByIP, bool) {
	dhcpIPToLeaseLock.Lock()
	v, ok := dhcpIPToLeases.Load(ip.String())
	dhcpIPToLeaseLock.Unlock()

	if ok {
		return v.([]dbschema.LeaseDetailsByIP), ok
	} else {
		return []dbschema.LeaseDetailsByIP{}, ok
	}
}

func respondByHostname(q *queryContext) (ret bool, rcode int, err error) {
	var hostname string

	if q._suffix != nil {
		for _, searchDomain := range system.SearchDomains {
			system.Console.Infof("=== %s %s %s", searchDomain, *q._suffix, q._domain)

			if searchDomain == *q._suffix {
				hostname = strings.TrimSuffix(q._domain, "."+*q._suffix)
				q._matchedSearchDomain = &searchDomain
				system.Console.Infof("matched question '%s' with search domain '%s'; new hostname: '%s'", q._domain, searchDomain, hostname)
				goto found
			}
		}

		// Hostname has a suffix, but it was not matched in the database
		return

	} else {
		hostname = q._domain
		system.Console.Infof("question '%s' does not have suffix; hostname: '%s'", q._domain, hostname)
	}
found:

	if v, ok := getLeasesByHostname(hostname); ok {
		var m *dns.Msg

		for _, lease := range v {
			if lease.IP.To4() != nil {
				if q._question.Qtype == dns.TypeA {
					system.Console.Infof("answering query %d with value '%s'", q.r.Id, lease.IP.To4().String())
					q.action = ActionDDNSHostname
					//q.Action = "ddns-hostname" // TODO: Special action for RcodeServerFailure?
					m = genResponse(q.r, q._question.Qtype, lease.IP.To4().String())
					err = q.respond(m)
					return true, dns.RcodeSuccess, err
				} else if q._question.Qtype == dns.TypeAAAA {
					system.Console.Debug("lease IP is IPv4, question is AAAA")
					continue
				}

			} else if lease.IP.To16() != nil {
				if q._question.Qtype == dns.TypeAAAA {
					system.Console.Infof("answering query %d with value '%s'", q.r.Id, lease.IP.To16().String())
					q.action = ActionDDNSHostname
					//q.Action = "ddns-hostname" // TODO: Special action for RcodeServerFailure?
					m = genResponse(q.r, q._question.Qtype, lease.IP.To16().String())
					err = q.respond(m)
					return true, dns.RcodeSuccess, err
				} else if q._question.Qtype == dns.TypeA {
					system.Console.Debug("lease IP is IPv6, question is A")
					continue
				}
			}
		}

		system.Console.Debugf("lease with hostname '%s' exists but query type '%s' is not A or AAAA", q._domain, dns.TypeToString[q._question.Qtype])
	}

	return
}

var (
	matchReverse = regexp.MustCompile(`^(?:(?:1?[0-9]{1,2}|2[0-4][0-9]|25[0-5])\.){4}in-addr\.arpa$`)
	getReverse   = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)\.(\d+)\.in-addr.arpa$`)
)

func respondByPTR(q *queryContext) (ret bool, rcode int, err error) {
	if q._question.Qtype != dns.TypePTR {
		return
	}

	if strings.HasSuffix(q._domain, ".in-addr.arpa") && matchReverse.MatchString(q._domain) {
		bits := getReverse.FindStringSubmatch(q._domain)

		ip := net.ParseIP(bits[4] + "." + bits[3] + "." + bits[2] + "." + bits[1])

		if leases, ok := getLeasesByIP(ip); ok {
			q.action = ActionDDNSPTR

			var values []string

			for _, lease := range leases {
				if lease.Hostname == nil {
					return
				}

				values = append(values, *lease.Hostname)
			}

			m := genResponse(q.r, q._question.Qtype, values...)
			err := q.respond(m)

			return true, dns.RcodeSuccess, err
		}
	}

	return
}
