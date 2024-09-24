package ethereum_rpc_call

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type TransactionObject struct {
	From *string `json:"from"`
	To   string  `json:"to"`
	Data string  `json:"data"`
}

type Result struct {
	Error string          `json:"error,omitempty"`
	Data  json.RawMessage `json:"data,omitempty"`
}

func BlockNumberToString(n int64) string {
	switch n {
	case -4:
		return "safe"
	case -3:
		return "finalized"
	case -2:
		return "latest"
	case -1:
		return "pending"
	case 0:
		return "earliest"
	default:
		return hexutil.EncodeBig(big.NewInt(n))
	}
}
