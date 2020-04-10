package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	. "github.com/ajruckman/xlib"
)

func main() {
	fmt.Print("Attempting to connect to server... ")

	conn, err := net.DialTimeout("tcp", "10.3.0.16:64417", time.Second*3)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			fmt.Println("failed.")
			return
		} else {
			Err(err)
		}
	}
	fmt.Println("connected.")

	_, err = conn.Write([]byte("gen_oui\n"))
	Err(err)

	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Server stopped.")
				conn.Close()
				break
			} else {
				conn.Close()
				Err(err)
			}
		}

		fmt.Print("+>", data)
	}
}
