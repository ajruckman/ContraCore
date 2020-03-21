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
	for {
		fmt.Print("Attempting to connect to server... ")

		conn, err := net.DialTimeout("tcp", "127.0.0.1:64417", time.Second*3)
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				fmt.Println("failed.")
				continue
			} else {
				Err(err)
			}
		}
		fmt.Println("connected.")

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

			fmt.Println("Event received:", data)
		}
	}
}
