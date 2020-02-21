package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/consensus/misc"
	"github.com/ledgerwatch/turbo-geth/core/state"
	"github.com/ledgerwatch/turbo-geth/eth"
	"github.com/ledgerwatch/turbo-geth/trie"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ledgerwatch/turbo-geth/consensus/ethash"
	"github.com/ledgerwatch/turbo-geth/core"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/core/vm"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/params"
)

func main() {
	sigs := make(chan os.Signal, 1)
	interruptCh := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

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

	bc, err := core.NewBlockChain(ethDb, &core.CacheConfig{TrieCleanLimit: 0, Disabled: true}, chainConfig, ethash.NewFullFaker(), vmConfig, nil)
	check(err)

	lastBlock := bc.CurrentHeader().Number.Uint64()
	blockNum := uint64(1)

	interrupt := false

	fmt.Println("started from", blockNum)
	fmt.Println("started to", lastBlock)
	stateDb, err := ethdb.NewBoltDatabase("/media/b00ris/ssd/ethchain/jump")
	check(err)

	for !interrupt {
		block := bc.GetBlockByNumber(blockNum)
		if block == nil {
			break
		}

		batch := stateDb.NewBatch()

		dbstate, err := state.NewTrieDbState(block.Root(), batch, block.NumberU64())
		check(err)
		dbstate.SetResolveReads(false)
		dbstate.SetNoHistory(true)
		dbstate.Rebuild()

		//dbstate := state.NewDbState(ethDb, block.NumberU64()-1)
		statedb := state.New(dbstate)
		signer := types.MakeSigner(chainConfig, block.Number())
		for _, tx := range block.Transactions() {
			// Assemble the transaction call message and return if the requested offset
			msg, _ := tx.AsMessage(signer)
			ctx := core.NewEVMContext(msg, block.Header(), bc, nil)
			// Not yet the searched for transaction, execute on top of the current state
			vmenv := vm.NewEVM(ctx, statedb, chainConfig, vmConfig)
			if _, _, _, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(tx.Gas())); err != nil {
				panic(fmt.Errorf("tx %x failed: %v", tx.Hash(), err))
			}

			err = statedb.FinalizeTx(vmenv.ChainConfig().WithEIPsFlags(context.Background(), block.Number()), dbstate)
			check(err)
		}
		if _, err = batch.Commit(); err != nil {
			fmt.Printf("Failed to commit batch: %v\n", err)
		}

		blockNum++
		if blockNum%1000 == 0 {
			push, nonpush := vm.GetJumps()
			fmt.Printf("Processed %d blocks, jumps from PUSHes %d, not from PUSHes %d\n", blockNum, push, nonpush)
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

	for !interrupt {
		trace := blockNum == 50492 // false // blockNum == 545080
		tds.SetResolveReads(blockNum >= witnessThreshold)
		block := bcb.GetBlockByNumber(blockNum)
		if block == nil {
			break
		}
		execStart := time.Now()
		statedb := state.New(tds)
		gp := new(core.GasPool).AddGas(block.GasLimit())
		usedGas := new(uint64)
		header := block.Header()
		tds.StartNewBuffer()
		var receipts types.Receipts
		if chainConfig.DAOForkSupport && chainConfig.DAOForkBlock != nil && chainConfig.DAOForkBlock.Cmp(block.Number()) == 0 {
			misc.ApplyDAOHardFork(statedb)
		}
		for i, tx := range block.Transactions() {
			statedb.Prepare(tx.Hash(), block.Hash(), i)
			receipt, err := core.ApplyTransaction(chainConfig, bcb, nil, gp, statedb, tds.TrieStateWriter(), header, tx, usedGas, vmConfig)
			if err != nil {
				fmt.Printf("tx %x failed: %v\n", tx.Hash(), err)
				return
			}
			if !chainConfig.IsByzantium(header.Number) {
				tds.StartNewBuffer()
			}
			receipts = append(receipts, receipt)
		}
		// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
		if _, err = engine.FinalizeAndAssemble(chainConfig, header, statedb, block.Transactions(), block.Uncles(), receipts); err != nil {
			fmt.Printf("Finalize of block %d failed: %v\n", blockNum, err)
			return
		}

		execTime1 := time.Since(execStart)
		execStart = time.Now()

		ctx := chainConfig.WithEIPsFlags(context.Background(), header.Number)
		if err := statedb.FinalizeTx(ctx, tds.TrieStateWriter()); err != nil {
			fmt.Printf("FinalizeTx of block %d failed: %v\n", blockNum, err)
			return
		}

		if witnessDBReader != nil {
			tds.SetBlockNr(blockNum)
			err = tds.ResolveStateTrieStateless(witnessDBReader)
			if err != nil {
				fmt.Printf("Failed to statelessly resolve state trie: %v\n", err)
				return
			}
		} else {
			var resolveWitnesses []*trie.Witness
			if resolveWitnesses, err = tds.ResolveStateTrie(witnessDBWriter != nil); err != nil {
				fmt.Printf("Failed to resolve state trie: %v\n", err)
				return
			}

			if len(resolveWitnesses) > 0 {
				witnessDBWriter.MustUpsert(blockNum, state.MaxTrieCacheGen, resolveWitnesses)
			}
		}
		execTime2 := time.Since(execStart)
		blockWitness = nil
		if blockNum >= witnessThreshold {
			// Witness has to be extracted before the state trie is modified
			var blockWitnessStats *trie.BlockWitnessStats
			bw, err = tds.ExtractWitness(trace, binary /* is binary */)
			if err != nil {
				fmt.Printf("error extracting witness for block %d: %v\n", blockNum, err)
				return
			}

			var buf bytes.Buffer
			blockWitnessStats, err = bw.WriteTo(&buf)
			if err != nil {
				fmt.Printf("error extracting witness for block %d: %v\n", blockNum, err)
				return
			}
			blockWitness = buf.Bytes()
			err = stats.AddRow(blockNum, blockWitnessStats)
			check(err)
		}
		finalRootFail := false
		execStart = time.Now()
		if blockNum >= witnessThreshold && blockWitness != nil { // blockWitness == nil means the extraction fails
			var s *state.Stateless
			var w *trie.Witness
			w, err = trie.NewWitnessFromReader(bytes.NewReader(blockWitness), false)
			bw.WriteDiff(w, os.Stdout)
			if err != nil {
				fmt.Printf("error deserializing witness for block %d: %v\n", blockNum, err)
				return
			}
			if _, ok := starkBlocks[blockNum-1]; ok {
				err = starkData(w, starkStatsBase, blockNum-1)
				check(err)
			}
			s, err = state.NewStateless(preRoot, w, blockNum-1, trace, binary /* is binary */)
			if err != nil {
				fmt.Printf("Error making stateless2 for block %d: %v\n", blockNum, err)
				filename := fmt.Sprintf("right_%d.txt", blockNum-1)
				f, err1 := os.Create(filename)
				if err1 == nil {
					defer f.Close()
					tds.PrintTrie(f)
				}
				return
			}
			if _, ok := starkBlocks[blockNum-1]; ok {
				err = statePicture(s.GetTrie(), s.GetCodeMap(), blockNum-1)
				check(err)
			}
			if err = runBlock(s, chainConfig, bcb, header, block, trace, !binary); err != nil {
				fmt.Printf("Error running block %d through stateless2: %v\n", blockNum, err)
				finalRootFail = true
			}
		}
		execTime3 := time.Since(execStart)
		execStart = time.Now()
		var preCalculatedRoot common.Hash
		if tryPreRoot {
			preCalculatedRoot, err = tds.CalcTrieRoots(blockNum == 1)
			if err != nil {
				fmt.Printf("failed to calculate preRoot for block %d: %v\n", blockNum, err)
				return
			}
		}
		execTime4 := time.Since(execStart)
		execStart = time.Now()
		roots, err := tds.UpdateStateTrie()
		if err != nil {
			fmt.Printf("failed to calculate IntermediateRoot: %v\n", err)
			return
		}
		execTime5 := time.Since(execStart)
		if tryPreRoot && tds.LastRoot() != preCalculatedRoot {
			filename := fmt.Sprintf("right_%d.txt", blockNum)
			f, err1 := os.Create(filename)
			if err1 == nil {
				defer f.Close()
				tds.PrintTrie(f)
			}
			fmt.Printf("block %d, preCalculatedRoot %x != lastRoot %x\n", blockNum, preCalculatedRoot, tds.LastRoot())
			return
		}
		if finalRootFail {
			filename := fmt.Sprintf("right_%d.txt", blockNum)
			f, err1 := os.Create(filename)
			if err1 == nil {
				defer f.Close()
				tds.PrintTrie(f)
			}
			return
		}
		if !chainConfig.IsByzantium(header.Number) {
			for i, receipt := range receipts {
				receipt.PostState = roots[i].Bytes()
			}
		}
		nextRoot := roots[len(roots)-1]
		if nextRoot != block.Root() {
			fmt.Printf("Root hash does not match for block %d, expected %x, was %x\n", blockNum, block.Root(), nextRoot)
			return
		}
		tds.SetBlockNr(blockNum)

		err = statedb.CommitBlock(ctx, tds.DbStateWriter())
		if err != nil {
			fmt.Printf("Commiting block %d failed: %v", blockNum, err)
			return
		}

		if batch.BatchSize() >= 100000 {
			if _, err := batch.Commit(); err != nil {
				fmt.Printf("Failed to commit batch: %v\n", err)
				return
			}
			tds.PruneTries(false)
		}

		select {
		case interrupt = <-interruptCh:
			fmt.Println("interrupted, please wait for cleanup...")
		default:
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
