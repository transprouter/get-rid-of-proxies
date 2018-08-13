package testutilities

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
	go func() {
		srv, err := l.Accept()
		if err != nil {
			t.Errorf("Can't accept connection: %s", err.Error())
		}
		go func() {
			fmt.Fprint(srv, response)
			srv.Close()
		}()
	}()
	client, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Errorf("Can't open connection to %s: %s", l.Addr().String(), err.Error())
	}
	return client.(*net.TCPConn), func() error {
		err = client.Close()
		if err != nil {
			return err
		}
		return l.Close()
	}
}
