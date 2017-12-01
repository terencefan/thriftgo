package rpc

import (
	"errors"
	"net/rpc"
	"strings"
	"sync"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
)

type serverCodec struct {
	iprot      Protocol
	oprot      Protocol
	nameCache  map[string]string // incoming name -> registered name
	methodName map[uint64]string // sequence ID -> method name
	mu         sync.Mutex
}

// NewServerCodec returns a new rpc.ServerCodec using Thrift RPC on conn using the specified protocol.
func NewServerCodec(iprot, oprot Protocol) rpc.ServerCodec {
	return &serverCodec{
		iprot:      iprot,
		oprot:      oprot,
		nameCache:  make(map[string]string, 8),
		methodName: make(map[uint64]string, 8),
	}
}

func (c *serverCodec) ReadRequestHeader(request *rpc.Request) error {
	name, messageType, seq, err := c.iprot.ReadMessageBegin()
	if err != nil {
		return err
	}
	if messageType != T_CALL { // Currently don't support one way
		return errors.New("thrift: expected Call message type")
	}

	// TODO: should use a limited size cache for the nameCache to avoid a possible
	//       memory overflow from nefarious or broken clients
	newName := c.nameCache[name]
	if newName == "" {
		newName = CamelCase(name)
		if !strings.ContainsRune(newName, '.') {
			newName = "Thrift." + newName
		}
		c.nameCache[name] = newName
	}

	c.mu.Lock()
	c.methodName[uint64(seq)] = name
	c.mu.Unlock()

	request.ServiceMethod = newName
	request.Seq = uint64(seq)

	return nil
}

func (c *serverCodec) ReadRequestBody(thriftStruct interface{}) error {
	if thriftStruct == nil {
		if err := c.iprot.Skip(T_STRUCT); err != nil {
			return err
		}
	} else {
		if err := DecodeStruct(c.iprot, thriftStruct); err != nil {
			return err
		}
	}
	return c.iprot.ReadMessageEnd()
}

func (c *serverCodec) WriteResponse(response *rpc.Response, thriftStruct interface{}) error {
	c.mu.Lock()
	methodName := c.methodName[response.Seq]
	delete(c.methodName, response.Seq)
	c.mu.Unlock()
	response.ServiceMethod = methodName

	mtype := byte(T_REPLY)
	if response.Error != "" {
		mtype = T_EXCEPTION
		etype := int32(ExceptionInternalError)
		if strings.HasPrefix(response.Error, "rpc: can't find") {
			etype = ExceptionUnknownMethod
		}
		thriftStruct = &TApplicationException{response.Error, etype}
	}
	if err := c.oprot.WriteMessageBegin(response.ServiceMethod, mtype, int32(response.Seq)); err != nil {
		return err
	}
	if err := EncodeStruct(c.oprot, thriftStruct); err != nil {
		return err
	}
	if err := c.oprot.WriteMessageEnd(); err != nil {
		return err
	}
	return c.oprot.Flush()
}

func (c *serverCodec) Close() error {
	// TODO why close connection ?
	c.iprot.Close()
	c.oprot.Close()
	return nil
}
