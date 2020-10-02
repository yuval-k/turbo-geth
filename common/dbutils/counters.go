package dbutils

import (
	"github.com/ledgerwatch/turbo-geth/ethdb/cbor"
)

const KeyIDs = "ids"

// IDs - store id of last inserted entity to db - increment before use
type IDs struct {
	Topic uint64
}

const KeyAggregates = "aggregates"

// Aggregates - store some statistical aggregates of data: for example min/max of values in some bucket
type Aggregates struct {
}

func (c *IDs) Unmarshal(data []byte) error {
	return cbor.Unmarshal(c, data)
}

func (c *IDs) Marshal() (data []byte, err error) {
	data = []byte{}
	err = cbor.Marshal(&data, c)
	return data, err
}

func (c *Aggregates) Unmarshal(data []byte) error {
	return cbor.Unmarshal(c, data)
}

func (c *Aggregates) Marshal() (data []byte, err error) {
	data = []byte{}
	err = cbor.Marshal(&data, c)
	return data, err
}
