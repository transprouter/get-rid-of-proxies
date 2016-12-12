package proxy

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/jeremiehuchet/go-through-proxies/xnet"
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
	port int
}

func (p DirectProxy) Forward(conn *xnet.Connection) {
	defer conn.Close()
	remoteConn, err := net.Dial("tcp", conn.Destination().String())
	defer remoteConn.Close()
	if err != nil {
		fmt.Printf("ERROR opening connection with %s\n", conn.Destination())
		return
	}
	pipe(conn, remoteConn)
}

func (p HTTPProxy) Forward(conn *xnet.Connection) {
	proxyConn, err := net.Dial("tcp", conn.Destination().String())
	if err != nil {
		fmt.Printf("ERROR opening connection with proxy at %s:%d\n", p.host, p.port)
	}
	// http ? or connect ?
	proxyConn.Write(nil)
	pipe(conn, proxyConn)
}

func pipe(local *xnet.Connection, remote net.Conn) {
	defer local.Close()
	defer remote.Close()

	var wg sync.WaitGroup
	wg.Add(2)
	// pipe downstream
	go cp(local.WriteCloser(), remote, &wg)
	// pipe upstream
	closeableRemoteDest := xnet.NewTCPWriteCloser(remote.(*net.TCPConn))
	go cp(closeableRemoteDest, local.Reader(), &wg)
	wg.Wait()
}

func cp(dst io.WriteCloser, src io.Reader, wg *sync.WaitGroup) (written int64, err error) {
	defer dst.Close()
	defer wg.Done()
	written, err = io.Copy(dst, src)
	if err != nil {
		fmt.Printf("Error forwarding data: %s\n", err)
	}
	return
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
