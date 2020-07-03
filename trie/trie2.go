// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty off
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package trie implements Merkle Patricia Tries.
package trie

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/hexutil"
	"github.com/ledgerwatch/turbo-geth/core/types/accounts"
	"github.com/ledgerwatch/turbo-geth/ethdb"
)

type Trie2 struct {
	rl      *RetainList
	rl2     *RetainList
	values  [][]byte
	values2 [][]byte
	accs    map[common.Hash]*accounts.Account
}

func NewTrie2() *Trie2 {
	return &Trie2{
		rl:  NewRetainList(0),
		rl2: NewRetainList(0),
	}
}

func (t *Trie2) Reset() {
	t.rl, t.rl2 = t.rl2, t.rl
	t.values, t.values2 = t.values2, t.values
	t.rl.Rewind()
	t.rl2.Rewind()
	t.rl2.hexes = t.rl2.hexes[:0]
	t.values2 = t.values2[:0]

	fmt.Printf("t2: %d %d\n", len(t.values), len(t.rl.hexes))

	sort.Sort(t)
}

func (t *Trie2) Len() int           { return len(t.values) }
func (t *Trie2) Less(i, j int) bool { return bytes.Compare(t.rl.hexes[i], t.rl.hexes[j]) < 0 }
func (t *Trie2) Swap(i, j int) {
	t.rl.hexes[i], t.rl.hexes[j] = t.rl.hexes[j], t.rl.hexes[i]
	t.values[i], t.values[j] = t.values[j], t.values[i]
}

func (t *Trie2) Seek(prefix []byte) ([]byte, []byte) {
	rl := t.rl

	// Adjust "GT" if necessary
	var gtAdjusted bool
	for rl.lteIndex < len(rl.hexes)-1 && bytes.Compare(rl.hexes[rl.lteIndex+1], prefix) <= 0 {
		rl.lteIndex++
		gtAdjusted = true
	}

	// Adjust "LTE" if necessary (normally will not be necessary)
	for !gtAdjusted && rl.lteIndex > 0 && bytes.Compare(rl.hexes[rl.lteIndex], prefix) > 0 {
		rl.lteIndex--
	}

	if rl.lteIndex < len(rl.hexes) {
		if bytes.Compare(prefix, rl.hexes[rl.lteIndex]) <= 0 {
			return rl.hexes[rl.lteIndex], t.values[rl.lteIndex]
		}
	}

	if rl.lteIndex < len(rl.hexes)-1 {
		if bytes.Compare(prefix, rl.hexes[rl.lteIndex+1]) <= 0 {
			rl.lteIndex++
			return rl.hexes[rl.lteIndex], t.values[rl.lteIndex]
		}
	}

	return nil, nil
}

// _next - assume to be called after at least 1 .Seek call
func (t *Trie2) _next() ([]byte, []byte) {
	rl := t.rl
	if rl.lteIndex >= len(rl.hexes)-1 {
		return nil, nil
	}

	rl.lteIndex++
	return rl.hexes[rl.lteIndex], t.values[rl.lteIndex]
}

func (t *Trie2) wrapHashCollector(hc HashCollector) HashCollector {
	return func(keyHex []byte, hash []byte) error {
		if len(keyHex) < 6 {
			if hash != nil {
				//t.rl2.AddHex(common.CopyBytes(keyHex))
				//t.values2 = append(t.values2, common.CopyBytes(hash))
			}
		}

		if hc != nil {
			if err := hc(keyHex, hash); err != nil {
				return err
			}
		}

		return nil
	}
}

// IHCursor - skip all elements for which RetainDecider returns true
type IHCursor struct {
	c *TwoAs1Cursor

	k, v   []byte
	filter func(k []byte) bool
}

func SkipRetain(c *TwoAs1Cursor) *IHCursor { return &IHCursor{c: c} }

func (c *IHCursor) Filter(filter func(k []byte) bool) *IHCursor {
	c.filter = filter
	return c
}

func (c *IHCursor) _seekTo(seek []byte) (err error) {
	c.k, c.v, err = c.c.SeekTo(seek)
	if err != nil {
		return err
	}
	if c.k == nil {
		return nil
	}

	if c.filter(c.k) {
		return nil
	}

	return c._next()
}

func (c *IHCursor) _next() (err error) {
	c.k, c.v, err = c.c.Next()
	if err != nil {
		return err
	}
	for {
		if c.k == nil {
			return nil
		}

		if c.filter(c.k) {
			return nil
		}

		c.k, c.v, err = c.c.Next()
		if err != nil {
			return err
		}
	}
}

func (c *IHCursor) SeekTo(seek []byte) ([]byte, []byte, bool, error) {
	if err := c._seekTo(seek); err != nil {
		return []byte{}, nil, false, err
	}

	isSequence := false
	if c.k == nil {
		return c.k, c.v, false, nil
	}

	if len(seek) == common.HashLength*2+common.IncarnationLength*2 {
		isSequence = bytes.Equal(c.k, seek)
	} else {
		if bytes.HasPrefix(c.k, seek) {
			tail := c.k[len(seek):] // if tail has only zeroes, then no state records can be between fstl.nextHex and fstl.ihK
			isSequence = true
			for _, n := range tail {
				if n != 0 {
					isSequence = false
					break
				}
			}
		} else {
			if bytes.HasPrefix(seek, c.k) {
				//fmt.Printf("1: %x, %x\n", seek, c.k)
			}
		}
	}

	return c.k, c.v, isSequence, nil
}

func IHIsValid(nibbles []byte) bool {
	return len(nibbles) != 0 && len(nibbles)%2 == 0
}

// TowAs1 - manage 2 cursors - make it looks as 1
type TwoAs1Cursor struct {
	c1               *Trie2
	c2               *IHDecompressCursor
	used1            bool
	used2            bool
	skipSeek2        bool
	skipSeek2Counter int
	k1, v1, k2, v2   []byte
}

func TwoAs1(c1 *Trie2, c2 *IHDecompressCursor) *TwoAs1Cursor {
	return &TwoAs1Cursor{
		c1: c1,
		c2: c2,
	}
}

func (c *TwoAs1Cursor) SeekTo(seek []byte) ([]byte, []byte, error) {
	var err error
	if len(seek) == 0 || (c.k1 == nil && !c.used1) || keyIsBefore(c.k1, seek) {
		c.k1, c.v1 = c.c1.Seek(seek)
		if c.k1 != nil && len(seek) == len(c.k1) && bytes.Equal(seek, c.k1) {
			c.used1 = true
			c.skipSeek2 = true
			c.skipSeek2Counter++
			return c.k1, c.v1, nil
		}
		c.used1 = false
	}

	if len(seek) == 0 || (c.k2 == nil && !c.used2) || keyIsBefore(c.k2, seek) {
		c.k2, c.v2, err = c.c2.SeekTo(seek)
		if err != nil {
			return []byte{}, nil, err
		}
		c.used2 = false
	}

	if c.k1 != nil && keyIsBefore(c.k1, c.k2) {
		if c.k1 != nil {
			//fmt.Printf("2: %x, %x\n", seek, c.k1)
		}
		c.used1 = true
		return c.k1, c.v1, nil
	}

	if c.k1 != nil {
		//fmt.Printf("3: %x, %x, %x\n", seek, c.k1, c.k2)
	}
	c.used2 = true
	return c.k2, c.v2, nil
}

func (c *TwoAs1Cursor) Next() ([]byte, []byte, error) {
	var err error
	if c.skipSeek2 {
		c.k2, c.v2, err = c.c2.SeekTo(c.k1)
		if err != nil {
			return []byte{}, nil, err
		}
		c.skipSeek2 = false
		c.used2 = true
	}

	if c.used1 {
		c.k1, c.v1 = c.c1._next()
		c.used1 = false
	}

	if c.used2 {
		c.k2, c.v2, err = c.c2.Next()
		if err != nil {
			return []byte{}, nil, err
		}
		c.used2 = false
	}

	if c.k1 != nil && keyIsBefore(c.k1, c.k2) {
		c.used1 = true
		return c.k1, c.v1, nil
	}

	c.used2 = true
	return c.k2, c.v2, nil
}

// StateCursor - does decompression of keys
type IHDecompressCursor struct {
	c          ethdb.Cursor
	k, kHex, v []byte
	seek       []byte
}

func IHDecompress(c ethdb.Cursor) *IHDecompressCursor {
	return &IHDecompressCursor{
		c:    c,
		seek: make([]byte, 0, common.HashLength*3),
		kHex: make([]byte, 0, common.HashLength*3),
	}
}

func (c *IHDecompressCursor) SeekTo(seek []byte) ([]byte, []byte, error) {
	hexutil.FromNibbles2(seek, &c.seek)
	var err error
	c.k, c.v, err = c.c.SeekTo(c.seek)
	if err != nil {
		return []byte{}, nil, err
	}
	if c.k == nil {
		return nil, nil, nil
	}

	hexutil.ToNibbles(c.k, &c.kHex)
	return c.kHex, c.v, nil
}

func (c *IHDecompressCursor) Next() ([]byte, []byte, error) {
	var err error
	c.k, c.v, err = c.c.Next()
	if err != nil {
		return []byte{}, nil, err
	}
	if c.k == nil {
		return nil, nil, nil
	}
	hexutil.ToNibbles(c.k, &c.kHex)
	return c.kHex, c.v, nil
}

// StateCursor - does decompression of keys
type StateCursor struct {
	c          cursor
	k, kHex, v []byte
}

func NewStateCursor(c cursor) *StateCursor {
	return &StateCursor{
		c: c,
	}
}

func (dc *StateCursor) SeekTo(seek []byte) ([]byte, []byte, []byte, error) {
	var err error
	dc.k, dc.v, err = dc.c.SeekTo(seek)
	if err != nil {
		return []byte{}, nil, nil, err
	}
	if dc.k == nil {
		return nil, nil, nil, nil
	}
	hexutil.ToNibbles(dc.k, &dc.kHex)
	return dc.k, dc.kHex, dc.v, nil
}

func (dc *StateCursor) Next() ([]byte, []byte, []byte, error) {
	var err error
	dc.k, dc.v, err = dc.c.Next()
	if err != nil {
		return []byte{}, nil, nil, err
	}
	if dc.k == nil {
		return nil, nil, nil, nil
	}
	hexutil.ToNibbles(dc.k, &dc.kHex)
	return dc.k, dc.kHex, dc.v, nil
}
