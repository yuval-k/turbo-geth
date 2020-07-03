package dbutils

import (
	"testing"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/stretchr/testify/assert"
)

func TestNextSubtree(t *testing.T) {
	var res []byte
	cases := []struct {
		in  string
		out string
	}{
		{"0102", "0103"},
		{"010f", "02"},
		{"0f0d", "0f0e"},
		{"010203040506070809", "01020304050607080a"},
		{"01020304050607080f", "0102030405060709"},
	}
	for _, c := range cases {
		NextSubtreeHex2(common.FromHex(c.in), &res)
		assert.Equal(t, c.out, common.Bytes2Hex(res))
	}
}
