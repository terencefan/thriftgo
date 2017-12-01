package rpc

import "github.com/stdrickforce/thriftgo/protocol"

type Processor interface {
	Process(iprot, oprot protocol.Protocol) error
}
