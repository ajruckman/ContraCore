package netmgr

//// Sets up a new client by sending the client every log in the cache and adding
//// them to the client list.
//func Onboard(conn net.Conn, initial []schema.Log) {
//	clients[conn.RemoteAddr().String()] = conn
//
//	system.Console.Infof("Sending client %d rows", len(initial))
//
//	e := func(err error) {
//		if err != nil {
//			if _, ok := err.(*net.OpError); ok {
//				system.Console.Info("Deleting disconnected client:", conn.RemoteAddr().String())
//				delete(clients, conn.RemoteAddr().String())
//			} else {
//				Err(err)
//			}
//		}
//	}
//
//	_, err := conn.Write([]byte("initial\n"))
//	e(err)
//
//	for _, v := range initial {
//		content := marshal(v)
//
//		_, err := conn.Write(content)
//		e(err)
//	}
//
//	_, err = conn.Write([]byte("!initial\n"))
//	e(err)
//}

//// Marshals a log to a JSON byte slice.
//func marshal(log schema.Log) []byte {
//	content, err := json.Marshal(log)
//	Err(err)
//
//	return append(content, '\n')
//}

