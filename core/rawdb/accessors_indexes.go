// Copyright 2018 The go-ethereum Authors
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

package rawdb

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/log"
)

// TxLookupEntry is a positional metadata to help looking up the data content of
// a transaction or receipt given only its hash.
type TxLookupEntry struct {
	BlockHash  common.Hash
	BlockIndex uint64
	Index      uint64
}

// ReadTxLookupEntry retrieves the positional metadata associated with a transaction
// hash to allow retrieving the transaction or receipt by hash.
func ReadTxLookupEntry(db DatabaseReader, hash common.Hash) *uint64 {
	data, _ := db.Get(dbutils.TxLookupPrefix, hash.Bytes())
	if len(data) == 0 {
		return nil
	}
	number := new(big.Int).SetBytes(data).Uint64()
	return &number
}

// WriteTxLookupEntries stores a positional metadata for every transaction from
// a block, enabling hash based transaction and receipt lookups.
func WriteTxLookupEntries(db DatabaseWriter, block *types.Block) {
	for _, tx := range block.Transactions() {
		data := block.Number().Bytes()
		if err := db.Put(dbutils.TxLookupPrefix, tx.Hash().Bytes(), data); err != nil {
			log.Crit("Failed to store transaction lookup entry", "err", err)
		}
	}
}

// DeleteTxLookupEntry removes all transaction data associated with a hash.
func DeleteTxLookupEntry(db DatabaseDeleter, hash common.Hash) error {
	return db.Delete(dbutils.TxLookupPrefix, hash.Bytes(), nil)
}

// ReadTransaction retrieves a specific transaction from the database, along with
// its added positional metadata.
func ReadTransaction(db ethdb.Database, hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	blockNumber := ReadTxLookupEntry(db, hash)
	if blockNumber == nil {
		return nil, common.Hash{}, 0, 0
	}

	txs, err := CanonicalTransactions(db, *blockNumber)
	if err != nil {
		log.Error("ReadCanonicalHash failed", "err", err)
		return nil, common.Hash{}, 0, 0
	}

	for txIndex, tx := range txs {
		if tx.Hash() == hash {
			blockHash, err1 := ReadCanonicalHash(db, *blockNumber)
			if err1 != nil {
				log.Error("ReadCanonicalHash failed", "err", err1)
				return nil, common.Hash{}, 0, 0
			}
			return tx, blockHash, *blockNumber, uint64(txIndex)
		}
	}
	log.Error("Transaction not found", "number", blockNumber, "txhash", hash)
	return nil, common.Hash{}, 0, 0
}

func ReOrgTransactions2(db ethdb.Database, fromBlockN uint64, srcForkId, dstForkId uint8) error {
	if srcForkId < dstForkId {
		return fmt.Errorf("unexpected srcForkId: %d and dstForkId: %d", srcForkId, dstForkId)
	}

	blockN := fromBlockN

	kk := make([]byte, 8+4+1)
	binary.BigEndian.PutUint64(kk, blockN)
	kk[12] = 255 // canonical forkID

	_ = db.Walk(dbutils.EthTx, kk, 0, func(k, v []byte) (bool, error) {
		copy(kk[8:12], k[8:12]) // copy txID
		kk[12] = dstForkId

		if err := db.Delete(dbutils.EthTx, k, nil); err != nil {
			return false, err
		}
		if err := db.Put(dbutils.EthTx, kk, v); err != nil {
			return false, err
		}
		return true, nil
	})

	kk[12] = srcForkId
	binary.BigEndian.PutUint32(kk[8:], 0)

	return db.Walk(dbutils.EthTx, kk, 0, func(k, v []byte) (bool, error) {
		if k[12] != srcForkId {
			return false, nil
		}

		copy(kk[8:12], k[8:12]) // copy txID
		kk[12] = 255            // canonical forkID

		if err := db.Delete(dbutils.EthTx, k, nil); err != nil {
			return false, err
		}
		if err := db.Append(dbutils.EthTx, kk, v); err != nil {
			return false, err
		}
		return true, nil
	})
}

func ReOrgTransactions(tx ethdb.Tx, fromBlockN uint64, srcForkId, dstForkId uint8) error {
	if srcForkId < dstForkId {
		return fmt.Errorf("unexpected srcForkId: %d and dstForkId: %d", srcForkId, dstForkId)
	}

	c := tx.Cursor(dbutils.EthTx)
	defer c.Close()
	rwC := tx.Cursor(dbutils.EthTx)
	defer rwC.Close()

	blockN := fromBlockN

	kk := make([]byte, 8+4+1)
	binary.BigEndian.PutUint64(kk, blockN)
	kk[12] = 255 // canonical forkID

	//rename all last blocks to new forkId
	for k, v, _ := c.Seek(kk); k != nil; k, v, _ = c.Next() {
		copy(kk[8:12], k[8:12]) // copy txID
		kk[12] = dstForkId

		_, _ = rwC.SeekExact(k)
		if err := rwC.DeleteCurrent(); err != nil {
			return err
		}
		if err := rwC.Put(kk, v); err != nil {
			return err
		}
	}

	//Append new for as canonical
	kk[12] = srcForkId
	binary.BigEndian.PutUint32(kk[8:], 0)
	for k, v, _ := c.Seek(kk); k != nil; k, v, _ = c.Next() {
		if k[12] != srcForkId {
			break
		}

		copy(kk[8:12], k[8:12]) // copy txID
		kk[12] = 255            // canonical forkID

		_, _ = rwC.SeekExact(k)
		if err := rwC.DeleteCurrent(); err != nil {
			return err
		}
		if err := rwC.Append(kk, v); err != nil {
			return err
		}
	}
	return nil
}

// ReadReceipt retrieves a specific transaction receipt from the database, along with
// its added positional metadata.
func ReadReceipt(db ethdb.Database, hash common.Hash) (*types.Receipt, common.Hash, uint64, uint64) {
	// Retrieve the context of the receipt based on the transaction hash
	blockNumber := ReadTxLookupEntry(db, hash)
	if blockNumber == nil {
		return nil, common.Hash{}, 0, 0
	}
	blockHash, err := ReadCanonicalHash(db, *blockNumber)
	if err != nil {
		log.Error("ReadCanonicalHash failed", "err", err)
		return nil, common.Hash{}, 0, 0
	}
	if blockHash == (common.Hash{}) {
		return nil, common.Hash{}, 0, 0
	}
	// Read all the receipts from the block and return the one with the matching hash
	receipts := ReadReceipts(db, blockHash, *blockNumber)
	for receiptIndex, receipt := range receipts {
		if receipt.TxHash == hash {
			return receipt, blockHash, *blockNumber, uint64(receiptIndex)
		}
	}
	log.Error("Receipt not found", "number", blockNumber, "hash", blockHash, "txhash", hash)
	return nil, common.Hash{}, 0, 0
}

// ReadBloomBits retrieves the compressed bloom bit vector belonging to the given
// section and bit index from the.
func ReadBloomBits(db DatabaseReader, bit uint, section uint64, head common.Hash) ([]byte, error) {
	return db.Get(dbutils.BloomBitsPrefix, dbutils.BloomBitsKey(bit, section, head))
}

// WriteBloomBits stores the compressed bloom bits vector belonging to the given
// section and bit index.
func WriteBloomBits(db DatabaseWriter, bit uint, section uint64, head common.Hash, bits []byte) {
	if err := db.Put(dbutils.BloomBitsPrefix, dbutils.BloomBitsKey(bit, section, head), bits); err != nil {
		log.Crit("Failed to store bloom bits", "err", err)
	}
}
