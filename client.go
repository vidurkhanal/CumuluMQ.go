package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		return
	}
	defer conn.Close()

	// Construct the message to send
	// messageType := byte(0x02) // Publish message
	// topic := []byte("topic")
	messageBody := []byte("{\"foo\": \"bar\"}")

	// Calculate the lengths
	// topicLength := len(topic)
	messageBodyLength := len(messageBody)
	totalLength := 4 + messageBodyLength

	// Create the message buffer
	message := make([]byte, totalLength)

	// Set the message fields
	binary.LittleEndian.PutUint32(message[0:4], uint32(totalLength))
	// message[4] = messageType
	// binary.LittleEndian.PutUint32(message[5:9], uint32(topicLength))
	// copy(message[9:9+topicLength], topic)
	// binary.LittleEndian.PutUint32(message[9+topicLength:9+topicLength+4],
	// uint32(messageBodyLength))
	copy(message[4:4+messageBodyLength],
		messageBody)

	// Send the message to the server
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}

	fmt.Println("Message sent successfully")
}
