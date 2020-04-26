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

// WithIfContinueABI is the input ABI used to generate the binding from.
const WithIfContinueABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"create\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"update\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// WithIfContinueBin is the compiled bytecode used for deploying new contracts.
const WithIfContinueBin = `608060405234801561001057600080fd5b506002600111156100305761002b600561004660201b60201c565b610041565b610040600661004960201b60201c565b5b61004c565b50565b50565b60c88061005a6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063780900dc14603757806382ab890a146062575b600080fd5b606060048036036020811015604b57600080fd5b8101908080359060200190929190505050608d565b005b608b60048036036020811015607657600080fd5b81019080803590602001909291905050506090565b005b50565b5056fea265627a7a723058205602a2ed3f9f3c167be892953389ecf83ecd3f75fff569afe151bca9315cb11b64736f6c63430005090032`

// DeployWithIfContinue deploys a new Ethereum contract, binding an instance of WithIfContinue to it.
func DeployWithIfContinue(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WithIfContinue, error) {
	parsed, err := abi.JSON(strings.NewReader(WithIfContinueABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(WithIfContinueBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WithIfContinue{WithIfContinueCaller: WithIfContinueCaller{contract: contract}, WithIfContinueTransactor: WithIfContinueTransactor{contract: contract}, WithIfContinueFilterer: WithIfContinueFilterer{contract: contract}}, nil
}

// WithIfContinue is an auto generated Go binding around an Ethereum contract.
type WithIfContinue struct {
	WithIfContinueCaller     // Read-only binding to the contract
	WithIfContinueTransactor // Write-only binding to the contract
	WithIfContinueFilterer   // Log filterer for contract events
}

// WithIfContinueCaller is an auto generated read-only Go binding around an Ethereum contract.
type WithIfContinueCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithIfContinueTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WithIfContinueTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithIfContinueFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WithIfContinueFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithIfContinueSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WithIfContinueSession struct {
	Contract     *WithIfContinue   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WithIfContinueCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WithIfContinueCallerSession struct {
	Contract *WithIfContinueCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// WithIfContinueTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WithIfContinueTransactorSession struct {
	Contract     *WithIfContinueTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// WithIfContinueRaw is an auto generated low-level Go binding around an Ethereum contract.
type WithIfContinueRaw struct {
	Contract *WithIfContinue // Generic contract binding to access the raw methods on
}

// WithIfContinueCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WithIfContinueCallerRaw struct {
	Contract *WithIfContinueCaller // Generic read-only contract binding to access the raw methods on
}

// WithIfContinueTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WithIfContinueTransactorRaw struct {
	Contract *WithIfContinueTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWithIfContinue creates a new instance of WithIfContinue, bound to a specific deployed contract.
func NewWithIfContinue(address common.Address, backend bind.ContractBackend) (*WithIfContinue, error) {
	contract, err := bindWithIfContinue(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WithIfContinue{WithIfContinueCaller: WithIfContinueCaller{contract: contract}, WithIfContinueTransactor: WithIfContinueTransactor{contract: contract}, WithIfContinueFilterer: WithIfContinueFilterer{contract: contract}}, nil
}

// NewWithIfContinueCaller creates a new read-only instance of WithIfContinue, bound to a specific deployed contract.
func NewWithIfContinueCaller(address common.Address, caller bind.ContractCaller) (*WithIfContinueCaller, error) {
	contract, err := bindWithIfContinue(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WithIfContinueCaller{contract: contract}, nil
}

// NewWithIfContinueTransactor creates a new write-only instance of WithIfContinue, bound to a specific deployed contract.
func NewWithIfContinueTransactor(address common.Address, transactor bind.ContractTransactor) (*WithIfContinueTransactor, error) {
	contract, err := bindWithIfContinue(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WithIfContinueTransactor{contract: contract}, nil
}

// NewWithIfContinueFilterer creates a new log filterer instance of WithIfContinue, bound to a specific deployed contract.
func NewWithIfContinueFilterer(address common.Address, filterer bind.ContractFilterer) (*WithIfContinueFilterer, error) {
	contract, err := bindWithIfContinue(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WithIfContinueFilterer{contract: contract}, nil
}

// bindWithIfContinue binds a generic wrapper to an already deployed contract.
func bindWithIfContinue(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(WithIfContinueABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WithIfContinue *WithIfContinueRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _WithIfContinue.Contract.WithIfContinueCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WithIfContinue *WithIfContinueRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WithIfContinue.Contract.WithIfContinueTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WithIfContinue *WithIfContinueRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WithIfContinue.Contract.WithIfContinueTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WithIfContinue *WithIfContinueCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _WithIfContinue.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WithIfContinue *WithIfContinueTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WithIfContinue.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WithIfContinue *WithIfContinueTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WithIfContinue.Contract.contract.Transact(opts, method, params...)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_WithIfContinue *WithIfContinueTransactor) Create(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _WithIfContinue.contract.Transact(opts, "create", newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_WithIfContinue *WithIfContinueSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _WithIfContinue.Contract.Create(&_WithIfContinue.TransactOpts, newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_WithIfContinue *WithIfContinueTransactorSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _WithIfContinue.Contract.Create(&_WithIfContinue.TransactOpts, newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_WithIfContinue *WithIfContinueTransactor) Update(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _WithIfContinue.contract.Transact(opts, "update", newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_WithIfContinue *WithIfContinueSession) Update(newBalance *big.Int) (*types.Transaction, error) {
	return _WithIfContinue.Contract.Update(&_WithIfContinue.TransactOpts, newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_WithIfContinue *WithIfContinueTransactorSession) Update(newBalance *big.Int) (*types.Transaction, error) {
	return _WithIfContinue.Contract.Update(&_WithIfContinue.TransactOpts, newBalance)
}
