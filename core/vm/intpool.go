// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"math/big"
	"sync"
)

var checkVal = big.NewInt(-42)

const poolLimit = 256

// IntPool is a pool of big integers that
// can be reused for all big.Int operations.
type IntPool struct {
	poolLimit int
	pool      *Stack
}

func NewIntPool(poolLimit int) *IntPool {
	return &IntPool{poolLimit: poolLimit, pool: NewStack(poolLimit)}
}

// get retrieves a big int from the pool, allocating one if the pool is empty.
// Note, the returned int's value is arbitrary and will not be zeroed!
func (p *IntPool) Get() *big.Int {
	if p.pool.len() > 0 {
		return p.pool.pop()
	}
	return new(big.Int)
}

// getZero retrieves a big int from the pool, setting it to zero or allocating
// a new one if the pool is empty.
func (p *IntPool) GetZero() *big.Int {
	if p.pool.len() > 0 {
		return p.pool.pop().SetUint64(0)
	}
	return new(big.Int)
}

// put returns an allocated big int to the pool to be later reused by get calls.
// Note, the values as saved as is; neither put nor get zeroes the ints out!
func (p *IntPool) Put(is ...*big.Int) {
	if len(p.pool.data) > poolLimit {
		return
	}
	for _, i := range is {
		if i == nil {
			continue
		}

		// verifyPool is a build flag. Pool verification makes sure the integrity
		// of the integer pool by comparing values to a default value.
		if verifyPool {
			i.Set(checkVal)
		}
		p.pool.push(i)
	}
}

// The IntPool pool's default capacity
const poolDefaultCap = 25

// intPoolPool manages a pool of intPools.
type intPoolPool struct {
	pools     []*IntPool
	poolLimit int
	lock      sync.Mutex
}

var PoolOfIntPools *intPoolPool
var poolOfIntPoolsOnce sync.Once

func GetPoolOfIntPools(limit ...int) *intPoolPool {
	poolOfIntPoolsOnce.Do(func() {
		l := poolLimit
		if len(limit) > 0 && limit[0] > 0 {
			l = limit[0]
		}

		PoolOfIntPools = &intPoolPool{
			pools:     make([]*IntPool, 0, poolDefaultCap),
			poolLimit: l,
		}
	})

	return PoolOfIntPools
}

// get is looking for an available pool to return.
func (ipp *intPoolPool) Get() *IntPool {
	ipp.lock.Lock()
	defer ipp.lock.Unlock()

	if len(PoolOfIntPools.pools) > 0 {
		ip := ipp.pools[len(ipp.pools)-1]
		ipp.pools = ipp.pools[:len(ipp.pools)-1]
		return ip
	}
	return NewIntPool(poolLimit)
}

// put a pool that has been allocated with get.
func (ipp *intPoolPool) Put(ip *IntPool) {
	ipp.lock.Lock()
	defer ipp.lock.Unlock()

	if len(ipp.pools) < cap(ipp.pools) {
		ipp.pools = append(ipp.pools, ip)
	}
}
