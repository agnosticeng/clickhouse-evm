package function

import (
	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function/ethereum_rpc_call"
	ethereum_rpc "github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function/ethreum_rpc"
	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function/evm_decode_call"
	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function/evm_decode_event"
	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function/evm_descriptor_from_fullsig"
	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function/evm_signature_from_descriptor"
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
			ethereum_rpc_call.Command(),
			evm_signature_from_descriptor.Command(),
			evm_descriptor_from_fullsig.Command(),
		},
	}
}
