package rule

import (
    "fmt"
    "strings"
)

func genRegex(domain string) string {
    return fmt.Sprintf(`(?:^|.+\.)%s$`, strings.Replace(domain, ".", `\.`, -1))
}
