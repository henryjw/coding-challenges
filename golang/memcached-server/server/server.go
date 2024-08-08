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
	// Ideally, the type for the cache should be an interface instead of a concrete type to allow flexibility of using different implementations.
	// However, this is fine for the purpose of this project
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
			sendMessage(fmt.Sprint("Unexpected error parsing the command: ", parseCommandErr, "\r\n"), conn)
			continue
		}

		var data string

		if command.Name != "get" {
			var dataFetchErr error
			// NOTE: there's no validation to check that the data size matches value of `byte count` in the command
			data, dataFetchErr = reader.ReadString('\n')

			if dataFetchErr != nil {
				log.Println("Error reading data: ", dataFetchErr)
			}
		}

		result, processCommandErr := receiver.processCommand(*command, strings.TrimSpace(data))

		if processCommandErr != nil {
			sendMessage(fmt.Sprint("Error processing command: ", processCommandErr, "\r\n"), conn)
			continue
		}

		_, writeErr := conn.Write([]byte(result + "\r\n"))
		if writeErr != nil {
			log.Println("Error sending message: ", writeErr)
			continue
		}
	}
}

// processCommand returns status for command. If an error occurs, it returns the status as an empty string
func (receiver *Server) processCommand(command utils.Command, value string) (string, error) {
	switch command.Name {
	case "set":
		return receiver.processSet(command, value)
	case "get":
		return receiver.processGet(command)
	case "add":
		return receiver.processAdd(command, value)
	case "replace":
		return receiver.processReplace(command, value)
	case "append":
		return receiver.processAppend(command, value)
	case "prepend":
		return receiver.processPrepend(command, value)
	}

	return "", fmt.Errorf("unexpected command name '%s'", command.Name)
}

func (receiver *Server) processSet(command utils.Command, value string) (string, error) {
	// TODO: set key expiration once supported by the cache
	err := receiver.cache.Set(command.Key, cache.Data{
		Flags:     command.Flags,
		ByteCount: command.ByteCount,
		Value:     value,
	})
	if err != nil {
		return "", err
	}

	return "STORED", nil
}

func (receiver *Server) processGet(command utils.Command) (string, error) {
	data, err := receiver.cache.Get(command.Key)

	keyNotFoundError := &cache.KeyNotFoundError{}
	if errors.As(err, &keyNotFoundError) {
		return "END", nil
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("VALUE %s %d %d", data.Value, data.Flags, data.ByteCount), nil
}

func (receiver *Server) processAdd(command utils.Command, value string) (string, error) {
	err := receiver.cache.Add(command.Key, cache.Data{
		Value:     value,
		Flags:     command.Flags,
		ByteCount: command.ByteCount,
	})

	keyExistsError := &cache.KeyAlreadyExistsError{}

	if errors.As(err, &keyExistsError) {
		return "NOT_STORED", nil
	}

	if err != nil {
		return "", err
	}

	return "STORED", nil
}

func (receiver *Server) processReplace(command utils.Command, value string) (string, error) {
	err := receiver.cache.Replace(command.Key, cache.Data{
		Value:     value,
		Flags:     command.Flags,
		ByteCount: command.ByteCount,
	})

	keyNotFoundError := &cache.KeyNotFoundError{}

	if errors.As(err, &keyNotFoundError) {
		return "NOT_STORED", nil
	}

	if err != nil {
		return "", err
	}

	return "STORED", nil
}

func (receiver *Server) processAppend(command utils.Command, value string) (string, error) {
	err := receiver.cache.Append(command.Key, cache.Data{
		Value:     value,
		Flags:     command.Flags,
		ByteCount: command.ByteCount,
	})

	keyNotFoundError := &cache.KeyNotFoundError{}

	if errors.As(err, &keyNotFoundError) {
		return "NOT_STORED", nil
	}

	if err != nil {
		return "", err
	}

	return "STORED", nil
}

func (receiver *Server) processPrepend(command utils.Command, value string) (string, error) {
	err := receiver.cache.Prepend(command.Key, cache.Data{
		Value:     value,
		Flags:     command.Flags,
		ByteCount: command.ByteCount,
	})

	keyNotFoundError := &cache.KeyNotFoundError{}

	if errors.As(err, &keyNotFoundError) {
		return "NOT_STORED", nil
	}

	if err != nil {
		return "", err
	}

	return "STORED", nil
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
