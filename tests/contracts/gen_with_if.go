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

// WithIfABI is the input ABI used to generate the binding from.
const WithIfABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"create\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// WithIfBin is the compiled bytecode used for deploying new contracts.
const WithIfBin = `608060405234801561001057600080fd5b5060026001111561002c5761002b600561003160201b60201c565b5b610034565b50565b6090806100426000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063780900dc14602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506058565b005b5056fea265627a7a72305820aaeaf2e2f752de718ca490639eecfa1a6e0bae597dba5ab64d643e3ce717891564736f6c63430005090032`

// DeployWithIf deploys a new Ethereum contract, binding an instance of WithIf to it.
func DeployWithIf(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WithIf, error) {
	parsed, err := abi.JSON(strings.NewReader(WithIfABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(WithIfBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WithIf{WithIfCaller: WithIfCaller{contract: contract}, WithIfTransactor: WithIfTransactor{contract: contract}, WithIfFilterer: WithIfFilterer{contract: contract}}, nil
}

// WithIf is an auto generated Go binding around an Ethereum contract.
type WithIf struct {
	WithIfCaller     // Read-only binding to the contract
	WithIfTransactor // Write-only binding to the contract
	WithIfFilterer   // Log filterer for contract events
}

// WithIfCaller is an auto generated read-only Go binding around an Ethereum contract.
type WithIfCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithIfTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WithIfTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithIfFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WithIfFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithIfSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WithIfSession struct {
	Contract     *WithIf           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WithIfCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WithIfCallerSession struct {
	Contract *WithIfCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// WithIfTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WithIfTransactorSession struct {
	Contract     *WithIfTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WithIfRaw is an auto generated low-level Go binding around an Ethereum contract.
type WithIfRaw struct {
	Contract *WithIf // Generic contract binding to access the raw methods on
}

// WithIfCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WithIfCallerRaw struct {
	Contract *WithIfCaller // Generic read-only contract binding to access the raw methods on
}

// WithIfTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WithIfTransactorRaw struct {
	Contract *WithIfTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWithIf creates a new instance of WithIf, bound to a specific deployed contract.
func NewWithIf(address common.Address, backend bind.ContractBackend) (*WithIf, error) {
	contract, err := bindWithIf(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WithIf{WithIfCaller: WithIfCaller{contract: contract}, WithIfTransactor: WithIfTransactor{contract: contract}, WithIfFilterer: WithIfFilterer{contract: contract}}, nil
}

// NewWithIfCaller creates a new read-only instance of WithIf, bound to a specific deployed contract.
func NewWithIfCaller(address common.Address, caller bind.ContractCaller) (*WithIfCaller, error) {
	contract, err := bindWithIf(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WithIfCaller{contract: contract}, nil
}

// NewWithIfTransactor creates a new write-only instance of WithIf, bound to a specific deployed contract.
func NewWithIfTransactor(address common.Address, transactor bind.ContractTransactor) (*WithIfTransactor, error) {
	contract, err := bindWithIf(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WithIfTransactor{contract: contract}, nil
}

// NewWithIfFilterer creates a new log filterer instance of WithIf, bound to a specific deployed contract.
func NewWithIfFilterer(address common.Address, filterer bind.ContractFilterer) (*WithIfFilterer, error) {
	contract, err := bindWithIf(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WithIfFilterer{contract: contract}, nil
}

// bindWithIf binds a generic wrapper to an already deployed contract.
func bindWithIf(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(WithIfABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WithIf *WithIfRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _WithIf.Contract.WithIfCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WithIf *WithIfRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WithIf.Contract.WithIfTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WithIf *WithIfRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WithIf.Contract.WithIfTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WithIf *WithIfCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _WithIf.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WithIf *WithIfTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WithIf.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WithIf *WithIfTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WithIf.Contract.contract.Transact(opts, method, params...)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_WithIf *WithIfTransactor) Create(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _WithIf.contract.Transact(opts, "create", newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_WithIf *WithIfSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _WithIf.Contract.Create(&_WithIf.TransactOpts, newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_WithIf *WithIfTransactorSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _WithIf.Contract.Create(&_WithIf.TransactOpts, newBalance)
}
