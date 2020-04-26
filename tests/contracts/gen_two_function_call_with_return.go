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

// TwoFunctionCallWithReturnABI is the input ABI used to generate the binding from.
const TwoFunctionCallWithReturnABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"create\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"update\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// TwoFunctionCallWithReturnBin is the compiled bytecode used for deploying new contracts.
const TwoFunctionCallWithReturnBin = `608060405234801561001057600080fd5b5060646000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506101b3806100646000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806327e235e314610046578063780900dc1461009e57806382ab890a146100cc575b600080fd5b6100886004803603602081101561005c57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061010e565b6040518082815260200191505060405180910390f35b6100ca600480360360208110156100b457600080fd5b8101908080359060200190929190505050610126565b005b6100f8600480360360208110156100e257600080fd5b8101908080359060200190929190505050610174565b6040518082815260200191505060405180910390f35b60006020528060005260406000206000915090505481565b61012f81610174565b6000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555050565b600081905091905056fea265627a7a723058202c468fd937aebf4700bc0e2d6784cc2801eba25b97ffad6c6c4de1ece0861fba64736f6c63430005090032`

// DeployTwoFunctionCallWithReturn deploys a new Ethereum contract, binding an instance of TwoFunctionCallWithReturn to it.
func DeployTwoFunctionCallWithReturn(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TwoFunctionCallWithReturn, error) {
	parsed, err := abi.JSON(strings.NewReader(TwoFunctionCallWithReturnABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TwoFunctionCallWithReturnBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TwoFunctionCallWithReturn{TwoFunctionCallWithReturnCaller: TwoFunctionCallWithReturnCaller{contract: contract}, TwoFunctionCallWithReturnTransactor: TwoFunctionCallWithReturnTransactor{contract: contract}, TwoFunctionCallWithReturnFilterer: TwoFunctionCallWithReturnFilterer{contract: contract}}, nil
}

// TwoFunctionCallWithReturn is an auto generated Go binding around an Ethereum contract.
type TwoFunctionCallWithReturn struct {
	TwoFunctionCallWithReturnCaller     // Read-only binding to the contract
	TwoFunctionCallWithReturnTransactor // Write-only binding to the contract
	TwoFunctionCallWithReturnFilterer   // Log filterer for contract events
}

// TwoFunctionCallWithReturnCaller is an auto generated read-only Go binding around an Ethereum contract.
type TwoFunctionCallWithReturnCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoFunctionCallWithReturnTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TwoFunctionCallWithReturnTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoFunctionCallWithReturnFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TwoFunctionCallWithReturnFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TwoFunctionCallWithReturnSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TwoFunctionCallWithReturnSession struct {
	Contract     *TwoFunctionCallWithReturn // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// TwoFunctionCallWithReturnCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TwoFunctionCallWithReturnCallerSession struct {
	Contract *TwoFunctionCallWithReturnCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// TwoFunctionCallWithReturnTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TwoFunctionCallWithReturnTransactorSession struct {
	Contract     *TwoFunctionCallWithReturnTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// TwoFunctionCallWithReturnRaw is an auto generated low-level Go binding around an Ethereum contract.
type TwoFunctionCallWithReturnRaw struct {
	Contract *TwoFunctionCallWithReturn // Generic contract binding to access the raw methods on
}

// TwoFunctionCallWithReturnCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TwoFunctionCallWithReturnCallerRaw struct {
	Contract *TwoFunctionCallWithReturnCaller // Generic read-only contract binding to access the raw methods on
}

// TwoFunctionCallWithReturnTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TwoFunctionCallWithReturnTransactorRaw struct {
	Contract *TwoFunctionCallWithReturnTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTwoFunctionCallWithReturn creates a new instance of TwoFunctionCallWithReturn, bound to a specific deployed contract.
func NewTwoFunctionCallWithReturn(address common.Address, backend bind.ContractBackend) (*TwoFunctionCallWithReturn, error) {
	contract, err := bindTwoFunctionCallWithReturn(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionCallWithReturn{TwoFunctionCallWithReturnCaller: TwoFunctionCallWithReturnCaller{contract: contract}, TwoFunctionCallWithReturnTransactor: TwoFunctionCallWithReturnTransactor{contract: contract}, TwoFunctionCallWithReturnFilterer: TwoFunctionCallWithReturnFilterer{contract: contract}}, nil
}

// NewTwoFunctionCallWithReturnCaller creates a new read-only instance of TwoFunctionCallWithReturn, bound to a specific deployed contract.
func NewTwoFunctionCallWithReturnCaller(address common.Address, caller bind.ContractCaller) (*TwoFunctionCallWithReturnCaller, error) {
	contract, err := bindTwoFunctionCallWithReturn(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionCallWithReturnCaller{contract: contract}, nil
}

// NewTwoFunctionCallWithReturnTransactor creates a new write-only instance of TwoFunctionCallWithReturn, bound to a specific deployed contract.
func NewTwoFunctionCallWithReturnTransactor(address common.Address, transactor bind.ContractTransactor) (*TwoFunctionCallWithReturnTransactor, error) {
	contract, err := bindTwoFunctionCallWithReturn(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionCallWithReturnTransactor{contract: contract}, nil
}

// NewTwoFunctionCallWithReturnFilterer creates a new log filterer instance of TwoFunctionCallWithReturn, bound to a specific deployed contract.
func NewTwoFunctionCallWithReturnFilterer(address common.Address, filterer bind.ContractFilterer) (*TwoFunctionCallWithReturnFilterer, error) {
	contract, err := bindTwoFunctionCallWithReturn(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TwoFunctionCallWithReturnFilterer{contract: contract}, nil
}

// bindTwoFunctionCallWithReturn binds a generic wrapper to an already deployed contract.
func bindTwoFunctionCallWithReturn(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TwoFunctionCallWithReturnABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TwoFunctionCallWithReturn.Contract.TwoFunctionCallWithReturnCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TwoFunctionCallWithReturn.Contract.TwoFunctionCallWithReturnTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TwoFunctionCallWithReturn.Contract.TwoFunctionCallWithReturnTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TwoFunctionCallWithReturn.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TwoFunctionCallWithReturn.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TwoFunctionCallWithReturn.Contract.contract.Transact(opts, method, params...)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TwoFunctionCallWithReturn.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _TwoFunctionCallWithReturn.Contract.Balances(&_TwoFunctionCallWithReturn.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _TwoFunctionCallWithReturn.Contract.Balances(&_TwoFunctionCallWithReturn.CallOpts, arg0)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnTransactor) Create(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionCallWithReturn.contract.Transact(opts, "create", newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionCallWithReturn.Contract.Create(&_TwoFunctionCallWithReturn.TransactOpts, newBalance)
}

// Create is a paid mutator transaction binding the contract method 0x780900dc.
//
// Solidity: function create(uint256 newBalance) returns()
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnTransactorSession) Create(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionCallWithReturn.Contract.Create(&_TwoFunctionCallWithReturn.TransactOpts, newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns(uint256)
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnTransactor) Update(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionCallWithReturn.contract.Transact(opts, "update", newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns(uint256)
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnSession) Update(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionCallWithReturn.Contract.Update(&_TwoFunctionCallWithReturn.TransactOpts, newBalance)
}

// Update is a paid mutator transaction binding the contract method 0x82ab890a.
//
// Solidity: function update(uint256 newBalance) returns(uint256)
func (_TwoFunctionCallWithReturn *TwoFunctionCallWithReturnTransactorSession) Update(newBalance *big.Int) (*types.Transaction, error) {
	return _TwoFunctionCallWithReturn.Contract.Update(&_TwoFunctionCallWithReturn.TransactOpts, newBalance)
}
