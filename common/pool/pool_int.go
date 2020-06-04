package pool

import (
	"github.com/holiman/uint256"
)

var Unt256 = NewUint256Pool(1024)

type uint256Pool struct {
	ints []*uint256.Int
}

func NewUint256Pool(c int) *uint256Pool {
	p := new(uint256Pool)
	for i := 0; i < c; i++ {
		p.ints = append(p.ints, new(uint256.Int))

	}
	return p
}

func (pool *uint256Pool) Get() *uint256.Int {
	if pool == nil {
		return new(uint256.Int)
	}

	l := len(pool.ints)
	if l == 0 {
		return new(uint256.Int)
	}

	n := pool.ints[l-1]
	pool.ints = pool.ints[:l-1]
	return n
}

func (pool *uint256Pool) GetZero() *uint256.Int {
	if pool == nil {
		return new(uint256.Int)
	}

	l := len(pool.ints)
	if l == 0 {
		return new(uint256.Int)
	}

	n := pool.ints[l-1]
	pool.ints = pool.ints[:l-1]
	return n.SetUint64(0)
}

func (pool *uint256Pool) Put(ns ...*uint256.Int) {
	if pool == nil || len(ns) == 0 {
		return
	}
	pool.ints = append(pool.ints, ns...)
}
