package extendedrpc

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Response struct {
	Message string
}

type Request struct {
	Name    string
	Message []byte
}

type Protocol struct{}

func (h *Protocol) Execute(req Request, res *Response) (err error) {
	res.Message = handleMessage(req.Name, string(req.Message))
	return
}

func handleMessage(name string, message string) string {
	assert.Equal(_tInstance, "test", name)
	return ""
}

var _tInstance *testing.T

type PingArgs struct{}
type PingReply struct {
	Message string
}

type wrapConn struct {
	net.Conn
	errChan chan error
}

func (w *wrapConn) Read(p []byte) (n int, err error) {
	n, err = w.Conn.Read(p)
	if err != nil {
		w.errChan <- err
	}
	return
}

func (w *wrapConn) Write(p []byte) (n int, err error) {
	n, err = w.Conn.Write(p)
	if err != nil {
		w.errChan <- err
	}
	return
}
