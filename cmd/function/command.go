package function

import (
	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function/evm_decode_event"
	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function/keccak256"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name: "function",
		Subcommands: []*cli.Command{
			evm_decode_event.Command(),
			keccak256.Command(),
		},
	}
}
