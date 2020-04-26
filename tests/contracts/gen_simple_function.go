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

// SimpleFunctionABI is the input ABI used to generate the binding from.
const SimpleFunctionABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"create\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// SimpleFunctionBin is the compiled bytecode used for deploying new contracts.
const SimpleFunctionBin = `608060405234801561001057600080fd5b5060646000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610154806100646000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806327e235e31461003b578063780900dc14610093575b600080fd5b61007d6004803603602081101561005157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506100c1565b6040518082815260200191505060405180910390f35b6100bf600480360360208110156100a957600080fd5b81019080803590602001909291905050506100d9565b005b60006020528060005260406000206000915090505481565b806000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055505056fea265627a7a72305820c0cd8276a80abe4b0c0f6cdec073d24f4dde53876707419236126af738ec737b64736f6c63430005090032`

// DeploySimpleFunction deploys a new Ethereum contract, binding an instance of SimpleFunction to it.
func DeploySimpleFunction(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SimpleFunction, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleFunctionABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SimpleFunctionBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimpleFunction{SimpleFunctionCaller: SimpleFunctionCaller{contract: contract}, SimpleFunctionTransactor: SimpleFunctionTransactor{contract: contract}, SimpleFunctionFilterer: SimpleFunctionFilterer{contract: contract}}, nil
}

// SimpleFunction is an auto generated Go binding around an Ethereum contract.
type SimpleFunction struct {
	SimpleFunctionCaller     // Read-only binding to the contract
	SimpleFunctionTransactor // Write-only binding to the contract
	SimpleFunctionFilterer   // Log filterer for contract events
}

// SimpleFunctionCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimpleFunctionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleFunctionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimpleFunctionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleFunctionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimpleFunctionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleFunctionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimpleFunctionSession struct {
	Contract     *SimpleFunction   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SimpleFunctionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimpleFunctionCallerSession struct {
	Contract *SimpleFunctionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// SimpleFunctionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimpleFunctionTransactorSession struct {
	Contract     *SimpleFunctionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// SimpleFunctionRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimpleFunctionRaw struct {
	Contract *SimpleFunction // Generic contract binding to access the raw methods on
}

// SimpleFunctionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimpleFunctionCallerRaw struct {
	Contract *SimpleFunctionCaller // Generic read-only contract binding to access the raw methods on
}

// SimpleFunctionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimpleFunctionTransactorRaw struct {
	Contract *SimpleFunctionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimpleFunction creates a new instance of SimpleFunction, bound to a specific deployed contract.
func NewSimpleFunction(address common.Address, backend bind.ContractBackend) (*SimpleFunction, error) {
	contract, err := bindSimpleFunction(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleFunction{SimpleFunctionCaller: SimpleFunctionCaller{contract: contract}, SimpleFunctionTransactor: SimpleFunctionTransactor{contract: contract}, SimpleFunctionFilterer: SimpleFunctionFilterer{contract: contract}}, nil
}

// NewSimpleFunctionCaller creates a new read-only instance of SimpleFunction, bound to a specific deployed contract.
func NewSimpleFunctionCaller(address common.Address, caller bind.ContractCaller) (*SimpleFunctionCaller, error) {
	contract, err := bindSimpleFunction(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleFunctionCaller{contract: contract}, nil
}

// NewSimpleFunctionTransactor creates a new write-only instance of SimpleFunction, bound to a specific deployed contract.
func NewSimpleFunctionTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleFunctionTransactor, error) {
	contract, err := bindSimpleFunction(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleFunctionTransactor{contract: contract}, nil
}

// NewSimpleFunctionFilterer creates a new log filterer instance of SimpleFunction, bound to a specific deployed contract.
func NewSimpleFunctionFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleFunctionFilterer, error) {
	contract, err := bindSimpleFunction(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleFunctionFilterer{contract: contract}, nil
}

// bindSimpleFunction binds a generic wrapper to an already deployed contract.
func bindSimpleFunction(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleFunctionABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleFunction *SimpleFunctionRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SimpleFunction.Contract.SimpleFunctionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleFunction *SimpleFunctionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleFunction.Contract.SimpleFunctionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleFunction *SimpleFunctionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleFunction.Contract.SimpleFunctionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleFunction *SimpleFunctionCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SimpleFunction.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleFunction *SimpleFunctionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleFunction.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleFunction *SimpleFunctionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleFunction.Contract.contract.Transact(opts, method, params...)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_SimpleFunction *SimpleFunctionCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleFunction.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_SimpleFunction *SimpleFunctionSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _SimpleFunction.Contract.Balances(&_SimpleFunction.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_SimpleFunction *SimpleFunctionCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _SimpleFunction.Contract.Balances(&_SimpleFunction.CallOpts, arg0)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_SimpleFunction *SimpleFunctionTransactor) Create(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _SimpleFunction.contract.Transact(opts, "create", newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_SimpleFunction *SimpleFunctionSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _SimpleFunction.Contract.Create(&_SimpleFunction.TransactOpts, newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_SimpleFunction *SimpleFunctionTransactorSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _SimpleFunction.Contract.Create(&_SimpleFunction.TransactOpts, newBalance)
}
