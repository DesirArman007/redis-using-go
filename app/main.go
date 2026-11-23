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

	var (
		//Your map is defined as map[string]string, which strictly requires a string key and a string value. If you try to skip .(string):
		store = make(map[string]string)
	)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		// takes the value out of the buffer
		val, err := utils.Decode(buf[:n])
		if err != nil {
			return
		}

		// utils.Decode returns interface{} because RESP can be any type (string, int, array, etc.).
		// A RESP array decodes to []interface{} (not []string) since its elements can be mixed types,
		// so we type-assert it to a slice first, then assert each element (like args[0]) to string.
		args, ok := val.([]interface{})
		if !ok || len(args) == 0 {
			return
		}

		// Extract the command name (first element) as a string
		cmd, ok := args[0].(string)

		_ = buf[:n]
		switch strings.ToUpper(cmd) {
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		case "ECHO":
			if len(args) >= 2 {
				value, ok := args[1].(string)
				if ok {
					resp := utils.EncodeBulkString(value)
					// resp := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
					conn.Write([]byte(resp))
				}
			}

		case "SET":
			store[args[1].(string)] = args[2].(string)
			conn.Write(utils.EncodeBulkString("OK"))

		case "GET":
			if args[1].(string) == store[args[1].(string)] {
				conn.Write(utils.EncodeBulkString(store[args[1].(string)]))
			}
		}

	}
}
