package rulegen

type Node struct {
    Value    string
    Children *map[string]*Node
    Blocked  bool
    FQDN     string
}

var (
    IncV1 int
    IncV2 int
    IncV3 int
)

func BlockV1(n *Node, fqdn string, path []string) {
    IncV1++
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
            (*cur.Children)[dc] = newNode
            return
        }

        if v, ok := (*cur.Children)[dc]; ok {
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
            (*cur.Children)[dc] = newNode
            cur = newNode
        }
    }
}

func BlockV2(n *Node, fqdn string, path []string) {
    IncV2++
    cur := n

    for i, dc := range path {
        isLast := i == len(path)-1

        if v, ok := (*cur.Children)[dc]; ok {
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
                Blocked:  isLast,
            }
            if isLast {
                newNode.FQDN = fqdn
            }
            (*cur.Children)[dc] = newNode
            cur = newNode
        }
    }
}

func BlockV3(n *Node, fqdn string, path []string) {
    IncV3++
    cur := n

    for _, dc := range path[:len(path)-2] {
        if v, ok := (*cur.Children)[dc]; ok {
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
            (*cur.Children)[dc] = newNode
            cur = newNode
        }
    }

    cur.Blocked = true
    cur.FQDN = fqdn
}

func Read(n *Node, result *[]string, depth int) {
    cur := n

    if cur.Blocked {
        *result = append(*result, cur.FQDN)
        //xlib.PrintIndent(depth, cur.FQDN)

        return // This should have no effect if #A works.
    }

    for _, v := range *cur.Children {
        Read(v, result, depth+1)
    }
}
