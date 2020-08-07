package ethdb

import (
	"context"
	"fmt"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"time"
)

// TxDb - Database interface around TX
// It's not thread-safe!
type TxDb struct {
	db      Database
	tx      Tx
	cursors map[string]*LmdbCursor
	len     uint64
}

func (m *TxDb) Close() {
	panic("don't call me")
}

func (m *TxDb) Begin() (DbWithPendingMutations, error) {
	batch := &TxDb{db: m.db, cursors: map[string]*LmdbCursor{}}
	if err := batch.begin(); err != nil {
		return nil, err
	}
	return batch, nil
}

func (m *TxDb) Put(bucket, key []byte, value []byte) error {
	m.len += uint64(len(key) + len(value))
	if value == nil {
		return m.cursors[string(bucket)].Delete(key)
	}
	return m.cursors[string(bucket)].Put(key, value)
}

func (m *TxDb) Delete(bucket, key []byte) error {
	m.len += uint64(len(key))
	return m.cursors[string(bucket)].Delete(key)
}

func (m *TxDb) NewBatch() DbWithPendingMutations {
	panic("don't call me")
}

func (m *TxDb) begin() error {
	tx, err := m.db.(HasKV).KV().Begin(context.Background(), true)
	if err != nil {
		return err
	}
	m.tx = tx
	for i := range dbutils.Buckets {
		m.cursors[string(dbutils.Buckets[i])] = tx.Bucket(dbutils.Buckets[i]).Cursor().(*LmdbCursor)
		if err := m.cursors[string(dbutils.Buckets[i])].initCursor(); err != nil {
			return err
		}
	}
	return nil
}

func (m *TxDb) KV() KV {
	if casted, ok := m.db.(HasKV); ok {
		return casted.KV()
	}
	return nil
}

// Can only be called from the worker thread
func (m *TxDb) Get(bucket, key []byte) ([]byte, error) {
	v, err := m.cursors[string(bucket)].SeekExact(key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrKeyNotFound
	}
	return v, nil
}

func (m *TxDb) GetIndexChunk(bucket, key []byte, timestamp uint64) ([]byte, error) {
	if m.db != nil {
		return m.db.GetIndexChunk(bucket, key, timestamp)
	}
	return nil, ErrKeyNotFound
}

func (m *TxDb) Has(bucket, key []byte) (bool, error) {
	v, err := m.Get(bucket, key)
	if err != nil {
		return false, err
	}
	return v != nil, nil
}

func (m *TxDb) DiskSize(ctx context.Context) (common.StorageSize, error) {
	if m.db == nil {
		return 0, nil
	}
	sz, err := m.db.(HasStats).DiskSize(ctx)
	if err != nil {
		return 0, err
	}
	return common.StorageSize(sz), nil
}

func (m *TxDb) MultiPut(tuples ...[]byte) (uint64, error) {
	panic("don't use me")
}

func (m *TxDb) BatchSize() int {
	return int(m.len)
}

// IdealBatchSize defines the size of the data batches should ideally add in one write.
func (m *TxDb) IdealBatchSize() int {
	return m.db.IdealBatchSize() * 100
}

// WARNING: Merged mem/DB walk is not implemented
func (m *TxDb) Walk(bucket, startkey []byte, fixedbits int, walker func([]byte, []byte) (bool, error)) error {
	m.panicOnEmptyDB()
	return m.db.Walk(bucket, startkey, fixedbits, walker)
}

// WARNING: Merged mem/DB walk is not implemented
func (m *TxDb) MultiWalk(bucket []byte, startkeys [][]byte, fixedbits []int, walker func(int, []byte, []byte) error) error {
	m.panicOnEmptyDB()
	return m.db.MultiWalk(bucket, startkeys, fixedbits, walker)
}

func (m *TxDb) Commit() (uint64, error) {
	defer func(t time.Time) { fmt.Printf("%s\n", time.Since(t)) }(time.Now())
	if m.db == nil {
		return 0, nil
	}
	if err := m.tx.Commit(context.Background()); err != nil {
		return 0, err
	}
	m.len = 0
	if err := m.begin(); err != nil {
		return 0, err
	}
	return 0, nil
}

func (m *TxDb) Rollback() {
	m.tx.Rollback()
	m.len = 0
}

func (m *TxDb) Keys() ([][]byte, error) {
	panic("don't use me")
}

func (m *TxDb) panicOnEmptyDB() {
	if m.db == nil {
		panic("Not implemented")
	}
}

func (m *TxDb) MemCopy() Database {
	m.panicOnEmptyDB()
	return m.db
}

// [TURBO-GETH] Freezer support (not implemented yet)
// Ancients returns an error as we don't have a backing chain freezer.
func (m *TxDb) Ancients() (uint64, error) {
	return 0, errNotSupported
}

// TruncateAncients returns an error as we don't have a backing chain freezer.
func (m *TxDb) TruncateAncients(items uint64) error {
	return errNotSupported
}
