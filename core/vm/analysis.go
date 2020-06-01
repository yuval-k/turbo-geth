// Copyright 2014 The go-ethereum Authors
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
	"sort"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/pool"
)

var jumpdests cache = newJumpDests(50000, 10, 1)

type cache interface {
	Set(hash common.Hash, v *pool.ByteBuffer)
	Get(hash common.Hash) (*pool.ByteBuffer, bool)
	Clear(codeHash common.Hash, local *pool.ByteBuffer)
}

type jumpDests struct {
	maps       map[common.Hash]*item
	lru        []map[common.Hash]struct{}
	chunks     []int
	minToClear int // 1..1000
	maxSize    int
}

type item struct {
	m    *pool.ByteBuffer
	used int
}

func newJumpDests(maxSize, nChunks, perMilleToClear int) *jumpDests {
	lru := make([]map[common.Hash]struct{}, nChunks)
	chunks := make([]int, nChunks)
	chunkSize := maxSize / nChunks
	for i := 0; i < nChunks; i++ {
		lru[i] = make(map[common.Hash]struct{}, chunkSize)
		chunks[i] = 1 << (1 + i*2)
	}

	return &jumpDests{
		make(map[common.Hash]*item, maxSize),
		lru,
		chunks,
		maxSize * perMilleToClear / 1000,
		maxSize,
	}
}

func (j *jumpDests) Set(hash common.Hash, v *pool.ByteBuffer) {
	_, ok := j.maps[hash]
	if ok {
		return
	}

	if len(j.maps) >= j.maxSize {
		j.gc()
	}

	j.maps[hash] = &item{v, 1}
	j.lru[0][hash] = struct{}{}
	return
}

func (j *jumpDests) Get(hash common.Hash) (*pool.ByteBuffer, bool) {
	jumps, ok := j.maps[hash]
	if !ok {
		return nil, false
	}

	jumps.used++
	idx := sort.SearchInts(j.chunks, jumps.used)

	// everything greater than j.chunks[len(chunks)-1] should be stored in the last chunk
	if idx >= 0 && idx < len(j.chunks)-1 {
		max := j.chunks[idx]
		if jumps.used >= max {
			// moving to the next chunk
			j.lru[idx+1][hash] = struct{}{}
			delete(j.lru[idx], hash)
		}
	}

	return jumps.m, true
}

func (j *jumpDests) gc() {
	n := 0
	for _, chunk := range j.lru {
		for hash := range chunk {
			delete(chunk, hash)

			item := j.maps[hash]
			delete(j.maps, hash)
			pool.PutBuffer(item.m)

			n++
			if n >= j.minToClear {
				return
			}
		}
	}
}

func (j *jumpDests) Clear(codeHash common.Hash, local *pool.ByteBuffer) {
	if codeHash == (common.Hash{}) {
		return
	}
	_, ok := j.maps[codeHash]
	if ok {
		return
	}
	// analysis is a local one
	pool.PutBuffer(local)
}

// codeBitmap collects data locations in code.
func codeBitmap(code []byte) *pool.ByteBuffer {
	// The bitmap is 4 bytes longer than necessary, in case the code
	// ends with a PUSH32, the algorithm will push zeroes onto the
	// bitvector outside the bounds of the actual code.
	bits := pool.GetBuffer(uint(len(code)/8 + 1 + 4))
	for pc := uint64(0); pc < uint64(len(code)); {
		op := OpCode(code[pc])

		if op >= PUSH1 && op <= PUSH32 {
			numbits := op - PUSH1 + 1
			pc++
			for ; numbits >= 8; numbits -= 8 {
				bits.SetBit8Pos(pc) // 8
				pc += 8
			}
			for ; numbits > 0; numbits-- {
				bits.SetBitPos(pc)
				pc++
			}
		} else {
			pc++
		}
	}
	return bits
}
