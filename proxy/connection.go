package proxy

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

// Connection holds addional information about the output connection request.
type Connection struct {
	io.ReadWriteCloser
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
	c.Dest = *new(Dest)
	c.originalConn, c.Dest = originalDestination(conn)
	c.Protocol, c.reader = inspectProtocol(c.originalConn)
	c.writer = c.originalConn
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

const soOriginalDest = 80

func originalDestination(conn *net.TCPConn) (newConn *net.TCPConn, dest Dest) {
	defer conn.Close()
	connFile, err := conn.File()
	defer connFile.Close()
	if err != nil {
		fmt.Printf("Unable to obtain unterdying os.File: %s\n", err)
		return
	}
	addr, err := syscall.GetsockoptIPv6Mreq(int(connFile.Fd()), syscall.IPPROTO_IP, soOriginalDest)
	if err != nil {
		fmt.Printf("Unable to obtain original destination: %s\n", err)
		return
	}
	newFileConn, err := net.FileConn(connFile)
	if err != nil {
		fmt.Printf("Unable to obtain new connection to os.File: %s\n", err)
		return
	}
	newConn = newFileConn.(*net.TCPConn)
	host := fmt.Sprintf("%d.%d.%d.%d", uint(addr.Multiaddr[4]), uint(addr.Multiaddr[5]), uint(addr.Multiaddr[6]), uint(addr.Multiaddr[7]))
	port := uint16(addr.Multiaddr[2])<<8 + uint16(addr.Multiaddr[3])
	dest = Dest{
		Host: host,
		Port: port,
	}
	return
}

func inspectProtocol(r io.Reader) (string, io.Reader) {
	consumed := new(bytes.Buffer)
	safe := io.TeeReader(r, consumed)
	http := isHTTP(safe)
	if http {
		return "HTTP", io.MultiReader(consumed, r)
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
