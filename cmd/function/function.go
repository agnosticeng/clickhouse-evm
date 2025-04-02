package function

import (
	"github.com/agnosticeng/clickhouse-evm/cmd/function/convert_format"
	"github.com/agnosticeng/clickhouse-evm/cmd/function/ethereum_decode_tx"
	"github.com/agnosticeng/clickhouse-evm/cmd/function/ethereum_rpc_call"
	ethereum_rpc "github.com/agnosticeng/clickhouse-evm/cmd/function/ethreum_rpc"
	"github.com/agnosticeng/clickhouse-evm/cmd/function/evm_decode_call"
	"github.com/agnosticeng/clickhouse-evm/cmd/function/evm_decode_calldata"
	"github.com/agnosticeng/clickhouse-evm/cmd/function/evm_decode_event"
	"github.com/agnosticeng/clickhouse-evm/cmd/function/evm_descriptor_from_fullsig"
	"github.com/agnosticeng/clickhouse-evm/cmd/function/evm_signature_from_descriptor"
	"github.com/agnosticeng/clickhouse-evm/cmd/function/keccak256"
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
			ethereum_decode_tx.Command(),
			evm_decode_calldata.Command(),
			convert_format.Command(),
		},
	}
}
