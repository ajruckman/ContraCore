package log

//// Initializes the ContraWeb query log event server.
//func listen() {
//	netmgr.loadCache()
//
//	ln, err := net.Listen("tcp", "0.0.0.0:64417")
//	Err(err)
//
//	for {
//		conn, err := ln.Accept()
//		Err(err)
//
//		system.Console.Info("New client: ", conn.RemoteAddr().String())
//
//		queryBufferLock.Lock()
//		buffer := make([]schema.Log, len(queryBuffer))
//		copy(buffer, queryBuffer)
//		queryBufferLock.Unlock()
//
//		go netmgr.Onboard(conn, netmgr.cache)
//	}
//}
