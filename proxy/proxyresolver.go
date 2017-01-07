package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/transprouter/transprouter/xnet"

	"bufio"

	"strings"
)

// Proxy describes how to contact proxy server
type Proxy interface {
	Forward(conn *xnet.Connection)
}

// DirectProxy directly forwards connections
type DirectProxy struct {
	Proxy
}

// HTTPProxy helps communicate with HTTP proxies
type HTTPProxy struct {
	Proxy
	host string
	port uint16
}

func (p DirectProxy) Forward(conn *xnet.Connection) {
	defer conn.Close()
	remoteConn, err := net.Dial("tcp", conn.Dest.String())
	defer remoteConn.Close()
	if err != nil {
		fmt.Printf("ERROR opening connection with %s\n", conn.Dest)
		return
	}
	pipe(conn, remoteConn.(*net.TCPConn))
}

func NewHTTPProxy(host string, port uint16) *HTTPProxy {
	p := new(HTTPProxy)
	p.host = host
	p.port = port
	return p
}

func (p HTTPProxy) Forward(conn *xnet.Connection) {
	proxyConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", p.host, p.port))
	if err != nil {
		fmt.Printf("ERROR opening connection with proxy at %s:%d: %s\n", p.host, p.port, err)
		return
	}
	fmt.Println(".")
	if conn.Protocol == "HTTP" {
		// forward raw request
		req, err := http.ReadRequest(bufio.NewReader(conn))
		if err != nil {
			fmt.Printf("ERROR Parsing HTTP request: %s\n", err)
		}
		req.URL.Scheme = "http"
		req.WriteProxy(proxyConn)
		pipe(conn, proxyConn.(*net.TCPConn))
		fmt.Println("Forwarded as HTTP")
	} else {
		// use CONNECT
		connect := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", conn.Dest, conn.Dest)
		fmt.Printf("forwarding using %s", connect)
		_, err := fmt.Fprintf(proxyConn, connect)
		if err != nil {
			fmt.Printf("ERROR sending CONNECT command to proxy at %s:%d: %s\n", p.host, p.port, err)
			return
		}
		status, err := bufio.NewReader(proxyConn).ReadString('\n')
		if err != nil {
			fmt.Printf("ERROR unreadable response from proxy at %s:%d: %s\n", p.host, p.port, err)
			return
		}
		if !strings.Contains(status, "200") {
			fmt.Printf("ERROR CONNECT command refused by proxy at %s:%d: %s\n", p.host, p.port, status)
			return
		}
		pipe(conn, proxyConn.(*net.TCPConn))
		fmt.Println("Forwarded using CONNECT")
	}
}

func pipe(local *xnet.Connection, remote *net.TCPConn) {
	defer local.Close()
	defer remote.Close()

	var wg sync.WaitGroup
	cp := func(dst xnet.TCPWriteCloser, src io.Reader) (written int64, err error) {
		defer wg.Done()
		defer dst.CloseWrite()
		written, err = io.Copy(dst, src)
		if err != nil {
			fmt.Printf("Error forwarding data: %s\n", err)
		}
		return
	}

	wg.Add(2)

	// pipe downstream
	go cp(local, remote)
	// pipe upstream
	go cp(remote, local)

	wg.Wait()
}

// Resolver help finding the proxy server to use to contact a given host
type Resolver interface {
	resolve(string) Proxy
}

// FixedResolver returns always the same proxy
type FixedResolver struct {
	Resolver
	proxy      Proxy
	exceptions []string
}

func (r FixedResolver) resolve(url string) Proxy {
	return r.proxy
}
