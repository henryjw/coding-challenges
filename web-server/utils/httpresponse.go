package utils

import (
	"fmt"
)

const HttpVersion = "1.1"

var httpStatusText = map[int]string{
	200: "OK",
	201: "No Content",
	400: "Bad Request",
	404: "Not Found",
	500: "Internal Server Error",
}

type HttpResponse struct {
	StatusCode int
	Body       string
}

func (response HttpResponse) Format() (string, error) {
	statusText, ok := httpStatusText[response.StatusCode]

	if !ok {
		return "", fmt.Errorf("unexpected error code: %d", response.StatusCode)
	}

	return fmt.Sprintf("HTTP/%s %d %s\r\n\r\n%s\r\n", HttpVersion, response.StatusCode, statusText, response.Body), nil
}
