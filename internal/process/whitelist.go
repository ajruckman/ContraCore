package process

import (
	"regexp"
	"sync"
	"time"

	"github.com/ajruckman/ContraCore/internal/db/contradb"
	"github.com/ajruckman/ContraCore/internal/db/contradb/dbschema"
	"github.com/ajruckman/ContraCore/internal/log"
	"github.com/ajruckman/ContraCore/internal/system"
)

var whitelistCache whitelistTree

type whitelistRule struct {
	dbschema.Whitelist

	compiledPattern   *regexp.Regexp
	compiledVendors   *[]*regexp.Regexp
	compiledHostnames *[]*regexp.Regexp
}

type whitelistTree struct {
	rules []whitelistRule

	lock sync.Mutex
}

func (t *whitelistTree) find(domain string) (res []whitelistRule) {
	for _, rule := range t.rules {
		if rule.compiledPattern.MatchString(domain) {
			res = append(res, rule)
		}
	}
	return
}

func readWhitelistRules() {
	rules, err := contradb.GetWhitelistRules()
	if _, ok := err.(*contradb.ErrContraDBOffline); ok {
		system.Console.Warning("not loading rules because ContraDB is not connected")
		return
	} else if err != nil {
		system.Console.Error("failed to load rules from ContraDB with error:")
		system.Console.Error(err.Error())
		return
	}

	whitelistCache.lock.Lock()

	began := time.Now()

	for _, rule := range rules {
		w := whitelistRule{
			Whitelist:       rule,
			compiledPattern: regexp.MustCompile(rule.Pattern),
		}

		if rule.Vendors != nil {
			w.compiledVendors = &[]*regexp.Regexp{}
			for _, vendor := range *rule.Vendors {
				*w.compiledVendors = append(*w.compiledVendors, regexp.MustCompile(vendor))
			}
		}

		if rule.Hostnames != nil {
			w.compiledHostnames = &[]*regexp.Regexp{}
			for _, hostname := range *rule.Hostnames {
				*w.compiledHostnames = append(*w.compiledHostnames, regexp.MustCompile(hostname))
			}
		}

		whitelistCache.rules = append(whitelistCache.rules, w)
	}

	system.Console.Infof("%d whitelist rules compiled and stored in %v", len(rules), time.Since(began))

	whitelistCache.lock.Unlock()
}

func whitelist(q *queryContext) (found bool) {

	if log.LogRuleLookupDurations {
		began := time.Now()

		defer func() {
			system.Console.Debugf("checked %d whitelist rules against query %d in %v", len(whitelistCache.rules), q.r.Id, time.Since(began))
		}()
	}

	rules := whitelistCache.find(q._domain)

	for _, rule := range rules {

		if rule.IPs != nil {
			for _, ip := range *rule.IPs {
				if q._client.Equal(ip) {
					system.Console.Infof("query %d is whitelisted by rule #%d: IP '%s' = '%s'", q.r.Id, rule.ID, q._client.String(), ip.String())
					return true
				}
			}
		}

		if rule.Subnets != nil {
			for _, subnet := range *rule.Subnets {
				if subnet.Contains(q._client) {
					system.Console.Infof("query %d is whitelisted by rule #%d: IP '%s' is in subnet '%s'", q.r.Id, rule.ID, q._client.String(), subnet.String())
					return true
				}
			}
		}

		if rule.MACs != nil && q.mac != nil {
			for _, mac := range *rule.MACs {
				if *q.mac == mac.String() {
					system.Console.Infof("query %d is whitelisted by rule #%d: MAC '%s' = '%s'", q.r.Id, rule.ID, *q.mac, mac.String())
					return true
				}
			}
		}

		if rule.Vendors != nil && q.vendor != nil {
			for _, vendor := range *rule.compiledVendors {
				if vendor.MatchString(*q.vendor) {
					system.Console.Infof("query %d is whitelisted by rule #%d: vendor '%s' matches '%s'", q.r.Id, rule.ID, *q.vendor, vendor.String())
					return true
				}
			}
		}

		if rule.Hostnames != nil && q.hostname != nil {
			for _, hostname := range *rule.compiledHostnames {
				if hostname.MatchString(*q.hostname) {
					system.Console.Infof("query %d is whitelisted by rule #%d: hostname '%s' matches '%s'", q.r.Id, rule.ID, *q.hostname, hostname.String())
					return true
				}
			}
		}

	}

	return
}
