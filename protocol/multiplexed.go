package protocol

import "fmt"

type MultiplexedProtocol struct {
	Protocol
	service   string
	delimeter string
}

func NewMultiplexedProtocol(p Protocol, service string) *MultiplexedProtocol {
	return &MultiplexedProtocol{
		Protocol:  p,
		service:   service,
		delimeter: ":",
	}
}

func (p MultiplexedProtocol) WriteMessageBegin(
	name string, mtype byte, seqid int32,
) error {
	name = fmt.Sprintf("%s%s%s", p.service, p.delimeter, name)
	return p.Protocol.WriteMessageBegin(name, mtype, seqid)
}
