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

// TwoFunctionCallWithoutReturnABI is the input ABI used to generate the binding from.
const TwoFunctionCallWithoutReturnABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"create\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"update\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// TwoFunctionCallWithoutReturnBin is the compiled bytecode used for deploying new contracts.
const TwoFunctionCallWithoutReturnBin = `608060405234801561001057600080fd5b5060646000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610199806100646000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806327e235e314610046578063780900dc1461009e57806382ab890a146100cc575b600080fd5b6100886004803603602081101561005c57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506100fa565b6040518082815260200191505060405180910390f35b6100ca600480360360208110156100b457600080fd5b8101908080359060200190929190505050610112565b005b6100f8600480360360208110156100e257600080fd5b810190808035906020019092919050505061011e565b005b60006020528060005260406000206000915090505481565b61011b8161011e565b50565b806000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055505056fea265627a7a723058202830aa36d9cb855f2dabbc9feada6de45d3b02814f4c724f1751d64b4caaeff464736f6c63430005090032`

// DeployTwoFunctionCallWithoutReturn deploys a new Ethereum contract, binding an instance of TwoFunctionCallWithoutReturn to it.
func DeployTwoFunctionCallWithoutReturn(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TwoFunctionCallWithoutReturn, error) {
	parsed, err := abi.JSON(strings.NewReader(TwoFunctionCallWithoutReturnABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TwoFunctionCallWithoutReturnBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TwoFunctionCallWithoutReturn{TwoFunctionCallWithoutReturnCaller: TwoFunctionCallWithoutReturnCaller{contract: contract}, TwoFunctionCallWithoutReturnTransactor: TwoFunctionCallWithoutReturnTransactor{contract: contract}, TwoFunctionCallWithoutReturnFilterer: TwoFunctionCallWithoutReturnFilterer{contract: contract}}, nil
}

// TwoFunctionCallWithoutReturn is an auto generated Go binding around an Ethereum contract.
type TwoFunctionCallWithoutReturn struct {
	TwoFunctionCallWithoutReturnCaller     // Read-only binding to the contract
	TwoFunctionCallWithoutReturnTransactor // Write-only binding to the contract
	TwoFunctionCallWithoutReturnFilterer   // Log filterer for contract events
}

// TwoFunctionCallWithoutReturnCaller is an auto generated read-only Go binding around an Ethereum contract.
type TwoFunctionCallWithoutReturnCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoFunctionCallWithoutReturnTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TwoFunctionCallWithoutReturnTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoFunctionCallWithoutReturnFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TwoFunctionCallWithoutReturnFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoFunctionCallWithoutReturnSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TwoFunctionCallWithoutReturnSession struct {
	Contract     *TwoFunctionCallWithoutReturn // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// TwoFunctionCallWithoutReturnCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TwoFunctionCallWithoutReturnCallerSession struct {
	Contract *TwoFunctionCallWithoutReturnCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// TwoFunctionCallWithoutReturnTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TwoFunctionCallWithoutReturnTransactorSession struct {
	Contract     *TwoFunctionCallWithoutReturnTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// TwoFunctionCallWithoutReturnRaw is an auto generated low-level Go binding around an Ethereum contract.
type TwoFunctionCallWithoutReturnRaw struct {
	Contract *TwoFunctionCallWithoutReturn // Generic contract binding to access the raw methods on
}

// TwoFunctionCallWithoutReturnCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TwoFunctionCallWithoutReturnCallerRaw struct {
	Contract *TwoFunctionCallWithoutReturnCaller // Generic read-only contract binding to access the raw methods on
}

// TwoFunctionCallWithoutReturnTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TwoFunctionCallWithoutReturnTransactorRaw struct {
	Contract *TwoFunctionCallWithoutReturnTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTwoFunctionCallWithoutReturn creates a new instance of TwoFunctionCallWithoutReturn, bound to a specific deployed contract.
func NewTwoFunctionCallWithoutReturn(address common.Address, backend bind.ContractBackend) (*TwoFunctionCallWithoutReturn, error) {
	contract, err := bindTwoFunctionCallWithoutReturn(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionCallWithoutReturn{TwoFunctionCallWithoutReturnCaller: TwoFunctionCallWithoutReturnCaller{contract: contract}, TwoFunctionCallWithoutReturnTransactor: TwoFunctionCallWithoutReturnTransactor{contract: contract}, TwoFunctionCallWithoutReturnFilterer: TwoFunctionCallWithoutReturnFilterer{contract: contract}}, nil
}

// NewTwoFunctionCallWithoutReturnCaller creates a new read-only instance of TwoFunctionCallWithoutReturn, bound to a specific deployed contract.
func NewTwoFunctionCallWithoutReturnCaller(address common.Address, caller bind.ContractCaller) (*TwoFunctionCallWithoutReturnCaller, error) {
	contract, err := bindTwoFunctionCallWithoutReturn(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionCallWithoutReturnCaller{contract: contract}, nil
}

// NewTwoFunctionCallWithoutReturnTransactor creates a new write-only instance of TwoFunctionCallWithoutReturn, bound to a specific deployed contract.
func NewTwoFunctionCallWithoutReturnTransactor(address common.Address, transactor bind.ContractTransactor) (*TwoFunctionCallWithoutReturnTransactor, error) {
	contract, err := bindTwoFunctionCallWithoutReturn(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionCallWithoutReturnTransactor{contract: contract}, nil
}

// NewTwoFunctionCallWithoutReturnFilterer creates a new log filterer instance of TwoFunctionCallWithoutReturn, bound to a specific deployed contract.
func NewTwoFunctionCallWithoutReturnFilterer(address common.Address, filterer bind.ContractFilterer) (*TwoFunctionCallWithoutReturnFilterer, error) {
	contract, err := bindTwoFunctionCallWithoutReturn(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionCallWithoutReturnFilterer{contract: contract}, nil
}

// bindTwoFunctionCallWithoutReturn binds a generic wrapper to an already deployed contract.
func bindTwoFunctionCallWithoutReturn(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TwoFunctionCallWithoutReturnABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TwoFunctionCallWithoutReturn.Contract.TwoFunctionCallWithoutReturnCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TwoFunctionCallWithoutReturn.Contract.TwoFunctionCallWithoutReturnTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TwoFunctionCallWithoutReturn.Contract.TwoFunctionCallWithoutReturnTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TwoFunctionCallWithoutReturn.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TwoFunctionCallWithoutReturn.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TwoFunctionCallWithoutReturn.Contract.contract.Transact(opts, method, params...)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TwoFunctionCallWithoutReturn.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _TwoFunctionCallWithoutReturn.Contract.Balances(&_TwoFunctionCallWithoutReturn.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _TwoFunctionCallWithoutReturn.Contract.Balances(&_TwoFunctionCallWithoutReturn.CallOpts, arg0)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnTransactor) Create(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionCallWithoutReturn.contract.Transact(opts, "create", newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionCallWithoutReturn.Contract.Create(&_TwoFunctionCallWithoutReturn.TransactOpts, newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnTransactorSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionCallWithoutReturn.Contract.Create(&_TwoFunctionCallWithoutReturn.TransactOpts, newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnTransactor) Update(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionCallWithoutReturn.contract.Transact(opts, "update", newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnSession) Update(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionCallWithoutReturn.Contract.Update(&_TwoFunctionCallWithoutReturn.TransactOpts, newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns()
func (_TwoFunctionCallWithoutReturn *TwoFunctionCallWithoutReturnTransactorSession) Update(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionCallWithoutReturn.Contract.Update(&_TwoFunctionCallWithoutReturn.TransactOpts, newBalance)
}
