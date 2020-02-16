package functions

import (
	"strings"
)

// GenPath returns a slice of strings created by splitting domain into its
// domain components and reversing the result.
func GenPath(domain string) []string {
	dcs := strings.Split(domain, ".")

	for i := len(dcs)/2 - 1; i >= 0; i-- {
		opp := len(dcs) - 1 - i
		dcs[i], dcs[opp] = dcs[opp], dcs[i]
	}

	return dcs
}

// RT returns in without a trailing dot.
func RT(in string) string {
	return strings.TrimSuffix(in, ".")
}
