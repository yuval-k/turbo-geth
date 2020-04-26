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

// WithIfElseABI is the input ABI used to generate the binding from.
const WithIfElseABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"create\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"update\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// WithIfElseBin is the compiled bytecode used for deploying new contracts.
const WithIfElseBin = `608060405234801561001057600080fd5b506002600111156100305761002b600561004660201b60201c565b610041565b610040600661004960201b60201c565b5b61004c565b50565b50565b60c88061005a6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063780900dc14603757806382ab890a146062575b600080fd5b606060048036036020811015604b57600080fd5b8101908080359060200190929190505050608d565b005b608b60048036036020811015607657600080fd5b81019080803590602001909291905050506090565b005b50565b5056fea265627a7a72305820cff6884153e082aa261b5d4df6027c25d1c9cd418c8727a247bb6d3c586fde4b64736f6c63430005090032`

// DeployWithIfElse deploys a new Ethereum contract, binding an instance of WithIfElse to it.
func DeployWithIfElse(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WithIfElse, error) {
	parsed, err := abi.JSON(strings.NewReader(WithIfElseABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(WithIfElseBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WithIfElse{WithIfElseCaller: WithIfElseCaller{contract: contract}, WithIfElseTransactor: WithIfElseTransactor{contract: contract}, WithIfElseFilterer: WithIfElseFilterer{contract: contract}}, nil
}

// WithIfElse is an auto generated Go binding around an Ethereum contract.
type WithIfElse struct {
	WithIfElseCaller     // Read-only binding to the contract
	WithIfElseTransactor // Write-only binding to the contract
	WithIfElseFilterer   // Log filterer for contract events
}

// WithIfElseCaller is an auto generated read-only Go binding around an Ethereum contract.
type WithIfElseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithIfElseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WithIfElseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithIfElseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WithIfElseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithIfElseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WithIfElseSession struct {
	Contract     *WithIfElse       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WithIfElseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WithIfElseCallerSession struct {
	Contract *WithIfElseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// WithIfElseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WithIfElseTransactorSession struct {
	Contract     *WithIfElseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// WithIfElseRaw is an auto generated low-level Go binding around an Ethereum contract.
type WithIfElseRaw struct {
	Contract *WithIfElse // Generic contract binding to access the raw methods on
}

// WithIfElseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WithIfElseCallerRaw struct {
	Contract *WithIfElseCaller // Generic read-only contract binding to access the raw methods on
}

// WithIfElseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WithIfElseTransactorRaw struct {
	Contract *WithIfElseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWithIfElse creates a new instance of WithIfElse, bound to a specific deployed contract.
func NewWithIfElse(address common.Address, backend bind.ContractBackend) (*WithIfElse, error) {
	contract, err := bindWithIfElse(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WithIfElse{WithIfElseCaller: WithIfElseCaller{contract: contract}, WithIfElseTransactor: WithIfElseTransactor{contract: contract}, WithIfElseFilterer: WithIfElseFilterer{contract: contract}}, nil
}

// NewWithIfElseCaller creates a new read-only instance of WithIfElse, bound to a specific deployed contract.
func NewWithIfElseCaller(address common.Address, caller bind.ContractCaller) (*WithIfElseCaller, error) {
	contract, err := bindWithIfElse(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WithIfElseCaller{contract: contract}, nil
}

// NewWithIfElseTransactor creates a new write-only instance of WithIfElse, bound to a specific deployed contract.
func NewWithIfElseTransactor(address common.Address, transactor bind.ContractTransactor) (*WithIfElseTransactor, error) {
	contract, err := bindWithIfElse(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WithIfElseTransactor{contract: contract}, nil
}

// NewWithIfElseFilterer creates a new log filterer instance of WithIfElse, bound to a specific deployed contract.
func NewWithIfElseFilterer(address common.Address, filterer bind.ContractFilterer) (*WithIfElseFilterer, error) {
	contract, err := bindWithIfElse(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WithIfElseFilterer{contract: contract}, nil
}

// bindWithIfElse binds a generic wrapper to an already deployed contract.
func bindWithIfElse(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(WithIfElseABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WithIfElse *WithIfElseRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _WithIfElse.Contract.WithIfElseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WithIfElse *WithIfElseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WithIfElse.Contract.WithIfElseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WithIfElse *WithIfElseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WithIfElse.Contract.WithIfElseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WithIfElse *WithIfElseCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _WithIfElse.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WithIfElse *WithIfElseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WithIfElse.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WithIfElse *WithIfElseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WithIfElse.Contract.contract.Transact(opts, method, params...)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_WithIfElse *WithIfElseTransactor) Create(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _WithIfElse.contract.Transact(opts, "create", newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_WithIfElse *WithIfElseSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _WithIfElse.Contract.Create(&_WithIfElse.TransactOpts, newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_WithIfElse *WithIfElseTransactorSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _WithIfElse.Contract.Create(&_WithIfElse.TransactOpts, newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_WithIfElse *WithIfElseTransactor) Update(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _WithIfElse.contract.Transact(opts, "update", newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_WithIfElse *WithIfElseSession) Update(newBalance *big.Int) (*types.Transaction, error) {
	return _WithIfElse.Contract.Update(&_WithIfElse.TransactOpts, newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_WithIfElse *WithIfElseTransactorSession) Update(newBalance *big.Int) (*types.Transaction, error) {
	return _WithIfElse.Contract.Update(&_WithIfElse.TransactOpts, newBalance)
}
