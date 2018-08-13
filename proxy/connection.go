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
}

// Inspect the given net.TCPConn.
// returns a Connection holding informations about the request.
func Inspect(conn *net.TCPConn) *Connection {
	c := new(Connection)
	c.initOriginalDestination(conn)
	c.inspectProtocol()
	return c
}

func (c *Connection) String() string {
	return fmt.Sprintf("%s://%s", c.Protocol, c.Dest.String())
}

func (c *Connection) Read(p []byte) (n int, err error) {
	return c.reader.Read(p)
}

func (c *Connection) Write(p []byte) (n int, err error) {
	return c.originalConn.Write(p)
}

// Close the underlying resources.
func (c *Connection) Close() error {
	return c.originalConn.Close()
}

const soOriginalDest = 80

func (c *Connection) initOriginalDestination(conn *net.TCPConn) {
	defer conn.Close()
	connFile, err := conn.File()
	defer connFile.Close()
	if err != nil {
		log.Fatalf("Unable to obtain unterdying os.File: %s", err)
		return
	}
	addr, err := syscall.GetsockoptIPv6Mreq(int(connFile.Fd()), syscall.IPPROTO_IP, soOriginalDest)
	if err != nil {
		log.Fatalf("Unable to obtain original destination: %s", err)
		return
	}
	newFileConn, err := net.FileConn(connFile)
	if err != nil {
		log.Fatalf("Unable to obtain new connection to os.File: %s", err)
		return
	}
	c.originalConn = newFileConn.(*net.TCPConn)
	c.Dest = Dest{
		Host: fmt.Sprintf("%d.%d.%d.%d", uint(addr.Multiaddr[4]), uint(addr.Multiaddr[5]), uint(addr.Multiaddr[6]), uint(addr.Multiaddr[7])),
		Port: uint16(addr.Multiaddr[2])<<8 + uint16(addr.Multiaddr[3]),
	}
}

func (c *Connection) inspectProtocol() {
	consumedBuffer := new(bytes.Buffer)
	backupedReader := io.TeeReader(c.originalConn, consumedBuffer)
	http := isHTTP(backupedReader)
	if http {
		c.Protocol = "HTTP"
	} else {
		c.Protocol = "unknown"
	}
	c.reader = io.MultiReader(consumedBuffer, c.originalConn)
}

var httpMethodRegexp = regexp.MustCompile(`^(GET|HEAD|POST|OPTIONS|CONNECT|TRACE|PUT|PATCH|DELETE) `)
var httpRegexp = regexp.MustCompile(`^(GET|HEAD|POST|OPTIONS|CONNECT|TRACE|PUT|PATCH|DELETE) .+ HTTP\/(0\.9|1\.0|1\.1)$`)

func isHTTP(r io.Reader) bool {
	buf := make([]byte, 10)
	n, err := io.ReadFull(r, buf)
	if err == io.EOF {
		// no data, it can't be HTTP
		return false
	}
	if err == io.ErrUnexpectedEOF {
		// not enougth data, it can't be HTTP
		return false
	}
	if err != nil {
		log.Fatalf("Can't verify if protocol matches HTTP: %s", err)
		return false
	}
	s := string(buf[:n])
	if !httpMethodRegexp.MatchString(s) {
		// it doesn't match the begining of an HTTP request
		return false
	}
	return true
}
