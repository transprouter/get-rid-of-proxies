package proxy

import (
	"bufio"
	"strings"
	"testing"
)

func TestInspectCanDetectHTTP(t *testing.T) {
	data := "GET /home HTTP/1.1"
	r := strings.NewReader(data)
	protocol, _ := inspectProtocol(r)
	if protocol != "HTTP" {
		t.Error("Expecting HTTP protocol")
	}
}

func TestInspectCanDetectUnknownWhenNoData(t *testing.T) {
	data := ""
	r := strings.NewReader(data)
	protocol, _ := inspectProtocol(r)
	if protocol != "unknown" {
		t.Error("Expecting unknown protocol")
	}
}

func TestInspectCanDetectUnknownWhenVeryFewData(t *testing.T) {
	data := "few"
	r := strings.NewReader(data)
	protocol, _ := inspectProtocol(r)
	if protocol != "unknown" {
		t.Error("Expecting unknown protocol")
	}
}

func TestReaderAfterInspect(t *testing.T) {
	r := strings.NewReader("GET /home HTTP/1.1")
	_, rr := inspectProtocol(r)

	s, _ := bufio.NewReader(rr).ReadString('\n')
	if s != "GET /home HTTP/1.1" {
		t.Errorf("expecting '%s' to be equal to 'GET /home HTTP/1.1'", s)
	}
}
