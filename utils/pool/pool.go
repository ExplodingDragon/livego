package pool

import (
	"sync"
)

type Pool struct {
	mu  sync.Mutex
	pos int
	buf []byte
}

const maxPoolSize = 1000 * 1024

func (p *Pool) Get(size int) []byte {
	if size > maxPoolSize {
		return make([]byte, size)
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	if maxPoolSize-p.pos < size {
		p.pos = 0
	}

	result := p.buf[p.pos : p.pos+size]
	p.pos += size
	return result
}

func NewPool() *Pool {
	return &Pool{
		buf: make([]byte, maxPoolSize),
	}
}
