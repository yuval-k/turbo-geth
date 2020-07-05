package trie

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/core/types/accounts"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/metrics"
	"github.com/ledgerwatch/turbo-geth/trie/rlphacks"
)

var (
	trieFlatDbSubTrieLoaderTimer = metrics.NewRegisteredTimer("trie/subtrieloader/flatdb", nil)
)

type StreamReceiver interface {
	Receive(
		itemType StreamItem,
		accountKey []byte,
		storageKey []byte,
		accountValue *accounts.Account,
		storageValue []byte,
		hash []byte,
		cutoff int,
	) error

	Result() SubTries
}

type FlatDbSubTrieLoader struct {
	trace              bool
	rl                 RetainDecider
	rangeIdx           int
	accAddrHashWithInc [40]byte // Concatenation of addrHash of the currently build account with its incarnation encoding
	dbPrefixes         [][]byte
	fixedbytes         []int
	masks              []byte
	cutoffs            []int
	kv                 ethdb.KV
	nextAccountKey     [32]byte
	k, v               []byte
	ihK, ihV           []byte
	minKeyAsNibbles    []byte

	itemPresent bool
	itemType    StreamItem

	// Storage item buffer
	storageKey   []byte
	storageValue []byte

	// Acount item buffer
	accountKey   []byte
	accountValue accounts.Account
	hashValue    []byte
	streamCutoff int

	receiver         StreamReceiver
	defaultReceiver  *DefaultReceiver
	hc               HashCollector
	t2               *Trie2
	seekIHCounter    int
	seekIHCounter2   int
	seekIHCounter3   int
	seekIHCounter4   int
	seekIHCounter5   int
	seekIHCounter6   int
	seekIHCounter7   int
	seekCounter      int
	seekCounter1     int
	seekCounter2     int
	seekCounter3     int
	seekCounter4     int
	seekCounter5     int
	seekCounterSkip5 int
	seekCounterSkip6 int
	seekCounter6     int
	seekCounter7     int
	nextCounter      int
	nextIHCounter    int
}

type DefaultReceiver struct {
	trace        bool
	rl           RetainDecider
	hc           HashCollector
	subTries     SubTries
	currStorage  bytes.Buffer // Current key for the structure generation algorithm, as well as the input tape for the hash builder
	succStorage  bytes.Buffer
	valueStorage bytes.Buffer // Current value to be used as the value tape for the hash builder
	curr         bytes.Buffer // Current key for the structure generation algorithm, as well as the input tape for the hash builder
	succ         bytes.Buffer
	value        bytes.Buffer // Current value to be used as the value tape for the hash builder
	groups       []uint16
	hb           *HashBuilder
	trie2        *Trie2
	wasIH        bool
	wasIHStorage bool
	hashData     GenStructStepHashData
	a            accounts.Account
	leafData     GenStructStepLeafData
	accData      GenStructStepAccountData
}

func NewDefaultReceiver() *DefaultReceiver {
	return &DefaultReceiver{hb: NewHashBuilder(false)}
}

func NewFlatDbSubTrieLoader() *FlatDbSubTrieLoader {
	fstl := &FlatDbSubTrieLoader{
		defaultReceiver: NewDefaultReceiver(),
	}
	return fstl
}

// Reset prepares the loader for reuse
func (fstl *FlatDbSubTrieLoader) Reset(db ethdb.Database, trie2 *Trie2, rl RetainDecider, receiverDecider RetainDecider, hc HashCollector, dbPrefixes [][]byte, fixedbits []int, trace bool) error {
	fstl.defaultReceiver.Reset(receiverDecider, hc, trie2, trace)
	fstl.hc = trie2.wrapHashCollector(hc)
	fstl.receiver = fstl.defaultReceiver
	fstl.rangeIdx = 0

	fstl.minKeyAsNibbles = fstl.minKeyAsNibbles[:0]
	fstl.trace = trace
	fstl.rl = rl
	fstl.t2 = trie2
	fstl.t2.Reset()
	fstl.dbPrefixes = dbPrefixes
	fstl.itemPresent = false
	if fstl.trace {
		fmt.Printf("----------\n")
		fmt.Printf("RebuildTrie\n")
	}
	if fstl.trace {
		fmt.Printf("fstl.rl: %s\n", fstl.rl)
		fmt.Printf("fixedbits: %d\n", fixedbits)
		fmt.Printf("dbPrefixes(%d): %x\n", len(dbPrefixes), dbPrefixes)
	}
	if len(dbPrefixes) == 0 {
		return nil
	}
	if hasKV, ok := db.(ethdb.HasKV); ok {
		fstl.kv = hasKV.KV()
	} else {
		return fmt.Errorf("database doest not implement KV: %T", db)
	}
	fixedbytes := make([]int, len(fixedbits))
	masks := make([]byte, len(fixedbits))
	cutoffs := make([]int, len(fixedbits))
	for i, bits := range fixedbits {
		cutoffs[i] = bits / 4
		fixedbytes[i], masks[i] = ethdb.Bytesmask(bits)
	}
	fstl.fixedbytes = fixedbytes
	fstl.masks = masks
	fstl.cutoffs = cutoffs

	fstl.seekIHCounter = 0
	fstl.seekIHCounter2 = 0
	fstl.seekIHCounter3 = 0
	fstl.seekIHCounter4 = 0
	fstl.seekIHCounter5 = 0
	fstl.seekCounter1 = 0
	fstl.seekCounter2 = 0
	fstl.seekCounter3 = 0
	fstl.seekCounter4 = 0
	fstl.seekCounter5 = 0
	fstl.seekCounter6 = 0
	fstl.seekCounter7 = 0
	fstl.nextCounter = 0
	fstl.nextIHCounter = 0
	return nil
}

func (fstl *FlatDbSubTrieLoader) SetStreamReceiver(receiver StreamReceiver) {
	fstl.receiver = receiver
}

type cursor interface {
	SeekTo(seek []byte) ([]byte, []byte, error)
	Next() ([]byte, []byte, error)
}

// iteration moves through the database buckets and creates at most
// one stream item, which is indicated by setting the field fstl.itemPresent to true
func (fstl *FlatDbSubTrieLoader) iteration(c ethdb.Cursor, ih *IHCursor, first bool) error {
	var isIH, isIHSequence bool
	var minKey []byte
	var err error
	if !first {
		isIH, minKey = keyIsBeforeOrEqual(fstl.ihK, fstl.k)
	}
	fixedbytes := fstl.fixedbytes[fstl.rangeIdx]
	cutoff := fstl.cutoffs[fstl.rangeIdx]
	dbPrefix := fstl.dbPrefixes[fstl.rangeIdx]
	mask := fstl.masks[fstl.rangeIdx]
	// Adjust rangeIdx if needed
	var cmp int = -1
	for cmp != 0 {
		if minKey == nil {
			if !first {
				cmp = 1
			}
		} else if fixedbytes > 0 { // In the first iteration, we do not have valid minKey, so we skip this part
			if len(minKey) < fixedbytes {
				cmp = bytes.Compare(minKey, dbPrefix[:len(minKey)])
				if cmp == 0 {
					cmp = -1
				}
			} else {
				cmp = bytes.Compare(minKey[:fixedbytes-1], dbPrefix[:fixedbytes-1])
				if cmp == 0 {
					k1 := minKey[fixedbytes-1] & mask
					k2 := dbPrefix[fixedbytes-1] & mask
					if k1 < k2 {
						cmp = -1
					} else if k1 > k2 {
						cmp = 1
					}
				}
			}
		} else {
			cmp = 0
		}
		if fstl.trace {
			fmt.Printf("minKey %x, dbPrefix %x, cmp %d, fstl.rangeIdx %d, %x\n", minKey, dbPrefix, cmp, fstl.rangeIdx, fstl.dbPrefixes)
		}
		if cmp == 0 && fstl.itemPresent {
			return nil
		}
		if cmp < 0 {
			// This happens after we have just incremented rangeIdx or on the very first iteration
			if first && len(dbPrefix) > common.HashLength {
				// Looking for storage sub-tree
				copy(fstl.accAddrHashWithInc[:], dbPrefix[:common.HashLength+common.IncarnationLength])
			}
			fstl.seekIHCounter++
			if fstl.ihK, fstl.ihV, isIHSequence, err = ih.SeekTo(dbPrefix); err != nil {
				return err
			}
			if len(dbPrefix) <= common.HashLength && len(fstl.ihK) > common.HashLength {
				// Advance to the first account
				fstl.seekIHCounter2++
				if nextAccount(fstl.ihK, fstl.nextAccountKey[:]) {
					if fstl.ihK, fstl.ihV, isIHSequence, err = ih.SeekTo(fstl.nextAccountKey[:]); err != nil {
						return err
					}
				} else {
					fstl.ihK = nil
				}
			}

			if isIHSequence {
				fstl.k = common.CopyBytes(fstl.ihK)
			} else {
				fstl.seekCounter++
				if fstl.k, fstl.v, err = c.SeekTo(dbPrefix); err != nil {
					return err
				}
				if len(dbPrefix) <= common.HashLength && len(fstl.k) > common.HashLength {
					// Advance past the storage to the first account
					if nextAccount(fstl.k, fstl.nextAccountKey[:]) {
						fstl.seekCounter2++
						if fstl.k, fstl.v, err = c.SeekTo(fstl.nextAccountKey[:]); err != nil {
							return err
						}
					} else {
						fstl.k = nil
					}
				}
			}

			isIH, minKey = keyIsBeforeOrEqual(fstl.ihK, fstl.k)
			if fixedbytes == 0 {
				cmp = 0
			}
		} else if cmp > 0 {
			if !first {
				fstl.rangeIdx++
			}
			if !first {
				fstl.itemPresent = true
				fstl.itemType = CutoffStreamItem
				fstl.streamCutoff = cutoff
				fstl.accountKey = nil
				fstl.storageKey = nil
				fstl.storageValue = nil
				fstl.hashValue = nil
				if fstl.trace {
					fmt.Printf("Inserting cutoff %d\n", cutoff)
				}
			}
			if fstl.rangeIdx == len(fstl.dbPrefixes) {
				return nil
			}
			fixedbytes = fstl.fixedbytes[fstl.rangeIdx]
			mask = fstl.masks[fstl.rangeIdx]
			dbPrefix = fstl.dbPrefixes[fstl.rangeIdx]
			if len(dbPrefix) > common.HashLength {
				// Looking for storage sub-tree
				copy(fstl.accAddrHashWithInc[:], dbPrefix[:common.HashLength+common.IncarnationLength])
			}
		}
	}

	if !isIH {
		if len(fstl.k) > common.HashLength && !bytes.HasPrefix(fstl.k, fstl.accAddrHashWithInc[:]) {
			if bytes.Compare(fstl.k, fstl.accAddrHashWithInc[:]) < 0 {
				// Skip all the irrelevant storage in the middle
				fstl.seekCounter3++
				if fstl.k, fstl.v, err = c.SeekTo(fstl.accAddrHashWithInc[:]); err != nil {
					return err
				}
			} else {
				if nextAccount(fstl.k, fstl.nextAccountKey[:]) {
					fstl.seekCounter4++
					if fstl.k, fstl.v, err = c.SeekTo(fstl.nextAccountKey[:]); err != nil {
						return err
					}
				} else {
					fstl.k = nil
				}
			}
			return nil
		}
		fstl.itemPresent = true
		if len(fstl.k) > common.HashLength {
			fstl.itemType = StorageStreamItem
			fstl.accountKey = nil
			fstl.storageKey = append(fstl.storageKey[:0], fstl.k...)
			fstl.hashValue = nil
			fstl.storageValue = append(fstl.storageValue[:0], fstl.v...)
			fstl.nextCounter++
			if fstl.k, fstl.v, err = c.Next(); err != nil {
				return err
			}
			if fstl.trace {
				fmt.Printf("k after storageWalker and Next: %x\n", fstl.k)
			}
		} else {
			fstl.itemType = AccountStreamItem
			fstl.accountKey = append(fstl.accountKey[:0], fstl.k...)
			fstl.storageKey = nil
			fstl.storageValue = nil
			fstl.hashValue = nil
			if err := fstl.accountValue.DecodeForStorage(fstl.v); err != nil {
				return fmt.Errorf("fail DecodeForStorage: %w", err)
			}
			copy(fstl.accAddrHashWithInc[:], fstl.k)
			binary.BigEndian.PutUint64(fstl.accAddrHashWithInc[32:], ^fstl.accountValue.Incarnation)

			if keyIsBefore(fstl.ihK, fstl.accAddrHashWithInc[:]) {
				fstl.seekIHCounter3++
				if fstl.ihK, fstl.ihV, isIHSequence, err = ih.SeekTo(fstl.accAddrHashWithInc[:]); err != nil {
					return err
				}
				if isIHSequence {
					// If can use next IH, then no reason to move `c` cursor.
					// But move .k forward - then `keyIsBeforeOrEqual` of new iteration will chose ihK.
					// .v not expected been used.
					fstl.k = common.CopyBytes(fstl.accAddrHashWithInc[:])
					fstl.seekCounterSkip5++
					return nil
				}
			}

			// Now we know the correct incarnation of the account, and we can skip all irrelevant storage records
			// Since 0 incarnation if 0xfff...fff, and we do not expect any records like that, this automatically
			// skips over all storage items
			if fstl.k, fstl.v, err = c.SeekTo(fstl.accAddrHashWithInc[:]); err != nil {
				return err
			}
			fstl.seekCounter5++
			if fstl.trace {
				fmt.Printf("k after accountWalker and SeekTo: %x\n", fstl.k)
			}
		}
		return nil
	}

	// ih part
	fstl.itemPresent = true
	if len(fstl.ihK) > common.HashLength {
		fstl.itemType = SHashStreamItem
		fstl.accountKey = nil
		fstl.storageKey = append(fstl.storageKey[:0], fstl.ihK...)
		fstl.hashValue = append(fstl.hashValue[:0], fstl.ihV...)
		fstl.storageValue = nil
	} else {
		fstl.itemType = AHashStreamItem
		fstl.accountKey = append(fstl.accountKey[:0], fstl.ihK...)
		fstl.storageKey = nil
		fstl.storageValue = nil
		fstl.hashValue = append(fstl.hashValue[:0], fstl.ihV...)
	}

	// skip subtree
	next, ok := dbutils.NextSubtree(fstl.ihK)
	if !ok { // no siblings left
		fstl.k, fstl.ihK, fstl.ihV = nil, nil, nil
		return nil
	}
	if fstl.trace {
		fmt.Printf("next: %x\n", next)
	}
	fstl.seekIHCounter6++

	if fstl.ihK, fstl.ihV, isIHSequence, err = ih.SeekTo(next); err != nil {
		return err
	}

	// If can use next IH, then no reason to move `c` cursor.
	// But move .k forward - then `keyIsBeforeOrEqual` of new iteration will chose ihK.
	// .v not expected been used.
	if isIHSequence {
		// if IH is prefix of next sub-trie, then IH will used on next iteration
		fstl.k = common.CopyBytes(fstl.ihK)
		return nil
	}

	if keyIsBefore(fstl.k, next) {
		fstl.seekCounter6++
		if fstl.k, fstl.v, err = c.SeekTo(next); err != nil {
			return err
		}
	}
	if len(next) <= common.HashLength && len(fstl.k) > common.HashLength {
		// Advance past the storage to the first account
		if nextAccount(fstl.k, fstl.nextAccountKey[:]) {
			if fstl.k, fstl.v, err = c.SeekTo(fstl.nextAccountKey[:]); err != nil {
				return err
			}
		} else {
			fstl.k = nil
		}
	}
	return nil
}

func (dr *DefaultReceiver) Reset(rl RetainDecider, hc HashCollector, trie2 *Trie2, trace bool) {
	dr.rl = rl
	dr.trie2 = trie2
	dr.hc = trie2.wrapHashCollector(hc)
	dr.curr.Reset()
	dr.succ.Reset()
	dr.value.Reset()
	dr.groups = dr.groups[:0]
	dr.a.Reset()
	dr.hb.Reset()
	dr.wasIH = false
	dr.currStorage.Reset()
	dr.succStorage.Reset()
	dr.valueStorage.Reset()
	dr.wasIHStorage = false
	dr.subTries = SubTries{}
	dr.trace = trace
	dr.hb.trace = trace
}

func (dr *DefaultReceiver) Receive(itemType StreamItem,
	accountKey []byte,
	storageKey []byte,
	accountValue *accounts.Account,
	storageValue []byte,
	hash []byte,
	cutoff int,
) error {
	switch itemType {
	case StorageStreamItem:
		dr.advanceKeysStorage(storageKey, true /* terminator */)
		if dr.currStorage.Len() > 0 {
			if err := dr.genStructStorage(); err != nil {
				return err
			}
		}
		dr.saveValueStorage(false, storageValue, hash)
	case SHashStreamItem:
		dr.advanceKeysStorage(storageKey, false /* terminator */)
		if dr.currStorage.Len() > 0 {
			if err := dr.genStructStorage(); err != nil {
				return err
			}
		}
		dr.saveValueStorage(true, storageValue, hash)
	case AccountStreamItem:
		dr.advanceKeysAccount(accountKey, true /* terminator */)
		if dr.curr.Len() > 0 && !dr.wasIH {
			dr.cutoffKeysStorage(2 * (common.HashLength + common.IncarnationLength))
			if dr.currStorage.Len() > 0 {
				if err := dr.genStructStorage(); err != nil {
					return err
				}
			}
			if dr.currStorage.Len() > 0 {
				if len(dr.groups) >= 2*common.HashLength {
					dr.groups = dr.groups[:2*common.HashLength-1]
				}
				for len(dr.groups) > 0 && dr.groups[len(dr.groups)-1] == 0 {
					dr.groups = dr.groups[:len(dr.groups)-1]
				}
				dr.currStorage.Reset()
				dr.succStorage.Reset()
				dr.wasIHStorage = false
				// There are some storage items
				dr.accData.FieldSet |= AccountFieldStorageOnly
			}
		}
		if dr.curr.Len() > 0 {
			if err := dr.genStructAccount(); err != nil {
				return err
			}
		}
		if err := dr.saveValueAccount(false, accountValue, hash); err != nil {
			return err
		}
	case AHashStreamItem:
		dr.advanceKeysAccount(accountKey, false /* terminator */)
		if dr.curr.Len() > 0 && !dr.wasIH {
			dr.cutoffKeysStorage(2 * (common.HashLength + common.IncarnationLength))
			if dr.currStorage.Len() > 0 {
				if err := dr.genStructStorage(); err != nil {
					return err
				}
			}
			if dr.currStorage.Len() > 0 {
				if len(dr.groups) >= 2*common.HashLength {
					dr.groups = dr.groups[:2*common.HashLength-1]
				}
				for len(dr.groups) > 0 && dr.groups[len(dr.groups)-1] == 0 {
					dr.groups = dr.groups[:len(dr.groups)-1]
				}
				dr.currStorage.Reset()
				dr.succStorage.Reset()
				dr.wasIHStorage = false
				// There are some storage items
				dr.accData.FieldSet |= AccountFieldStorageOnly
			}
		}
		if dr.curr.Len() > 0 {
			if err := dr.genStructAccount(); err != nil {
				return err
			}
		}
		if err := dr.saveValueAccount(true, accountValue, hash); err != nil {
			return err
		}
	case CutoffStreamItem:
		if dr.trace {
			fmt.Printf("storage cuttoff %d\n", cutoff)
		}
		if cutoff >= 2*(common.HashLength+common.IncarnationLength) {
			dr.cutoffKeysStorage(cutoff)
			if dr.currStorage.Len() > 0 {
				if err := dr.genStructStorage(); err != nil {
					return err
				}
			}
			if dr.currStorage.Len() > 0 {
				if len(dr.groups) >= cutoff {
					dr.groups = dr.groups[:cutoff-1]
				}
				for len(dr.groups) > 0 && dr.groups[len(dr.groups)-1] == 0 {
					dr.groups = dr.groups[:len(dr.groups)-1]
				}
				dr.currStorage.Reset()
				dr.succStorage.Reset()
				dr.wasIHStorage = false
				dr.subTries.roots = append(dr.subTries.roots, dr.hb.root())
				dr.subTries.Hashes = append(dr.subTries.Hashes, dr.hb.rootHash())
			} else {
				dr.subTries.roots = append(dr.subTries.roots, nil)
				dr.subTries.Hashes = append(dr.subTries.Hashes, common.Hash{})
			}
		} else {
			dr.cutoffKeysAccount(cutoff)
			if dr.curr.Len() > 0 && !dr.wasIH {
				dr.cutoffKeysStorage(2 * (common.HashLength + common.IncarnationLength))
				if dr.currStorage.Len() > 0 {
					if err := dr.genStructStorage(); err != nil {
						return err
					}
				}
				if dr.currStorage.Len() > 0 {
					if len(dr.groups) >= 2*common.HashLength {
						dr.groups = dr.groups[:2*common.HashLength-1]
					}
					for len(dr.groups) > 0 && dr.groups[len(dr.groups)-1] == 0 {
						dr.groups = dr.groups[:len(dr.groups)-1]
					}
					dr.currStorage.Reset()
					dr.succStorage.Reset()
					dr.wasIHStorage = false
					// There are some storage items
					dr.accData.FieldSet |= AccountFieldStorageOnly
				}
			}
			if dr.curr.Len() > 0 {
				if err := dr.genStructAccount(); err != nil {
					return err
				}
			}
			if dr.curr.Len() > 0 {
				if len(dr.groups) > cutoff {
					dr.groups = dr.groups[:cutoff]
				}
				for len(dr.groups) > 0 && dr.groups[len(dr.groups)-1] == 0 {
					dr.groups = dr.groups[:len(dr.groups)-1]
				}
			}
			dr.subTries.roots = append(dr.subTries.roots, dr.hb.root())
			dr.subTries.Hashes = append(dr.subTries.Hashes, dr.hb.rootHash())
			dr.groups = dr.groups[:0]
			dr.hb.Reset()
			dr.wasIH = false
			dr.wasIHStorage = false
			dr.curr.Reset()
			dr.succ.Reset()
			dr.currStorage.Reset()
			dr.succStorage.Reset()
		}
	}
	return nil
}

func (dr *DefaultReceiver) Result() SubTries {
	return dr.subTries
}

func (fstl *FlatDbSubTrieLoader) LoadSubTries() (SubTries, error) {
	defer trieFlatDbSubTrieLoaderTimer.UpdateSince(time.Now())
	if len(fstl.dbPrefixes) == 0 {
		return SubTries{}, nil
	}

	if rl, ok := fstl.rl.(*RetainList); ok {
		fmt.Printf("RetainList size: %dK\n", len(rl.hexes)/1000)
	}
	if err := fstl.kv.View(context.Background(), func(tx ethdb.Tx) error {
		c := tx.Bucket(dbutils.CurrentStateBucket).Cursor()
		var filter = func(k []byte) bool {

			if fstl.rl.Retain(k) {
				if fstl.hc != nil {
					_ = fstl.hc(k, nil)
				}
				return false
			}

			if len(k) < fstl.cutoffs[fstl.rangeIdx] {
				return false
			}

			return true
		}

		ihc := Filter(filter, tx.Bucket(dbutils.IntermediateTrieHashBucket).Cursor())
		t2c := Filter(filter, fstl.t2)
		ih := IH(t2c, ihc)

		//sz, _ := tx.Bucket(dbutils.IntermediateTrieHashBucket).Size()
		//ii := 0
		//total := 0
		//agg := map[string]int{}
		//_ = tx.Bucket(dbutils.IntermediateTrieHashBucket).Cursor().Walk(func(k, v []byte) (bool, error) {
		//	total++
		//	agg[string(v)]++
		//	if len(v) == 0 {
		//		ii++
		//	}
		//	return true, nil
		//})

		//fmt.Printf("DB IH size: %dMb, total: %d, empty roots: %d\n", sz/1024/1024, total, ii)
		//kk := 0
		//for k := range agg {
		//	if agg[k] > 5_000 {
		//		kk++
		//		fmt.Printf("Often hashes example: %d, %x\n", agg[k], k)
		//	}
		//}
		//fmt.Printf("Often hashes: %d\n", kk)

		i := 1
		if err := fstl.iteration(c, ih, true /* first */); err != nil {
			return err
		}
		for fstl.rangeIdx < len(fstl.dbPrefixes) {
			for !fstl.itemPresent {
				i++
				if err := fstl.iteration(c, ih, false /* first */); err != nil {
					return err
				}
			}
			if fstl.itemPresent {
				if err := fstl.receiver.Receive(fstl.itemType, fstl.accountKey, fstl.storageKey, &fstl.accountValue, fstl.storageValue, fstl.hashValue, fstl.streamCutoff); err != nil {
					return err
				}
				fstl.itemPresent = false
			}
		}

		fmt.Printf(".iteration() calls: %d\n", i)
		fmt.Printf("ih.SeekTo called: 1=%d, 2=%d, 3=%d, 4=%d, 5=%d, 6=%d, 7=%d\n", fstl.seekIHCounter, fstl.seekIHCounter2, fstl.seekIHCounter3, fstl.seekIHCounter4, fstl.seekIHCounter5, fstl.seekIHCounter6, fstl.seekIHCounter7)
		fmt.Printf("c.SeekTo called: 1=%d, 2=%d,3=%d,4=%d,5=%d,6=%d,7=%d\n", fstl.seekCounter, fstl.seekCounter2, fstl.seekCounter3, fstl.seekCounter4, fstl.seekCounter5, fstl.seekCounter6, fstl.seekCounter7)
		fmt.Printf("c.Next called: 1=%d, counterSkip5=%d, counterSkip6=%d\n", fstl.nextCounter, fstl.seekCounterSkip5, fstl.seekCounterSkip6)
		fmt.Printf("2As1: skipSeek2Counter=%d, db.seekAccCouner=%d, db.seekStorageCouner=%d, db.nextCounter=%d\n", ih.skipSeek2Counter, ihc.seekAccCouner, ihc.seekStorageCouner, ihc.nextCounter)
		return nil
	}); err != nil {
		return SubTries{}, err
	}
	return fstl.receiver.Result(), nil
}

func (fstl *FlatDbSubTrieLoader) AttachRequestedCode(db ethdb.Getter, requests []*LoadRequestForCode) error {
	for _, req := range requests {
		codeHash := req.codeHash
		code, err := db.Get(dbutils.CodeBucket, codeHash[:])
		if err != nil {
			return err
		}
		if req.bytecode {
			if err := req.t.UpdateAccountCode(req.addrHash[:], codeNode(code)); err != nil {
				return err
			}
		} else {
			if err := req.t.UpdateAccountCodeSize(req.addrHash[:], len(code)); err != nil {
				return err
			}
		}
	}
	return nil
}

func keyToNibbles(k []byte, w io.ByteWriter) {
	for _, b := range k {
		//nolint:errcheck
		w.WriteByte(b / 16)
		//nolint:errcheck
		w.WriteByte(b % 16)
	}
}

func (dr *DefaultReceiver) advanceKeysStorage(k []byte, terminator bool) {
	dr.currStorage.Reset()
	dr.currStorage.Write(dr.succStorage.Bytes())
	dr.succStorage.Reset()
	// Transform k to nibbles, but skip the incarnation part in the middle
	keyToNibbles(k, &dr.succStorage)

	if terminator {
		dr.succStorage.WriteByte(16)
	}
}

func (dr *DefaultReceiver) cutoffKeysStorage(cutoff int) {
	dr.currStorage.Reset()
	dr.currStorage.Write(dr.succStorage.Bytes())
	dr.succStorage.Reset()
	if dr.currStorage.Len() > 0 {
		dr.succStorage.Write(dr.currStorage.Bytes()[:cutoff-1])
		dr.succStorage.WriteByte(dr.currStorage.Bytes()[cutoff-1] + 1) // Modify last nibble in the incarnation part of the `currStorage`
	}
}

func (dr *DefaultReceiver) genStructStorage() error {
	var err error
	var data GenStructStepData
	if dr.wasIHStorage {
		dr.hashData.Hash = common.BytesToHash(dr.valueStorage.Bytes())
		data = &dr.hashData
	} else {
		dr.leafData.Value = rlphacks.RlpSerializableBytes(dr.valueStorage.Bytes())
		data = &dr.leafData
	}
	dr.groups, err = GenStructStep(dr.rl.Retain, dr.currStorage.Bytes(), dr.succStorage.Bytes(), dr.hb, dr.hc, data, dr.groups, dr.trace)
	if err != nil {
		return err
	}
	return nil
}

func (dr *DefaultReceiver) saveValueStorage(isIH bool, v, h []byte) {
	// Remember the current value
	dr.wasIHStorage = isIH
	dr.valueStorage.Reset()
	if isIH {
		dr.valueStorage.Write(h)
	} else {
		dr.valueStorage.Write(v)
	}
}

func (dr *DefaultReceiver) advanceKeysAccount(k []byte, terminator bool) {
	dr.curr.Reset()
	dr.curr.Write(dr.succ.Bytes())
	dr.succ.Reset()
	for _, b := range k {
		dr.succ.WriteByte(b / 16)
		dr.succ.WriteByte(b % 16)
	}
	if terminator {
		dr.succ.WriteByte(16)
	}
}

func (dr *DefaultReceiver) cutoffKeysAccount(cutoff int) {
	dr.curr.Reset()
	dr.curr.Write(dr.succ.Bytes())
	dr.succ.Reset()
	if dr.curr.Len() > 0 && cutoff > 0 {
		dr.succ.Write(dr.curr.Bytes()[:cutoff-1])
		dr.succ.WriteByte(dr.curr.Bytes()[cutoff-1] + 1) // Modify last nibble before the cutoff point
	}
}

func (dr *DefaultReceiver) genStructAccount() error {
	var data GenStructStepData
	if dr.wasIH {
		copy(dr.hashData.Hash[:], dr.value.Bytes())
		data = &dr.hashData
	} else {
		dr.accData.Balance.Set(&dr.a.Balance)
		if dr.a.Balance.Sign() != 0 {
			dr.accData.FieldSet |= AccountFieldBalanceOnly
		}
		dr.accData.Nonce = dr.a.Nonce
		if dr.a.Nonce != 0 {
			dr.accData.FieldSet |= AccountFieldNonceOnly
		}
		dr.accData.Incarnation = dr.a.Incarnation
		data = &dr.accData
	}
	dr.wasIHStorage = false
	dr.currStorage.Reset()
	dr.succStorage.Reset()
	var err error
	if dr.groups, err = GenStructStep(dr.rl.Retain, dr.curr.Bytes(), dr.succ.Bytes(), dr.hb, dr.hc, data, dr.groups, dr.trace); err != nil {
		return err
	}
	dr.accData.FieldSet = 0
	return nil
}

func (dr *DefaultReceiver) saveValueAccount(isIH bool, v *accounts.Account, h []byte) error {
	dr.wasIH = isIH
	if isIH {
		dr.value.Reset()
		dr.value.Write(h)
		return nil
	}
	dr.a.Copy(v)
	// Place code on the stack first, the storage will follow
	if !dr.a.IsEmptyCodeHash() {
		// the first item ends up deepest on the stack, the second item - on the top
		dr.accData.FieldSet |= AccountFieldCodeOnly
		if err := dr.hb.hash(dr.a.CodeHash[:]); err != nil {
			return err
		}
	}
	return nil
}

func nextAccount(in, out []byte) bool {
	copy(out, in)
	for i := len(out) - 1; i >= 0; i-- {
		if out[i] != 255 {
			out[i]++
			return true
		}
		out[i] = 0
	}
	return false
}

// keyIsBefore - kind of bytes.Compare, but nil is the last key. And return
func keyIsBeforeOrEqual(k1, k2 []byte) (bool, []byte) {
	if k1 == nil {
		return false, k2
	}

	if k2 == nil {
		return true, k1
	}

	if bytes.Compare(k1, k2) <= 0 {
		return true, k1
	}
	return false, k2
}

// keyIsBefore - kind of bytes.Compare, but nil is the last key. And return
func keyIsBefore(k1, k2 []byte) bool {
	if k1 == nil {
		return false
	}

	if k2 == nil {
		return true
	}

	return bytes.Compare(k1, k2) < 0
}
