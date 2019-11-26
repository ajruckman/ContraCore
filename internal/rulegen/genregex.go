package rulegen

import (
    "fmt"
    "strings"
)

func GenRegex(domain string) string {
    return fmt.Sprintf(`(?:^|.+\.)%s$`, strings.Replace(domain, ".", `\.`, -1))
}
