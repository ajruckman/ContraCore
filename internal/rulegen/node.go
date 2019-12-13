package rulegen

import (
    "sync"

    "github.com/tevino/abool"
)

type Node struct {
    Value   string
    Blocked *abool.AtomicBool
    FQDN    string

    Children sync.Map
}

//func BlockV2(n *Node, fqdn string, path []string) {
//   cur := n
//
//   for i, dc := range path {
//
//       isLast := i == len(path)-1
//
//       if isLast {
//           newNode := &Node{
//               Value:    dc,
//               Children: &(map[string]*Node{}),
//               Blocked:  true,
//               FQDN:     fqdn,
//           }
//           cur.lock.Lock()
//           (*cur.Children)[dc] = newNode
//           cur.lock.Unlock()
//           return
//       }
//
//       cur.lock.Lock()
//       v, ok := (*cur.Children)[dc]
//       cur.lock.Unlock()
//
//       if ok {
//           cur = v
//
//           // This node is already blocked and a complete rule already exists
//           // for this path. #A
//           if cur.Blocked {
//               return
//           }
//       } else {
//           newNode := &Node{
//               Value:    dc,
//               Children: &(map[string]*Node{}),
//           }
//
//           cur.lock.Lock()
//           (*cur.Children)[dc] = newNode
//           cur.lock.Unlock()
//
//           cur = newNode
//       }
//   }
//}
//
//func BlockV3(n *Node, fqdn string, path []string) {
//   cur := n
//
//   for _, dc := range path {
//       cur.lock.Lock()
//       v, ok := (*cur.Children)[dc]
//       cur.lock.Unlock()
//
//       if ok {
//           cur = v
//
//           // This node is already blocked and a complete rule already exists
//           // for this path. #A
//           if cur.Blocked {
//               return
//           }
//       } else {
//           newNode := &Node{
//               Value:    dc,
//               Children: &(map[string]*Node{}),
//           }
//           cur.lock.Lock()
//           (*cur.Children)[dc] = newNode
//           cur.lock.Unlock()
//
//           cur = newNode
//       }
//   }
//
//   cur.Blocked = true
//   cur.FQDN = fqdn
//}

func BlockV4(n *Node, fqdn string, path []string) {
    cur := n

    for _, dc := range path {
        v, ok := cur.Children.Load(dc)

        if ok {
            cur = v.(*Node)

            // This node is already blocked and a complete rule already exists
            // for this path. #A
            if cur.Blocked.IsSet() {
                return
            }
        } else {
            newNode := Node{
                Value:   dc,
                Blocked: abool.New(),
            }

            cur.Children.Store(dc, &newNode)
            cur = &newNode
        }
    }

    cur.Blocked.Set()
    cur.FQDN = fqdn
}

//func BlockV5(n *Node, fqdn string, path []string) {
//  cur := n
//
//  for _, dc := range path {
//      v, ok := cur.Children.Load(dc)
//
//      if ok {
//          // This node is already blocked and a complete rule already exists
//          // for this path. #A
//          if v.(*Node).Blocked {
//              return
//          }
//
//          cur = v.(*Node)
//      } else {
//          newNode := &Node{
//              Value: dc,
//          }
//
//          cur.Children.Store(dc, newNode)
//          cur = newNode
//      }
//  }
//
//  cur.Blocked = true
//  cur.FQDN = fqdn
//}

func Read(n *Node, result *[]string) {
    cur := n

    if cur.Blocked.IsSet() {
        *result = append(*result, cur.FQDN)

        return // This should have no effect if #A works.
    }

    (cur.Children).Range(func(key, value interface{}) bool {
        Read(value.(*Node), result)
        return true
    })

    //for _, v := range *cur.Children {
    //    Read(v, result)
    //}
}
