package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
)

func split(data []string, size int) [][]string {

	var chunks [][]string

	for i := 0; i < len(data); i += size {

		end := i + size
		if end > len(data) {
			end = len(data)
		}

		chunks = append(chunks, data[i:end])
	}

	return chunks
}

func send(worker string, chunk []string, results chan map[string]int, wg *sync.WaitGroup) {

	defer wg.Done()

	conn, err := net.Dial("tcp", worker)
	if err != nil {
		fmt.Println("Connection error:", worker, err)
		return
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	encoder.Encode(chunk)

	var result map[string]int
	decoder.Decode(&result)

	results <- result
}

func main() {

	// 1. Read file
	file, err := os.Open("passwords.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var passwords []string

	for scanner.Scan() {
		passwords = append(passwords, scanner.Text())
	}

	// 2. Split into 3 chunks
	chunks := split(passwords, len(passwords)/3)

	// 3. Workers IPs (غيّرهم حسب أجهزةك)
	workers := []string{
		"192.168.1.10:9001",
		"192.168.1.11:9001",
		"192.168.1.12:9001",
	}

	results := make(chan map[string]int, 3)

	var wg sync.WaitGroup

	// 4. Send tasks
	for i, chunk := range chunks {

		wg.Add(1)

		go send(workers[i], chunk, results, &wg)
	}

	// 5. Close channel after finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// 6. Reduce phase
	final := make(map[string]int)

	for res := range results {
		for k, v := range res {
			final[k] += v
		}
	}

	// 7. Output
	fmt.Println("\nFINAL RESULT:")
	for k, v := range final {
		fmt.Println(k, "=>", v)
	}
}