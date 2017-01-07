package xnet

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"syscall"
)

// Dest holds connection destination information
type Dest struct {
	Host string
	Port uint16
}

func (d Dest) String() string {
	return fmt.Sprintf("%s:%d", d.Host, d.Port)
}

type TCPReadCloser interface {
	io.Reader
	CloseRead() error
}

type TCPWriteCloser interface {
	io.Writer
	CloseWrite() error
}

// Connection holds addional information about the output connection request.
type Connection struct {
	io.ReadWriteCloser
	TCPReadCloser
	TCPWriteCloser
	originalConn *net.TCPConn
	Protocol     string
	Dest         Dest
	reader       io.Reader
	writer       io.Writer
}

// FIXME how to mock without shipping this test code?
func MockConnection(host string, port uint16, r io.Reader, w io.Writer) *Connection {
	c := new(Connection)
	c.originalConn = new(net.TCPConn)
	c.Dest = Dest{
		Host: host,
		Port: port,
	}
	c.Protocol, c.reader = inspectProtocol(r)
	c.writer = w
	return c
}

// Inspect the given net.TCPConn.
// returns a Connection holding informations about the request.
func Inspect(conn *net.TCPConn) *Connection {
	c := new(Connection)
	c.originalConn = conn
	c.Dest = *new(Dest)
	c.Dest.Host, c.Dest.Port = originalDestination(conn)
	c.Protocol, c.reader = inspectProtocol(conn)
	c.writer = conn
	return c
}

func (c Connection) String() string {
	return fmt.Sprintf("%s://%s", c.Protocol, c.Dest.String())
}

func (c Connection) Read(p []byte) (n int, err error) {
	return c.reader.Read(p)
}

func (c Connection) Write(p []byte) (n int, err error) {
	return c.writer.Write(p)
}

func (c Connection) Close() error {
	return c.originalConn.Close()
}

func (c Connection) CloseRead() error {
	return c.originalConn.CloseRead()
}

func (c Connection) CloseWrite() error {
	return c.originalConn.CloseWrite()
}

const soOriginalDest = 80

func originalDestination(conn *net.TCPConn) (host string, port uint16) {
	connFile, err := conn.File()
	//defer connFile.Close()
	if err != nil {
		fmt.Printf("Unable to obtain unterdying os.File: %s\n", err)
		return
	}
	addr, err := syscall.GetsockoptIPv6Mreq(int(connFile.Fd()), syscall.IPPROTO_IP, soOriginalDest)
	if err != nil {
		fmt.Printf("Unable to obtain original destination: %s\n", err)
		return
	}
	host = fmt.Sprintf("%d.%d.%d.%d", uint(addr.Multiaddr[4]), uint(addr.Multiaddr[5]), uint(addr.Multiaddr[6]), uint(addr.Multiaddr[7]))
	port = uint16(addr.Multiaddr[2])<<8 + uint16(addr.Multiaddr[3])
	return host, port
}

func inspectProtocol(r io.Reader) (string, io.Reader) {
	consumed := new(bytes.Buffer)
	safe := io.TeeReader(r, consumed)
	http := isHTTP(safe)
	if http {
		return "http", io.MultiReader(consumed, r)
	}
	return "unknown", io.MultiReader(consumed, r)
}

var httpMethodRegexp = regexp.MustCompile(`^(GET|HEAD|POST|OPTIONS|CONNECT|TRACE|PUT|PATCH|DELETE) `)
var httpRegexp = regexp.MustCompile(`^(GET|HEAD|POST|OPTIONS|CONNECT|TRACE|PUT|PATCH|DELETE) .+ HTTP\/(0\.9|1\.0|1\.1)$`)

func isHTTP(r io.Reader) bool {
	buf := make([]byte, 10)
	n, err := io.ReadFull(r, buf)
	if err == io.EOF {
		// no data, it can't be HTTP
		return false
	} else if err == io.ErrUnexpectedEOF {
		// not enougth data, it can't be HTTP
		return false
	} else if err != nil {
		log.Fatalf("Can't verify if protocol matches HTTP: %s", err)
	}
	s := string(buf[:n])
	if !httpMethodRegexp.MatchString(s) {
		// it doesn't match the begining of an HTTP request
		return false
	}
	return true
}
