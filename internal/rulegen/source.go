package rulegen

import (
    "bufio"
    "net/http"
    "strings"
    "sync"

    . "github.com/ajruckman/xlib"
)

var (
    MaxPar    = 4
    ChunkSize = 100
)

var (
    wg    = sync.WaitGroup{}
    total int
    guard = make(chan struct{}, MaxPar)
    root  Node

    linesInB  = make(chan []string)
    linesInBP = make(chan *[]string)
)

func ProcessFromURLs(urls []string) ([]string, int) {
    linesInB = make(chan []string)

    total = 0
    root = Node{}

    for i := 0; i < MaxPar; i++ {
        wg.Add(1)
        go processWorker()
    }

    c := 0
    var batch []string

    for _, url := range urls {
        //fmt.Print("Reading ", url, "... ")
        resp, err := http.Get(url)
        Err(err)
        //fmt.Print(resp.StatusCode, " ")

        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            batch = append(batch, scanner.Text())

            if c > ChunkSize {
                linesInB <- batch
                batch = []string{}
                c = 0
            } else {
                c++
            }
        }

        linesInB <- batch

        //fmt.Println("done", c)
    }
    close(linesInB)

    var res []string
    wg.Wait()
    Read(&root, &res)
    //fmt.Println("/linesOut", len(res))
    return res, total
}

func ProcessFromURLsPointers(urls []string) ([]string, int) {
    linesInBP = make(chan *[]string)

    total = 0
    root = Node{}

    for i := 0; i < MaxPar; i++ {
        wg.Add(1)
        go processWorkerPointers()
    }

    c := 0
    var batch *[]string
    batch = &([]string{})

    for _, url := range urls {
        //fmt.Print("Reading ", url, "... ")
        resp, err := http.Get(url)
        Err(err)
        //fmt.Print(resp.StatusCode, " ")

        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
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

    var res []string
    wg.Wait()
    Read(&root, &res)
    //fmt.Println("/linesOut", len(res))
    return res, total
}

func processWorker() {
    for set := range linesInB {
        for _, t := range set {
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
            BlockV4(&root, t, path)
        }
    }

    wg.Done()
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

            total++

            // You might be inclined to use a 'map[string]struct{}' variable to
            // check whether the current value of 't' has already been seen, but
            // it is actually ~11% faster to let the tree structure handle
            // deduplication.

            ////////// Deduplication

            path := GenPath(t)
            BlockV4(&root, t, path)
        }

        set = nil // Clean up slice from memory after processing
    }

    wg.Done()
}

func ProcessDomainSources(evaluator func(*Node, string, []string), contents []string) ([]string, int) {
    total = 0
    root = Node{}

    // Chunk code from: https://stackoverflow.com/a/35179941/9911189
    chunkSize := (len(contents) + MaxPar - 1) / MaxPar

    for i := 0; i < len(contents); i += chunkSize {
        end := i + chunkSize

        if end > len(contents) {
            end = len(contents)
        }

        guard <- struct{}{}
        wg.Add(1)

        //fmt.Println(fmt.Sprintf("Goroutine %-10d -> %-10d of %d starting", i, end, len(contents)))

        go work(evaluator, contents[i:end], i)
    }
    //fmt.Println()

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

    //fmt.Println("Goroutine", fmt.Sprintf("%-10d", i), "done @", time.Now().Format("04:05.000"))
    wg.Done()
    <-guard
}
