package ethereum_rpc_filter

import (
	"encoding/json"
	"errors"
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
		&cli.StringFlag{Name: "poll-method", Value: "eth_getFilterChanges"},
		&cli.DurationFlag{Name: "poll-interval", Value: time.Second},
	}
}

func Command() *cli.Command {
	return &cli.Command{
		Name:  "ethereum-rpc-filter",
		Flags: lo.Flatten([][]cli.Flag{jsonrpc_cli.Flags(), Flags()}),
		Action: func(ctx *cli.Context) error {
			var (
				method       = ctx.Args().Get(0)
				pollMethod   = ctx.String("poll-method")
				pollInterval = ctx.Duration("poll-interval")
				endpoint     = ctx.String("endpoint")

				inputFilterCol  = new(proto.ColBytes)
				outputResultCol = new(proto.ColBytes)

				input = proto.Results{
					{Name: "filter", Data: inputFilterCol},
				}

				output = proto.Input{
					{Name: "result", Data: outputResultCol},
				}

				buf proto.Buffer
			)

			for {
				var (
					inputBlock proto.Block
					err        = inputBlock.DecodeRawBlock(
						proto.NewReader(os.Stdin),
						54451,
						input,
					)
				)

				if errors.Is(err, io.EOF) {
					return nil
				}

				if err != nil {
					return err
				}

				var filters []json.RawMessage

				for i := 0; i < input.Rows(); i++ {
					filters = append(filters, inputFilterCol.Row(i))
				}

				params, err := json.Marshal(filters)

				if err != nil {
					return err
				}

				client, err := jsonrpc.NewHTTPClient(ctx.Context)

				if err != nil {
					return err
				}

				defer client.Close()

				if pollInterval == 0 {
					pollInterval = time.Second
				}

				var msg = jsonrpc.NewRequest(method, params)

				if err := client.Call(ctx.Context, endpoint, msg); err != nil {
					return err
				}

				if msg.Error != nil {
					return err
				}

				var filterId string

				if err := json.Unmarshal(msg.Result, &filterId); err != nil {
					return err
				}

				defer client.Call(
					ctx.Context,
					endpoint,
					jsonrpc.NewRequest(
						"eth_uninstallFilter",
						lo.Must(json.Marshal([]interface{}{filterId})),
					),
				)

				for {
					msg := jsonrpc.NewRequest(
						pollMethod,
						lo.Must(json.Marshal([]interface{}{filterId})),
					)

					if err := client.Call(ctx.Context, endpoint, msg); err != nil {
						return err
					}

					if msg.Error != nil {
						return msg.Error
					}

					var rows []json.RawMessage

					if err := json.Unmarshal(msg.Result, &rows); err != nil {
						return err
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

					time.Sleep(pollInterval)
				}

				proto.Reset(
					&buf,
					inputFilterCol,
				)
			}
		},
	}
}
