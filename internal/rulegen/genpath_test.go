package rulegen

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestGenPath(t *testing.T) {
    assert.ElementsMatch(t, GenPath("com"), []string{"com"})
    assert.ElementsMatch(t, GenPath("google.com"), []string{"com", "google"})
    assert.ElementsMatch(t, GenPath("ads.google.com"), []string{"com", "google", "ads"})
}
