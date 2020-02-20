package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ledgerwatch/turbo-geth/consensus/ethash"
	"github.com/ledgerwatch/turbo-geth/core"
	"github.com/ledgerwatch/turbo-geth/core/state"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/core/vm"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/params"
)

func main() {
	sigs := make(chan os.Signal, 1)
	interruptCh := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		interruptCh <- true
	}()

	ethDb, err := ethdb.NewBoltDatabase("/media/b00ris/ssd/ethchain/thin_1/geth/chaindata")
	check(err)
	defer ethDb.Close()

	chainConfig := params.MainnetChainConfig
	srwFile, err := os.OpenFile("storage_read_writes.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	check(err)
	defer srwFile.Close()

	w := bufio.NewWriter(srwFile)
	defer w.Flush()

	vmConfig := vm.Config{}
	bc, err := core.NewBlockChain(ethDb, nil, chainConfig, ethash.NewFaker(), vmConfig, nil)
	check(err)

	lastBlock := bc.CurrentHeader().Number.Uint64()
	blockNum := uint64(1)

	interrupt := false
	for !interrupt {
		block := bc.GetBlockByNumber(blockNum)
		if block == nil {
			break
		}
		dbstate := state.NewDbState(ethDb, block.NumberU64()-1)
		statedb := state.New(dbstate)
		signer := types.MakeSigner(chainConfig, block.Number())
		for _, tx := range block.Transactions() {
			// Assemble the transaction call message and return if the requested offset
			msg, _ := tx.AsMessage(signer)
			context := core.NewEVMContext(msg, block.Header(), bc, nil)
			// Not yet the searched for transaction, execute on top of the current state
			vmenv := vm.NewEVM(context, statedb, chainConfig, vmConfig)
			if _, _, _, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(tx.Gas())); err != nil {
				panic(fmt.Errorf("tx %x failed: %v", tx.Hash(), err))
			}
		}

		blockNum++
		if blockNum%1000 == 0 {
			push, nonpush := vm.GetJumps()
			fmt.Fprintf(w, "Processed %d blocks, jumps from PUSHes %d, not from PUSHes %d\n", blockNum, push, nonpush)
		}

		// Check for interrupts
		select {
		case interrupt = <-interruptCh:
			fmt.Println("interrupted, please wait for cleanup...")
		default:
		}

		if lastBlock < blockNum {
			interrupt = true
		}
	}

	push, nonpush := vm.GetJumps()
	fmt.Fprintf(w,"Processed %d blocks, jumps from PUSHes %d, not from PUSHes %d\n", blockNum, push, nonpush)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
