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

// TwoFunctionEmptyABI is the input ABI used to generate the binding from.
const TwoFunctionEmptyABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"create\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"update\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// TwoFunctionEmptyBin is the compiled bytecode used for deploying new contracts.
const TwoFunctionEmptyBin = `608060405234801561001057600080fd5b5060c88061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063780900dc14603757806382ab890a146062575b600080fd5b606060048036036020811015604b57600080fd5b8101908080359060200190929190505050608d565b005b608b60048036036020811015607657600080fd5b81019080803590602001909291905050506090565b005b50565b5056fea265627a7a72305820032f00a3ee76729c80d7d01ede94a8512dbe4fa393011cc36cbc9eba4d9e48a664736f6c63430005090032`

// DeployTwoFunctionEmpty deploys a new Ethereum contract, binding an instance of TwoFunctionEmpty to it.
func DeployTwoFunctionEmpty(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TwoFunctionEmpty, error) {
	parsed, err := abi.JSON(strings.NewReader(TwoFunctionEmptyABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TwoFunctionEmptyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TwoFunctionEmpty{TwoFunctionEmptyCaller: TwoFunctionEmptyCaller{contract: contract}, TwoFunctionEmptyTransactor: TwoFunctionEmptyTransactor{contract: contract}, TwoFunctionEmptyFilterer: TwoFunctionEmptyFilterer{contract: contract}}, nil
}

// TwoFunctionEmpty is an auto generated Go binding around an Ethereum contract.
type TwoFunctionEmpty struct {
	TwoFunctionEmptyCaller     // Read-only binding to the contract
	TwoFunctionEmptyTransactor // Write-only binding to the contract
	TwoFunctionEmptyFilterer   // Log filterer for contract events
}

// TwoFunctionEmptyCaller is an auto generated read-only Go binding around an Ethereum contract.
type TwoFunctionEmptyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoFunctionEmptyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TwoFunctionEmptyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoFunctionEmptyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TwoFunctionEmptyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoFunctionEmptySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TwoFunctionEmptySession struct {
	Contract     *TwoFunctionEmpty // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TwoFunctionEmptyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TwoFunctionEmptyCallerSession struct {
	Contract *TwoFunctionEmptyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// TwoFunctionEmptyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TwoFunctionEmptyTransactorSession struct {
	Contract     *TwoFunctionEmptyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// TwoFunctionEmptyRaw is an auto generated low-level Go binding around an Ethereum contract.
type TwoFunctionEmptyRaw struct {
	Contract *TwoFunctionEmpty // Generic contract binding to access the raw methods on
}

// TwoFunctionEmptyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TwoFunctionEmptyCallerRaw struct {
	Contract *TwoFunctionEmptyCaller // Generic read-only contract binding to access the raw methods on
}

// TwoFunctionEmptyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TwoFunctionEmptyTransactorRaw struct {
	Contract *TwoFunctionEmptyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTwoFunctionEmpty creates a new instance of TwoFunctionEmpty, bound to a specific deployed contract.
func NewTwoFunctionEmpty(address common.Address, backend bind.ContractBackend) (*TwoFunctionEmpty, error) {
	contract, err := bindTwoFunctionEmpty(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionEmpty{TwoFunctionEmptyCaller: TwoFunctionEmptyCaller{contract: contract}, TwoFunctionEmptyTransactor: TwoFunctionEmptyTransactor{contract: contract}, TwoFunctionEmptyFilterer: TwoFunctionEmptyFilterer{contract: contract}}, nil
}

// NewTwoFunctionEmptyCaller creates a new read-only instance of TwoFunctionEmpty, bound to a specific deployed contract.
func NewTwoFunctionEmptyCaller(address common.Address, caller bind.ContractCaller) (*TwoFunctionEmptyCaller, error) {
	contract, err := bindTwoFunctionEmpty(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionEmptyCaller{contract: contract}, nil
}

// NewTwoFunctionEmptyTransactor creates a new write-only instance of TwoFunctionEmpty, bound to a specific deployed contract.
func NewTwoFunctionEmptyTransactor(address common.Address, transactor bind.ContractTransactor) (*TwoFunctionEmptyTransactor, error) {
	contract, err := bindTwoFunctionEmpty(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionEmptyTransactor{contract: contract}, nil
}

// NewTwoFunctionEmptyFilterer creates a new log filterer instance of TwoFunctionEmpty, bound to a specific deployed contract.
func NewTwoFunctionEmptyFilterer(address common.Address, filterer bind.ContractFilterer) (*TwoFunctionEmptyFilterer, error) {
	contract, err := bindTwoFunctionEmpty(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionEmptyFilterer{contract: contract}, nil
}

// bindTwoFunctionEmpty binds a generic wrapper to an already deployed contract.
func bindTwoFunctionEmpty(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TwoFunctionEmptyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TwoFunctionEmpty *TwoFunctionEmptyRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TwoFunctionEmpty.Contract.TwoFunctionEmptyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TwoFunctionEmpty *TwoFunctionEmptyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TwoFunctionEmpty.Contract.TwoFunctionEmptyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TwoFunctionEmpty *TwoFunctionEmptyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TwoFunctionEmpty.Contract.TwoFunctionEmptyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TwoFunctionEmpty *TwoFunctionEmptyCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TwoFunctionEmpty.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TwoFunctionEmpty *TwoFunctionEmptyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TwoFunctionEmpty.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TwoFunctionEmpty *TwoFunctionEmptyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TwoFunctionEmpty.Contract.contract.Transact(opts, method, params...)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_TwoFunctionEmpty *TwoFunctionEmptyTransactor) Create(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionEmpty.contract.Transact(opts, "create", newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_TwoFunctionEmpty *TwoFunctionEmptySession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionEmpty.Contract.Create(&_TwoFunctionEmpty.TransactOpts, newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_TwoFunctionEmpty *TwoFunctionEmptyTransactorSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionEmpty.Contract.Create(&_TwoFunctionEmpty.TransactOpts, newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_TwoFunctionEmpty *TwoFunctionEmptyTransactor) Update(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionEmpty.contract.Transact(opts, "update", newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_TwoFunctionEmpty *TwoFunctionEmptySession) Update(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionEmpty.Contract.Update(&_TwoFunctionEmpty.TransactOpts, newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_TwoFunctionEmpty *TwoFunctionEmptyTransactorSession) Update(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionEmpty.Contract.Update(&_TwoFunctionEmpty.TransactOpts, newBalance)
}
