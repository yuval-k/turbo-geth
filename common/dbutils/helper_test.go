package dbutils_test

import (
	"bytes"
	"fmt"
	"sort"
	"testing"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/crypto"
)

func BenchmarkName(b *testing.B) {
	N := 1_000_000

	keys := make([][]byte, N)
	keys2 := make([][]byte, N)
	keys3 := make([][]byte, N)
	keys4 := make([][]byte, N)
	values := make([][]byte, N)
	m := make(map[string][]byte, N)
	for i := 0; i < N; i++ {
		keys[i] = crypto.Keccak256(common.FromHex(fmt.Sprintf("%x", i)))
		keys2[i] = keys[i]
		keys3[i] = keys[i]
		keys4[i] = keys[i]
	}

	b.Run("1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dbutils.Unique(keys, values)
		}
	})

	b.Run("2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, k := range keys2 {
				m[string(k)] = k
			}
		}
	})

	b.Run("3", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for k, v := range m {
				_ = []byte(k)
				_ = v
			}
		}
	})

	b.Run("4", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b := &BB{
				keys:   keys2,
				values: values,
			}
			sort.Stable(b)
		}
	})

	b.Run("5", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a, c := dbutils.Unique(keys, values)
			b := &BB{
				keys:   a,
				values: c,
			}
			sort.Sort(b)
		}
	})
}

type BB struct {
	keys   [][]byte
	values [][]byte
}

func (b *BB) Len() int {
	return len(b.keys)
}

func (b *BB) Less(i, j int) bool {
	return bytes.Compare(b.keys[i], b.keys[j]) < 0
}

func (b *BB) Swap(i, j int) {
	b.keys[i], b.keys[j] = b.keys[j], b.keys[i]
	b.values[i], b.values[j] = b.values[j], b.values[i]
}
