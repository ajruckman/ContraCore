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
    "github.com/tevino/abool"

    "github.com/ajruckman/ContraCore/internal/db"
)

var (
    MaxPar    = 4
    ChunkSize = 10000
    SaveSize  = 10000
)

var (
    wg = sync.WaitGroup{}
    //total int
    //guard = make(chan struct{}, MaxPar)
    root Node

    //linesInB  = make(chan []string)
    linesInBP = make(chan *[]string)

    rulesIn = make(chan *[][]interface{})
)

func ProcessFromURLsPointers(urls []string) (res []string, total int) {
    linesInBP = make(chan *[]string)

    //total = 0
    root = Node{
        Blocked: abool.New(),
    }

    for i := 0; i < MaxPar; i++ {
        wg.Add(1)
        go processWorkerPointers()
    }

    c := 0
    //var batch *[]string
    //batch = &([]string{})

    for _, url := range urls {
        //fmt.Print("Reading ", url, "... ")
        resp, err := http.Get(url)
        Err(err)
        //fmt.Print(resp.StatusCode, " ")

        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            total++

            *batch = append(*batch, scanner.Text())

            if c > ChunkSize {
                linesInBP <- batch
                batch = &([]string{})
                c = 0
            } else {
                c++
            }
        }

        linesInBP <- batch

        //fmt.Println("done", c)
    }
    close(linesInBP)

    wg.Wait()

    Read(&root, &res)

    //save(&res)

    //fmt.Println("/linesOut", len(res))
    return
}

func save(res *[]string) {
    _, err := db.XDB.Exec(`TRUNCATE TABLE rule;`)
    Err(err)

    go dbSaveWorker()

    c := 0
    var batch *[][]interface{}
    batch = &([][]interface{}{})

    for _, rule := range *res {
        *batch = append(*batch, []interface{}{
            rule,
            GenRegex(rule),
        })

        if c > SaveSize {
            rulesIn <- batch
            batch = &([][]interface{}{})
            c = 0
        } else {
            c++
        }
    }

}

func dbSaveWorker() {
    for set := range rulesIn {
        copied, err := db.PDB.CopyFrom(context.Background(), pgx.Identifier{"rule"}, []string{"domain", "pattern"}, pgx.CopyFromRows(*set))
        Err(err)

        fmt.Println(copied)
    }
}

func processWorkerPointers() {
    for set := range linesInBP {
        for _, t := range *set {
            if strings.HasPrefix(t, "#") {
                continue
            }

            t = strings.TrimSpace(t)
            t = strings.TrimSuffix(t, ".") // Some lists have trailing dots

            // Match lines like '0.0.0.0 ads.google.com'
            if strings.Contains(t, " ") {

                // Skip lines with more than 1 space
                if strings.Count(t, " ") > 1 {
                    continue
                }

                for _, prefix := range prefixes {
                    if strings.HasPrefix(t, prefix+" ") {
                        t = strings.TrimPrefix(t, prefix+" ")
                        goto next
                    }
                }
                continue // Skip if none matched

            next:
            }

            c := strings.Count(t, ".")

            // Preserve domains like 'www.com'
            if c >= 2 && strings.HasPrefix(t, "www.") {
                t = strings.TrimPrefix(t, "www.")
            } else if c == 0 {
                continue
            }

            // You might be inclined to use a 'map[string]struct{}' variable to
            // check whether the current value of 't' has already been seen, but
            // it is ~11% faster to let the tree structure handle deduplication.

            ////////// Deduplication

            path := GenPath(t)
            BlockV4(&root, t, path)
        }

        //set = nil // Clean up slice from memory after processing
    }

    wg.Done()
}

// List of identifiers to match for before domains in the domain scanners.
var prefixes = []string{"0.0.0.0", "127.0.0.1", "::", "::0", "::1"}
