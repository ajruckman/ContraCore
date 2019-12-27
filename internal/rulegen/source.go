package rulegen

import (
    "bufio"
    "context"
    "fmt"
    "net/http"
    "strings"
    "sync"

    . "github.com/ajruckman/xlib"
    "github.com/jackc/pgx/v4"
    "go.uber.org/atomic"
    "golang.org/x/text/encoding/charmap"

    "github.com/ajruckman/ContraCore/internal/db"
)

var (
    MaxPar    = 1
    ChunkSize = 10000
    SaveSize  = 10000
)

var (
    loadWG = sync.WaitGroup{}
    saveWG = sync.WaitGroup{}

    root Node

    linesInBP chan []string
    rulesIn   chan [][]interface{}

    distinct atomic.Int32
    seen     sync.Map
)

// TODO: check domains from URLs against manual blacklist domains. For example,
// if a rule like '^0as.*\.win' exists and there are domains in a list like
// '0as24865347578835677.win', skip those domains.

func GenFromURLs(urls []string) ([]string, int) {
    var res []string
    linesInBP = make(chan []string)

    root = Node{
        Blocked: atomic.NewBool(false),
    }

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

            if c >= ChunkSize {
                linesInBP <- batch
                batch = []string{}
                c = 0
            } else {
                c++
            }
        }

        linesInBP <- batch

        fmt.Println("done, read", l, "lines")
    }

    close(linesInBP)
    loadWG.Wait()

    Read(&root, &res)

    return res, int(distinct.Load())
}

const naiveMode = true

func SaveRules(res []string) {
    rulesIn = make(chan [][]interface{})

    _, err := db.XDB.Exec(`TRUNCATE TABLE rule;`)
    Err(err)

    saveWG.Add(1)
    go dbSaveWorker()

    c := 0
    var batch [][]interface{}

    for _, rule := range res {
        p := GenPath(rule)

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
                GenRegex(rule),
                rule,
                d,
                tld,
                sld, // Should be safe is #B works.
            })
        } else {
            // This program wil always generate class-2 rules because it omits domains without periods.
            // This means that the value of 'class' is always 2 and 'p[1]' is always safe (#B).
            batch = append(batch, []interface{}{
                GenRegex(rule),
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

func dbSaveWorker() {
    for set := range rulesIn {
        _, err := db.PDB.CopyFrom(context.Background(), pgx.Identifier{"rule"}, []string{"pattern", "domain", "class", "tld", "sld"}, pgx.CopyFromRows(set))
        Err(err)
    }

    saveWG.Done()
}

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

            path := GenPath(t)
            BlockV4(&root, t, path)
        }
    }

    loadWG.Done()
}

// List of identifiers to match for before domains in the domain scanners.
var prefixes = []string{"0.0.0.0", "127.0.0.1", "::", "::0", "::1"}
