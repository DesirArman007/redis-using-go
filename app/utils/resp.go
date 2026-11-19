package utils

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

func Decode(data []byte) (interface{}, error) {

	if len(data) == 0 {
		return nil, errors.New("no data")
	}

	value, _, err := DecodeRESP(data)
	return value, err
}

func readLine(data []byte) (line []byte, nextIndex int, err error) {
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
	line, nextIndex, err := readLine(data[1:])
	if err != nil {
		return nil, 0, err
	}

	length, err := strconv.Atoi(string(line))
	if length == -1 {
		return nil, 1 + nextIndex, nil
	}
	if err != nil {
		return nil, 0, err
	}

	start := 1 + nextIndex
	end := start + length

	if len(data) < end+2 {
		return nil, 0, errors.New("incomplete bulk string")
	}
	return string(data[start:end]), end + 2, nil

}

func decodeArray(data []byte) (interface{}, int, error) {
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

		val, bytesRead, err := DecodeRESP(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elements[i] = val
		pos += bytesRead
	}

	return elements, pos, nil
}

func encodeBulkString(s string) (data []byte) {
	length := len(s)
	res := fmt.Sprintf("$%d\r\n%s\r\n", length, s)

	return []byte(res)
}
