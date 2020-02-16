package rule

import (
	"sync"

	"go.uber.org/atomic"
)

// A node containing information about a domain component.
type node struct {
	// The domain component's value.
	Value string

	// Whether or not this domain component is blocked.
	Blocked *atomic.Bool

	// The complete domain which this domain component is a part of.
	// Nil if this is not the most-specific domain component in the domain.
	FQDN *atomic.String

	// All children (subdomains) of this node.
	Children sync.Map
}

// Stores a domain into the passed domain component tree.
func block(n *node, fqdn string, path []string) {
	cur := n

	for _, dc := range path {
		v, ok := cur.Children.Load(dc)

		if ok {
			cur = v.(*node)

			// This node is already blocked and a complete rule already exists
			// for this path. #A
			if cur.Blocked.Load() {
				return
			}
		} else {
			newNode := node{
				Value:   dc,
				Blocked: atomic.NewBool(false),
				FQDN:    &atomic.String{},
			}

			cur.Children.Store(dc, &newNode)
			cur = &newNode
		}
	}

	cur.Blocked.Store(true)
	cur.FQDN.Store(fqdn)
}

// Returns the least-specific blocked domains in the domain component tree.
func read(n *node, result *[]string) {
	cur := n

	if cur.Blocked.Load() {
		*result = append(*result, cur.FQDN.Load())

		return // This should have no effect if #A works.
	}

	(cur.Children).Range(func(key, value interface{}) bool {
		read(value.(*node), result)
		return true
	})
}
