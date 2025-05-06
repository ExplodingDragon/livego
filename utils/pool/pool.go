package pool

type Pool struct {
	pos int
	buf []byte
}

const maxPoolSize = 500 * 1024

func (pool *Pool) Get(size int) []byte {
	if maxPoolSize-pool.pos < size {
		pool.pos = 0
		pool.buf = make([]byte, maxPoolSize)
	}
	b := pool.buf[pool.pos : pool.pos+size]
	pool.pos += size
	return b
}

func NewPool() *Pool {
	return &Pool{
		buf: make([]byte, maxPoolSize),
	}
}
