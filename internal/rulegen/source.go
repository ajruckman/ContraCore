package rulegen

import (
    "fmt"
    "strings"
    "sync"
    "time"
)

const (
    maxPar = 4
)

var (
    wg    = sync.WaitGroup{}
    total int
    guard = make(chan struct{}, maxPar)
    root  Node
)

func ProcessDomainSources(evaluator func(*Node, string, []string), contents []string) ([]string, int) {
    total = 0

    // Chunk code from: https://stackoverflow.com/a/35179941/9911189
    chunkSize := (len(contents) + maxPar - 1) / maxPar

    for i := 0; i < len(contents); i += chunkSize {
        end := i + chunkSize

        if end > len(contents) {
            end = len(contents)
        }

        guard <- struct{}{}
        wg.Add(1)

        fmt.Println(fmt.Sprintf("Goroutine %-10d -> %-10d of %d starting", i, end, len(contents)))

        go work(evaluator, contents[i:end], i)
    }
    fmt.Println()

    wg.Wait()

    var res []string
    Read(&root, &res)

    return res, total
}

// List of identifiers to match for before domains in the domain scanners.
var prefixes = []string{"0.0.0.0", "127.0.0.1", "::", "::0", "::1"}

func work(evaluator func(*Node, string, []string), content []string, i int) {

    for _, t := range content {
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

        total++

        // You might be inclined to use a 'map[string]struct{}' variable to
        // check whether the current value of 't' has already been seen, but
        // it is actually ~11% faster to let the tree structure handle
        // deduplication.

        ////////// Deduplication

        path := GenPath(t)
        evaluator(&root, t, path)
    }

    fmt.Println("Goroutine", fmt.Sprintf("%-10d", i), "done @", time.Now().Format("04:05.000"))
    wg.Done()
    <-guard
}
