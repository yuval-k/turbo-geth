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

const LevelsInMem = 5

type Trie2 struct {
	rl     *RetainList
	values [][]byte
	batch  map[string][]byte
	accs   map[common.Hash]*accounts.Account
}

func NewTrie2() *Trie2 {
	return &Trie2{
		rl:    NewRetainList(0),
		batch: map[string][]byte{},
	}
}

func (t *Trie2) Reset() {
	t.rl.hexes = t.rl.hexes[:0]
	t.values = t.values[:0]
	for k, v := range t.batch {
		t.rl.hexes = append(t.rl.hexes, []byte(k))
		t.values = append(t.values, v)
	}
	t.rl.Rewind()

	sort.Sort(t)

	fmt.Printf("t2: %d %d\n", len(t.values), len(t.rl.hexes))
}

func (t *Trie2) Len() int           { return len(t.values) }
func (t *Trie2) Less(i, j int) bool { return bytes.Compare(t.rl.hexes[i], t.rl.hexes[j]) < 0 }
func (t *Trie2) Swap(i, j int) {
	t.rl.hexes[i], t.rl.hexes[j] = t.rl.hexes[j], t.rl.hexes[i]
	t.values[i], t.values[j] = t.values[j], t.values[i]
}

func (t *Trie2) SeekTo(prefix []byte) ([]byte, []byte, error) {
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
			return rl.hexes[rl.lteIndex], t.values[rl.lteIndex], nil
		}
	}

	if rl.lteIndex < len(rl.hexes)-1 {
		if bytes.Compare(prefix, rl.hexes[rl.lteIndex+1]) <= 0 {
			rl.lteIndex++
			return rl.hexes[rl.lteIndex], t.values[rl.lteIndex], nil
		}
	}

	return nil, nil, nil
}

// Next - assume to be called after at least 1 .Seek call
func (t *Trie2) Next() ([]byte, []byte, error) {
	rl := t.rl
	if rl.lteIndex >= len(rl.hexes)-1 {
		return nil, nil, nil
	}

	rl.lteIndex++
	return rl.hexes[rl.lteIndex], t.values[rl.lteIndex], nil
}

func (t *Trie2) wrapHashCollector(hc HashCollector) HashCollector {
	return func(keyHex []byte, hash []byte) error {
		if len(keyHex) <= LevelsInMem {
			if hash == nil {
				delete(t.batch, string(keyHex))
			} else {
				t.batch[string(keyHex)] = common.CopyBytes(hash)
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

// FilterCursor - skip all elements for which RetainDecider returns true
type FilterCursor struct {
	c cursor

	k, v              []byte
	filter            func(k []byte) bool
	seekAccCouner     int
	seekStorageCouner int
	nextCounter       int
}

func Filter(filter func(k []byte) bool, c cursor) *FilterCursor {
	return &FilterCursor{c: c, filter: filter}
}

func (c *FilterCursor) _seekTo(seek []byte) (err error) {
	if len(seek) > 64 {
		c.seekStorageCouner++
	} else {
		c.seekAccCouner++
	}

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

func (c *FilterCursor) _next() (err error) {
	c.nextCounter++
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

		c.nextCounter++
		c.k, c.v, err = c.c.Next()
		if err != nil {
			return err
		}
	}
}

func (c *FilterCursor) SeekTo(seek []byte) ([]byte, []byte, error) {
	if err := c._seekTo(seek); err != nil {
		return []byte{}, nil, err
	}

	return c.k, c.v, nil
}

func IHIsValid(nibbles []byte) bool {
	return len(nibbles) != 0 && len(nibbles)%2 == 0
}

// IHCursor - manage 2 cursors - make it looks as 1
type IHCursor struct {
	c1                    *FilterCursor
	c2                    *FilterCursor
	used1                 bool
	used2                 bool
	skipSeek2             bool
	skipSeek2Counter      int
	k1, k1Old, v1, k2, v2 []byte
}

func IH(c1 *FilterCursor, c2 *FilterCursor) *IHCursor {
	return &IHCursor{
		c1: c1,
		c2: c2,
	}
}

func (c *IHCursor) _seek1(seek []byte) error {
	var err error
	c.k1, c.v1, err = c.c1.SeekTo(seek)
	if err != nil {
		return err
	}

	return nil
}

func (c *IHCursor) _seek2(seek []byte) error {
	var err error
	c.k2, c.v2, err = c.c2.SeekTo(seek)
	if err != nil {
		return err
	}
	c.used2 = false
	c.skipSeek2 = false
	return nil
}

func (c *IHCursor) SeekTo(seek []byte) ([]byte, []byte, bool, error) {
	if err := c._seek1(seek); err != nil {
		return []byte{}, nil, false, err
	}
	if c.k1 != nil {
		isSequence := false
		if bytes.HasPrefix(c.k1, seek) {
			tail := c.k1[len(seek):] // if tail has only zeroes, then no state records can be between fstl.nextHex and fstl.ihK
			isSequence = true
			for _, n := range tail {
				if n != 0 {
					isSequence = false
					break
				}
			}
		}

		if isSequence {
			c.skipSeek2 = true
			c.skipSeek2Counter++
			return c.k1, c.v1, true, nil
		}
	}

	//fmt.Printf("Before: %x, k1=%x, k2=%x\n", seek, c.k1, c.k2)
	if err := c._seek2(seek); err != nil {
		return []byte{}, nil, false, err
	}
	//fmt.Printf("After: k1=%x, k2=%x, %t\n", c.k1, c.k2, keyIsBefore(c.k1, c.k2))

	if c.k1 != nil && keyIsBefore(c.k1, c.k2) {
		return c.k1, c.v1, false, nil
	}

	c.used2 = true

	if c.k2 != nil {
		isSequence := false
		if bytes.HasPrefix(c.k2, seek) {
			tail := c.k2[len(seek):] // if tail has only zeroes, then no state records can be between fstl.nextHex and fstl.ihK
			isSequence = true
			for _, n := range tail {
				if n != 0 {
					isSequence = false
					break
				}
			}
		}

		if isSequence {
			return c.k2, c.v2, true, nil
		}
	}
	return c.k2, c.v2, false, nil
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
