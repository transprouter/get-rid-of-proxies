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
	host string
	port uint16
}

func (d Dest) String() string {
	return fmt.Sprintf("%s:%d", d.host, d.port)
}

// Connection holds addional information about the output connection request.
type Connection struct {
	io.Closer
	originalConn *net.TCPConn
	protocol     string
	dest         Dest
	reader       io.Reader
	writer       io.WriteCloser
}

// Inspect the given net.TCPConn.
// returns a Connection holding informations about the request.
func Inspect(conn *net.TCPConn) *Connection {
	c := new(Connection)
	c.originalConn = conn
	c.dest = *new(Dest)
	c.dest.host, c.dest.port = originalDestination(conn)
	c.protocol, c.reader = inspectProtocol(conn)
	c.writer = NewTCPWriteCloser(conn)
	return c
}

func (c Connection) String() string {
	return fmt.Sprintf("%s://%s", c.protocol, c.dest.String())
}

// Destination returns the original destination of the connection.
func (c Connection) Destination() Dest {
	return c.dest
}

// Protocol return the connection protocol the data seems to refer to.
func (c Connection) Protocol() string {
	return c.protocol
}

// Reader returns an io.Reader to read the data the connection want to transmit.
func (c Connection) Reader() io.Reader {
	return c.reader
}

// WriteCloser returns an io.WriteCloser to write the response data.
func (c Connection) WriteCloser() io.WriteCloser {
	return c.writer
}

// Close the opened resources.
func (c Connection) Close() error {
	return c.originalConn.Close()
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
