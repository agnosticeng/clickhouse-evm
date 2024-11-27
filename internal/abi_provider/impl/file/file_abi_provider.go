package file

import (
	"net/http"
	"os"
	"strings"

	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/indexed_abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type FileABIProvider struct {
	idx indexed_abi.IndexedABI
}

func (p *FileABIProvider) Events(selector string) ([]*abi.Event, error) {
	var evt = p.idx.EventsSigHashIndex[selector]

	if evt == nil {
		return nil, nil
	} else {
		return []*abi.Event{evt}, nil
	}
}

func (p *FileABIProvider) Methods(selector string) ([]*abi.Method, error) {
	var meth = p.idx.MethodsSigHashIndex[selector]

	if meth == nil {
		return nil, nil
	} else {
		return []*abi.Method{meth}, nil
	}
}

func (p *FileABIProvider) Close() error {
	return nil
}

func FromPath(path string) (*FileABIProvider, error) {
	path = strings.TrimPrefix(path, "file://")

	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	idx, err := indexed_abi.JSON(f)

	if err != nil {
		return nil, err
	}

	return &FileABIProvider{
		idx: idx,
	}, nil
}

func FromURL(path string) (*FileABIProvider, error) {
	resp, err := http.Get(path)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	idx, err := indexed_abi.JSON(resp.Body)

	if err != nil {
		return nil, err
	}

	return &FileABIProvider{
		idx: idx,
	}, nil
}
