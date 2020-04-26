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

// OnlyConstructorABI is the input ABI used to generate the binding from.
const OnlyConstructorABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// OnlyConstructorBin is the compiled bytecode used for deploying new contracts.
const OnlyConstructorBin = `6080604052348015600f57600080fd5b50603e80601d6000396000f3fe6080604052600080fdfea265627a7a72305820134c4b2e067be235c3336e38c81f5243fbcfbfabbd99d04562b666d2a3be885864736f6c63430005090032`

// DeployOnlyConstructor deploys a new Ethereum contract, binding an instance of OnlyConstructor to it.
func DeployOnlyConstructor(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OnlyConstructor, error) {
	parsed, err := abi.JSON(strings.NewReader(OnlyConstructorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OnlyConstructorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OnlyConstructor{OnlyConstructorCaller: OnlyConstructorCaller{contract: contract}, OnlyConstructorTransactor: OnlyConstructorTransactor{contract: contract}, OnlyConstructorFilterer: OnlyConstructorFilterer{contract: contract}}, nil
}

// OnlyConstructor is an auto generated Go binding around an Ethereum contract.
type OnlyConstructor struct {
	OnlyConstructorCaller     // Read-only binding to the contract
	OnlyConstructorTransactor // Write-only binding to the contract
	OnlyConstructorFilterer   // Log filterer for contract events
}

// OnlyConstructorCaller is an auto generated read-only Go binding around an Ethereum contract.
type OnlyConstructorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnlyConstructorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OnlyConstructorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnlyConstructorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OnlyConstructorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnlyConstructorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OnlyConstructorSession struct {
	Contract     *OnlyConstructor  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OnlyConstructorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OnlyConstructorCallerSession struct {
	Contract *OnlyConstructorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// OnlyConstructorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OnlyConstructorTransactorSession struct {
	Contract     *OnlyConstructorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// OnlyConstructorRaw is an auto generated low-level Go binding around an Ethereum contract.
type OnlyConstructorRaw struct {
	Contract *OnlyConstructor // Generic contract binding to access the raw methods on
}

// OnlyConstructorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OnlyConstructorCallerRaw struct {
	Contract *OnlyConstructorCaller // Generic read-only contract binding to access the raw methods on
}

// OnlyConstructorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OnlyConstructorTransactorRaw struct {
	Contract *OnlyConstructorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOnlyConstructor creates a new instance of OnlyConstructor, bound to a specific deployed contract.
func NewOnlyConstructor(address common.Address, backend bind.ContractBackend) (*OnlyConstructor, error) {
	contract, err := bindOnlyConstructor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OnlyConstructor{OnlyConstructorCaller: OnlyConstructorCaller{contract: contract}, OnlyConstructorTransactor: OnlyConstructorTransactor{contract: contract}, OnlyConstructorFilterer: OnlyConstructorFilterer{contract: contract}}, nil
}

// NewOnlyConstructorCaller creates a new read-only instance of OnlyConstructor, bound to a specific deployed contract.
func NewOnlyConstructorCaller(address common.Address, caller bind.ContractCaller) (*OnlyConstructorCaller, error) {
	contract, err := bindOnlyConstructor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OnlyConstructorCaller{contract: contract}, nil
}

// NewOnlyConstructorTransactor creates a new write-only instance of OnlyConstructor, bound to a specific deployed contract.
func NewOnlyConstructorTransactor(address common.Address, transactor bind.ContractTransactor) (*OnlyConstructorTransactor, error) {
	contract, err := bindOnlyConstructor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OnlyConstructorTransactor{contract: contract}, nil
}

// NewOnlyConstructorFilterer creates a new log filterer instance of OnlyConstructor, bound to a specific deployed contract.
func NewOnlyConstructorFilterer(address common.Address, filterer bind.ContractFilterer) (*OnlyConstructorFilterer, error) {
	contract, err := bindOnlyConstructor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OnlyConstructorFilterer{contract: contract}, nil
}

// bindOnlyConstructor binds a generic wrapper to an already deployed contract.
func bindOnlyConstructor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OnlyConstructorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OnlyConstructor *OnlyConstructorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OnlyConstructor.Contract.OnlyConstructorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OnlyConstructor *OnlyConstructorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnlyConstructor.Contract.OnlyConstructorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OnlyConstructor *OnlyConstructorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnlyConstructor.Contract.OnlyConstructorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OnlyConstructor *OnlyConstructorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OnlyConstructor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OnlyConstructor *OnlyConstructorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OnlyConstructor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OnlyConstructor *OnlyConstructorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OnlyConstructor.Contract.contract.Transact(opts, method, params...)
}
