package rule

import (
	"fmt"
	"strings"
)

// Creates a regular expression from a domain that will block the domain and all
// subdomains of that domain.
func genRegex(domain string) string {
	return fmt.Sprintf(`(?:^|.+\.)%s$`, strings.Replace(domain, ".", `\.`, -1))
}
