package trie

import (
	"bytes"
	"math/big"

	"github.com/ledgerwatch/turbo-geth/common"
)

const (
	flagCode = 1 << iota
	flagStorage
	flagNonce
	flagBalance
)

// OperatorKindCode is "enum" type for defining the opcodes of the stack machine that reconstructs the structure of tries from Structure tape
type OperatorKindCode uint8

const (
	// OpLeaf creates leaf node and pushes it onto the node stack, its hash onto the hash stack
	OpLeaf OperatorKindCode = iota
	// OpExtension pops a node from the node stack, constructs extension node from it and its operand's key, and pushes this extension node onto
	// the node stack, its hash onto the hash stack
	OpExtension
	// OpBranch has operand, which is a bitset representing digits in the branch node. Pops the children nodes from the node stack (the number of
	// children is equal to the number of bits in the bitset), constructs branch node and pushes it onto the node stack, its hash onto the hash stack
	OpBranch
	// OpHash and pushes the hash them onto the stack.
	OpHash
	// OpCode constructs code node and pushes it onto the node stack, its hash onto the hash stack.
	OpCode
	// OpAccountLeaf constructs an account node (without any storage and code) and pushes it onto the node stack, its hash onto the hash stack.
	OpAccountLeaf
	// OpEmptyRoot places nil onto the node stack, and empty root hash onto the hash stack.
	OpEmptyRoot

	// OpNewTrie stops the processing, because another trie is encoded into the witness.
	OpNewTrie = OperatorKindCode(0xBB)
)

// WitnessOperator is a single operand in the block witness. It knows how to serialize/deserialize itself.
type WitnessOperator interface {
	GetKey() []byte
	WriteTo(m *OperatorMarshaller, previousNibbles []byte) error

	// LoadFrom always assumes that the opcode value was already read
	LoadFrom(u *OperatorUnmarshaller, previousNibbles []byte) error
}

type OperatorIntermediateHash struct {
	Key  []byte
	Hash common.Hash
}

func (o *OperatorIntermediateHash) GetKey() []byte {
	return o.Key
}

func (o *OperatorIntermediateHash) WriteTo(output *OperatorMarshaller, previousNibbles []byte) error {
	if err := output.WriteOpCode(OpHash); err != nil {
		return nil
	}

	if err := output.WriteKey(o.Key, previousNibbles); err != nil {
		return err
	}

	return output.WriteHash(o.Hash)
}

func (o *OperatorIntermediateHash) LoadFrom(loader *OperatorUnmarshaller, previousNibbles []byte) error {
	if key, err := loader.ReadKey(previousNibbles); err == nil {
		o.Key = key
	} else {
		return err
	}

	if hash, err := loader.ReadHash(); err == nil {
		o.Hash = hash
	} else {
		return err
	}
	return nil
}

type OperatorLeafValue struct {
	Key   []byte
	Value []byte
}

func (o *OperatorLeafValue) GetKey() []byte {
	return o.Key
}

func (o *OperatorLeafValue) WriteTo(output *OperatorMarshaller, previousNibbles []byte) error {
	if err := output.WriteOpCode(OpLeaf); err != nil {
		return err
	}

	if err := output.WriteKey(o.Key, previousNibbles); err != nil {
		return err
	}

	return output.WriteByteArrayValue(o.Value)
}

func (o *OperatorLeafValue) LoadFrom(loader *OperatorUnmarshaller, previousNibbles []byte) error {
	key, err := loader.ReadKey(previousNibbles)
	if err != nil {
		return err
	}

	o.Key = key

	value, err := loader.ReadByteArray()
	if err != nil {
		return err
	}

	o.Value = value
	return nil
}

type OperatorLeafAccount struct {
	Key      []byte
	Nonce    uint64
	Balance  *big.Int
	Root     common.Hash
	CodeHash common.Hash
	CodeSize uint64
}

func (o *OperatorLeafAccount) GetKey() []byte {
	return o.Key
}

func (o *OperatorLeafAccount) WriteTo(output *OperatorMarshaller, previousNibbles []byte) error {
	if err := output.WriteOpCode(OpAccountLeaf); err != nil {
		return err
	}

	if err := output.WriteKey(o.Key, previousNibbles); err != nil {
		return err
	}

	flags := byte(0)
	if o.Nonce > 0 {
		flags |= flagNonce
	}

	if o.Balance.Sign() != 0 {
		flags |= flagBalance
	}

	if !bytes.Equal(o.Root[:], EmptyRoot[:]) {
		flags |= flagStorage
	}

	if !bytes.Equal(o.CodeHash[:], EmptyCodeHash[:]) {
		flags |= flagCode
	}

	if err := output.WriteByteValue(flags); err != nil {
		return err
	}

	if o.Nonce > 0 {
		if err := output.WriteUint64Value(o.Nonce); err != nil {
			return err
		}
	}

	if o.Balance.Sign() != 0 {
		if err := output.WriteByteArrayValue(o.Balance.Bytes()); err != nil {
			return err
		}
	}

	if flags&flagStorage != 0 {
		if err := output.WriteHash(o.Root); err != nil {
			return err
		}
	}

	if flags&flagCode != 0 {
		if err := output.WriteHash(o.CodeHash); err != nil {
			return err
		}
		if err := output.WriteUint64Value(o.CodeSize); err != nil {
			return err
		}
	}

	return nil
}

func (o *OperatorLeafAccount) LoadFrom(loader *OperatorUnmarshaller, previousNibbles []byte) error {
	key, err := loader.ReadKey(previousNibbles)
	if err != nil {
		return err
	}

	o.Key = common.CopyBytes(key)

	flags, err := loader.ReadByte()
	if err != nil {
		return err
	}

	if flags&flagNonce != 0 {
		o.Nonce, err = loader.ReadUInt64()
		if err != nil {
			return err
		}
	}

	balance := big.NewInt(0)

	if flags&flagBalance != 0 {
		var balanceBytes []byte
		balanceBytes, err = loader.ReadByteArray()
		if err != nil {
			return err
		}
		balance.SetBytes(balanceBytes)
	}

	o.Balance = balance

	if flags&flagStorage != 0 {
		if root, err := loader.ReadHash(); err == nil {
			o.Root = root
		} else {
			return err
		}
	} else {
		o.Root = EmptyRoot
	}

	if flags&flagCode != 0 {
		if codeHash, err := loader.ReadHash(); err == nil {
			o.CodeHash = codeHash
		} else {
			return err
		}
		if codeSize, err := loader.ReadUInt64(); err == nil {
			o.CodeSize = codeSize
		} else {
			return err
		}
	} else {
		o.CodeHash = EmptyCodeHash
	}

	return nil
}

type OperatorCode struct {
	Key  []byte
	Code []byte
}

func (o *OperatorCode) GetKey() []byte {
	return o.Key
}

func (o *OperatorCode) WriteTo(output *OperatorMarshaller, previousNibbles []byte) error {
	if err := output.WriteOpCode(OpCode); err != nil {
		return err
	}

	if err := output.WriteKey(o.Key, previousNibbles); err != nil {
		return err
	}

	return output.WriteCode(o.Code)
}

func (o *OperatorCode) LoadFrom(loader *OperatorUnmarshaller, previousNibbles []byte) error {
	if key, err := loader.ReadKey(previousNibbles); err == nil {
		o.Key = key
	} else {
		return err
	}

	if code, err := loader.ReadByteArray(); err == nil {
		o.Code = code
	} else {
		return err
	}
	return nil
}
