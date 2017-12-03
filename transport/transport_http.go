package transport

import (
	"bytes"
	"errors"
	"net"
	"net/http"
	"time"
)

type THttpTransport struct {
	addr    string
	buf     *bytes.Buffer
	timeout time.Duration
}

type THttpTransportFactory struct {
	addr string
}

func (g *THttpTransport) Read(message []byte) (int, error) {
	return g.buf.Read(message)
}

func (g *THttpTransport) Write(message []byte) (int, error) {
	return g.buf.Write(message)
}

func (g *THttpTransport) Open() error {
	return nil
}

func (g *THttpTransport) Close() error {
	return nil
}

func (g *THttpTransport) SetTimeout(d time.Duration) {
	g.timeout = d
}

func (g *THttpTransport) Flush() (err error) {
	client := &http.Client{
		Timeout: g.timeout,
	}
	resp, err := client.Post(g.addr, "application/thrift", g.buf)

	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		err = errors.New("[THttpTransport] time limit exceeded")
	}
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	if _, err := g.buf.ReadFrom(resp.Body); err != nil {
		return err
	}
	return nil
}

func (g *THttpTransportFactory) GetTransport() Transport {
	return NewTHttpTransport(g.addr)
}

func NewTHttpTransport(addr string) *THttpTransport {
	return &THttpTransport{
		addr: addr,
		buf:  bytes.NewBuffer([]byte{}),
	}
}

func NewTHttpTransportFactory(addr string) *THttpTransportFactory {
	return &THttpTransportFactory{
		addr: addr,
	}
}
