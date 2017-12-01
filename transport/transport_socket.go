package transport

import (
	"net"
	"time"
)

const (
	sockTypeTcp  = "tcp"
	sockTypeUnix = "unix"
)

type TSocket struct {
	conn     net.Conn
	addr     string
	sockType string
	timeout  time.Duration
}

type TSocketFactory struct {
	addr     string
	sockType string
}

func (t *TSocket) SetTimeout(d time.Duration) {
	t.timeout = d
}

func (t *TSocket) Read(message []byte) (int, error) {
	if t.timeout > 0 {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}
	if t.conn == nil {
		return 0, NewTransportError("TSocket", "read on unopened transport")
	}
	return t.conn.Read(message)
}

func (t *TSocket) Write(message []byte) (int, error) {
	if t.timeout > 0 {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}
	if t.conn == nil {
		return 0, NewTransportError("TSocket", "write on unopened transport")
	}
	return t.conn.Write(message)
}

func (t *TSocket) Open() (err error) {
	if t.conn != nil {
		return NewTransportError("TSocket", "transport has been already opened")
	}
	switch t.sockType {
	case sockTypeTcp:
		t.conn, err = net.Dial(sockTypeTcp, t.addr)
	case sockTypeUnix:
		t.conn, err = net.Dial(sockTypeUnix, t.addr)
	default:
		err = NewTransportError("TSocket", "invalid socket type "+t.sockType)
	}
	return err
}

func (t *TSocket) Close() error {
	return t.conn.Close()
}

func (t *TSocket) Flush() error {
	return nil
}

func (t *TSocketFactory) GetTransport() Transport {
	return &TSocket{
		addr:     t.addr,
		sockType: t.sockType,
	}
}

func NewTSocket(addr string) *TSocket {
	return &TSocket{
		addr:     addr,
		sockType: sockTypeTcp,
		timeout:  defaultTimeout,
	}
}

func NewTSocketTimeout(addr string, timeout time.Duration) *TSocket {
	return &TSocket{
		addr:     addr,
		sockType: sockTypeTcp,
		timeout:  timeout,
	}
}

func NewTUnixSocket(addr string) *TSocket {
	return &TSocket{
		addr:     addr,
		sockType: sockTypeUnix,
		timeout:  defaultTimeout,
	}
}

func NewTSocketConn(conn net.Conn) *TSocket {
	return &TSocket{
		conn:    conn,
		timeout: defaultTimeout,
	}
}

func NewTSocketFactory(addr string) *TSocketFactory {
	return &TSocketFactory{
		addr:     addr,
		sockType: sockTypeTcp,
	}
}

func NewTUnixSocketFactory(addr string) *TSocketFactory {
	return &TSocketFactory{
		addr:     addr,
		sockType: sockTypeUnix,
	}
}
