package protocol

import (
	"testing"

	. "github.com/stdrickforce/thriftgo/thrift"
	. "github.com/stdrickforce/thriftgo/transport"
)

func TestNewProtocol(t *testing.T) {
	var proto Protocol
	var wrapper Protocol
	var trans = NewTMemoryBuffer()

	proto = TNullProtocol
	proto = NewTBinaryProtocol(trans, true, true)
	wrapper = NewStoredProtocol(proto, "ping", T_CALL, 1)
	wrapper = NewMultiplexedProtocol(proto, "ping")

	_, _ = proto, wrapper

}
