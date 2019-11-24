package rulegen

import (
    "bufio"
    "strings"
    "time"
)

func ReadDomainList(evaluator func(*Node, string, []string), scanner *bufio.Scanner) (res []string, total int, kept int, ratio float64, duration time.Duration) {
    root := Node{
        Children: &map[string]*Node{},
    }

    //total := 0
    //began := time.Now()
    //var began time.Time

    //for _, url := range urls {
        //resp, err := http.Get(url)
        //Err(err)
        //began = time.Now()

        //scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            t := scanner.Text()

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

            //evaluator(&root)

            path := GenPath(t)
            evaluator(&root, t, path)
            //root.BlockV3(t, path)

            _ = evaluator

            // V1: 1269353
            // V2: 1269353
            // V3: 1269353

            //res = append(res, t)
        }
    //}

    Read(&root, &res, 0)

    //duration = time.Since(began)
    //kept = len(res)
    //ratio = float64(kept) / float64(total)

    return
}

var prefixes = []string{"0.0.0.0", "127.0.0.1", "::", "::0", "::1"}
