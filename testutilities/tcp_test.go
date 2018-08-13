package testutilities

import (
	"bufio"
	"testing"
)

func CanOpenTCPConnection(t *testing.T) {
	conn, close := TCPConnection(t, "Hello world!")
	defer close()
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Error("Can't read connection", err)
	}
	if response != "Hello world!" {
		t.Error("Expecting 'Hello world!' reponse")
	}
}
