package file

import (
	"os"
	"strings"

	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/indexed_abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type FileABIProvider struct {
	idx indexed_abi.IndexedABI
}

func NewFileABIProvider(path string) (*FileABIProvider, error) {
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

func (p *FileABIProvider) Event(selector string) (*abi.Event, error) {
	return p.idx.EventsSigHashIndex[selector], nil
}

func (p *FileABIProvider) Method(selector string) (*abi.Method, error) {
	return p.idx.MethodsSigHashIndex[selector], nil
}

func (p *FileABIProvider) Close() error {
	return nil
}
