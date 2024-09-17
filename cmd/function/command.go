package function

import (
	ethereum_rpc "github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function/ethreum_rpc"
	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function/evm_decode_call"
	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function/evm_decode_event"
	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function/keccak256"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name: "function",
		Subcommands: []*cli.Command{
			keccak256.Command(),
			evm_decode_event.Command(),
			evm_decode_call.Command(),
			ethereum_rpc.Command(),
		},
	}
}
