package cache

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ajruckman/ContraCore/internal/db/contradb"
	"github.com/ajruckman/ContraCore/internal/functions"
	"github.com/ajruckman/ContraCore/internal/system"
)

var BlacklistCache = blacklistTree{}

type blacklistTree struct {
	class0Rules []*regexp.Regexp
	class1Rules map[string][]*regexp.Regexp
	class2Rules map[string]map[string][]*regexp.Regexp

	lock sync.RWMutex
}

func ReadBlacklist(callback functions.ProgressCallback) {
	BlacklistCache.lock.Lock()
	defer BlacklistCache.lock.Unlock()

	rules, err := contradb.GetBlacklistRules()
	if _, ok := err.(*contradb.ErrContraDBOffline); ok {
		system.Console.Warning("not loading blacklist rules because ContraDB is not connected")
		return
	} else if err != nil {
		system.Console.Error("failed to load blacklist rules from ContraDB with error:")
		system.Console.Error(err.Error())
		return
	}

	BlacklistCache.class1Rules = map[string][]*regexp.Regexp{}
	BlacklistCache.class2Rules = map[string]map[string][]*regexp.Regexp{}

	l := len(rules)
	began := time.Now()

	for i, rule := range rules {
		if i%10000 == 0 {
			if ret := callback(fmt.Sprintf("Caching blacklist rule %d of %d", i, l)); ret {
				return
			}
		}

		switch rule.Class {
		case 0:
			BlacklistCache.class0Rules = append(BlacklistCache.class0Rules, regexp.MustCompile(rule.Pattern))

		case 1:
			if _, ok := BlacklistCache.class1Rules[rule.TLD]; !ok {
				BlacklistCache.class1Rules[rule.TLD] = []*regexp.Regexp{}
			}

			BlacklistCache.class1Rules[rule.TLD] = append(BlacklistCache.class1Rules[rule.TLD], regexp.MustCompile(rule.Pattern))

		case 2:
			if _, ok := BlacklistCache.class2Rules[rule.TLD]; !ok {
				BlacklistCache.class2Rules[rule.TLD] = map[string][]*regexp.Regexp{}
			}

			if _, ok := BlacklistCache.class2Rules[rule.TLD][rule.SLD]; !ok {
				BlacklistCache.class2Rules[rule.TLD][rule.SLD] = []*regexp.Regexp{}
			}

			BlacklistCache.class2Rules[rule.TLD][rule.SLD] = append(BlacklistCache.class2Rules[rule.TLD][rule.SLD], regexp.MustCompile(rule.Pattern))
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

	switch c {

	case 2:
		tld := path[0]
		sld := path[1]

		began := time.Now()
		for _, rule := range t.class2Rules[tld][sld] {
			if rule.MatchString(domain) {
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
			if rule.MatchString(domain) {
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
			if rule.MatchString(domain) {
				return true
			}
		}
		if system.LogRuleLookupDurations {
			system.Console.Debugf("%v: 0 > %v", domain, time.Since(began))
		}
	}

	return false
}
