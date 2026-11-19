package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connections:", err)
			break
		}

		// passes the connection string to handleconn goroutine
		go handleConn(conn) // one goroutine per client
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		val, err := utils.Decode(buf[:n])
		if err != nil {
			return
		}
		args, ok := val.([]interface{})
		if !ok || len(args) == 0 {
			return
		}

		cmd, ok := args[0].(string)

		_ = buf[:n]
		switch strings.ToUpper(cmd) {
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		case "ECHO":
			if len(args) >= 2 {
				value, ok := args[1].(string)
				if ok {
					resp := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
					conn.Write([]byte(resp))
				}
			}
		}

	}
}
