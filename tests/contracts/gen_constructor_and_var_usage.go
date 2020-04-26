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

// ConstructorAndVarUsageABI is the input ABI used to generate the binding from.
const ConstructorAndVarUsageABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// ConstructorAndVarUsageBin is the compiled bytecode used for deploying new contracts.
const ConstructorAndVarUsageBin = `608060405234801561001057600080fd5b5060646000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555060cf806100636000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806327e235e314602d575b600080fd5b606c60048036036020811015604157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506082565b6040518082815260200191505060405180910390f35b6000602052806000526040600020600091509050548156fea265627a7a72305820b7ee699ba48acc615b2b7c88503bf66fbcbcfb510dd1ad919d564756b75e261c64736f6c63430005090032`

// DeployConstructorAndVarUsage deploys a new Ethereum contract, binding an instance of ConstructorAndVarUsage to it.
func DeployConstructorAndVarUsage(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ConstructorAndVarUsage, error) {
	parsed, err := abi.JSON(strings.NewReader(ConstructorAndVarUsageABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ConstructorAndVarUsageBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConstructorAndVarUsage{ConstructorAndVarUsageCaller: ConstructorAndVarUsageCaller{contract: contract}, ConstructorAndVarUsageTransactor: ConstructorAndVarUsageTransactor{contract: contract}, ConstructorAndVarUsageFilterer: ConstructorAndVarUsageFilterer{contract: contract}}, nil
}

// ConstructorAndVarUsage is an auto generated Go binding around an Ethereum contract.
type ConstructorAndVarUsage struct {
	ConstructorAndVarUsageCaller     // Read-only binding to the contract
	ConstructorAndVarUsageTransactor // Write-only binding to the contract
	ConstructorAndVarUsageFilterer   // Log filterer for contract events
}

// ConstructorAndVarUsageCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConstructorAndVarUsageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstructorAndVarUsageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConstructorAndVarUsageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstructorAndVarUsageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConstructorAndVarUsageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstructorAndVarUsageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConstructorAndVarUsageSession struct {
	Contract     *ConstructorAndVarUsage // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ConstructorAndVarUsageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConstructorAndVarUsageCallerSession struct {
	Contract *ConstructorAndVarUsageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// ConstructorAndVarUsageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConstructorAndVarUsageTransactorSession struct {
	Contract     *ConstructorAndVarUsageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ConstructorAndVarUsageRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConstructorAndVarUsageRaw struct {
	Contract *ConstructorAndVarUsage // Generic contract binding to access the raw methods on
}

// ConstructorAndVarUsageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConstructorAndVarUsageCallerRaw struct {
	Contract *ConstructorAndVarUsageCaller // Generic read-only contract binding to access the raw methods on
}

// ConstructorAndVarUsageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConstructorAndVarUsageTransactorRaw struct {
	Contract *ConstructorAndVarUsageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConstructorAndVarUsage creates a new instance of ConstructorAndVarUsage, bound to a specific deployed contract.
func NewConstructorAndVarUsage(address common.Address, backend bind.ContractBackend) (*ConstructorAndVarUsage, error) {
	contract, err := bindConstructorAndVarUsage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConstructorAndVarUsage{ConstructorAndVarUsageCaller: ConstructorAndVarUsageCaller{contract: contract}, ConstructorAndVarUsageTransactor: ConstructorAndVarUsageTransactor{contract: contract}, ConstructorAndVarUsageFilterer: ConstructorAndVarUsageFilterer{contract: contract}}, nil
}

// NewConstructorAndVarUsageCaller creates a new read-only instance of ConstructorAndVarUsage, bound to a specific deployed contract.
func NewConstructorAndVarUsageCaller(address common.Address, caller bind.ContractCaller) (*ConstructorAndVarUsageCaller, error) {
	contract, err := bindConstructorAndVarUsage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConstructorAndVarUsageCaller{contract: contract}, nil
}

// NewConstructorAndVarUsageTransactor creates a new write-only instance of ConstructorAndVarUsage, bound to a specific deployed contract.
func NewConstructorAndVarUsageTransactor(address common.Address, transactor bind.ContractTransactor) (*ConstructorAndVarUsageTransactor, error) {
	contract, err := bindConstructorAndVarUsage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConstructorAndVarUsageTransactor{contract: contract}, nil
}

// NewConstructorAndVarUsageFilterer creates a new log filterer instance of ConstructorAndVarUsage, bound to a specific deployed contract.
func NewConstructorAndVarUsageFilterer(address common.Address, filterer bind.ContractFilterer) (*ConstructorAndVarUsageFilterer, error) {
	contract, err := bindConstructorAndVarUsage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConstructorAndVarUsageFilterer{contract: contract}, nil
}

// bindConstructorAndVarUsage binds a generic wrapper to an already deployed contract.
func bindConstructorAndVarUsage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConstructorAndVarUsageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstructorAndVarUsage *ConstructorAndVarUsageRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ConstructorAndVarUsage.Contract.ConstructorAndVarUsageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstructorAndVarUsage *ConstructorAndVarUsageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstructorAndVarUsage.Contract.ConstructorAndVarUsageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstructorAndVarUsage *ConstructorAndVarUsageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstructorAndVarUsage.Contract.ConstructorAndVarUsageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstructorAndVarUsage *ConstructorAndVarUsageCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ConstructorAndVarUsage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstructorAndVarUsage *ConstructorAndVarUsageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstructorAndVarUsage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstructorAndVarUsage *ConstructorAndVarUsageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstructorAndVarUsage.Contract.contract.Transact(opts, method, params...)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_ConstructorAndVarUsage *ConstructorAndVarUsageCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ConstructorAndVarUsage.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_ConstructorAndVarUsage *ConstructorAndVarUsageSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _ConstructorAndVarUsage.Contract.Balances(&_ConstructorAndVarUsage.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_ConstructorAndVarUsage *ConstructorAndVarUsageCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _ConstructorAndVarUsage.Contract.Balances(&_ConstructorAndVarUsage.CallOpts, arg0)
}
