package process

import (
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"

	"github.com/ajruckman/ContraCore/internal/db/contradb"
	"github.com/ajruckman/ContraCore/internal/functions"
	"github.com/ajruckman/ContraCore/internal/log"
	"github.com/ajruckman/ContraCore/internal/system"
)

var blacklistCache = blacklistTree{}

type blacklistTree struct {
	class0Rules []*regexp.Regexp
	class1Rules map[string][]*regexp.Regexp
	class2Rules map[string]map[string][]*regexp.Regexp

	lock sync.Mutex
}

func (t *blacklistTree) check(domain string) bool {
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
		if log.LogRuleLookupDurations {
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
		if log.LogRuleLookupDurations {
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
		if log.LogRuleLookupDurations {
			system.Console.Debugf("%v: 0 > %v", domain, time.Since(began))
		}
	}

	return false
}

func readBlacklistRules() {
	rules, err := contradb.GetBlacklistRules()
	if _, ok := err.(*contradb.ErrContraDBOffline); ok {
		system.Console.Warning("not loading blacklist rules because ContraDB is not connected")
		return
	} else if err != nil {
		system.Console.Error("failed to load blacklist rules from ContraDB with error:")
		system.Console.Error(err.Error())
		return
	}

	blacklistCache.lock.Lock()

	blacklistCache.class1Rules = map[string][]*regexp.Regexp{}
	blacklistCache.class2Rules = map[string]map[string][]*regexp.Regexp{}

	l := len(rules)
	began := time.Now()

	for i, rule := range rules {
		if i%10000 == 0 {
			system.Console.Infof("Caching blacklist rule %d of %d", i, l)
		}

		switch rule.Class {
		case 0:
			blacklistCache.class0Rules = append(blacklistCache.class0Rules, regexp.MustCompile(rule.Pattern))

		case 1:
			if _, ok := blacklistCache.class1Rules[rule.TLD]; !ok {
				blacklistCache.class1Rules[rule.TLD] = []*regexp.Regexp{}
			}

			blacklistCache.class1Rules[rule.TLD] = append(blacklistCache.class1Rules[rule.TLD], regexp.MustCompile(rule.Pattern))

		case 2:
			if _, ok := blacklistCache.class2Rules[rule.TLD]; !ok {
				blacklistCache.class2Rules[rule.TLD] = map[string][]*regexp.Regexp{}
			}

			if _, ok := blacklistCache.class2Rules[rule.TLD][rule.SLD]; !ok {
				blacklistCache.class2Rules[rule.TLD][rule.SLD] = []*regexp.Regexp{}
			}

			blacklistCache.class2Rules[rule.TLD][rule.SLD] = append(blacklistCache.class2Rules[rule.TLD][rule.SLD], regexp.MustCompile(rule.Pattern))
		}
	}

	system.Console.Infof("%d blacklist rules compiled and stored in %v", l, time.Since(began))

	blacklistCache.lock.Unlock()
}

func blacklist(q *queryContext) (ret bool, rcode int, err error) {
	if blacklistCache.check(q._domain) {
		q.action = ActionBlock
		var m *dns.Msg
		var v string

		switch q._question.Qtype {
		case dns.TypeA:
			v = system.SpoofedA
			m = genResponse(q.r, q._question.Qtype, v)

		case dns.TypeAAAA:
			v = system.SpoofedAAAA
			m = genResponse(q.r, q._question.Qtype, v)

		case dns.TypeCNAME:
			v = system.SpoofedCNAME
			m = genResponse(q.r, q._question.Qtype, v)

		default:
			v = system.SpoofedDefault
		}

		system.Console.Infof("blocking query %d with value '%s'", q.r.Id, v)

		m = genResponse(q.r, q._question.Qtype, v)
		err = q.respond(m)

		return true, dns.RcodeSuccess, err
	}

	return
}
