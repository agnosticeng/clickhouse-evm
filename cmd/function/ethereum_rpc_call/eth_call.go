package ethereum_rpc_call

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type TransactionObject struct {
	From *string `json:"from"`
	To   string  `json:"to"`
	Data string  `json:"data"`
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
