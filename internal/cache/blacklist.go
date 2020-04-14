package cache

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ajruckman/ContraCore/internal/db/contradb"
	"github.com/ajruckman/ContraCore/internal/db/contradb/dbschema"
	"github.com/ajruckman/ContraCore/internal/functions"
	"github.com/ajruckman/ContraCore/internal/system"
)

var BlacklistCache = blacklistTree{}

type blacklistRule struct {
	Pattern *regexp.Regexp
	Expires *time.Time
}

type blacklistTree struct {
	class0Rules []blacklistRule
	class1Rules map[string][]blacklistRule
	class2Rules map[string]map[string][]blacklistRule

	lock sync.RWMutex
}

func ReadBlacklist(callback functions.ProgressCallback) {
	BlacklistCache.lock.Lock()
	defer BlacklistCache.lock.Unlock()

	rules, err := contradb.GetBlacklistRules()
	if _, ok := err.(*contradb.ErrContraDBOffline); ok {
		system.Console.Warning("Not loading blacklist rules because ContraDB is not connected")
		return
	} else if err != nil {
		system.Console.Error("Failed to load blacklist rules from ContraDB with error:")
		system.Console.Error(err.Error())
		return
	}

	BlacklistCache.class1Rules = map[string][]blacklistRule{}
	BlacklistCache.class2Rules = map[string]map[string][]blacklistRule{}

	l := len(rules)
	began := time.Now()

	addRule := func(dbRule dbschema.Blacklist) (newRule blacklistRule, ret bool) {
		newRule.Pattern, err = regexp.Compile(dbRule.Pattern)
		if err != nil {
			system.Console.Warningf("Failed to compile regular expression '%s' in rule %D", dbRule.Pattern, dbRule.ID)
			if ret := callback(fmt.Sprintf("Failed to compile regular expression '%s' in rule %d", dbRule.Pattern, dbRule.ID)); ret {
				return blacklistRule{}, true
			}
		}
		newRule.Expires = dbRule.Expires
		return
	}

	for i, rule := range rules {
		if i%10000 == 0 {
			if ret := callback(fmt.Sprintf("Caching blacklist rule %d of %d", i, l)); ret {
				return
			}
		}

		switch rule.Class {
		case 0:
			newRule, ret := addRule(rule)
			if ret {
				return
			}

			BlacklistCache.class0Rules = append(BlacklistCache.class0Rules, newRule)

		case 1:
			if rule.TLD == nil {
				system.Console.Warningf("Rule '%s' is class 1 but has a null TLD; skipping", rule.Pattern)
				continue
			}

			if _, ok := BlacklistCache.class1Rules[*rule.TLD]; !ok {
				BlacklistCache.class1Rules[*rule.TLD] = []blacklistRule{}
			}

			newRule, ret := addRule(rule)
			if ret {
				return
			}

			BlacklistCache.class1Rules[*rule.TLD] = append(BlacklistCache.class1Rules[*rule.TLD], newRule)

			//BlacklistCache.class1Rules[rule.TLD] = append(BlacklistCache.class1Rules[rule.TLD], regexp.MustCompile(rule.Pattern))

		case 2:
			if rule.TLD == nil {
				system.Console.Warningf("Rule '%s' is class 2 but has a null TLD; skipping", rule.Pattern)
				continue
			} else if rule.SLD == nil {
				system.Console.Warningf("Rule '%s' is class 2 but has a null SLD; skipping", rule.Pattern)
				continue
			}

			if _, ok := BlacklistCache.class2Rules[*rule.TLD]; !ok {
				BlacklistCache.class2Rules[*rule.TLD] = map[string][]blacklistRule{}
			}

			if _, ok := BlacklistCache.class2Rules[*rule.TLD][*rule.SLD]; !ok {
				BlacklistCache.class2Rules[*rule.TLD][*rule.SLD] = []blacklistRule{}
			}

			newRule, ret := addRule(rule)
			if ret {
				return
			}

			BlacklistCache.class2Rules[*rule.TLD][*rule.SLD] = append(BlacklistCache.class2Rules[*rule.TLD][*rule.SLD], newRule)

			//blacklistRule, err := newRule(rule.Pattern, rule.Expires)
			//
			////pattern, err := regexp.Compile(rule.Pattern)
			//if err != nil {
			//	system.Console.Warningf("Failed to compile regular expression '%s' in rule %D", rule.Pattern, rule.ID)
			//	if ret := callback(fmt.Sprintf("Failed to compile regular expression '%s' in rule %d", rule.Pattern, rule.ID)); ret {
			//		return
			//	}
			//	continue
			//}
			//
			//BlacklistCache.class2Rules[rule.TLD][rule.SLD] = append(BlacklistCache.class2Rules[rule.TLD][rule.SLD], pattern)
		}
	}

	callback(fmt.Sprintf("%d blacklist rules compiled and stored in %v", l, time.Since(began)))
}

func (t *blacklistTree) Check(domain string) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()

	path := functions.GenPath(domain)
	if strings.Count(domain, ".") == 0 {
		return false
	}

	c := len(strings.Split(domain, "."))

	if c > 2 {
		c = 2
	}

	now := time.Now()

	switch c {

	case 2:
		tld := path[0]
		sld := path[1]

		began := time.Now()
		for _, rule := range t.class2Rules[tld][sld] {
			if rule.Expires != nil {
				if rule.Expires.Before(now) {
					system.Console.Infof("Blacklist rule has expired: %v at %v", rule.Pattern, rule.Expires)
					return false
				}
			}
			if rule.Pattern.MatchString(domain) {
				return true
			}
		}
		if system.LogRuleLookupDurations {
			system.Console.Debugf("%v: 2 > %v", domain, time.Since(began))
		}

		fallthrough

	case 1:
		tld := path[0]

		began := time.Now()
		for _, rule := range t.class1Rules[tld] {
			if rule.Expires != nil {
				if rule.Expires.Before(now) {
					system.Console.Infof("Blacklist rule has expired: %v at %v", rule.Pattern, rule.Expires)
					return false
				}
			}
			if rule.Pattern.MatchString(domain) {
				return true
			}
		}
		if system.LogRuleLookupDurations {
			system.Console.Debugf("%v: 1 > %v", domain, time.Since(began))
		}

		fallthrough

	case 0:
		began := time.Now()
		for _, rule := range t.class0Rules {
			if rule.Expires != nil {
				if rule.Expires.Before(now) {
					system.Console.Debugf("Blacklist rule has expired: %v at %v", rule.Pattern, rule.Expires)
					return false
				}
			}
			if rule.Pattern.MatchString(domain) {
				return true
			}
		}
		if system.LogRuleLookupDurations {
			system.Console.Debugf("%v: 0 > %v", domain, time.Since(began))
		}
	}

	return false
}
