package server

import (
	"sync"
)

type bufPool interface {
	Get() []byte
	Put([]byte)
}

type pool struct {
	*sync.Pool
	size int
}

func newPool(size int) bufPool {
	return &pool{
		size: size,
		Pool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, size)
			},
		},
	}
}

func (p *pool) Get() []byte {
	return p.Pool.Get().([]byte)
}

func (p *pool) Put(b []byte) {
	if cap(b) != p.size {
		panic("invalid buffer size")
	}
	p.Pool.Put(b[:0])
}
