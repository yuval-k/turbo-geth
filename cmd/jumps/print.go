package main

import (
	"context"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledgerwatch/turbo-geth/accounts/abi/bind"
	"github.com/ledgerwatch/turbo-geth/accounts/abi/bind/backends"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/consensus/ethash"
	"github.com/ledgerwatch/turbo-geth/core"
	"github.com/ledgerwatch/turbo-geth/core/asm"
	"github.com/ledgerwatch/turbo-geth/core/types/accounts"
	"github.com/ledgerwatch/turbo-geth/core/vm"
	"github.com/ledgerwatch/turbo-geth/crypto"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/params"
	"github.com/ledgerwatch/turbo-geth/tests/contracts"
)

func mainPrint() {
	var files []string
	binPath := "tests/contracts/build"
	err := filepath.Walk(binPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), "bin") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	type bin struct {
		code string
		name string
	}

	bins := make([]bin, len(files))
	for i, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		name := strings.Split(filepath.Base(file), ".")
		bins[i] = bin{string(content), name[0]}
	}

	for _, code := range bins {
		res, err := asm.Disassembled(code.code)
		if err != nil {
			//panic(fmt.Sprintf("file %q: %v", code.name, err))
		}

		err = ioutil.WriteFile("cmd/jumps/prints/"+code.name+".print", []byte(res), 0666)
		if err != nil {
			panic(err)
		}
	}

	// deployed code
	var (
		db       = ethdb.NewMemDatabase()
		key, _   = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address  = crypto.PubkeyToAddress(key.PublicKey)
		funds    = big.NewInt(1000000000)
		gspec    = &core.Genesis{
			Config: &params.ChainConfig{
				ChainID:             big.NewInt(1),
				HomesteadBlock:      new(big.Int),
				EIP150Block:         new(big.Int),
				EIP155Block:         new(big.Int),
				EIP158Block:         big.NewInt(1),
				EIP2027Block:        big.NewInt(4),
				ByzantiumBlock:      big.NewInt(1),
				ConstantinopleBlock: big.NewInt(1),
			},
			Alloc: core.GenesisAlloc{
				address:  {Balance: funds},
			},
		}
		genesis   = gspec.MustCommit(db)
		genesisDb = db.MemCopy()
	)
	engine := ethash.NewFaker()
	blockchain, err := core.NewBlockChain(db, nil, gspec.Config, engine, vm.Config{}, nil)
	if err != nil {
		panic(err)
	}
	blockchain.EnableReceipts(true)
	contractBackend := backends.NewSimulatedBackendWithConfig(gspec.Alloc, gspec.Config, gspec.GasLimit)
	transactOpts := bind.NewKeyedTransactor(key)

	type deployed struct{
		addr common.Address
		name string
	}
	var contractAddresses []deployed
	blockNum := 2
	ctx := blockchain.WithContext(context.Background(), big.NewInt(int64(genesis.NumberU64())+1))
	blocks, _ := core.GenerateChain(ctx, gspec.Config, genesis, engine, genesisDb, blockNum, func(i int, block *core.BlockGen) {
		if i == 1 {
			contractAddress, tx, _, err := contracts.DeployConstructorAndVar(transactOpts, contractBackend)
			if err != nil {
				panic(err)
			}
			contractAddresses = append(contractAddresses, deployed{contractAddress, "constructor_and_var"})
			block.AddTx(tx)

			contractAddress, tx, _, err = contracts.DeployConstructorAndVarUsage(transactOpts, contractBackend)
			if err != nil {
				panic(err)
			}
			contractAddresses = append(contractAddresses, deployed{contractAddress, "constructor_and_var_usage"})
			block.AddTx(tx)

			contractAddress, tx, _, err = contracts.DeployOnlyConstructor(transactOpts, contractBackend)
			if err != nil {
				panic(err)
			}
			contractAddresses = append(contractAddresses, deployed{contractAddress, "only_constructor"})
			block.AddTx(tx)

			contractAddress, tx, _, err = contracts.DeploySimpleFunction(transactOpts, contractBackend)
			if err != nil {
				panic(err)
			}
			contractAddresses = append(contractAddresses, deployed{contractAddress, "simple_function"})
			block.AddTx(tx)

			contractAddress, tx, _, err = contracts.DeployTwoFunction(transactOpts, contractBackend)
			if err != nil {
				panic(err)
			}
			contractAddresses = append(contractAddresses, deployed{contractAddress, "two_function"})
			block.AddTx(tx)

			contractAddress, tx, _, err = contracts.DeployTwoFunctionCallWithReturn(transactOpts, contractBackend)
			if err != nil {
				panic(err)
			}
			contractAddresses = append(contractAddresses, deployed{contractAddress, "two_functionCallWithReturn"})
			block.AddTx(tx)

			contractAddress, tx, _, err = contracts.DeployTwoFunctionCallWithoutReturn(transactOpts, contractBackend)
			if err != nil {
				panic(err)
			}
			contractAddresses = append(contractAddresses, deployed{contractAddress, "two_functionCallWithoutReturn"})
			block.AddTx(tx)

			contractAddress, tx, _, err = contracts.DeployTwoFunctionEmpty(transactOpts, contractBackend)
			if err != nil {
				panic(err)
			}
			contractAddresses = append(contractAddresses, deployed{contractAddress, "two_functionEmpty"})
			block.AddTx(tx)

			contractAddress, tx, _, err = contracts.DeployWithIf(transactOpts, contractBackend)
			if err != nil {
				panic(err)
			}
			contractAddresses = append(contractAddresses, deployed{contractAddress, "with_if"})
			block.AddTx(tx)

			contractAddress, tx, _, err = contracts.DeployWithIfContinue(transactOpts, contractBackend)
			if err != nil {
				panic(err)
			}
			contractAddresses = append(contractAddresses, deployed{contractAddress, "with_if_continue"})
			block.AddTx(tx)

			contractAddress, tx, _, err = contracts.DeployWithIfElse(transactOpts, contractBackend)
			if err != nil {
				panic(err)
			}
			contractAddresses = append(contractAddresses, deployed{contractAddress, "with_if_else"})
			block.AddTx(tx)
		}
		contractBackend.Commit()
	})

	_, err = blockchain.InsertChain(context.Background(), blocks)
	if err != nil {
		panic(err)
	}

	for _, contract := range contractAddresses {
		h, _ := common.HashData(contract.addr[:])

		acc, err := db.Get(dbutils.AccountsBucket, h.Bytes())
		if err != nil {
			panic(err)
		}

		myAcc := accounts.NewAccount()
		err = myAcc.DecodeForStorage(acc)
		if err != nil {
			panic(err)
		}

		code, err := db.Get(dbutils.CodeBucket, myAcc.CodeHash.Bytes())
		if err != nil {
			panic(err)
		}

		res, err := asm.DisassembledBytes(code)
		if err != nil {
			//panic(fmt.Sprintf("file %q: %v", code.name, err))
		}

		err = ioutil.WriteFile("cmd/jumps/prints/"+contract.name+"_deployed.print", []byte(res), 0666)
		if err != nil {
			panic(err)
		}
	}
}