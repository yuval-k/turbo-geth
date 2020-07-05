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
	"context"
	"fmt"
	"testing"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrie2Seek(t *testing.T) {
	require, assert, t2 := require.New(t), assert.New(t), NewTrie2()

	hc := t2.wrapHashCollector(nil)
	require.NoError(hc(common.FromHex("010203"), []byte{1}))
	require.NoError(hc(common.FromHex("0102"), []byte{2}))         // this key must not be visible, because trie2 stores only Odd prefixes
	require.NoError(hc(common.FromHex("010200000000"), []byte{3})) // this key must not be visible, because trie2 stores only prefixes <= 5
	require.NoError(hc(common.FromHex("01"), []byte{4}))
	require.NoError(hc(common.FromHex("03"), []byte{5}))
	require.NoError(hc(common.FromHex("020908"), []byte{6}))
	require.NoError(hc(common.FromHex("0a"), []byte{7}))
	require.NoError(hc(common.FromHex("0a0d000000"), []byte{8}))
	t2.Reset()

	_, _, _ = t2.SeekTo([]byte{})
	res, v, _ := t2.Next()
	assert.Equal("0102", common.Bytes2Hex(res))
	assert.Equal([]byte{2}, v)

	res, v, _ = t2.Next()
	assert.Equal("010203", common.Bytes2Hex(res))
	assert.Equal([]byte{1}, v)

	_, _, _ = t2.SeekTo(common.FromHex("0a"))
	res, v, _ = t2.Next()
	assert.Equal("0a0d000000", common.Bytes2Hex(res))
	assert.Equal([]byte{8}, v)

	cases := []struct {
		in, expect string
		expectV    byte
	}{
		{"", "01", 4},
		{"0100", "0102", 2},
		{"04", "0a", 7},
		{"0a0d00", "0a0d000000", 8},
	}

	for _, c := range cases {
		res, v, _ := t2.SeekTo(common.FromHex(c.in))
		assert.Equal(c.expect, common.Bytes2Hex(res), "seek to "+c.in)
		assert.Equal([]byte{c.expectV}, v, "seek to "+c.in)
	}

	res, v, _ = t2.SeekTo(common.FromHex("0f"))
	require.Nil(res)
	require.Nil(v)
	res, v, _ = t2.SeekTo(common.FromHex("0f0e"))
	require.Nil(res)
	require.Nil(v)
	res, v, _ = t2.SeekTo(common.FromHex("0f0f"))
	require.Nil(res)
	require.Nil(v)
	res, v, _ = t2.SeekTo(common.FromHex("0f0f0f"))
	require.Nil(res)
	require.Nil(v)
	res, v, _ = t2.Next()
	require.Nil(res)
	require.Nil(v)
	res, v, _ = t2.Next()
	require.Nil(res)
	require.Nil(v)
}

func TestTwoAs1(t *testing.T) {
	require, assert, db, t2 := require.New(t), assert.New(t), ethdb.NewMemDatabase(), NewTrie2()
	defer db.Close()

	hc := t2.wrapHashCollector(nil)
	require.NoError(hc(common.FromHex("01"), []byte{4}))
	require.NoError(hc(common.FromHex("0102"), []byte{2}))         // this key must not be visible, because trie2 stores only Odd prefixes
	require.NoError(hc(common.FromHex("010200000000"), []byte{3})) // this key must not be visible, because trie2 stores only prefixes <= 5
	require.NoError(hc(common.FromHex("010203"), []byte{1}))
	require.NoError(hc(common.FromHex("020908"), []byte{6}))
	require.NoError(hc(common.FromHex("03"), []byte{5}))
	require.NoError(hc(common.FromHex("ad"), []byte{7}))
	require.NoError(hc(common.FromHex("ad0000"), []byte{8}))
	t2.Reset()

	require.NoError(db.Put(dbutils.IntermediateTrieHashBucket, common.FromHex("00"), []byte{1, 1}))
	require.NoError(db.Put(dbutils.IntermediateTrieHashBucket, common.FromHex("0001"), []byte{1, 2}))
	require.NoError(db.Put(dbutils.IntermediateTrieHashBucket, common.FromHex("11"), []byte{1, 3}))
	require.NoError(db.Put(dbutils.IntermediateTrieHashBucket, common.FromHex("12"), []byte{1, 4}))
	require.NoError(db.Put(dbutils.IntermediateTrieHashBucket, common.FromHex("1234"), []byte{1, 5}))
	require.NoError(db.Put(dbutils.IntermediateTrieHashBucket, common.FromHex("22"), []byte{1, 6}))

	cases := []struct {
		in, expect string
		expectV    []byte
	}{
		{"", "0000", []byte{1, 1}},
		{"0100", "0101", []byte{1, 3}},
		{"04", "ad", []byte{7}},
		{"05", "ad", []byte{7}},
		{"06", "ad", []byte{7}},
		{"ad00", "ad0000", []byte{8}},
	}

	require.NoError(db.KV().View(context.Background(), func(tx ethdb.Tx) error {
		var filter = func(k []byte) bool {
			return true
		}

		ihc := Filter(filter, IHDecompress(tx.Bucket(dbutils.IntermediateTrieHashBucket).Cursor()))
		t2c := Filter(filter, t2)
		ih := IH(t2c, ihc)

		for _, c := range cases {
			res, v, _, err := ih.SeekTo(common.FromHex(c.in))
			require.NoError(err)
			assert.Equal(c.expect, common.Bytes2Hex(res), "seek to "+c.in)
			assert.Equal(c.expectV, v, "seek to "+c.in)
		}

		res, v, _, err := ih.SeekTo(common.FromHex("fe"))
		require.NoError(err)
		require.Nil(res)
		require.Nil(v)
		res, v, _, err = ih.SeekTo(common.FromHex("ff"))
		require.NoError(err)
		require.Nil(res)
		require.Nil(v)

		ih.SeekTo([]byte{})
		ih.SeekTo(common.FromHex("0102"))
		fmt.Printf("1: %d\n", ih.skipSeek2Counter)
		fmt.Printf("2: seeksAcc=%d, seeksSt=%d, next=%d\n", ihc.seekAccCouner, ihc.seekStorageCouner, ihc.nextCounter)
		return nil
	}))
}
