package rule

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"sync"

	. "github.com/ajruckman/xlib"
	"github.com/jackc/pgx/v4"
	"go.uber.org/atomic"
	"golang.org/x/text/encoding/charmap"

	"github.com/ajruckman/ContraCore/internal/db/contradb"
	"github.com/ajruckman/ContraCore/internal/functions"
)

var (
	// Controls the maximum number of concurrent rule gen workers.
	MaxPar = 1

	// Controls the maximum number of rule source lines to send to a rule gen
	// worker.
	ChunkSize = 10000

	// Controls the maximum number of generated rules to save to ContraDB in a
	// batch.
	SaveSize = 10000
)

var (
	loadWG = sync.WaitGroup{}
	saveWG = sync.WaitGroup{}

	root node

	linesInBP chan []string
	rulesIn   chan [][]interface{}

	// The number of distinct domains seen before deduplication.
	distinct atomic.Int32
	seen     sync.Map
)

// TODO: check domains from URLs against manual blacklist domains. For example,
// if a rule like '^0as.*\.win' exists and there are domains in a list like
// '0as24865347578835677.win', skip those domains.

// Generates a list of regular expressions that will block all the domains found
// in the passed URLs. Each regular expression will block 1 or more domains
// found in the passed URLs, and it will not return redundant rules.
//
// For example, if a rule exists for 'bbb.aaa' and domains 'ccc.bbb.aaa' and
// 'ddd.bbb.aaa' are found in the passed URLs, only a rule blocking 'bbb.aaa'
// and its subdomains will be returned, as rules blocking its subdomains would
// be redundant.
func GenFromURLs(urls []string) ([]string, int) {
	var res []string
	linesInBP = make(chan []string)

	root = node{
		Blocked: atomic.NewBool(false),
	}

	// Spawn multiple rule processor workers.
	for i := 0; i < MaxPar; i++ {
		loadWG.Add(1)
		go ruleGenWorker()
	}

	c := 0
	var batch []string

	for _, url := range urls {
		fmt.Print("Reading ", url, "... ")
		resp, err := http.Get(url)
		Err(err)
		fmt.Print(resp.StatusCode, " ")

		conv := charmap.Windows1252.NewDecoder().Reader(resp.Body)
		scanner := bufio.NewScanner(conv)
		l := 0

		for scanner.Scan() {
			l++

			batch = append(batch, scanner.Text())

			// Send batches of lines onto the rule source channel.
			if c >= ChunkSize {
				linesInBP <- batch
				batch = []string{}
				c = 0
			} else {
				c++
			}
		}

		// Send remaining lines onto the rule source channel.
		linesInBP <- batch

		fmt.Println("done, read", l, "lines")
	}

	close(linesInBP)
	loadWG.Wait()

	read(&root, &res)

	return res, int(distinct.Load())
}

// If true, the rule generator will not always generate class-2 rules.
const naiveMode = false

// Processes batches of rule source lines pushed onto the rule source channel.
func ruleGenWorker() {
	for set := range linesInBP {
		for _, t := range set {
			if strings.HasPrefix(t, "#") {
				continue
			}

			t = strings.TrimSpace(t)
			t = strings.TrimSuffix(t, ".") // Some lists have trailing dots.

			// Match lines like '0.0.0.0 ads.google.com'.
			if strings.Contains(t, " ") || strings.Contains(t, "\t") {

				// Skip lines with more than 1 space.
				if strings.Count(t, " ") > 1 {
					continue
				}

				for _, prefix := range prefixes {
					if strings.HasPrefix(t, prefix+" ") || strings.HasPrefix(t, prefix+"\t") {
						t = strings.TrimPrefix(t, prefix)
						t = strings.TrimSpace(t)
						goto next
					}
				}
				continue // Skip if none matched

			next:
			}

			c := strings.Count(t, ".")

			// Preserve domains like 'www.com'.
			if c >= 2 && strings.HasPrefix(t, "www.") {
				t = strings.TrimPrefix(t, "www.")
			} else if c == 0 {
				continue // Skip single-word domains like 'country'. #B
			}

			// You might be inclined to use a 'map[string]struct{}' variable to
			// check whether the current value of 't' has already been seen, but
			// it is ~11% faster to let the tree structure handle deduplication.

			if _, ok := seen.Load(t); !ok {
				seen.Store(t, struct{}{})
				distinct.Add(1)
			}

			////////// Deduplication

			path := functions.GenPath(t)
			block(&root, t, path)
		}
	}

	loadWG.Done()
}

// Identifiers to match for before domains in the domain scanners.
var prefixes = [...]string{"0.0.0.0", "127.0.0.1", "::", "::0", "::1"}

// Saves a slice of rules (regular expressions) to ContraDB in batches.
func SaveRules(res []string) {
	rulesIn = make(chan [][]interface{})

	_, err := contradb.Exec(`TRUNCATE TABLE rule;`)
	Err(err)

	saveWG.Add(1)
	go dbSaveWorker()

	c := 0
	var batch [][]interface{}

	for _, rule := range res {
		p := functions.GenPath(rule)

		if naiveMode {
			d := strings.Count(rule, ".")
			if d > 2 {
				d = 2
			}

			tld := p[0]
			var sld *string = nil

			if d == 2 {
				sld = &p[1]
			}

			batch = append(batch, []interface{}{
				genRegex(rule),
				rule,
				d,
				tld,
				sld, // Should be safe is #B works.
			})
		} else {
			// This program wil always generate class-2 rules because it omits
			// domains without periods. This means that the value of 'class' is
			// always 2 and 'p[1]' is always safe (#B).
			batch = append(batch, []interface{}{
				genRegex(rule),
				rule,
				2,
				p[0],
				p[1],
			})
		}

		if c >= SaveSize {
			rulesIn <- batch
			batch = [][]interface{}{}
			c = 0
		} else {
			c++
		}
	}

	rulesIn <- batch

	close(rulesIn)
	saveWG.Wait()
}

// Saves rule batches pushed onto the rule save channel to ContraDB.
func dbSaveWorker() {
	for set := range rulesIn {
		_, err := contradb.CopyFrom(pgx.Identifier{"rule"}, []string{"pattern", "domain", "class", "tld", "sld"}, pgx.CopyFromRows(set))
		Err(err)
	}

	saveWG.Done()
}
