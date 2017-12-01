package protocol

import (
	. "github.com/stdrickforce/thriftgo/thrift"
)

func WriteTApplicationException(proto Protocol, e *TApplicationException) (err error) {
	if err = proto.WriteStructBegin("TApplicationException"); err != nil {
		return
	}
	if err = proto.WriteFieldBegin("message", T_STRING, 1); err != nil {
		return
	}
	if err = proto.WriteString(e.Message); err != nil {
		return
	}
	if err = proto.WriteFieldEnd(); err != nil {
		return
	}
	if err = proto.WriteFieldBegin("type", T_I32, 2); err != nil {
		return
	}
	if err = proto.WriteI32(e.Type); err != nil {
		return
	}
	if err = proto.WriteFieldEnd(); err != nil {
		return
	}
	if err = proto.WriteFieldStop(); err != nil {
		return
	}
	if err = proto.WriteStructEnd(); err != nil {
		return
	}
	return
}

func ReadTApplicationException(proto Protocol) (e *TApplicationException, err error) {
	var (
		message string
		code    int32
	)
	if err = proto.ReadStructBegin(); err != nil {
		return
	}
	if _, _, err = proto.ReadFieldBegin(); err != nil {
		return
	}
	if message, err = proto.ReadString(); err != nil {
		return
	}
	if err = proto.ReadFieldEnd(); err != nil {
		return
	}
	if _, _, err = proto.ReadFieldBegin(); err != nil {
		return
	}
	if code, err = proto.ReadI32(); err != nil {
		return
	}
	if err = proto.ReadFieldEnd(); err != nil {
		return
	}
	if err = proto.ReadStructEnd(); err != nil {
		return
	}
	// NOTE read field stop
	if _, _, err = proto.ReadFieldBegin(); err != nil {
		return
	}
	e = NewTApplicationException(message, code)
	return
}
