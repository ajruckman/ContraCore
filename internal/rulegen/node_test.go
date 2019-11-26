package rulegen

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestBlockV4(t *testing.T) {
    var (
        root = Node{
            Children: &(map[string]*Node{}),
        }
        res []string
    )

    // Block all domains in 'testSet1'.
    for _, domain := range testSet1 {
        BlockV4(&root, domain, GenPath(domain))
    }
    Read(&root, &res)
    assert.ElementsMatch(t, testSet1Res1, res)

    // Block 'microsoft.com' manually and check that its subdomains are no
    // longer blocked.
    res = nil
    BlockV4(&root, "microsoft.com", GenPath("microsoft.com"))
    Read(&root, &res)
    assert.ElementsMatch(t, testSet1Res2, res)
}
