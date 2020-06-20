package trie

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/ledgerwatch/turbo-geth/common"
)

type HashNodeFunc func(node, bool, []byte) (int, error)

type MerklePathLimiter struct {
	RetainDecider RetainDecider
	HashFunc      HashNodeFunc
}

type WitnessBuilder struct {
	root     node
	trace    bool
	operands []WitnessOperator
}

func NewWitnessBuilder(root node, trace bool) *WitnessBuilder {
	return &WitnessBuilder{
		root:     root,
		trace:    trace,
		operands: make([]WitnessOperator, 0),
	}
}

func (b *WitnessBuilder) Build(limiter *MerklePathLimiter) (*Witness, error) {
	err := b.makeBlockWitness(b.root, []byte{}, limiter, true)
	sort.Slice(b.operands, func(i, j int) bool {
		keyI := b.operands[i].GetKey()
		keyJ := b.operands[j].GetKey()
		for z := 0; z < len(keyI); z++ {
			if z >= len(keyJ) {
				return false
			}
			if keyI[z] < keyJ[z] {
				return true
			}
			if keyI[z] > keyJ[z] {
				return false
			}
		}
		// check operands, code should always go last
		opI := b.operands[i]
		opJ := b.operands[j]

		if _, ok := opI.(*OperatorCode); ok {
			return false
		}

		if _, ok := opJ.(*OperatorCode); ok {
			return true
		}

		return true
	})
	/*
		for i, o := range b.operands {
			fmt.Printf("%d. K=%x\n", i, o.GetKey())
		}
	*/
	witness := NewWitness(b.operands)
	b.operands = nil
	return witness, err
}

func (b *WitnessBuilder) addLeafOp(key []byte, value []byte) error {
	if b.trace {
		fmt.Printf("LEAF_VALUE: k %x v:%x\n", key, value)
	}

	var op OperatorLeafValue

	op.Key = make([]byte, len(key))
	copy(op.Key[:], key)
	if value != nil {
		op.Value = make([]byte, len(value))
		copy(op.Value[:], value)
	}

	b.operands = append(b.operands, &op)
	return nil
}

func (b *WitnessBuilder) addAccountLeafOp(key []byte, accountNode *accountNode) error {
	if b.trace {
		fmt.Printf("LEAF_ACCOUNT: k %x acc:%x\n", key, accountNode.Hash())
	}

	var op OperatorLeafAccount
	op.Key = make([]byte, len(key))
	copy(op.Key[:], key)

	op.Nonce = accountNode.Nonce
	op.Balance = big.NewInt(0)
	op.Balance.SetBytes(accountNode.Balance.Bytes())
	copy(op.Root[:], accountNode.Root[:])
	copy(op.CodeHash[:], accountNode.CodeHash[:])

	b.operands = append(b.operands, &op)

	return nil
}

func (b *WitnessBuilder) makeHashNode(n node, force bool, hashNodeFunc HashNodeFunc) (hashNode, error) {
	switch n := n.(type) {
	case hashNode:
		return n, nil
	default:
		var hash common.Hash
		if _, err := hashNodeFunc(n, force, hash[:]); err != nil {
			return hashNode{}, err
		}
		return hashNode{hash: hash[:], iws: n.witnessSize()}, nil
	}
}

func (b *WitnessBuilder) addHashOp(key []byte, n hashNode) error {
	if b.trace {
		fmt.Printf("I_HASH: key=%x type: %T v %s\n", key, n, n)
	}

	var op OperatorIntermediateHash

	op.Key = common.CopyBytes(key)
	op.Hash = common.BytesToHash(n.hash)

	b.operands = append(b.operands, &op)
	return nil
}

func (b *WitnessBuilder) addCodeOp(key []byte, code []byte) error {
	if b.trace {
		fmt.Printf("CODE: key=%x len=%d\n", key, len(code))
	}

	var op OperatorCode

	op.Key = common.CopyBytes(key)
	op.Code = make([]byte, len(code))
	copy(op.Code, code)

	b.operands = append(b.operands, &op)
	return nil
}

func (b *WitnessBuilder) processAccountCode(key []byte, n *accountNode, retainDec RetainDecider) error {
	if n.IsEmptyRoot() && n.IsEmptyCodeHash() {
		return nil
	}

	if n.code == nil || (retainDec != nil && !retainDec.IsCodeTouched(n.CodeHash)) {
		return nil
	}

	return b.addCodeOp(key, n.code)
}

func (b *WitnessBuilder) processAccountStorage(n *accountNode, hex []byte, limiter *MerklePathLimiter) error {
	if n.IsEmptyRoot() && n.IsEmptyCodeHash() {
		return nil
	}

	// Here we substitute rs parameter for storageRs, because it needs to become the default
	return b.makeBlockWitness(n.storage, hex, limiter, true)
}

func (b *WitnessBuilder) makeBlockWitness(
	nd node, hex []byte, limiter *MerklePathLimiter, force bool) error {

	processAccountNode := func(key []byte, storageKey []byte, n *accountNode) error {
		if key[len(key)-1] == 0x10 {
			key = key[:len(key)-2]
		}
		var retainDec RetainDecider
		if limiter != nil {
			retainDec = limiter.RetainDecider
		}
		if err := b.processAccountCode(key, n, retainDec); err != nil {
			return err
		}
		if err := b.processAccountStorage(n, storageKey, limiter); err != nil {
			return err
		}
		return b.addAccountLeafOp(key, n)
	}

	switch n := nd.(type) {
	case nil:
		return nil
	case valueNode:
		return b.addLeafOp(hex, n)
	case *accountNode:
		return processAccountNode(hex, hex, n)
	case *shortNode:
		h := n.Key
		// Remove terminator
		if h[len(h)-1] == 16 {
			h = h[:len(h)-1]
		}
		hexVal := concat(hex, h...)
		switch v := n.Val.(type) {
		case valueNode:
			return b.addLeafOp(hexVal, v[:])
		case *accountNode:
			return processAccountNode(hexVal, hexVal, v)
		default:
			if err := b.makeBlockWitness(n.Val, hexVal, limiter, false); err != nil {
				return err
			}

			return nil
		}
	case *duoNode:
		hashOnly := limiter != nil && !limiter.RetainDecider.Retain(hex) // Save this because rl can move on to other keys during the recursive invocation
		if b.trace {
			fmt.Printf("b.retainDec.Retain(%x) -> %v\n", hex, !hashOnly)
		}
		if hashOnly {
			hn, err := b.makeHashNode(n, force, limiter.HashFunc)
			if err != nil {
				return err
			}
			return b.addHashOp(hex, hn)
		}

		i1, i2 := n.childrenIdx()

		if err := b.makeBlockWitness(n.child1, expandKeyHex(hex, i1), limiter, false); err != nil {
			return err
		}
		if err := b.makeBlockWitness(n.child2, expandKeyHex(hex, i2), limiter, false); err != nil {
			return err
		}
		return nil

	case *fullNode:
		hashOnly := limiter != nil && !limiter.RetainDecider.Retain(hex) // Save this because rs can move on to other keys during the recursive invocation
		if hashOnly {
			hn, err := b.makeHashNode(n, force, limiter.HashFunc)
			if err != nil {
				return err
			}
			return b.addHashOp(hex, hn)
		}

		var mask uint32
		for i, child := range n.Children {
			if child != nil {
				if err := b.makeBlockWitness(child, expandKeyHex(hex, byte(i)), limiter, false); err != nil {
					return err
				}
				mask |= (uint32(1) << uint(i))
			}
		}
		return nil

	case hashNode:
		hashOnly := limiter == nil || !limiter.RetainDecider.Retain(hex)
		if hashOnly {
			return b.addHashOp(hex, n)
		}
		return fmt.Errorf("unexpected hashNode: %s, at hex: %x, (%d), hashOnly: %t", n, hex, len(hex), hashOnly)
	default:
		return fmt.Errorf("unexpected node type: %T", nd)
	}
}

func expandKeyHex(hex []byte, nibble byte) []byte {
	result := make([]byte, len(hex)+1)
	copy(result, hex)
	result[len(hex)] = nibble
	return result
}
