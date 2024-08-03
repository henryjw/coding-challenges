package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
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
	writer := bufio.NewWriter(conn)
	//buffer := make([]byte, 0)

	for {
		message, readErr := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if readErr != nil {
			fmt.Println("Error reading from connection: ", readErr)
			continue
		}

		fmt.Printf("Message received: '%s'\n", message)

		// TODO: implement actual functionality. This just echoes the message back to the client
		_, writeErr := writer.WriteString(fmt.Sprintf("Echo: %v\n", message))
		if writeErr != nil {
			fmt.Println("Error sending message: ", writeErr)
			continue
		}

		writer.Flush()
	}
}
