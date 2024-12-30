package fullsig

import (
	"fmt"
	"strings"

	"github.com/agnosticeng/evmabi/fullsig"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type FullsigABIProvider struct {
	event  *abi.Event
	method *abi.Method
}

func NewFullsigABIProvider(s string) (*FullsigABIProvider, error) {
	var p FullsigABIProvider

	if len(s) == 0 {
		return nil, fmt.Errorf("fullsig must not be empty")
	}

	switch {
	case strings.HasPrefix(s, "event"):
		evt, err := fullsig.ParseEvent(s)

		if err != nil {
			return nil, err
		}

		p.event = &evt

	case strings.HasPrefix(s, "function"):
		meth, err := fullsig.ParseMethod(s)

		if err != nil {
			return nil, err
		}

		p.method = &meth

	default:
		return nil, fmt.Errorf("unknown fullsig kind: %s", s)
	}

	return &p, nil
}

func (p *FullsigABIProvider) Events(selector string) ([]*abi.Event, error) {
	if p.event == nil || selector != string(p.event.ID.Bytes()) {
		return nil, nil
	}

	return []*abi.Event{p.event}, nil
}

func (p *FullsigABIProvider) Methods(selector string) ([]*abi.Method, error) {
	if p.method == nil || selector != string(p.method.ID) {
		return nil, nil
	}

	return []*abi.Method{p.method}, nil
}

func (p *FullsigABIProvider) Close() error {
	return nil
}
