package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"web-server/utils"
)

const wwwDir = "./www"

var FileNotFoundError = errors.New("file not found")

type Server struct {
}

func New() *Server {
	return &Server{}
}

func (s *Server) Start(portNumber int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		log.Println("Error creating listener", err)
		return err
	}

	defer func() {
		err := listener.Close()
		if err != nil {
			log.Println("Error closing listener", err)
		}
	}()

	log.Printf("Server running on port %d\n", portNumber)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		err := conn.Close()
		log.Println("Error closing connection", err)
	}()

	log.Println("Received connection")

	httpRequest := make([]byte, 1024)
	_, err := conn.Read(httpRequest)

	if err != nil {
		log.Println("Error reading httpRequest: ", err)
		return
	}

	log.Printf("%s\n", httpRequest)

	req, err := utils.ParseRequest(string(httpRequest))

	if err != nil {
		// In a real server, an appropriate error HTML page would be returned
		// instead of plain text
		sendHttpResponse(utils.HttpResponse{
			StatusCode: 500,
			Body:       "Unexpected error parsing request",
		}, conn)
		return
	}

	htmlContents, err := getHtml(req.Path)

	if err != nil {
		if errors.Is(err, FileNotFoundError) {
			sendHttpResponse(utils.HttpResponse{
				StatusCode: 404,
			}, conn)
		} else {
			sendHttpResponse(utils.HttpResponse{
				StatusCode: 500,
				Body:       "Unexpected error occurred",
			}, conn)
		}
		return
	}

	sendHttpResponse(utils.HttpResponse{
		StatusCode: 200,
		Body:       htmlContents,
	}, conn)
}

// sendHttpResponse sends response. If there's an error sending the response,
// then the error is logged and ignored
func sendHttpResponse(response utils.HttpResponse, conn net.Conn) {
	formattedResponse, err := response.Format()
	if err != nil {
		errorResponse, err := utils.HttpResponse{
			StatusCode: 500,
			Body:       "An unexpected error occurred",
		}.Format()

		if err != nil {
			log.Printf("Error formatting error response: %v\n", err)
		}

		_, err = conn.Write([]byte(errorResponse))

		if err != nil {
			log.Printf("Error sending error response: %v\n", err)
		}

		return
	}

	_, err = conn.Write([]byte(formattedResponse))

	if err != nil {
		log.Printf("Error sending error response: %v\n", err)
	}
}

func getHtml(requestPath string) (string, error) {
	if requestPath == "/" {
		requestPath = "/index.html"
	}

	filePath := path.Join(wwwDir, requestPath)
	fileData, err := os.ReadFile(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			return "", FileNotFoundError
		}
		return "", err
	}

	return string(fileData), nil
}
