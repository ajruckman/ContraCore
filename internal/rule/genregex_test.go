package rule

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenRegexFormat(t *testing.T) {
	generated := genRegex("ads.google.com")
	assert.Equal(t, generated, `(?:^|.+\.)ads\.google\.com$`)
}

func TestGenRegex(t *testing.T) {
	var (
		generated string
		matches   bool
		err       error
	)

	for _, domain := range testSet1 {
		generated = genRegex(domain)

		// Should match

		matches, err = regexp.MatchString(generated, domain)
		assert.NoError(t, err)
		assert.True(t, matches)

		// Should not match

		matches, err = regexp.MatchString(generated, "."+domain)
		assert.NoError(t, err)
		assert.False(t, matches)

		matches, err = regexp.MatchString(generated, domain+".")
		assert.NoError(t, err)
		assert.False(t, matches)

		// Should match

		matches, err = regexp.MatchString(generated, "www."+domain)
		assert.NoError(t, err)
		assert.True(t, matches)

		matches, err = regexp.MatchString(generated, "www.ads."+domain)
		assert.NoError(t, err)
		assert.True(t, matches)
	}
}
