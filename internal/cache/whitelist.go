package cache

import (
	"fmt"
	"net"
	"regexp"
	"sync"
	"time"

	"github.com/ajruckman/ContraCore/internal/db/contradb"
	"github.com/ajruckman/ContraCore/internal/db/contradb/dbschema"
	"github.com/ajruckman/ContraCore/internal/functions"
	"github.com/ajruckman/ContraCore/internal/system"
)

var WhitelistCache whitelistTree

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

func ReadWhitelist(callback functions.ProgressStatusCallback) {
	WhitelistCache.lock.Lock()
	defer WhitelistCache.lock.Unlock()

	WhitelistCache.rules = []whitelistRule{}

	rules, err := contradb.GetWhitelistRules()
	if _, ok := err.(*contradb.ErrContraDBOffline); ok {
		system.Console.Warning("Not loading rules because ContraDB is not connected")
		callback("Not loading rules because ContraDB is not connected", err)
		return
	} else if err != nil {
		system.Console.Error("Failed to load rules from ContraDB with error:")
		system.Console.Error(err.Error())
		callback("Failed to load rules from ContraDB", err)
		return
	}

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

		WhitelistCache.rules = append(WhitelistCache.rules, w)
	}

	callback(fmt.Sprintf("%d whitelist rules compiled and stored in %v", len(rules), time.Since(began)), nil)
}

func (t *whitelistTree) Check(domain string, client net.IP, _mac *string, _hostname *string, _vendor *string) bool {
	WhitelistCache.lock.Lock()
	defer WhitelistCache.lock.Unlock()

	if system.LogRuleLookupDurations {
		began := time.Now()

		defer func() {
			system.Console.Debugf("Checked %d whitelist rules against query %s in %v", len(WhitelistCache.rules), domain, time.Since(began))
		}()
	}

	rules := WhitelistCache.find(domain)

	for _, rule := range rules {

		if rule.IPs != nil {
			for _, ip := range *rule.IPs {
				if client.Equal(ip) {
					system.Console.Infof("Query %s is whitelisted by rule #%d: IP '%s' = '%s'", domain, rule.ID, client.String(), ip.String())
					return true
				}
			}
		}

		if rule.Subnets != nil {
			for _, subnet := range *rule.Subnets {
				if subnet.Contains(client) {
					system.Console.Infof("Query %s is whitelisted by rule #%d: IP '%s' is in subnet '%s'", domain, rule.ID, client.String(), subnet.String())
					return true
				}
			}
		}

		if rule.MACs != nil && _mac != nil {
			for _, mac := range *rule.MACs {
				if *_mac == mac.String() {
					system.Console.Infof("Query %s is whitelisted by rule #%d: MAC '%s' = '%s'", domain, rule.ID, *_mac, mac.String())
					return true
				}
			}
		}

		if rule.Vendors != nil && _vendor != nil {
			for _, vendor := range *rule.compiledVendors {
				if vendor.MatchString(*_vendor) {
					system.Console.Infof("Query %s is whitelisted by rule #%d: vendor '%s' matches '%s'", domain, rule.ID, *_vendor, vendor.String())
					return true
				}
			}
		}

		if rule.Hostnames != nil && _hostname != nil {
			for _, hostname := range *rule.compiledHostnames {
				if hostname.MatchString(*_hostname) {
					system.Console.Infof("Query %s is whitelisted by rule #%d: hostname '%s' matches '%s'", domain, rule.ID, *_hostname, hostname.String())
					return true
				}
			}
		}

	}

	return false
}

func (t *whitelistTree) find(domain string) (res []whitelistRule) {
	for _, rule := range t.rules {
		if rule.compiledPattern.MatchString(domain) {
			res = append(res, rule)
		}
	}
	return
}
