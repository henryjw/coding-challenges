package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"memcached-server/cache"
	"memcached-server/utils"
	"net"
	"strings"
)

type Server struct {
	cache *cache.Cache
}

func New(cache *cache.Cache) *Server {
	// Ideally, the type for the cache be an interface instead of a concrete type to allow flexibility of using different implementations.
	//However, this is fine for the purpose of this project
	return &Server{
		cache: cache,
	}
}

// Run Runs server. This is a blocking call and will not return until the server is stopped
func (receiver *Server) Run(portNumber int) error {
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

	receiver.handleConnections(listener)

	return nil
}

// handleConnections Handles incoming connections
func (receiver *Server) handleConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}

		fmt.Println("Accepted new connection")

		go receiver.handleConnection(conn)
	}
}

func (receiver *Server) handleConnection(conn net.Conn) {
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

			if errors.Is(readErr, io.EOF) {
				break
			} else {
				log.Println("Error reading from connection: ", readErr)
			}
			continue
		}

		log.Printf("Message received: '%s'\n", message)

		command, parseCommandErr := utils.ParseCommand(message)

		if parseCommandErr != nil {
			sendMessage(fmt.Sprint("Unexpected error parsing the command: ", parseCommandErr), conn)
			continue
		}
		// TODO: get data from connection if command name isn't "get"

		var data string

		if command.Name != "get" {
			var dataFetchErr error
			data, dataFetchErr = reader.ReadString('\n')

			if dataFetchErr != nil {
				log.Println("Error reading data: ", dataFetchErr)
			}
		}

		result, processCommandErr := receiver.processCommand(*command, data)

		if processCommandErr != nil {
			sendMessage(fmt.Sprint("Error processing command: ", processCommandErr), conn)
			continue
		}

		_, writeErr := conn.Write([]byte(result))
		if writeErr != nil {
			log.Println("Error sending message: ", writeErr)
			continue
		}
	}
}

// processCommand returns status for command. If an error occurs, it returns the status as an empty string
func (receiver *Server) processCommand(command utils.Command, data string) (string, error) {
	switch command.Name {
	case "set":
		return receiver.processSet(command, data)
	case "get":
		return receiver.processGet(command)
	case "add":
	case "replace":
	case "append":
	case "prepend":
		return "", fmt.Errorf("command '%s' not yet implemented", command.Name)
	}

	return "", fmt.Errorf("unexpected command name '%s'", command.Name)
}

func (receiver *Server) processSet(command utils.Command, data string) (string, error) {
	// TODO: set key expiration once supported by the cache
	err := receiver.cache.Set(command.Key, data)
	if err != nil {
		return "", err
	}

	return "STORED", nil
}

func (receiver *Server) processGet(command utils.Command) (string, error) {
	// TODO: update to allow return of data, flags, and bytecount
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
