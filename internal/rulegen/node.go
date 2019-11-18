package rulegen

type Node struct {
    Value    string
    Children *map[string]*Node
    Blocked  bool
}

func (n *Node) Seek(path []string) (*Node, bool) {
    cur := n

    for i, dc := range path {
        isLast := i == len(path)-1

        if v, ok := (*cur.Children)[dc]; ok {
            cur = v

            // This node is already blocked and a complete rule already exists
            // for this path.
            if cur.Blocked && !isLast {
                return cur, true
            }
        } else {
            newNode := &Node{
                Value:    dc,
                Children: &(map[string]*Node{}),
                Blocked:  isLast,
            }
        }
    }
}
