package etl

import (
	"bytes"
	"github.com/ledgerwatch/turbo-geth/common"
	"sort"
	"strconv"
)

const (
	//SliceBuffer - just simple slice w
	SortableSliceBuffer = iota
	//SortableAppendBuffer - map[k] [v1 v2 v3]
	SortableAppendBuffer
	// SortableOldestAppearedBuffer - buffer that keeps only the oldest entries.
	// if first v1 was added under key K, then v2; only v1 will stay
	SortableOldestAppearedBuffer

	BufferOptimalSize = 256 * 1024 * 1024 /* 256 mb | var because we want to sometimes change it from tests */
	BufIOSize         = 64 * 4096         // 64 pages | default is 1 page | increasing further doesn't show speedup on SSD
)

type Buffer interface {
	Put(k, v []byte)
	Get(i int) ([]byte, []byte)
	Len() int
	Reset()
	GetEntries() [][]byte
	Sort()
	CheckFlushSize() bool
}

var (
	_ Buffer = &sortableBuffer{}
	_ Buffer = &appendSortableBuffer{}
	_ Buffer = &oldestEntrySortableBuffer{}
)

func NewSortableBuffer(bufferOptimalSize int) *sortableBuffer {
	return &sortableBuffer{
		entries:     make(map[string][]byte, 1024),
		size:        0,
		optimalSize: bufferOptimalSize,
	}
}

type sortableBuffer struct {
	sortedBuf   [][]byte
	entries     map[string][]byte
	size        int
	optimalSize int
}

func (b *sortableBuffer) Put(k, v []byte) {
	b.size += len(k) + len(v)
	b.entries[string(k)] = v
}

func (b *sortableBuffer) Size() int {
	return b.size
}

func (b *sortableBuffer) Len() int {
	return len(b.entries) / 2
}

func (b *sortableBuffer) Less(i, j int) bool {
	return bytes.Compare(b.sortedBuf[i*2], b.sortedBuf[j*2]) < 0
}

func (b *sortableBuffer) Swap(i, j int) {
	ki, kj := i*2, j*2
	b.sortedBuf[ki], b.sortedBuf[kj] = b.sortedBuf[kj], b.sortedBuf[ki]
	b.sortedBuf[ki+1], b.sortedBuf[kj+1] = b.sortedBuf[kj+1], b.sortedBuf[ki+1]
}

func (b *sortableBuffer) Get(i int) ([]byte, []byte) {
	return b.sortedBuf[i*2], b.sortedBuf[i*2+1]
}

func (b *sortableBuffer) Reset() {
	b.sortedBuf = make([][]byte, 0, 1024)
	b.entries = make(map[string][]byte, 1024)
	b.size = 0
}
func (b *sortableBuffer) Sort() {
	for k, v := range b.entries {
		b.sortedBuf = append(b.sortedBuf, []byte(k), v)
	}
	sort.Sort(b)
}

func (b *sortableBuffer) GetEntries() [][]byte {
	return b.sortedBuf
}

func (b *sortableBuffer) CheckFlushSize() bool {
	return b.size >= b.optimalSize
}

func NewAppendBuffer(bufferOptimalSize int) *appendSortableBuffer {
	return &appendSortableBuffer{
		entries:     make(map[string][]byte, 1024),
		size:        0,
		optimalSize: bufferOptimalSize,
	}
}

type appendSortableBuffer struct {
	entries     map[string][]byte
	size        int
	optimalSize int
	sortedBuf   [][]byte
}

func (b *appendSortableBuffer) Put(k, v []byte) {
	ks := string(k)
	stored, ok := b.entries[ks]
	if !ok {
		b.size += len(k)
	}
	b.size += len(v)
	stored = append(stored, v...)
	b.entries[ks] = stored
}

func (b *appendSortableBuffer) Size() int {
	return b.size
}

func (b *appendSortableBuffer) Len() int {
	return len(b.entries)
}
func (b *appendSortableBuffer) Sort() {
	for i, v := range b.entries {
		b.sortedBuf = append(b.sortedBuf, []byte(i), v)
	}
	sort.Sort(b)
}

func (b *appendSortableBuffer) Less(i, j int) bool {
	return bytes.Compare(b.sortedBuf[i*2], b.sortedBuf[j*2]) < 0
}

func (b *appendSortableBuffer) Swap(i, j int) {
	ki, kj := i*2, j*2
	b.sortedBuf[ki], b.sortedBuf[kj] = b.sortedBuf[kj], b.sortedBuf[ki]
	b.sortedBuf[ki+1], b.sortedBuf[kj+1] = b.sortedBuf[kj+1], b.sortedBuf[ki+1]
}

func (b *appendSortableBuffer) Get(i int) ([]byte, []byte) {
	return b.sortedBuf[i*2], b.sortedBuf[i*2+1]
}
func (b *appendSortableBuffer) Reset() {
	b.sortedBuf = make([][]byte, 0, 1024)
	b.entries = make(map[string][]byte, 1024)
	b.size = 0
}

func (b *appendSortableBuffer) GetEntries() [][]byte {
	return b.sortedBuf
}

func (b *appendSortableBuffer) CheckFlushSize() bool {
	return b.size >= b.optimalSize
}

func NewOldestEntryBuffer(bufferOptimalSize int) *oldestEntrySortableBuffer {
	return &oldestEntrySortableBuffer{
		entries:     make(map[string][]byte, 1024),
		size:        0,
		optimalSize: bufferOptimalSize,
	}
}

type oldestEntrySortableBuffer struct {
	entries     map[string][]byte
	size        int
	optimalSize int
	sortedBuf   [][]byte
}

func (b *oldestEntrySortableBuffer) Put(k, v []byte) {
	ks := string(k)
	_, ok := b.entries[ks]
	if ok {
		// if we already had this entry, we are going to keep it and ignore new value
		return
	}

	b.size += len(k)
	b.size += len(v)
	if v != nil {
		v = common.CopyBytes(v)
	}
	b.entries[ks] = v
}

func (b *oldestEntrySortableBuffer) Size() int {
	return b.size
}

func (b *oldestEntrySortableBuffer) Len() int {
	return len(b.entries)
}
func (b *oldestEntrySortableBuffer) Sort() {
	for k, v := range b.entries {
		b.sortedBuf = append(b.sortedBuf, []byte(k), v)
	}
	sort.Sort(b)
}

func (b *oldestEntrySortableBuffer) Less(i, j int) bool {
	return bytes.Compare(b.sortedBuf[i*2], b.sortedBuf[j*2]) < 0
}

func (b *oldestEntrySortableBuffer) Swap(i, j int) {
	ki, kj := i*2, j*2
	b.sortedBuf[ki], b.sortedBuf[kj] = b.sortedBuf[kj], b.sortedBuf[ki]
	b.sortedBuf[ki+1], b.sortedBuf[kj+1] = b.sortedBuf[kj+1], b.sortedBuf[ki+1]
}

func (b *oldestEntrySortableBuffer) Get(i int) ([]byte, []byte) {
	return b.sortedBuf[i*2], b.sortedBuf[i*2+1]
}

func (b *oldestEntrySortableBuffer) Reset() {
	b.sortedBuf = make([][]byte, 0, 1024)
	b.entries = make(map[string][]byte, 1024)
	b.size = 0
}

func (b *oldestEntrySortableBuffer) GetEntries() [][]byte {
	return b.sortedBuf
}

func (b *oldestEntrySortableBuffer) CheckFlushSize() bool {
	return b.size >= b.optimalSize
}

func getBufferByType(tp int, size int) Buffer {
	switch tp {
	case SortableSliceBuffer:
		return NewSortableBuffer(size)
	case SortableAppendBuffer:
		return NewAppendBuffer(size)
	case SortableOldestAppearedBuffer:
		return NewOldestEntryBuffer(size)
	default:
		panic("unknown buffer type " + strconv.Itoa(tp))
	}
}
