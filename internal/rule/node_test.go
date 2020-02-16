package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ajruckman/ContraCore/internal/functions"
)

func TestBlockV4(t *testing.T) {
	var res []string

	// Block all domains in 'testSet1'.
	for _, domain := range testSet1 {
		block(&root, domain, functions.GenPath(domain))
	}
	read(&root, &res)
	assert.ElementsMatch(t, testSet1Res1, res)

	// Block 'microsoft.com' manually and check that its subdomains are no
	// longer blocked.
	res = nil
	block(&root, "microsoft.com", functions.GenPath("microsoft.com"))
	read(&root, &res)
	assert.ElementsMatch(t, testSet1Res2, res)
}
