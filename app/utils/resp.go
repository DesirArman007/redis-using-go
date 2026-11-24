package utils

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

func Decode(data []byte) (interface{}, error) {
	// RESP requires at least one byte, because the first byte identifies the type.
	if len(data) == 0 {
		return nil, errors.New("no data")
	}

	// Decode the first RESP value and ignore the consumed-byte count.
	value, _, err := DecodeRESP(data)
	return value, err
}

func readLine(data []byte) (line []byte, nextIndex int, err error) {
	// Read a RESP header line, for example "4\r\n" or "PONG\r\n".
	// Each line ends with CRLF ("\r\n"). readLine returns:
	//   - line: the bytes before "\r\n", such as "4" or "PONG"
	//   - nextIndex: the number of bytes consumed, including "\r\n"
	// decodeBulkString uses line as the length of the string.
	// decodeArray uses line as the number of elements in the array.
	// nextIndex tells both decoders where the remaining RESP value starts.
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, 0, errors.New("line delimiter not found")
	}

	return data[:idx], idx + 2, nil
}

func DecodeRESP(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	// The first byte determines the RESP data type.
	switch data[0] {
	case '*':
		return decodeArray(data)
	case '$':
		return decodeBulkString(data)
	default:
		return nil, 0, fmt.Errorf("unknown RESP type: %c", data[0])
	}

}

func decodeBulkString(data []byte) (interface{}, int, error) {
	// Skip '$' and read the declared bulk-string length.
	line, nextIndex, err := readLine(data[1:])
	if err != nil {
		return nil, 0, err
	}

	length, err := strconv.Atoi(string(line))
	if length == -1 {
		// Null bulk strings in RESP are represented as $-1.
		return nil, 1 + nextIndex, nil
	}
	if err != nil {
		return nil, 0, err
	}

	// The bulk string starts after "$<length>\r\n".
	start := 1 + nextIndex
	end := start + length

	if len(data) < end+2 {
		return nil, 0, errors.New("incomplete bulk string")
	}

	// Return the string value and the total number of bytes consumed.
	// This includes "$", the length header, "\r\n", the string bytes, and the
	// trailing "\r\n". decodeArray adds this byte count to its pos pointer
	// so decoding continues at the start of the next array element.
	return string(data[start:end]), end + 2, nil

}

func decodeArray(data []byte) (interface{}, int, error) {
	// Skip '*' and read the number of elements in the array.
	line, nextIndex, err := readLine(data[1:])
	if err != nil {
		return nil, 0, err
	}

	length, err := strconv.Atoi(string(line))
	if err != nil {
		return nil, 0, err
	}

	elements := make([]interface{}, length)
	pos := 1 + nextIndex

	for i := 0; i < length; i++ {
		if pos >= len(data) {
			return nil, 0, errors.New("incomplete array")
		}

		// Arrays may contain bulk strings or nested arrays, so decode each value recursively.
		val, bytesRead, err := DecodeRESP(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elements[i] = val
		pos += bytesRead
	}

	return elements, pos, nil
}

func EncodeBulkString(s string) (data []byte) {
	length := len(s)
	// Encode a string as RESP bulk: "$<byte-length>\r\n<value>\r\n".
	res := fmt.Sprintf("$%d\r\n%s\r\n", length, s)

	return []byte(res)
}
