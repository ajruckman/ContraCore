package process

import (
    "regexp"
    "strings"
    "sync"
    "time"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/db/contradb"
    "github.com/ajruckman/ContraCore/internal/functions"
    "github.com/ajruckman/ContraCore/internal/log"
    "github.com/ajruckman/ContraCore/internal/system"
)

type ruleTree struct {
    class0Rules []*regexp.Regexp
    class1Rules map[string][]*regexp.Regexp
    class2Rules map[string]map[string][]*regexp.Regexp

    lock sync.Mutex
}

var (
    ruleCache = ruleTree{}
)

func (r *ruleTree) check(domain string) bool {
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
        for _, rule := range r.class2Rules[tld][sld] {
            if rule.MatchString(domain) {
                return true
            }
        }
        if log.LogDurations {
            system.Console.Debugf("%v: 2 > %v", domain, time.Since(began))
        }

        fallthrough

    case 1:
        tld := path[0]

        began := time.Now()
        for _, rule := range r.class1Rules[tld] {
            if rule.MatchString(domain) {
                return true
            }
        }
        if log.LogDurations {
            system.Console.Debugf("%v: 1 > %v", domain, time.Since(began))
        }

        fallthrough

    case 0:
        began := time.Now()
        for _, rule := range r.class0Rules {
            if rule.MatchString(domain) {
                return true
            }
        }
        if log.LogDurations {
            system.Console.Debugf("%v: 0 > %v", domain, time.Since(began))
        }
    }

    return false
}

func readRules() {
    rules, err := contradb.GetRules()
    Err(err)

    ruleCache.lock.Lock()

    ruleCache.class1Rules = map[string][]*regexp.Regexp{}
    ruleCache.class2Rules = map[string]map[string][]*regexp.Regexp{}

    l := len(rules)
    began := time.Now()

    for i, rule := range rules {
        if i%10000 == 0 {
            system.Console.Infof("Caching rule %d of %d", i, l)
        }

        switch rule.Class {
        case 0:
            ruleCache.class0Rules = append(ruleCache.class0Rules, regexp.MustCompile(rule.Pattern))

        case 1:
            if _, ok := ruleCache.class1Rules[rule.TLD]; !ok {
                ruleCache.class1Rules[rule.TLD] = []*regexp.Regexp{}
            }

            ruleCache.class1Rules[rule.TLD] = append(ruleCache.class1Rules[rule.TLD], regexp.MustCompile(rule.Pattern))

        case 2:
            if _, ok := ruleCache.class2Rules[rule.TLD]; !ok {
                ruleCache.class2Rules[rule.TLD] = map[string][]*regexp.Regexp{}
            }

            if _, ok := ruleCache.class2Rules[rule.TLD][rule.SLD]; !ok {
                ruleCache.class2Rules[rule.TLD][rule.SLD] = []*regexp.Regexp{}
            }

            ruleCache.class2Rules[rule.TLD][rule.SLD] = append(ruleCache.class2Rules[rule.TLD][rule.SLD], regexp.MustCompile(rule.Pattern))
        }
    }

    system.Console.Infof("%d rules compiled and stored in %v", len(rules), time.Since(began))

    ruleCache.lock.Unlock()
}
