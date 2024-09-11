package http

import (
	"net/http"

	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/indexed_abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type HTTPABIProvider struct {
	idx indexed_abi.IndexedABI
}

func NewHTTPABIProvider(path string) (*HTTPABIProvider, error) {
	resp, err := http.Get(path)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	idx, err := indexed_abi.JSON(resp.Body)

	if err != nil {
		return nil, err
	}

	return &HTTPABIProvider{
		idx: idx,
	}, nil
}

func (p *HTTPABIProvider) Event(selector string) (*abi.Event, error) {
	return p.idx.EventsSigHashIndex[selector], nil
}

func (p *HTTPABIProvider) Method(selector string) (*abi.Method, error) {
	return p.idx.MethodsSigHashIndex[selector], nil
}

func (p *HTTPABIProvider) Close() error {
	return nil
}
