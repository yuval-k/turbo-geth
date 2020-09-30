package bitmapdb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/RoaringBitmap/gocroaring"
	"github.com/c2h5oh/datasize"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/ethdb"
)

const ColdShardLimit = 512 * datasize.KB
const HotShardLimit = 4 * datasize.KB

// AppendMergeByOr - appending delta to existing data in db, merge by Or
// Method maintains sharding - because some bitmaps are >1Mb and when new incoming blocks process it
//	 updates ~300 of bitmaps - by append small amount new values. It cause much big writes (LMDB does copy-on-write).
//
// 3 terms are used: cold_shard, hot_shard, delta
//   delta - most recent changes (appendable)
//   hot_shard - merge delta here until hot_shard size < HotShardLimit, otherwise merge hot to cold
//   cold_shard - merge hot_shard here until cold_shard size < ColdShardLimit, otherwise mark hot as cold, create new hot from delta
// if no cold shard, create new hot from delta - it will turn hot to cold automatically
// never merges cold with cold for compaction - because it's expensive operation
func AppendMergeByOr(c ethdb.Cursor, key []byte, delta *gocroaring.Bitmap) error {
	shardNForDelta := uint16(0)

	hotK, hotV, seekErr := c.Seek(key)
	if seekErr != nil {
		return seekErr
	}
	if hotK != nil && bytes.HasPrefix(hotK, key) {
		hotShardN := ^binary.BigEndian.Uint16(hotK[len(hotK)-2:])
		if len(hotV) > int(HotShardLimit) { // merge hot to cold
			var err error
			shardNForDelta, err = hotShardOverflow(c, key, hotShardN, hotV)
			if err != nil {
				return err
			}
		} else { // merge delta to hot, write result to `hotShardN`
			hot, err := gocroaring.Read(hotV)
			if err != nil {
				return err
			}

			delta = gocroaring.FastOr(delta, hot)
			shardNForDelta = hotShardN
		}
	}

	newK := make([]byte, len(key)+2)
	copy(newK, key)
	binary.BigEndian.PutUint16(newK[len(newK)-2:], ^shardNForDelta)

	//delta.RunOptimize()
	newV := make([]byte, delta.SerializedSizeInBytes())
	err := delta.Write(newV)
	if err != nil {
		return err
	}
	err = c.Put(newK, newV)
	if err != nil {
		return err
	}
	return nil
}

func hotShardOverflow(c ethdb.Cursor, initialKey []byte, hotShardN uint16, hotV []byte) (shardNForDelta uint16, err error) {
	if hotShardN == 0 { // no cold shards, create new hot from delta - it will turn hot to cold automatically
		return 1, nil
	}

	coldK, coldV, err := c.Next() // get cold shard from db
	if err != nil {
		return 0, err
	}

	if coldK == nil || !bytes.HasPrefix(coldK, initialKey) {
		return 0, fmt.Errorf("cold shard not found in db: key=%x, hotShardN=%d, found=%x", initialKey, hotShardN, coldK)
	}

	if len(coldV) > int(ColdShardLimit) { // cold shard is too big, write delta to `hotShardN + 1` - it will turn hot to cold automatically
		return hotShardN + 1, nil
	}

	// merge hot to cold and replace hot by delta (by write delta to `hotShardN`)
	cold, err := gocroaring.Read(coldV)
	if err != nil {
		return 0, err
	}

	hot, err := gocroaring.Read(hotV)
	if err != nil {
		return 0, err
	}

	cold = gocroaring.FastOr(cold, hot)

	//cold.RunOptimize()
	coldBytes := make([]byte, cold.SerializedSizeInBytes())
	err = cold.Write(coldBytes)
	if err != nil {
		return 0, err
	}
	err = c.Put(common.CopyBytes(coldK), coldBytes) // copy 'coldK' if want replace c.PutCurrent() by c.Put()
	if err != nil {
		return 0, err
	}

	return hotShardN, nil
}

// TruncateRange - gets existing bitmap in db and call RemoveRange operator on it.
// starts from hot shard, stops when shard not overlap with [from-to)
// !Important: [from, to)
func TruncateRange(c ethdb.Cursor, key []byte, from, to uint64) error {
	updated := 0
	for k, v, err := c.Seek(key); k != nil; k, v, err = c.Next() {
		if err != nil {
			return err
		}

		if !bytes.HasPrefix(k, key) {
			break
		}

		bm, err := gocroaring.Read(v)
		if err != nil {
			return err
		}
		if uint64(bm.Maximum()) < from {
			break
		}
		noReasonToCheckNextShard := uint64(bm.Minimum()) <= from && uint64(bm.Maximum()) >= to

		updated++
		bm.RemoveRange(from, to)
		if bm.GetCardinality() == 0 { // don't store empty bitmaps
			err = c.DeleteCurrent()
			if err != nil {
				return err
			}
			continue
		}

		bm.RunOptimize()
		newV := make([]byte, bm.SerializedSizeInBytes())
		err = bm.Write(newV)
		if err != nil {
			return err
		}
		err = c.Put(common.CopyBytes(k), newV)
		if err != nil {
			return err
		}

		if noReasonToCheckNextShard {
			break
		}
	}

	return nil
}

func Get(c ethdb.Cursor, key []byte) (*gocroaring.Bitmap, error) {
	var shards []*gocroaring.Bitmap
	for k, v, err := c.Seek(key); k != nil; k, v, err = c.Next() {
		if err != nil {
			return nil, err
		}

		if !bytes.HasPrefix(k, key) {
			break
		}

		bm, err := gocroaring.Read(v)
		if err != nil {
			return nil, err
		}
		shards = append(shards, bm)
	}

	if len(shards) == 0 {
		return gocroaring.New(), nil
	}
	return gocroaring.FastOr(shards...), nil
}

//const ShardLimit = 4 * datasize.KB
const ShardLimit = 3 * datasize.KB

var blockNBytes = make([]byte, 4)

func AppendMergeByOr2(c ethdb.Cursor, key []byte, delta *gocroaring.Bitmap) error {
	lastShardKey := make([]byte, len(key)+4)
	copy(lastShardKey, key)
	binary.BigEndian.PutUint32(lastShardKey[len(lastShardKey)-4:], ^uint32(0))

	currentLastV, seekErr := c.SeekExact(lastShardKey)
	if seekErr != nil {
		return seekErr
	}

	if currentLastV == nil { // no existing shards, then just create one
		err := writeBitmapSharded(c, key, delta)
		if err != nil {
			return err
		}
		return nil
	}

	last, err := gocroaring.Read(currentLastV)
	if err != nil {
		return err
	}

	if len(currentLastV) < int(ShardLimit) { // merge delta to last shard
		delta = gocroaring.Or(delta, last)

		err = writeBitmapSharded(c, key, delta)
		if err != nil {
			return err
		}

		return nil
	}

	// rename existing last shard and create new last shard
	max := last.Maximum()
	binary.BigEndian.PutUint32(blockNBytes, max)
	if err := c.Put(append(common.CopyBytes(key), blockNBytes...), currentLastV); err != nil {
		return err
	}

	err = writeBitmapSharded(c, key, delta)
	if err != nil {
		return err
	}
	return nil
}

func writeBitmapSharded(c ethdb.Cursor, key []byte, delta *gocroaring.Bitmap) error {
	shardKey := make([]byte, len(key)+4)
	copy(shardKey, key)
	sz := delta.SerializedSizeInBytes()
	if sz <= int(ShardLimit) {
		newV := make([]byte, delta.SerializedSizeInBytes())
		err := delta.Write(newV)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint32(shardKey[len(shardKey)-4:], ^uint32(0))
		err = c.Put(common.CopyBytes(shardKey), newV)
		if err != nil {
			return err
		}
		return nil
	}

	shardsAmount := uint32(sz / int(ShardLimit))
	if shardsAmount == 0 {
		shardsAmount = 1
	}
	step := (delta.Maximum() - delta.Minimum()) / shardsAmount
	step = step / 16
	shard, tmp := gocroaring.New(), gocroaring.New() // shard will write to db, tmp will use to add data to shard
	for delta.Cardinality() > 0 {
		from := uint64(delta.Minimum())
		to := from + uint64(step)
		tmp.Clear()
		tmp.AddRange(from, to)
		tmp.And(delta)
		shard.Or(tmp)
		shard.RunOptimize()
		delta.RemoveRange(from, to)
		if delta.Cardinality() == 0 {
			break
		}
		if shard.SerializedSizeInBytes() >= int(ShardLimit) {
			newV := make([]byte, shard.SerializedSizeInBytes())
			err := shard.Write(newV)
			if err != nil {
				return err
			}
			binary.BigEndian.PutUint32(shardKey[len(shardKey)-4:], shard.Maximum())

			err = c.Put(common.CopyBytes(shardKey), newV)
			if err != nil {
				return err
			}
			shard.Clear()
		}
	}

	if shard.SerializedSizeInBytes() > 0 {
		newV := make([]byte, shard.SerializedSizeInBytes())
		err := shard.Write(newV)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint32(shardKey[len(shardKey)-4:], ^uint32(0))
		err = c.Put(common.CopyBytes(shardKey), newV)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

// TruncateRange - gets existing bitmap in db and call RemoveRange operator on it.
// starts from hot shard, stops when shard not overlap with [from-to)
// !Important: [from, to)
func TruncateRange2(c ethdb.Cursor, key []byte, from, to uint64) error {
	shardKey := make([]byte, len(key)+4)
	copy(shardKey, key)
	binary.BigEndian.PutUint32(shardKey[len(shardKey)-4:], uint32(from))

	for k, v, err := c.Seek(shardKey); k != nil; k, v, err = c.Next() {
		if err != nil {
			return err
		}

		if !bytes.HasPrefix(k, key) {
			break
		}

		bm, err := gocroaring.Read(v)
		if err != nil {
			return err
		}
		noReasonToCheckNextShard := (uint64(bm.Minimum()) <= from && uint64(bm.Maximum()) >= to) || binary.BigEndian.Uint32(k[len(k)-4:]) == ^uint32(0)

		bm.RemoveRange(from, to)
		if bm.GetCardinality() == 0 { // don't store empty bitmaps
			err = c.DeleteCurrent()
			if err != nil {
				return err
			}
			if noReasonToCheckNextShard {
				break
			}
			continue
		}

		bm.RunOptimize()
		newV := make([]byte, bm.SerializedSizeInBytes())
		err = bm.Write(newV)
		if err != nil {
			return err
		}
		err = c.Put(common.CopyBytes(k), newV)
		if err != nil {
			return err
		}

		if noReasonToCheckNextShard {
			break
		}
	}

	// rename last shard
	k, v, err := c.Current()
	if err != nil {
		return err
	}
	if k == nil { // if last shard was deleted, do 1 step back
		k, v, err = c.Prev()
		if err != nil {
			return err
		}
	}

	if binary.BigEndian.Uint32(k[len(k)-4:]) == ^uint32(0) { // nothing to return
		return nil
	}
	if !bytes.HasPrefix(k, key) {
		return nil
	}

	copyV := common.CopyBytes(v)
	err = c.DeleteCurrent()
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(shardKey[len(shardKey)-4:], ^uint32(0))
	err = c.Put(shardKey, copyV)
	if err != nil {
		return err
	}

	return nil
}

func Get2(c ethdb.Cursor, key []byte, from, to uint32) (*gocroaring.Bitmap, error) {
	var shards []*gocroaring.Bitmap

	fromKey := make([]byte, len(key)+4)
	copy(fromKey, key)
	binary.BigEndian.PutUint32(fromKey[len(fromKey)-4:], from)
	for k, v, err := c.Seek(fromKey); k != nil; k, v, err = c.Next() {
		if err != nil {
			return nil, err
		}

		if !bytes.HasPrefix(k, key) {
			break
		}

		bm, err := gocroaring.Read(v)
		if err != nil {
			return nil, err
		}
		shards = append(shards, bm)

		if binary.BigEndian.Uint32(k[len(k)-4:]) >= to {
			break
		}
	}

	if len(shards) == 0 {
		return gocroaring.New(), nil
	}
	return gocroaring.FastOr(shards...), nil
}
