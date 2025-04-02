package table_function

import (
	"github.com/agnosticeng/clickhouse-evm/cmd/table_function/ethereum_rpc_filter"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name: "table-function",
		Subcommands: []*cli.Command{
			ethereum_rpc_filter.Command(),
		},
	}
}
