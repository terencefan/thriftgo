package rpc

import (
	"fmt"
	"reflect"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
)

type SimpleProcessor struct {
	name    string
	rtype   reflect.Type
	rval    reflect.Value
	methods map[string]*methodType
}

func (p *SimpleProcessor) ReadMessageHeader(
	proto Protocol, name *string, mtype *byte, seqid *int32,
) (err error) {
	if *name, *mtype, *seqid, err = proto.ReadMessageBegin(); err != nil {
		return
	}
	if *mtype != T_CALL {
		return fmt.Errorf("message")
	}

	return
}

func (p *SimpleProcessor) ReadMessageBody(proto Protocol, req interface{}) (err error) {
	if err = proto.Skip(T_STRUCT); err != nil {
		return
	}
	if err = proto.ReadMessageEnd(); err != nil {
		return
	}
	return
}

func (p *SimpleProcessor) WriteMessageHeader(
	proto Protocol, name string, mtype byte, seqid int32,
) (err error) {
	if err = proto.WriteMessageBegin(name, mtype, seqid); err != nil {
		return
	}
	return
}

func (p *SimpleProcessor) WriteMessageBody(proto Protocol, res interface{}) (err error) {
	if err = proto.WriteByte(T_STOP); err != nil {
		return
	}
	if err = proto.WriteMessageEnd(); err != nil {
		return
	}
	return
}

func (p *SimpleProcessor) Process(iprot, oprot Protocol) (err error) {
	var (
		name   string
		mtype  byte
		seqid  int32
		method *methodType
	)

	if err = p.ReadMessageHeader(iprot, &name, &mtype, &seqid); err != nil {
		return
	}

	name = CamelCase(name)

	// determine if a method has been defined in processor.
	if method = p.methods[name]; method == nil {
		ae := NewTApplicationException(
			fmt.Sprintf("method %s not defined in %s", name, p.name),
			ExceptionUnknownMethod,
		)
		if err = iprot.Skip(T_STRUCT); err != nil {
			return
		}
		if err = p.WriteMessageHeader(oprot, name, T_EXCEPTION, seqid); err != nil {
			return
		}
		return WriteTApplicationException(oprot, ae)
	}

	// define request & response interface.
	var (
		req = reflect.New(method.RequestType)
		res = reflect.New(method.ResponseType)
	)

	if err = p.ReadMessageBody(iprot, req); err != nil {
		return
	}

	if err = p.WriteMessageHeader(oprot, name, T_REPLY, seqid); err != nil {
		return
	}
	if err = p.WriteMessageBody(iprot, res); err != nil {
		return
	}
	return
}

func NewSimpleProcessor(svr interface{}) *SimpleProcessor {

	rtype := reflect.TypeOf(svr)
	rval := reflect.ValueOf(svr)

	return &SimpleProcessor{
		name:    reflect.Indirect(rval).Type().Name(),
		methods: suitableMethods(rtype),
		rtype:   rtype,
		rval:    rval,
	}
}
