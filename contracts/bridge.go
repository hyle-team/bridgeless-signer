// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BridgeMetaData contains all meta data concerning the Bridge contract.
var BridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"srcAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"dstAddress\",\"type\":\"string\"}],\"name\":\"BridgeIn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dstAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"originTxId\",\"type\":\"string\"}],\"name\":\"BridgeOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowList\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_dst\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"_isWrapped\",\"type\":\"bool\"}],\"name\":\"bridgeIn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_dst\",\"type\":\"string\"}],\"name\":\"bridgeInNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_originTxId\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"_isWrapped\",\"type\":\"bool\"}],\"name\":\"bridgeOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_originTxId\",\"type\":\"string\"}],\"name\":\"bridgeOutNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"operator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// BridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeMetaData.ABI instead.
var BridgeABI = BridgeMetaData.ABI

// Bridge is an auto generated Go binding around an Ethereum contract.
type Bridge struct {
	BridgeCaller     // Read-only binding to the contract
	BridgeTransactor // Write-only binding to the contract
	BridgeFilterer   // Log filterer for contract events
}

// BridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeSession struct {
	Contract     *Bridge           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeCallerSession struct {
	Contract *BridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeTransactorSession struct {
	Contract     *BridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeRaw struct {
	Contract *Bridge // Generic contract binding to access the raw methods on
}

// BridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeCallerRaw struct {
	Contract *BridgeCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeTransactorRaw struct {
	Contract *BridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridge creates a new instance of Bridge, bound to a specific deployed contract.
func NewBridge(address common.Address, backend bind.ContractBackend) (*Bridge, error) {
	contract, err := bindBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bridge{BridgeCaller: BridgeCaller{contract: contract}, BridgeTransactor: BridgeTransactor{contract: contract}, BridgeFilterer: BridgeFilterer{contract: contract}}, nil
}

// NewBridgeCaller creates a new read-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeCaller(address common.Address, caller bind.ContractCaller) (*BridgeCaller, error) {
	contract, err := bindBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeCaller{contract: contract}, nil
}

// NewBridgeTransactor creates a new write-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeTransactor, error) {
	contract, err := bindBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeTransactor{contract: contract}, nil
}

// NewBridgeFilterer creates a new log filterer instance of Bridge, bound to a specific deployed contract.
func NewBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeFilterer, error) {
	contract, err := bindBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeFilterer{contract: contract}, nil
}

// bindBridge binds a generic wrapper to an already deployed contract.
func bindBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.BridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transact(opts, method, params...)
}

// AllowList is a free data retrieval call binding the contract method 0x2848aeaf.
//
// Solidity: function allowList(address ) view returns(bool)
func (_Bridge *BridgeCaller) AllowList(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "allowList", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowList is a free data retrieval call binding the contract method 0x2848aeaf.
//
// Solidity: function allowList(address ) view returns(bool)
func (_Bridge *BridgeSession) AllowList(arg0 common.Address) (bool, error) {
	return _Bridge.Contract.AllowList(&_Bridge.CallOpts, arg0)
}

// AllowList is a free data retrieval call binding the contract method 0x2848aeaf.
//
// Solidity: function allowList(address ) view returns(bool)
func (_Bridge *BridgeCallerSession) AllowList(arg0 common.Address) (bool, error) {
	return _Bridge.Contract.AllowList(&_Bridge.CallOpts, arg0)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_Bridge *BridgeCaller) Operator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "operator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_Bridge *BridgeSession) Operator() (common.Address, error) {
	return _Bridge.Contract.Operator(&_Bridge.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_Bridge *BridgeCallerSession) Operator() (common.Address, error) {
	return _Bridge.Contract.Operator(&_Bridge.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bridge *BridgeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bridge *BridgeSession) Owner() (common.Address, error) {
	return _Bridge.Contract.Owner(&_Bridge.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bridge *BridgeCallerSession) Owner() (common.Address, error) {
	return _Bridge.Contract.Owner(&_Bridge.CallOpts)
}

// BridgeIn is a paid mutator transaction binding the contract method 0x4a5368bf.
//
// Solidity: function bridgeIn(address _token, uint256 _amount, uint256 _chainId, string _dst, bool _isWrapped) returns()
func (_Bridge *BridgeTransactor) BridgeIn(opts *bind.TransactOpts, _token common.Address, _amount *big.Int, _chainId *big.Int, _dst string, _isWrapped bool) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "bridgeIn", _token, _amount, _chainId, _dst, _isWrapped)
}

// BridgeIn is a paid mutator transaction binding the contract method 0x4a5368bf.
//
// Solidity: function bridgeIn(address _token, uint256 _amount, uint256 _chainId, string _dst, bool _isWrapped) returns()
func (_Bridge *BridgeSession) BridgeIn(_token common.Address, _amount *big.Int, _chainId *big.Int, _dst string, _isWrapped bool) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeIn(&_Bridge.TransactOpts, _token, _amount, _chainId, _dst, _isWrapped)
}

// BridgeIn is a paid mutator transaction binding the contract method 0x4a5368bf.
//
// Solidity: function bridgeIn(address _token, uint256 _amount, uint256 _chainId, string _dst, bool _isWrapped) returns()
func (_Bridge *BridgeTransactorSession) BridgeIn(_token common.Address, _amount *big.Int, _chainId *big.Int, _dst string, _isWrapped bool) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeIn(&_Bridge.TransactOpts, _token, _amount, _chainId, _dst, _isWrapped)
}

// BridgeInNative is a paid mutator transaction binding the contract method 0x93ce02da.
//
// Solidity: function bridgeInNative(uint256 _chainId, string _dst) payable returns()
func (_Bridge *BridgeTransactor) BridgeInNative(opts *bind.TransactOpts, _chainId *big.Int, _dst string) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "bridgeInNative", _chainId, _dst)
}

// BridgeInNative is a paid mutator transaction binding the contract method 0x93ce02da.
//
// Solidity: function bridgeInNative(uint256 _chainId, string _dst) payable returns()
func (_Bridge *BridgeSession) BridgeInNative(_chainId *big.Int, _dst string) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeInNative(&_Bridge.TransactOpts, _chainId, _dst)
}

// BridgeInNative is a paid mutator transaction binding the contract method 0x93ce02da.
//
// Solidity: function bridgeInNative(uint256 _chainId, string _dst) payable returns()
func (_Bridge *BridgeTransactorSession) BridgeInNative(_chainId *big.Int, _dst string) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeInNative(&_Bridge.TransactOpts, _chainId, _dst)
}

// BridgeOut is a paid mutator transaction binding the contract method 0xb9fe6843.
//
// Solidity: function bridgeOut(address _token, address _receiver, uint256 _amount, string _originTxId, bool _isWrapped) returns()
func (_Bridge *BridgeTransactor) BridgeOut(opts *bind.TransactOpts, _token common.Address, _receiver common.Address, _amount *big.Int, _originTxId string, _isWrapped bool) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "bridgeOut", _token, _receiver, _amount, _originTxId, _isWrapped)
}

// BridgeOut is a paid mutator transaction binding the contract method 0xb9fe6843.
//
// Solidity: function bridgeOut(address _token, address _receiver, uint256 _amount, string _originTxId, bool _isWrapped) returns()
func (_Bridge *BridgeSession) BridgeOut(_token common.Address, _receiver common.Address, _amount *big.Int, _originTxId string, _isWrapped bool) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeOut(&_Bridge.TransactOpts, _token, _receiver, _amount, _originTxId, _isWrapped)
}

// BridgeOut is a paid mutator transaction binding the contract method 0xb9fe6843.
//
// Solidity: function bridgeOut(address _token, address _receiver, uint256 _amount, string _originTxId, bool _isWrapped) returns()
func (_Bridge *BridgeTransactorSession) BridgeOut(_token common.Address, _receiver common.Address, _amount *big.Int, _originTxId string, _isWrapped bool) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeOut(&_Bridge.TransactOpts, _token, _receiver, _amount, _originTxId, _isWrapped)
}

// BridgeOutNative is a paid mutator transaction binding the contract method 0x07bb49bc.
//
// Solidity: function bridgeOutNative(address _receiver, uint256 _amount, string _originTxId) returns()
func (_Bridge *BridgeTransactor) BridgeOutNative(opts *bind.TransactOpts, _receiver common.Address, _amount *big.Int, _originTxId string) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "bridgeOutNative", _receiver, _amount, _originTxId)
}

// BridgeOutNative is a paid mutator transaction binding the contract method 0x07bb49bc.
//
// Solidity: function bridgeOutNative(address _receiver, uint256 _amount, string _originTxId) returns()
func (_Bridge *BridgeSession) BridgeOutNative(_receiver common.Address, _amount *big.Int, _originTxId string) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeOutNative(&_Bridge.TransactOpts, _receiver, _amount, _originTxId)
}

// BridgeOutNative is a paid mutator transaction binding the contract method 0x07bb49bc.
//
// Solidity: function bridgeOutNative(address _receiver, uint256 _amount, string _originTxId) returns()
func (_Bridge *BridgeTransactorSession) BridgeOutNative(_receiver common.Address, _amount *big.Int, _originTxId string) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeOutNative(&_Bridge.TransactOpts, _receiver, _amount, _originTxId)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bridge *BridgeTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bridge *BridgeSession) RenounceOwnership() (*types.Transaction, error) {
	return _Bridge.Contract.RenounceOwnership(&_Bridge.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bridge *BridgeTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Bridge.Contract.RenounceOwnership(&_Bridge.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bridge *BridgeTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bridge *BridgeSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TransferOwnership(&_Bridge.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bridge *BridgeTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TransferOwnership(&_Bridge.TransactOpts, newOwner)
}

// BridgeBridgeInIterator is returned from FilterBridgeIn and is used to iterate over the raw logs and unpacked data for BridgeIn events raised by the Bridge contract.
type BridgeBridgeInIterator struct {
	Event *BridgeBridgeIn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeBridgeInIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeBridgeIn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeBridgeIn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeBridgeInIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeBridgeInIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeBridgeIn represents a BridgeIn event raised by the Bridge contract.
type BridgeBridgeIn struct {
	Token      common.Address
	SrcAddress common.Address
	Amount     *big.Int
	ChainId    *big.Int
	DstAddress string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBridgeIn is a free log retrieval operation binding the contract event 0xb682b0d4c459c3870ca6b0531dfa5ed978c7136fcafb37aaf1a123edb44ccd4f.
//
// Solidity: event BridgeIn(address indexed token, address srcAddress, uint256 amount, uint256 chainId, string dstAddress)
func (_Bridge *BridgeFilterer) FilterBridgeIn(opts *bind.FilterOpts, token []common.Address) (*BridgeBridgeInIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "BridgeIn", tokenRule)
	if err != nil {
		return nil, err
	}
	return &BridgeBridgeInIterator{contract: _Bridge.contract, event: "BridgeIn", logs: logs, sub: sub}, nil
}

// WatchBridgeIn is a free log subscription operation binding the contract event 0xb682b0d4c459c3870ca6b0531dfa5ed978c7136fcafb37aaf1a123edb44ccd4f.
//
// Solidity: event BridgeIn(address indexed token, address srcAddress, uint256 amount, uint256 chainId, string dstAddress)
func (_Bridge *BridgeFilterer) WatchBridgeIn(opts *bind.WatchOpts, sink chan<- *BridgeBridgeIn, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "BridgeIn", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeBridgeIn)
				if err := _Bridge.contract.UnpackLog(event, "BridgeIn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBridgeIn is a log parse operation binding the contract event 0xb682b0d4c459c3870ca6b0531dfa5ed978c7136fcafb37aaf1a123edb44ccd4f.
//
// Solidity: event BridgeIn(address indexed token, address srcAddress, uint256 amount, uint256 chainId, string dstAddress)
func (_Bridge *BridgeFilterer) ParseBridgeIn(log types.Log) (*BridgeBridgeIn, error) {
	event := new(BridgeBridgeIn)
	if err := _Bridge.contract.UnpackLog(event, "BridgeIn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeBridgeOutIterator is returned from FilterBridgeOut and is used to iterate over the raw logs and unpacked data for BridgeOut events raised by the Bridge contract.
type BridgeBridgeOutIterator struct {
	Event *BridgeBridgeOut // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeBridgeOutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeBridgeOut)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeBridgeOut)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeBridgeOutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeBridgeOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeBridgeOut represents a BridgeOut event raised by the Bridge contract.
type BridgeBridgeOut struct {
	Token      common.Address
	DstAddress common.Address
	Amount     *big.Int
	OriginTxId string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBridgeOut is a free log retrieval operation binding the contract event 0x0cc532d34ef13618ce9f3b733023adcb6dd0edd643d714296e414870f627e8b9.
//
// Solidity: event BridgeOut(address indexed token, address dstAddress, uint256 amount, string originTxId)
func (_Bridge *BridgeFilterer) FilterBridgeOut(opts *bind.FilterOpts, token []common.Address) (*BridgeBridgeOutIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "BridgeOut", tokenRule)
	if err != nil {
		return nil, err
	}
	return &BridgeBridgeOutIterator{contract: _Bridge.contract, event: "BridgeOut", logs: logs, sub: sub}, nil
}

// WatchBridgeOut is a free log subscription operation binding the contract event 0x0cc532d34ef13618ce9f3b733023adcb6dd0edd643d714296e414870f627e8b9.
//
// Solidity: event BridgeOut(address indexed token, address dstAddress, uint256 amount, string originTxId)
func (_Bridge *BridgeFilterer) WatchBridgeOut(opts *bind.WatchOpts, sink chan<- *BridgeBridgeOut, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "BridgeOut", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeBridgeOut)
				if err := _Bridge.contract.UnpackLog(event, "BridgeOut", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBridgeOut is a log parse operation binding the contract event 0x0cc532d34ef13618ce9f3b733023adcb6dd0edd643d714296e414870f627e8b9.
//
// Solidity: event BridgeOut(address indexed token, address dstAddress, uint256 amount, string originTxId)
func (_Bridge *BridgeFilterer) ParseBridgeOut(log types.Log) (*BridgeBridgeOut, error) {
	event := new(BridgeBridgeOut)
	if err := _Bridge.contract.UnpackLog(event, "BridgeOut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Bridge contract.
type BridgeOwnershipTransferredIterator struct {
	Event *BridgeOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeOwnershipTransferred represents a OwnershipTransferred event raised by the Bridge contract.
type BridgeOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Bridge *BridgeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BridgeOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BridgeOwnershipTransferredIterator{contract: _Bridge.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Bridge *BridgeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BridgeOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeOwnershipTransferred)
				if err := _Bridge.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Bridge *BridgeFilterer) ParseOwnershipTransferred(log types.Log) (*BridgeOwnershipTransferred, error) {
	event := new(BridgeOwnershipTransferred)
	if err := _Bridge.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
