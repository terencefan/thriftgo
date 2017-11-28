package transport

import (
	"fmt"
	"io"
)

type Transport interface {
	io.ReadWriteCloser
	Open() error
	Flush() error
}

type TransportFactory interface {
	GetTransport() Transport
}

type TransportError struct {
	transport string
	message   string
}

func (e *TransportError) Error() string {
	return fmt.Sprintf("[%s] %s", e.transport, e.message)
}

func NewTransportError(transport, message string) *TransportError {
	return &TransportError{
		transport: transport,
		message:   message,
	}
}

type TransportWrapper interface {
	GetTransport(Transport) Transport
}

type transportWrapper struct {
}

func (self *transportWrapper) GetTransport(trans Transport) Transport {
	return trans
}

var TTransportWrapper = &transportWrapper{}
