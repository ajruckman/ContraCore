package serve

import (
    "fmt"
    "regexp"
    "strings"
    "sync"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/rulegen"
)

type ruleTree struct {
    Tree map[string]map[string][]*regexp.Regexp
    lock sync.Mutex
}

var (
    ruleCache = ruleTree{}
)

func (r *ruleTree) check(domain string) bool {
    path := rulegen.GenPath(domain)

    if strings.Count(domain, ".") == 0 {
        return false
    }

    tld := path[0]
    sld := path[1]

    for _, match := range r.Tree[tld][sld] {
        if match.MatchString(domain) {
            return true
        }
    }

    return false
}

func cacheRules() {
    rules, err := db.GetRules()
    Err(err)

    ruleCache.lock.Lock()

    ruleCache.Tree = map[string]map[string][]*regexp.Regexp{}

    l := len(rules)

    for i, rule := range rules {
        if i % 10000 == 0 {
            fmt.Printf("Caching rule %d of %d\n", i, l)
        }

        if _, ok := ruleCache.Tree[rule.TLD]; !ok {
            ruleCache.Tree[rule.TLD] = map[string][]*regexp.Regexp{}
        }

        if _, ok := ruleCache.Tree[rule.TLD][rule.SLD]; !ok {
            ruleCache.Tree[rule.TLD][rule.SLD] = []*regexp.Regexp{}
        }

        ruleCache.Tree[rule.TLD][rule.SLD] = append(ruleCache.Tree[rule.TLD][rule.SLD], regexp.MustCompile(rule.Pattern))
    }

    ruleCache.lock.Unlock()
}
