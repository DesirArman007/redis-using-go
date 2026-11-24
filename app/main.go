package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/utils"
)

type Entry struct {
	Value     string
	ExpiresAt time.Time // zero value means no expiry
}

var (
	//Your map is defined as map[string]string, which strictly requires a string key and a string value. If you try to skip .(string):
	store = make(map[string]Entry)
)

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

	for {

		// reads the buffer
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		// takes the value out of the buffer
		// using the array decoder to decode the values
		// through the Custom resp parser in resp.go
		// decode() -> decoderesp() -> decodeArray()
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
		if !ok {
			continue
		}

		_ = buf[:n]

		// swtich-case start, here converting all commands
		//of switch to string in uppsercase to pass on to match the cases
		switch strings.ToUpper(cmd) {

		// returns a bulkstring of response format
		case "PING":
			conn.Write([]byte("+PONG\r\n"))

		// echo "string", it reads the string and then return the coded string
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

			key, ok1 := args[1].(string)
			value, ok2 := args[2].(string)
			if !ok1 || !ok2 {
				continue
			}

			var expiresAt time.Time

			if len(args) >= 5 {
				opt, ok3 := args[3].(string)
				ttlStr, ok4 := args[4].(string)
				if ok3 && ok4 {
					ttl, err := strconv.Atoi(ttlStr)
					if err == nil {
						switch strings.ToUpper(opt) {
						case "PX":
							expiresAt = time.Now().Add(time.Duration(ttl) * time.Millisecond)

						case "EX":
							expiresAt = time.Now().Add(time.Duration(ttl) * time.Second)
						}
					}

				}
			}
			store[key] = Entry{Value: value, ExpiresAt: expiresAt}

			conn.Write([]byte("+OK\r\n"))

		case "GET":
			if len(args) < 2 {
				conn.Write([]byte("-ERR wrong number of arguments\r\n"))
				continue
			}
			key, ok := args[1].(string)
			if !ok {
				continue
			}
			entry, exists := store[key]
			if !exists {
				conn.Write([]byte("$-1\r\n"))
				continue
			}

			if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
				delete(store, key)
				conn.Write([]byte("$-1\r\n"))
				continue
			}
			conn.Write([]byte(utils.EncodeBulkString(entry.Value)))
		}

	}
}
