package trie

import (
	"fmt"

	"github.com/holiman/uint256"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/core/types/accounts"
)

func BuildTrieFromWitness(witness *Witness, isBinary bool, trace bool) (*Trie, error) {
	t := New(common.Hash{})
	fmt.Println(t)
	for _, operator := range witness.Operators {
		switch op := operator.(type) {
		case *OperatorLeafValue:
			if trace {
				fmt.Printf("LEAF ")
			}
			t.Update(hexToKeybytes(op.Key), op.Value)

		case *OperatorIntermediateHash:
			if trace {
				fmt.Printf("HASH ")
			}
			_, t.root = t.insert(t.root, op.Key, hashNode{op.Hash[:], 0})

		case *OperatorCode:
			if trace {
				fmt.Printf("CODE ")
			}
			t.UpdateAccountCode(hexToKeybytes(op.Key), codeNode(op.Code))

		case *OperatorLeafAccount:
			if trace {
				fmt.Printf("ACCOUNTLEAF ")
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

			t.UpdateAccount(hexToKeybytes(op.Key), &acc)

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
