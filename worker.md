package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

func mapper(data []string) map[string]int {
	result := make(map[string]int)

	for _, p := range data {
		result[p]++
	}

	return result
}

func handle(conn net.Conn) {
	defer conn.Close()

	var chunk []string

	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	err := decoder.Decode(&chunk)
	if err != nil {
		fmt.Println("Decode error:", err)
		return
	}

	result := mapper(chunk)

	encoder.Encode(result)
}

func main() {

	port := ":9001" // غيّرها على كل جهاز

	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Worker running on", port)

	for {
		conn, _ := ln.Accept()
		go handle(conn)
	}
}