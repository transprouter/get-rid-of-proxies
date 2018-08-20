package testutils

import (
	"fmt"
	"net"
	"testing"
)

func TCPConnection(t *testing.T, response string) (*net.TCPConn, func() error) {
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 0))
	if err != nil {
		t.Errorf("Can't listen on 'localhost:%s': %s", l.Addr().String(), err.Error())
	}
	go listenSingleConnection(t, l, response)
	client, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Errorf("Can't open connection to %s: %s", l.Addr().String(), err.Error())
	}
	return client.(*net.TCPConn), func() error {
		client.Close()
		return l.Close()
	}
}

func listenSingleConnection(t *testing.T, l net.Listener, response string) {
	srv, err := l.Accept()
	defer srv.Close()
	if err != nil {
		t.Errorf("Can't accept connection: %s", err.Error())
	}
	fmt.Fprint(srv, response)
}
