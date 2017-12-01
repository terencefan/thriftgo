package transport

import (
	"fmt"
	"io"
	"time"
)

const defaultTimeout = time.Second * 30

type Transport interface {
	io.ReadWriteCloser
	Open() error
	Flush() error
	SetTimeout(time.Duration)
}

type TransportFactory interface {
	GetTransport() Transport
}

func NewTransportError(transport, message string) error {
	return fmt.Errorf("[%s] %s", transport, message)
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
