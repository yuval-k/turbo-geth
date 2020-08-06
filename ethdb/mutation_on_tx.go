package ethdb

import (
	"context"
	"fmt"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"time"
)

// Mutation alternative - which doesn't store anything in mem, but instead open KV transaction and call Get/Put/Commit/Rollback directly on it
// It's not thread-safe!
// please use defer batch.Close() when use this class
// it's not production ready, just experimental
type mutationOnTx struct {
	db      Database
	tx      Tx
	cursors map[string]*LmdbCursor
	len     uint64
}

func (m *mutationOnTx) begin() {
	tx, err := m.db.(HasKV).KV().Begin(context.Background(), true)
	if err != nil {
		panic(err)
	}
	m.tx = tx
	for i := range dbutils.Buckets {
		m.cursors[string(dbutils.Buckets[i])] = tx.Bucket(dbutils.Buckets[i]).Cursor().(*LmdbCursor)
		m.cursors[string(dbutils.Buckets[i])].initCursor()
	}
}

func NewMutationOnTx(db Database) DbWithPendingMutations {
	m := &mutationOnTx{db: db, cursors: map[string]*LmdbCursor{}}
	m.begin()
	return m
}

func (m *mutationOnTx) KV() KV {
	if casted, ok := m.db.(HasKV); ok {
		return casted.KV()
	}
	return nil
}

// Can only be called from the worker thread
func (m *mutationOnTx) Get(bucket, key []byte) ([]byte, error) {
	v, err := m.cursors[string(bucket)].SeekExact(key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrKeyNotFound
	}
	return v, nil
}

func (m *mutationOnTx) GetIndexChunk(bucket, key []byte, timestamp uint64) ([]byte, error) {
	if m.db != nil {
		return m.db.GetIndexChunk(bucket, key, timestamp)
	}
	return nil, ErrKeyNotFound
}

func (m *mutationOnTx) Has(bucket, key []byte) (bool, error) {
	v, err := m.Get(bucket, key)
	if err != nil {
		return false, err
	}
	return v != nil, nil

}

func (m *mutationOnTx) DiskSize(ctx context.Context) (common.StorageSize, error) {
	if m.db == nil {
		return 0, nil
	}
	sz, err := m.db.(HasStats).DiskSize(ctx)
	if err != nil {
		return 0, err
	}
	return common.StorageSize(sz), nil
}

func (m *mutationOnTx) Put(bucket, key []byte, value []byte) error {
	m.len += uint64(len(key) + len(value))
	//fmt.Printf("Put %s\n", common.StorageSize(m.len))
	if value == nil {
		return m.cursors[string(bucket)].Delete(key)
	}
	return m.cursors[string(bucket)].Put(key, value)
}

func (m *mutationOnTx) MultiPut(tuples ...[]byte) (uint64, error) {
	panic("don't use me")
}

func (m *mutationOnTx) BatchSize() int {
	return int(m.len)
}

// IdealBatchSize defines the size of the data batches should ideally add in one write.
func (m *mutationOnTx) IdealBatchSize() int {
	return m.db.IdealBatchSize()
}

// WARNING: Merged mem/DB walk is not implemented
func (m *mutationOnTx) Walk(bucket, startkey []byte, fixedbits int, walker func([]byte, []byte) (bool, error)) error {
	m.panicOnEmptyDB()
	return m.db.Walk(bucket, startkey, fixedbits, walker)
}

// WARNING: Merged mem/DB walk is not implemented
func (m *mutationOnTx) MultiWalk(bucket []byte, startkeys [][]byte, fixedbits []int, walker func(int, []byte, []byte) error) error {
	m.panicOnEmptyDB()
	return m.db.MultiWalk(bucket, startkeys, fixedbits, walker)
}

func (m *mutationOnTx) Delete(bucket, key []byte) error {
	m.len += uint64(len(key))
	return m.cursors[string(bucket)].Delete(key)
}

func (m *mutationOnTx) Commit() (uint64, error) {
	defer func(t time.Time) { fmt.Printf("%s\n", time.Since(t)) }(time.Now())
	if m.db == nil {
		return 0, nil
	}
	if err := m.tx.Commit(context.Background()); err != nil {
		return 0, err
	}
	m.len = 0
	m.begin()
	return 0, nil
}

func (m *mutationOnTx) Rollback() {
	fmt.Printf("Rollback\n")
	m.tx.Rollback()
	m.len = 0
	m.begin()
}

func (m *mutationOnTx) Keys() ([][]byte, error) {
	panic("don't use me")
}

func (m *mutationOnTx) Close() {
	fmt.Printf("Close\n")
	m.Rollback()
	m.tx.Rollback()
}

func (m *mutationOnTx) NewBatch() DbWithPendingMutations {
	fmt.Printf("NewBatch\n")
	tx, err := m.KV().Begin(context.Background(), true)
	if err != nil {
		panic(err)
	}
	return &mutationOnTx{
		db: m.db,
		tx: tx,
	}
}

func (m *mutationOnTx) panicOnEmptyDB() {
	if m.db == nil {
		panic("Not implemented")
	}
}

func (m *mutationOnTx) MemCopy() Database {
	m.panicOnEmptyDB()
	return m.db
}

func (m *mutationOnTx) ID() uint64 {
	return m.db.ID()
}

// [TURBO-GETH] Freezer support (not implemented yet)
// Ancients returns an error as we don't have a backing chain freezer.
func (m *mutationOnTx) Ancients() (uint64, error) {
	return 0, errNotSupported
}

// TruncateAncients returns an error as we don't have a backing chain freezer.
func (m *mutationOnTx) TruncateAncients(items uint64) error {
	return errNotSupported
}
