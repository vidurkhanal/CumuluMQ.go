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

	// DEBUG
	// fmt.Println("Message: %#v", message)

	// Send the message to the server
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}

	fmt.Println("Message sent successfully")
}
