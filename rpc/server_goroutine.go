package rpc

import (
	"fmt"
	"io"
	"net"

	"github.com/stdrickforce/thriftgo/protocol"
	"github.com/stdrickforce/thriftgo/transport"
)

type GoroutineServer struct {
	addr      string
	processor Processor
}

func (s *GoroutineServer) Serve() (err error) {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	fmt.Printf("server is listening on %s\n", s.addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			// TODO write log
		}
		fmt.Printf("new connection has been established: [%s]\n", conn.RemoteAddr().String())
		go s.process(conn)
	}
	return
}

func (s *GoroutineServer) process(conn net.Conn) {
	trans := transport.NewTSocketConn(conn)
	proto := protocol.NewTBinaryProtocol(trans, true, true)

	defer trans.Close()

	for {
		if err := s.processor.Process(proto, proto); err != nil {
			if err == io.EOF {
				fmt.Printf("connection has been closed: [%s]\n", conn.RemoteAddr().String())
			} else {
				fmt.Printf("unexpected error: [%s]", err.Error())
			}
			return
		}
	}
}

func NewGoroutineServer(processor Processor, addr string) *GoroutineServer {
	fmt.Println(processor)
	return &GoroutineServer{
		addr:      addr,
		processor: processor,
	}
}
