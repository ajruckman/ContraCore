package rulegen

import (
    "sync"
)

type Node struct {
    Value    string
    Children *map[string]*Node
    Blocked  bool
    FQDN     string

    //Children sync.Map

    lock sync.Mutex
}

func BlockV2(n *Node, fqdn string, path []string) {
    cur := n

    for i, dc := range path {

        isLast := i == len(path)-1

        if isLast {
            newNode := &Node{
                Value:    dc,
                Children: &(map[string]*Node{}),
                Blocked:  true,
                FQDN:     fqdn,
            }
            cur.lock.Lock()
            (*cur.Children)[dc] = newNode
            cur.lock.Unlock()
            return
        }

        cur.lock.Lock()
        v, ok := (*cur.Children)[dc]
        cur.lock.Unlock()

        if ok {
            cur = v

            // This node is already blocked and a complete rule already exists
            // for this path. #A
            if cur.Blocked {
                return
            }
        } else {
            newNode := &Node{
                Value:    dc,
                Children: &(map[string]*Node{}),
            }

            cur.lock.Lock()
            (*cur.Children)[dc] = newNode
            cur.lock.Unlock()

            cur = newNode
        }
    }
}

func BlockV3(n *Node, fqdn string, path []string) {
    cur := n

    for _, dc := range path {
        cur.lock.Lock()
        v, ok := (*cur.Children)[dc]
        cur.lock.Unlock()

        if ok {
            cur = v

            // This node is already blocked and a complete rule already exists
            // for this path. #A
            if cur.Blocked {
                return
            }
        } else {
            newNode := &Node{
                Value:    dc,
                Children: &(map[string]*Node{}),
            }
            cur.lock.Lock()
            (*cur.Children)[dc] = newNode
            cur.lock.Unlock()

            cur = newNode
        }
    }

    cur.Blocked = true
    cur.FQDN = fqdn
}

func BlockV4(n *Node, fqdn string, path []string) {
    cur := n
    l := len(path) - 1

    for i, dc := range path {
        cur.lock.Lock()
        v, ok := (*cur.Children)[dc]
        cur.lock.Unlock()

        if ok {
            cur = v

            // This node is already blocked and a complete rule already exists
            // for this path. #A
            if cur.Blocked {
                return
            }
        } else {
            newNode := &Node{
                Value: dc,
            }

            // This should conserve memory; don't create a bunch of dangling
            // maps
            if i != l {
                newNode.Children = &(map[string]*Node{})
            }

            cur.lock.Lock()
            (*cur.Children)[dc] = newNode
            cur.lock.Unlock()

            cur = newNode
        }
    }

    cur.Blocked = true
    cur.FQDN = fqdn
}

func BlockV5(n *Node, fqdn string, path []string) {
    cur := n
    l := len(path) - 1

    for i, dc := range path {
        if v, ok := (*cur.Children)[dc]; ok {
            // This node is already blocked and a complete rule already exists
            // for this path. #A
            if v.Blocked {
                return
            }

            cur = v
        } else {
            newNode := &Node{
                Value: dc,
            }

            // This should conserve memory; don't create a bunch of dangling
            // maps
            if i != l {
                newNode.Children = &(map[string]*Node{})
            }

            (*cur.Children)[dc] = newNode
            cur = newNode
        }
    }

    cur.Blocked = true
    cur.FQDN = fqdn
}

func BlockV6(n *Node, fqdn string, path []string) {
    cur := n
    l := len(path) - 1

    for i, dc := range path {
        if v, ok := (*cur.Children)[dc]; ok {
            // This node is already blocked and a complete rule already exists
            // for this path. #A
            if v.Blocked {
                return
            }

            cur = v
        } else {
            newNode := &Node{
                Value: dc,
            }

            // This should conserve memory; don't create a bunch of dangling
            // maps
            if i != l {
                newNode.Children = &(map[string]*Node{})
            }

            (*cur.Children)[dc] = newNode
            cur = newNode
        }
    }

    cur.Blocked = true
    cur.FQDN = fqdn
    cur.Children = nil
}

func Read(n *Node, result *[]string) {
    cur := n

    if cur.Blocked {
        *result = append(*result, cur.FQDN)

        return // This should have no effect if #A works.
    }

    for _, v := range *cur.Children {
        Read(v, result)
    }
}
