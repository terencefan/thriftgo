package transport

import (
	"net"
)

const (
	sockTypeTcp  = "tcp"
	sockTypeUnix = "unix"
)

type TSocket struct {
	conn     net.Conn
	addr     string
	sockType string
}

type TSocketFactory struct {
	addr     string
	sockType string
}

func (self *TSocket) Read(message []byte) (int, error) {
	if self.conn == nil {
		return 0, NewTransportError("TSocket", "read on unopened transport")
	}
	return self.conn.Read(message)
}

func (self *TSocket) Write(message []byte) (int, error) {
	if self.conn == nil {
		return 0, NewTransportError("TSocket", "write on unopened transport")
	}
	return self.conn.Write(message)
}

func (self *TSocket) Open() (err error) {
	if self.conn != nil {
		return NewTransportError("TSocket", "transport has been already opened")
	}
	switch self.sockType {
	case sockTypeTcp:
		self.conn, err = net.Dial(sockTypeTcp, self.addr)
	case sockTypeUnix:
		self.conn, err = net.Dial(sockTypeUnix, self.addr)
	default:
		err = NewTransportError("TSocket", "invalid socket type "+self.sockType)
	}
	return err
}

func (self *TSocket) Close() error {
	return self.conn.Close()
}

func (self *TSocket) Flush() error {
	return nil
}

func (self *TSocketFactory) GetTransport() Transport {
	return &TSocket{
		addr:     self.addr,
		sockType: self.sockType,
	}
}

func NewTSocket(addr string) *TSocket {
	return &TSocket{
		addr:     addr,
		sockType: sockTypeTcp,
	}
}

func NewTUnixSocket(addr string) *TSocket {
	return &TSocket{
		addr:     addr,
		sockType: sockTypeUnix,
	}
}

func NewTSocketConn(conn net.Conn) *TSocket {
	return &TSocket{
		conn: conn,
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
