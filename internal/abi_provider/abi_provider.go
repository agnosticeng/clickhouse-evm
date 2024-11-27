package abi_provider

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type ABIProvider interface {
	Events(selector string) ([]*abi.Event, error)
	Methods(selector string) ([]*abi.Method, error)
	Close() error
}
