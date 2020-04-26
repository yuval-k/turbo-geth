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

// TwoFunctionABI is the input ABI used to generate the binding from.
const TwoFunctionABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"create\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"update\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// TwoFunctionBin is the compiled bytecode used for deploying new contracts.
const TwoFunctionBin = `608060405234801561001057600080fd5b5060646000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506101d3806100646000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806327e235e314610046578063780900dc1461009e57806382ab890a146100cc575b600080fd5b6100886004803603602081101561005c57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506100fa565b6040518082815260200191505060405180910390f35b6100ca600480360360208110156100b457600080fd5b8101908080359060200190929190505050610112565b005b6100f8600480360360208110156100e257600080fd5b8101908080359060200190929190505050610158565b005b60006020528060005260406000206000915090505481565b806000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555050565b806000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055505056fea265627a7a72305820ef42a138e47ab8f8727bc093dabc59b2bd48407a340a7b39b295ea052068286364736f6c63430005090032`

// DeployTwoFunction deploys a new Ethereum contract, binding an instance of TwoFunction to it.
func DeployTwoFunction(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TwoFunction, error) {
	parsed, err := abi.JSON(strings.NewReader(TwoFunctionABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TwoFunctionBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TwoFunction{TwoFunctionCaller: TwoFunctionCaller{contract: contract}, TwoFunctionTransactor: TwoFunctionTransactor{contract: contract}, TwoFunctionFilterer: TwoFunctionFilterer{contract: contract}}, nil
}

// TwoFunction is an auto generated Go binding around an Ethereum contract.
type TwoFunction struct {
	TwoFunctionCaller     // Read-only binding to the contract
	TwoFunctionTransactor // Write-only binding to the contract
	TwoFunctionFilterer   // Log filterer for contract events
}

// TwoFunctionCaller is an auto generated read-only Go binding around an Ethereum contract.
type TwoFunctionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoFunctionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TwoFunctionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoFunctionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TwoFunctionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoFunctionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TwoFunctionSession struct {
	Contract     *TwoFunction      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TwoFunctionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TwoFunctionCallerSession struct {
	Contract *TwoFunctionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// TwoFunctionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TwoFunctionTransactorSession struct {
	Contract     *TwoFunctionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// TwoFunctionRaw is an auto generated low-level Go binding around an Ethereum contract.
type TwoFunctionRaw struct {
	Contract *TwoFunction // Generic contract binding to access the raw methods on
}

// TwoFunctionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TwoFunctionCallerRaw struct {
	Contract *TwoFunctionCaller // Generic read-only contract binding to access the raw methods on
}

// TwoFunctionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TwoFunctionTransactorRaw struct {
	Contract *TwoFunctionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTwoFunction creates a new instance of TwoFunction, bound to a specific deployed contract.
func NewTwoFunction(address common.Address, backend bind.ContractBackend) (*TwoFunction, error) {
	contract, err := bindTwoFunction(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TwoFunction{TwoFunctionCaller: TwoFunctionCaller{contract: contract}, TwoFunctionTransactor: TwoFunctionTransactor{contract: contract}, TwoFunctionFilterer: TwoFunctionFilterer{contract: contract}}, nil
}

// NewTwoFunctionCaller creates a new read-only instance of TwoFunction, bound to a specific deployed contract.
func NewTwoFunctionCaller(address common.Address, caller bind.ContractCaller) (*TwoFunctionCaller, error) {
	contract, err := bindTwoFunction(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionCaller{contract: contract}, nil
}

// NewTwoFunctionTransactor creates a new write-only instance of TwoFunction, bound to a specific deployed contract.
func NewTwoFunctionTransactor(address common.Address, transactor bind.ContractTransactor) (*TwoFunctionTransactor, error) {
	contract, err := bindTwoFunction(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionTransactor{contract: contract}, nil
}

// NewTwoFunctionFilterer creates a new log filterer instance of TwoFunction, bound to a specific deployed contract.
func NewTwoFunctionFilterer(address common.Address, filterer bind.ContractFilterer) (*TwoFunctionFilterer, error) {
	contract, err := bindTwoFunction(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionFilterer{contract: contract}, nil
}

// bindTwoFunction binds a generic wrapper to an already deployed contract.
func bindTwoFunction(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TwoFunctionABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TwoFunction *TwoFunctionRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TwoFunction.Contract.TwoFunctionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TwoFunction *TwoFunctionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TwoFunction.Contract.TwoFunctionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TwoFunction *TwoFunctionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TwoFunction.Contract.TwoFunctionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TwoFunction *TwoFunctionCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TwoFunction.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TwoFunction *TwoFunctionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TwoFunction.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TwoFunction *TwoFunctionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TwoFunction.Contract.contract.Transact(opts, method, params...)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TwoFunction *TwoFunctionCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TwoFunction.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TwoFunction *TwoFunctionSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _TwoFunction.Contract.Balances(&_TwoFunction.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TwoFunction *TwoFunctionCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _TwoFunction.Contract.Balances(&_TwoFunction.CallOpts, arg0)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_TwoFunction *TwoFunctionTransactor) Create(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunction.contract.Transact(opts, "create", newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_TwoFunction *TwoFunctionSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunction.Contract.Create(&_TwoFunction.TransactOpts, newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_TwoFunction *TwoFunctionTransactorSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunction.Contract.Create(&_TwoFunction.TransactOpts, newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_TwoFunction *TwoFunctionTransactor) Update(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunction.contract.Transact(opts, "update", newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_TwoFunction *TwoFunctionSession) Update(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunction.Contract.Update(&_TwoFunction.TransactOpts, newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_TwoFunction *TwoFunctionTransactorSession) Update(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunction.Contract.Update(&_TwoFunction.TransactOpts, newBalance)
}
