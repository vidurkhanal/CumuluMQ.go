package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	maxRetries = 3
	retryDelay = time.Second
)

func connectWithRetry(id int) (net.Conn, error) {
	var conn net.Conn
	var err error

	// Retry logic for establishing connection
	for i := 0; i < maxRetries; i++ {
		conn, err = net.Dial("tcp", "0.0.0.0:8080")
		if err == nil {
			return conn, nil
		}
		fmt.Printf("Request %d: Failed to connect to server (attempt %d): %v\n", id, i+1, err)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("Request %d: Failed to connect to server after %d attempts: %v", id, maxRetries, err)
}

func sendRequest(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	conn, err := connectWithRetry(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Construct the message to send
	topic := []byte("topic")
	messageBody := []byte("{\"foo\": \"bar\"}")
	// Calculate the lengths
	topicLength := len(topic)
	messageBodyLength := len(messageBody)
	totalLength := 1 + 4 + topicLength + 4 + messageBodyLength
	// Create the message buffer
	message := make([]byte, 4+totalLength)
	// Set the message fields
	binary.LittleEndian.PutUint32(message[0:4], uint32(totalLength))
	message[4] = byte(0x01)
	binary.LittleEndian.PutUint32(message[5:9], uint32(topicLength))
	copy(message[9:9+topicLength], topic)
	binary.LittleEndian.PutUint32(message[9+topicLength:9+topicLength+4], uint32(messageBodyLength))
	copy(message[9+topicLength+4:], messageBody)

	// Send the message to the server
	_, err = conn.Write(message)
	if err != nil {
		fmt.Printf("Request %d: Failed to send message: %v\n", id, err)
		return
	}

	fmt.Printf("Request %d: Message sent successfully\n", id)

	// Read the response length
	respLenBuf := make([]byte, 4)
	_, err = conn.Read(respLenBuf)
	if err != nil {
		fmt.Printf("Request %d: Failed to read response length: %v\n", id, err)
		return
	}
	respLength := binary.LittleEndian.Uint32(respLenBuf)

	// Read the response data in chunks
	response := make([]byte, respLength)
	totalRead := 0
	for totalRead < int(respLength) {
		n, err := conn.Read(response[totalRead:])
		if err != nil {
			fmt.Printf("Request %d: Failed to read response: %v\n", id, err)
			return
		}
		totalRead += n
	}

	// Print the response
	fmt.Printf("Request %d: Response received: %s\n", id, string(response))
}

func main() {
	var wg sync.WaitGroup
	numRequests := 200

	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go sendRequest(&wg, i)
	}

	// Wait for all requests to complete
	wg.Wait()
	fmt.Println("All requests completed")
}
