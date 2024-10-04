package ethereum_rpc_pending_transactions

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/jsonrpc"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/jsonrpc_cli"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.DurationFlag{Name: "poll-interval", Value: time.Second},
	}
}

func Command() *cli.Command {
	return &cli.Command{
		Name:  "ethereum-rpc-pending-transactions",
		Flags: lo.Flatten([][]cli.Flag{jsonrpc_cli.Flags(), Flags()}),
		Action: func(ctx *cli.Context) error {
			var (
				endpoint        = ctx.String("endpoint")
				pollInterval    = ctx.Duration("poll-interval")
				outputResultCol = new(proto.ColBytes)
				output          = proto.Input{
					{Name: "result", Data: outputResultCol},
				}

				buf proto.Buffer
			)

			client, err := jsonrpc.NewHTTPClient(ctx.Context)

			if err != nil {
				return err
			}

			defer client.Close()

			if pollInterval == 0 {
				pollInterval = time.Second
			}

			resp, err := client.Call(
				ctx.Context,
				endpoint,
				jsonrpc.NewRequest("eth_newPendingTransactionFilter", nil),
			)

			if err != nil {
				return err
			}

			var filterId string

			if err := json.Unmarshal(resp.Result, &filterId); err != nil {
				return err
			}

			defer client.Call(
				ctx.Context,
				endpoint,
				jsonrpc.NewRequest(
					"eth_uninstallFilter",
					jsonrpc.MustMarshalJSON([]interface{}{filterId}),
				),
			)

			for {
				resp, err := client.Call(
					ctx.Context,
					endpoint,
					jsonrpc.NewRequest(
						"eth_getFilterChanges",
						jsonrpc.MustMarshalJSON([]interface{}{filterId}),
					),
				)

				if err != nil {
					return err
				}

				if resp.Error != nil {
					return resp.Error
				}

				var rows []json.RawMessage

				if err := json.Unmarshal(resp.Result, &rows); err != nil {
					return err
				}

				if len(rows) == 0 {
					continue
				}

				for _, row := range rows {
					outputResultCol.Append(row)
				}

				var outputblock = proto.Block{
					Columns: 1,
					Rows:    len(rows),
				}

				if err := outputblock.EncodeRawBlock(&buf, 54451, output); err != nil {
					return err
				}

				if _, err := io.Copy(os.Stdout, buf.Reader()); err != nil {
					return err
				}

				proto.Reset(
					&buf,
					outputResultCol,
				)

			}
		},
	}
}
