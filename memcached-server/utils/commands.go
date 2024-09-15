package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Command struct {
	// "set", "get", or "delete"
	Name string

	Key string

	// If it is zero, the item never expires. If it's non-zero it is the number of seconds into the future in which the data expires
	ExpiresIn int

	// an arbitrary 16-bit unsigned integer (written out in decimal) that the server stores along with the data and sends back when the item is retrieved
	Flags uint16

	// optional parameter that instructs the server to not send the reply
	Noreply bool

	// the number of bytes is the number of bytes in the data block to follow, not including the delimiting
	ByteCount int
}

func ParseCommand(rawCommand string) (*Command, error) {
	if strings.HasPrefix(rawCommand, "get") {
		return parseGetCommand(rawCommand)
	}

	re := regexp.MustCompile(
		fmt.Sprintf(
			// Handling the possibility of multiple spaces. Although, multiple spaces probably aren't allowed by the
			// actual Memcached implementation
			`%s[ ]+%s[ ]+%s[ ]+%s[ ]+%s%s`,
			buildNamedCaptureGroupRegexp(`\S+`, "name"),
			buildNamedCaptureGroupRegexp(`\S+`, "key"),
			buildNamedCaptureGroupRegexp(`\S+`, "flags"),
			buildNamedCaptureGroupRegexp(`\S+`, "exptime"),
			buildNamedCaptureGroupRegexp(`\S+`, "byteCount"),
			buildNamedCaptureGroupRegexp("([ ]+noreply)?", "noReply"),
		),
	)
	matches := re.FindStringSubmatch(rawCommand)
	namedGroups := map[string]string{}

	for i, match := range matches {
		if len(re.SubexpNames()[i]) > 0 {
			namedGroups[re.SubexpNames()[i]] = strings.TrimSpace(match)
		}
	}

	flags, convertErr := strconv.Atoi(namedGroups["flags"])

	if convertErr != nil {
		return nil, errors.New("error parsing command: `flags` must be a number")
	}

	expiresIn, convertErr := strconv.Atoi(namedGroups["exptime"])

	if convertErr != nil {
		return nil, errors.New("error parsing command: `exptime` must be a number")
	}

	byteCount, convertErr := strconv.Atoi(namedGroups["byteCount"])
	if convertErr != nil {
		return nil, errors.New("error parsing command: `byteCount` must be a number")
	}

	noReply := false

	if namedGroups["noReply"] == "noreply" {
		noReply = true
	}

	command := &Command{
		Name:      namedGroups["name"],
		Flags:     uint16(flags),
		Key:       namedGroups["key"],
		ExpiresIn: expiresIn,
		Noreply:   noReply,
		ByteCount: byteCount,
	}

	return command, nil
}

func parseGetCommand(rawCommand string) (*Command, error) {
	split := strings.Split(rawCommand, " ")

	if len(split) != 2 {
		return nil, fmt.Errorf("unexpected command structure for: '%s'", rawCommand)
	}

	return &Command{
		Name:  split[0],
		Key:   split[1],
		Flags: uint16(0),
	}, nil
}

// buildNamedCaptureGroupRegexp wraps expression as a named capture group. E.g., "[a-z]+", "firstName" -> "(?P<firstName>[a-z]+)"
func buildNamedCaptureGroupRegexp(expression string, name string) string {
	return fmt.Sprintf("(?P<%s>%s)", name, expression)
}
