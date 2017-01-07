package proxy

import (
	"testing"

	"strings"

	"bufio"

	"bytes"

	"github.com/jeremiehuchet/get-rid-of-proxies/xnet"
)

func tTestCanForwardDirectly(t *testing.T) {
	req := bufio.NewReader(strings.NewReader("GET / HTTP/1.1\nHost: perdu.com\n\n"))
	resp := bytes.NewBuffer(make([]byte, 10))
	conn := xnet.MockConnection("perdu.com", 80, req, resp)
	p := DirectProxy{}
	p.Forward(conn)
	if !strings.Contains(resp.String(), "<h1>Perdu sur l'Internet ?</h1><h2>Pas de panique, on va vous aider</h2>") {
		t.Errorf("Response should contain perdu.com content but it was:\n%s\n", resp.String())
	}
}

func tTestCanForwardHTTPToProxy(t *testing.T) {
	req := bufio.NewReader(strings.NewReader("GET http://perdu.com/ HTTP/1.1\nHost: perdu.com\n\n"))
	resp := bytes.NewBuffer(make([]byte, 10))
	conn := xnet.MockConnection("perdu.com", 80, req, resp)
	p := NewHTTPProxy("localhost", 3128)
	p.Forward(conn)
	if !strings.Contains(resp.String(), "<h1>Perdu sur l'Internet ?</h1><h2>Pas de panique, on va vous aider</h2>") {
		t.Errorf("Response should contain perdu.com content but it was:\n%s\n", resp.String())
	}
}
