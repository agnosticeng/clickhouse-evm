package noop

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type NoopABIProvider struct{}

func NewNoopABIProvider() *NoopABIProvider {
	return &NoopABIProvider{}
}

func (p *NoopABIProvider) Event(selector string) (*abi.Event, error) {
	return nil, nil
}

func (p *NoopABIProvider) Method(selector string) (*abi.Method, error) {
	return nil, nil
}

func (p *NoopABIProvider) Close() error {
	return nil
}
