package transport

import (
	"testing"
)

func TestNewTransport(t *testing.T) {
	var transport Transport
	var wrapper Transport

	transport = TNullTransport
	transport = NewTMemoryBuffer()
	transport = NewTSocket(":6000")
	transport = NewTHttpTransport(":6000", "/")

	wrapper = NewTFramedTransport(transport, false, false)
	wrapper = NewTBufferedTransport(transport)
	wrapper = NewTBufferedTransportSize(transport, 1024, 1024)

	_, _ = transport, wrapper
}
