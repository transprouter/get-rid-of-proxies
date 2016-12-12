package xnet

import (
	"io"
	"net"
)

type TCPReadCloser struct {
	io.ReadCloser
	conn *net.TCPConn
}

func NewTCPReadCloser(conn *net.TCPConn) *TCPReadCloser {
	r := new(TCPReadCloser)
	r.conn = conn
	return r
}

func (r TCPReadCloser) Read(p []byte) (n int, err error) {
	return r.conn.Read(p)
}

func (r TCPReadCloser) Close() error {
	return r.conn.CloseRead()
}

type TCPWriteCloser struct {
	io.WriteCloser
	conn *net.TCPConn
}

func NewTCPWriteCloser(conn *net.TCPConn) *TCPWriteCloser {
	w := new(TCPWriteCloser)
	w.conn = conn
	return w
}

func (w TCPWriteCloser) Write(p []byte) (n int, err error) {
	return w.conn.Write(p)
}

func (w TCPWriteCloser) Close() error {
	return w.conn.CloseRead()
}
