package rpc

import (
	"log"
	"reflect"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type methodType struct {
	method       reflect.Method
	RequestType  reflect.Type
	ResponseType reflect.Type
}

func suitableMethods(typ reflect.Type) (methods map[string]*methodType) {
	methods = make(map[string]*methodType)

	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		mname := method.Name
		mtype := method.Type

		if mtype.NumIn() != 3 {
			log.Println("method", mname, "has wrong number of ins:", method.Type.NumIn())
			continue
		}

		reqType := mtype.In(1)
		resType := mtype.In(1)

		// request should be a pointer.
		if reqType.Kind() != reflect.Ptr {
			log.Println("method", mname, "request type should be a pointer:", reqType)
			continue
		}

		// response should be a pointer.
		if resType.Kind() != reflect.Ptr {
			log.Println("method", mname, "response type should be a pointer:", resType)
			continue
		}

		// return type should be error
		if mtype.NumOut() != 1 || mtype.Out(0) != typeOfError {
			log.Println("method", mname, "not return error")
		}

		methods[mname] = &methodType{
			method:       method,
			RequestType:  reqType,
			ResponseType: resType,
		}
	}
	return
}
