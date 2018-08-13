package proxy

import (
	"bufio"
	"testing"

	"github.com/transprouter/transprouter/testutilities"
)

func TestInspectCanDetectHTTP(t *testing.T) {
	tcpConn, close := testutilities.TCPConnection(t, "GET /home HTTP/1.1")
	defer close()
	c := Inspect(tcpConn)
	if c.Protocol != "HTTP" {
		t.Error("Expecting HTTP protocol")
	}
}

func TestInspectCanDetectUnknownWhenNoData(t *testing.T) {
	tcpConn, close := testutilities.TCPConnection(t, "")
	defer close()
	c := Inspect(tcpConn)
	if c.Protocol != "unknown" {
		t.Error("Expecting unknown protocol")
	}
}

func TestInspectCanDetectUnknownWhenVeryFewData(t *testing.T) {
	tcpConn, close := testutilities.TCPConnection(t, "few")
	defer close()
	c := Inspect(tcpConn)
	if c.Protocol != "unknown" {
		t.Error("Expecting unknown protocol")
	}
}

func TestReaderAfterInspect(t *testing.T) {
	tcpConn, close := testutilities.TCPConnection(t, "GET /home HTTP/1.1")
	defer close()
	c := Inspect(tcpConn)
	s, _ := bufio.NewReader(c).ReadString('\n')
	if s != "GET /home HTTP/1.1" {
		t.Errorf("expecting '%s' to be equal to 'GET /home HTTP/1.1'", s)
	}
}
