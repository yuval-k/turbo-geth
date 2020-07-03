// Copyright 2016 The go-ethereum Authors
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

package hexutil_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/hexutil"
	"github.com/ledgerwatch/turbo-geth/common/pool"
	"github.com/stretchr/testify/assert"
)

func TestCompressNibbles(t *testing.T) {
	cases := []struct {
		in     string
		expect string
	}{
		{in: "0000", expect: "00"},
		{in: "0102", expect: "12"},
		{in: "0102030405060708090f", expect: "123456789f"},
		{in: "0f000101", expect: "f011"},
		{in: "", expect: ""},
	}

	compressBuf := pool.GetBuffer(64)
	defer pool.PutBuffer(compressBuf)
	decompressBuf := pool.GetBuffer(64)
	defer pool.PutBuffer(decompressBuf)
	for _, tc := range cases {
		compressBuf.Reset()
		decompressBuf.Reset()

		in := common.Hex2Bytes(tc.in)
		hexutil.FromNibbles2(in, &compressBuf.B)
		compressed := compressBuf.Bytes()
		msg := "On: " + tc.in + " Len: " + strconv.Itoa(len(compressed))
		assert.Equal(t, tc.expect, fmt.Sprintf("%x", compressed), msg)
		hexutil.ToNibbles(compressed, &decompressBuf.B)
		decompressed := decompressBuf.Bytes()
		assert.Equal(t, tc.in, fmt.Sprintf("%x", decompressed), msg)
	}
}
