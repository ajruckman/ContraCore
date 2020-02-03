package rule

import (
    "sync"

    "go.uber.org/atomic"
)

type node struct {
    Value   string
    Blocked *atomic.Bool
    FQDN    *atomic.String

    Children sync.Map
}

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
