package utils

import "testing"

func TestParseRequest(t *testing.T) {
	request, err := ParseRequest("GET /index.html HTTP/1.1")

	if err != nil {
		t.Fatal(err)
	}

	if request.Path != "/index.html" {
		t.Errorf("invalid path. Expected /index.html, got %s", request.Path)
	}

	if request.HttpVersion != "1.1" {
		t.Errorf("invalid HTTP version. Expected 1.1, got %s\n", request.HttpVersion)
	}
}
