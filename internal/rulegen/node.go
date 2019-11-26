package rulegen

import (
    "sync"
)

type Node struct {
    Value    string
    //Children *map[string]*Node
    Blocked  bool
    FQDN     string

    Children *sync.Map
}

func BlockV2(n *Node, fqdn string, path []string) {
    cur := n

    for i, dc := range path {

        isLast := i == len(path)-1

        if isLast {
            newNode := &Node{
                Value:    dc,
                //Children: &(map[string]*Node{}),
                Children: &sync.Map{},
                Blocked:  true,
                FQDN:     fqdn,
            }
            (*cur.Children).Store(dc, newNode)
            //(*cur.Children)[dc] = newNode
            return
        }

        if v, ok := (*cur.Children).Load(dc); ok {
        //if v, ok := (*cur.Children)[dc]; ok {
            cur = v.(*Node)

            // This node is already blocked and a complete rule already exists
            // for this path. #A
            if cur.Blocked {
                return
            }
        } else {
            newNode := &Node{
                Value:    dc,
                //Children: &(map[string]*Node{}),
                Children: &sync.Map{},
            }

            //(*cur.Children)[dc] = newNode
            (*cur.Children).Store(dc, newNode)
            cur = newNode
        }
    }
}

func BlockV3(n *Node, fqdn string, path []string) {
    cur := n

    for _, dc := range path {
        //if v, ok := (*cur.Children)[dc]; ok {
        if v, ok := (*cur.Children).Load(dc); ok {
            cur = v.(*Node)

            // This node is already blocked and a complete rule already exists
            // for this path. #A
            if cur.Blocked {
                return
            }
        } else {
            newNode := &Node{
                Value:    dc,
                //Children: &(map[string]*Node{}),
                Children: &sync.Map{},
            }
            //(*cur.Children)[dc] = newNode
            (*cur.Children).Store(dc, newNode)
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
        if v, ok := (*cur.Children).Load(dc); ok {
        //if v, ok := (*cur.Children)[dc]; ok {
            cur = v.(*Node)

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
                newNode.Children = &sync.Map{}
            }

            //(*cur.Children)[dc] = newNode
            (*cur.Children).Store(dc, newNode)
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
        if v, ok := (*cur.Children).Load(dc); ok {
        //if v, ok := (*cur.Children)[dc]; ok {
            // This node is already blocked and a complete rule already exists
            // for this path. #A
            if v.(*Node).Blocked {
                return
            }

            cur = v.(*Node)
        } else {
            newNode := &Node{
                Value: dc,
            }

            // This should conserve memory; don't create a bunch of dangling
            // maps
            if i != l {
                //newNode.Children = &(map[string]*Node{})
                newNode.Children= &sync.Map{}
            }

            //(*cur.Children)[dc] = newNode
            (*cur.Children).Store(dc, newNode)
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
        if v, ok := (*cur.Children).Load(dc); ok {
        //if v, ok := (*cur.Children)[dc]; ok {
            // This node is already blocked and a complete rule already exists
            // for this path. #A
            if v.(*Node).Blocked {
                return
            }

            cur = v.(*Node)
        } else {
            newNode := &Node{
                Value: dc,
            }

            // This should conserve memory; don't create a bunch of dangling
            // maps
            if i != l {
                newNode.Children = &sync.Map{}
            }

            //(*cur.Children)[dc] = newNode
            (*cur.Children).Store(dc, newNode)
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

    cur.Children.Range(func(key interface{}, value interface{}) bool {
        Read(value.(*Node), result)
        return true
    })

    //for _, v := range (*cur.Children) {
    //    Read(v, result)
    //}
}
