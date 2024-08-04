package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"memcached-server/utils"
	"net"
	"strings"
)

// run Runs server. This is a blocking call and will not return until the server is stopped
func run(portNumber int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))

	if err != nil {
		return fmt.Errorf("error starting server: %v", err)
	}

	log.Printf("Server listening on port %d\n", portNumber)

	defer func() {
		closeErr := listener.Close()
		if closeErr != nil {
			log.Println("Error closing listener", err)
			return
		}

		log.Println("Listener closed")
	}()

	handleConnections(listener)

	return nil
}

// handleConnections Handles incoming connections
func handleConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}

		fmt.Println("Accepted new connection")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			log.Println("Error closing connection: ", closeErr)
			return
		}

		log.Printf("Successfully closed connection")
	}()

	reader := bufio.NewReader(conn)

	for {
		message, readErr := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if readErr != nil {
			log.Println("Error reading from connection: ", readErr)
			continue
		}

		log.Printf("Message received: '%s'\n", message)

		command, parseCommandErr := utils.ParseCommand(message)

		if parseCommandErr != nil {
			sendMessage(fmt.Sprint("Unexpected error parsing the command: ", parseCommandErr), conn)
			continue
		}

		result, processCommandErr := processCommand(*command)

		if processCommandErr != nil {
			sendMessage(fmt.Sprint("Error processing command: ", processCommandErr), conn)
			continue
		}

		_, writeErr := conn.Write([]byte(fmt.Sprintf("Echo: %v\n", result)))
		if writeErr != nil {
			log.Println("Error sending message: ", writeErr)
			continue
		}
	}
}

func processCommand(command utils.Command) (string, error) {
	return "", errors.New("not yet implemented")
}

func sendMessage(message string, writer io.Writer) error {
	log.Printf("Sending mesessage: '%s'\n", message)
	_, err := writer.Write([]byte(message))

	if err != nil {
		log.Println("Error sending message: ", err)
	} else {
		log.Println("Successfully sent message")
	}

	return err
}
