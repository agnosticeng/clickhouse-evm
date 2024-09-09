package evm_decode_event

import (
	"io"

	eth_abi "github.com/ethereum/go-ethereum/accounts/abi"
)

func ParseIndexedABI(r io.Reader) (IndexedABI, error) {
	_abi, err := eth_abi.JSON(r)

	if err != nil {
		return IndexedABI{}, err
	}

	var res IndexedABI
	res.ABI = _abi
	res.EventsSigHashIndex = make(map[string]*eth_abi.Event)
	res.MethodsSigHashIndex = make(map[string]*eth_abi.Method)

	for _, event := range _abi.Events {
		res.EventsSigHashIndex[string(event.ID.Bytes())] = &event
	}

	for _, method := range _abi.Methods {
		res.MethodsSigHashIndex[string(method.ID)] = &method
	}

	return res, nil
}

type IndexedABI struct {
	eth_abi.ABI
	EventsSigHashIndex  map[string]*eth_abi.Event
	MethodsSigHashIndex map[string]*eth_abi.Method
}
