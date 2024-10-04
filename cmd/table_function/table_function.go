package table_function

import (
	ethereum_rpc_pending_transactions "github.com/agnosticeng/agnostic-clickhouse-udf/cmd/table_function/ethreum_rpc_pending_transactions"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name: "table-function",
		Subcommands: []*cli.Command{
			ethereum_rpc_pending_transactions.Command(),
		},
	}
}
