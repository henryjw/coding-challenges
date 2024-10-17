package utils

import "testing"

func TestFormatResponse(t *testing.T) {
	response := HttpResponse{
		StatusCode: 404,
		Body:       "Hello",
	}

	formatted, err := response.Format()

	if err != nil {
		t.Fatal(err)
	}

	expected := "HTTP/1.1 404 Not Found\r\n\r\nHello\r\n"

	if formatted != expected {
		t.Errorf("unexpected response format. Expected '%s', got '%s'", expected, formatted)
	}
}

func TestFormatResponse_UnsupportedStatusCode(t *testing.T) {
	response := HttpResponse{
		StatusCode: 800,
	}

	_, err := HttpResponse.Format(response)

	if err == nil {
		t.Fatal("Expected error")
	}

	expectedErrorMessage := "unexpected error code: 800"

	if err.Error() != expectedErrorMessage {
		t.Errorf("expected error message '%s', got '%s'", expectedErrorMessage, err)
	}
}
