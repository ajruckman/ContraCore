package functions

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestGenPath(t *testing.T) {
    assert.ElementsMatch(t, GenPath("com"), []string{"com"})
    assert.ElementsMatch(t, GenPath("google.com"), []string{"com", "google"})
    assert.ElementsMatch(t, GenPath("ads.google.com"), []string{"com", "google", "ads"})
}

func TestRT(t *testing.T) {
    assert.Equal(t, RT("google.com."), "google.com")
    assert.Equal(t, RT("google.com"), "google.com")
}
