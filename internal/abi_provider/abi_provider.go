package abi_provider

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type ABIProvider interface {
	Event(selector string) (*abi.Event, error)
	Method(selector string) (*abi.Method, error)
	Close() error
}
