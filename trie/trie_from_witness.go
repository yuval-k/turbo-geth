package trie

import (
	"fmt"

	"github.com/holiman/uint256"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/core/types/accounts"
)

func BuildTrieFromWitness(witness *Witness, isBinary bool, trace bool) (*Trie, error) {
	t := New(common.Hash{})
	for _, operator := range witness.Operators {
		switch op := operator.(type) {
		case *OperatorLeafValue:
			if trace {
				fmt.Printf("LEAF ")
			}
			t.Update(hexToKeybytes(op.Key), op.Value)

		case *OperatorIntermediateHash:
			if trace {
				fmt.Printf("HASH %x\n", op.Key)
			}
			_, t.root = t.insert(t.root, op.Key, hashNode{op.Hash[:], 0})

		case *OperatorCode:
			if trace {
				fmt.Printf("CODE 0x%x->%v\n", op.Key, len(op.Code))
			}
			err := t.UpdateAccountCode(hexToKeybytes(op.Key), codeNode(op.Code))
			if err != nil {
				fmt.Printf("err while updating code: %v\n", err)
			}

		case *OperatorLeafAccount:
			if trace {
				fmt.Printf("ACCOUNTLEAF %x\n", op.Key)
			}
			balance := uint256.NewInt()
			balance.SetBytes(op.Balance.Bytes())
			nonce := op.Nonce

			acc := accounts.Account{
				Initialised: true,
				Nonce:       nonce,
				Balance:     *balance,
				Root:        op.Root,
				CodeHash:    op.CodeHash,
				Incarnation: 0,
			}

			k := hexToKeybytes(op.Key)
			fmt.Printf("inserting acc %x\n", k)

			t.UpdateAccount(k, &acc)
			if op.CodeSize > 0 {
				fmt.Printf("updating code size of %x -> %d\n", k, op.CodeSize)
				t.UpdateAccountCodeSize(k, int(op.CodeSize))
			}

		default:
			return nil, fmt.Errorf("unknown operand type: %T", operator)
		}
	}
	if trace {
		fmt.Printf("\n")
	}
	_ = t.Hash()
	return t, nil
}
