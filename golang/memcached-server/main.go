package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"memcached-server/utils"
	"net"
	"os"
	"strings"
)

func main() {
	app := &cli.App{
		Name: "Simple Memcached Server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "p",
				Value: 9999,
				Usage: "Port number to run the server",
			},
		},
		Action: func(context *cli.Context) error {
			return run(context.Int("p"))
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Error running app: %v\n", err)
	}
}

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

// handleConnections this is a blocking call
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
