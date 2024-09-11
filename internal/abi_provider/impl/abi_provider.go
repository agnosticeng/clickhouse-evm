package impl

import (
	"strings"

	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/abi_provider"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/abi_provider/impl/file"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/abi_provider/impl/fullsig"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/abi_provider/impl/http"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/abi_provider/impl/noop"
)

func NewABIProvider(s string) (abi_provider.ABIProvider, error) {
	switch {
	case strings.HasPrefix(s, "file://"):
		return file.NewFileABIProvider(s)
	case strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https"):
		return http.NewHTTPABIProvider(s)
	case len(s) > 0:
		return fullsig.NewFullsigABIProvider(s)
	default:
		return noop.NewNoopABIProvider(), nil
	}
}
