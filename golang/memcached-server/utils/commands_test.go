package utils

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseCommand(t *testing.T) {
	rawCommand := "set test 0 100 4 noreply"
	command, err := ParseCommand(rawCommand)

	if err != nil {
		t.Fatal(err)
	}

	expected := &Command{
		Name:      "set",
		Key:       "test",
		Flags:     0,
		ExpiresIn: 100,
		ByteCount: 4,
		Noreply:   true,
	}

	assertSame(*expected, *command, t)
}

func TestParseCommandNoreplyNotSet(t *testing.T) {
	rawCommand := "set test 0 100 4"
	command, err := ParseCommand(rawCommand)

	if err != nil {
		t.Fatal(err)
	}

	expected := &Command{
		Name:      "set",
		Key:       "test",
		Flags:     0,
		ExpiresIn: 100,
		ByteCount: 4,
		Noreply:   false,
	}

	assertSame(*expected, *command, t)
}

func TestNonNumericFlags_Error(t *testing.T) {
	rawCommand := "set test x 100 4"
	_, err := ParseCommand(rawCommand)

	if err == nil {
		t.Fatal("Expected error")
	}

	expected := errors.New("error parsing command: `flags` must be a number")
	if err.Error() != expected.Error() {
		t.Fatalf("Unexpected error. Expected: '%s', got: '%s'\n", expected.Error(), err.Error())
	}
}

func TestNonNumericExpTime_Error(t *testing.T) {
	rawCommand := "set test 0 x 4"
	_, err := ParseCommand(rawCommand)

	if err == nil {
		t.Fatal("Expected error")
	}

	expected := errors.New("error parsing command: `exptime` must be a number")
	if err.Error() != expected.Error() {
		t.Fatalf("Unexpected error. Expected: '%s', got: '%s'\n", expected.Error(), err.Error())
	}
}

func TestNonNumericByteCount_Error(t *testing.T) {
	rawCommand := "set test 0 100 x"
	_, err := ParseCommand(rawCommand)

	if err == nil {
		t.Fatal("Expected error")
	}

	expected := errors.New("error parsing command: `byteCount` must be a number")
	if err.Error() != expected.Error() {
		t.Fatalf("Unexpected error. Expected: '%s', got: '%s'\n", expected.Error(), err.Error())
	}
}

func assertSame(expected Command, actual Command, t *testing.T) {
	// In a production test, it would be more useful to output the fields that aren't equal to simplify troubleshooting.
	// But it's not worth the extra effort for this learning project
	if !reflect.DeepEqual(expected, actual) {
		t.Error("The commands are not the same")
	}
}
