// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"math/big"
	"strings"

	ethereum "github.com/ledgerwatch/turbo-geth"
	"github.com/ledgerwatch/turbo-geth/accounts/abi"
	"github.com/ledgerwatch/turbo-geth/accounts/abi/bind"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ConstructorAndVarABI is the input ABI used to generate the binding from.
const ConstructorAndVarABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// ConstructorAndVarBin is the compiled bytecode used for deploying new contracts.
const ConstructorAndVarBin = `608060405234801561001057600080fd5b5060cf8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806327e235e314602d575b600080fd5b606c60048036036020811015604157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506082565b6040518082815260200191505060405180910390f35b6000602052806000526040600020600091509050548156fea265627a7a72305820113c26735797eb6eab527afa6d248b3cad43c7c7d54867e50aca1cfa681f334564736f6c63430005090032`

// DeployConstructorAndVar deploys a new Ethereum contract, binding an instance of ConstructorAndVar to it.
func DeployConstructorAndVar(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ConstructorAndVar, error) {
	parsed, err := abi.JSON(strings.NewReader(ConstructorAndVarABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ConstructorAndVarBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConstructorAndVar{ConstructorAndVarCaller: ConstructorAndVarCaller{contract: contract}, ConstructorAndVarTransactor: ConstructorAndVarTransactor{contract: contract}, ConstructorAndVarFilterer: ConstructorAndVarFilterer{contract: contract}}, nil
}

// ConstructorAndVar is an auto generated Go binding around an Ethereum contract.
type ConstructorAndVar struct {
	ConstructorAndVarCaller     // Read-only binding to the contract
	ConstructorAndVarTransactor // Write-only binding to the contract
	ConstructorAndVarFilterer   // Log filterer for contract events
}

// ConstructorAndVarCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConstructorAndVarCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstructorAndVarTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConstructorAndVarTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstructorAndVarFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConstructorAndVarFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstructorAndVarSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConstructorAndVarSession struct {
	Contract     *ConstructorAndVar // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ConstructorAndVarCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConstructorAndVarCallerSession struct {
	Contract *ConstructorAndVarCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ConstructorAndVarTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConstructorAndVarTransactorSession struct {
	Contract     *ConstructorAndVarTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ConstructorAndVarRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConstructorAndVarRaw struct {
	Contract *ConstructorAndVar // Generic contract binding to access the raw methods on
}

// ConstructorAndVarCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConstructorAndVarCallerRaw struct {
	Contract *ConstructorAndVarCaller // Generic read-only contract binding to access the raw methods on
}

// ConstructorAndVarTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConstructorAndVarTransactorRaw struct {
	Contract *ConstructorAndVarTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConstructorAndVar creates a new instance of ConstructorAndVar, bound to a specific deployed contract.
func NewConstructorAndVar(address common.Address, backend bind.ContractBackend) (*ConstructorAndVar, error) {
	contract, err := bindConstructorAndVar(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConstructorAndVar{ConstructorAndVarCaller: ConstructorAndVarCaller{contract: contract}, ConstructorAndVarTransactor: ConstructorAndVarTransactor{contract: contract}, ConstructorAndVarFilterer: ConstructorAndVarFilterer{contract: contract}}, nil
}

// NewConstructorAndVarCaller creates a new read-only instance of ConstructorAndVar, bound to a specific deployed contract.
func NewConstructorAndVarCaller(address common.Address, caller bind.ContractCaller) (*ConstructorAndVarCaller, error) {
	contract, err := bindConstructorAndVar(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConstructorAndVarCaller{contract: contract}, nil
}

// NewConstructorAndVarTransactor creates a new write-only instance of ConstructorAndVar, bound to a specific deployed contract.
func NewConstructorAndVarTransactor(address common.Address, transactor bind.ContractTransactor) (*ConstructorAndVarTransactor, error) {
	contract, err := bindConstructorAndVar(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConstructorAndVarTransactor{contract: contract}, nil
}

// NewConstructorAndVarFilterer creates a new log filterer instance of ConstructorAndVar, bound to a specific deployed contract.
func NewConstructorAndVarFilterer(address common.Address, filterer bind.ContractFilterer) (*ConstructorAndVarFilterer, error) {
	contract, err := bindConstructorAndVar(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConstructorAndVarFilterer{contract: contract}, nil
}

// bindConstructorAndVar binds a generic wrapper to an already deployed contract.
func bindConstructorAndVar(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConstructorAndVarABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstructorAndVar *ConstructorAndVarRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ConstructorAndVar.Contract.ConstructorAndVarCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstructorAndVar *ConstructorAndVarRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstructorAndVar.Contract.ConstructorAndVarTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstructorAndVar *ConstructorAndVarRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstructorAndVar.Contract.ConstructorAndVarTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstructorAndVar *ConstructorAndVarCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ConstructorAndVar.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstructorAndVar *ConstructorAndVarTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstructorAndVar.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstructorAndVar *ConstructorAndVarTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstructorAndVar.Contract.contract.Transact(opts, method, params...)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_ConstructorAndVar *ConstructorAndVarCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConstructorAndVar.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_ConstructorAndVar *ConstructorAndVarSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _ConstructorAndVar.Contract.Balances(&_ConstructorAndVar.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_ConstructorAndVar *ConstructorAndVarCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _ConstructorAndVar.Contract.Balances(&_ConstructorAndVar.CallOpts, arg0)
}
