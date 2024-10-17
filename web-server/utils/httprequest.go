package utils

import (
	"fmt"
)

type HttpRequest struct {
	HttpVersion string
	Path        string
}

func ParseRequest(rawRequest string) (HttpRequest, error) {
	var httpVersion string
	var requestPath string

	_, err := fmt.Sscanf(rawRequest, "GET %s HTTP/%s", &requestPath, &httpVersion)

	if err != nil {
		return HttpRequest{}, fmt.Errorf("unexpected error parsing request: %s", err)
	}

	return HttpRequest{HttpVersion: httpVersion, Path: requestPath}, nil
}
