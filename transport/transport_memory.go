package transport

import (
	"bytes"
	"time"
)

type TMemoryBuffer struct {
	*bytes.Buffer
}

func (self *TMemoryBuffer) SetTimeout(d time.Duration) {
}

func (self *TMemoryBuffer) Open() error {
	return nil
}

func (self *TMemoryBuffer) Close() error {
	return nil
}

func (self *TMemoryBuffer) Flush() error {
	return nil
}

func (self *TMemoryBuffer) GetBytes() []byte {
	return self.Buffer.Bytes()
}

func NewTMemoryBuffer() *TMemoryBuffer {
	return &TMemoryBuffer{
		bytes.NewBuffer([]byte{}),
	}
}
