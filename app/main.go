package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment the code below to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error appecpting connections:", err)
			break
		}
		
		go handleConn(conn) // one goroutine per client
	}
}

func handleConn(conn net.Conn){
	defer conn.Close()

	buf := make([]byte, 1024)
	for(
		n, err := conn.Read(buf)
		if err != nil{
			return
		}

		_ =buf[:n]
		conn.Write([]byte("+PONG\r\n"))
	)
}
