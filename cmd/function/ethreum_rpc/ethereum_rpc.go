package ethereum_rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/jsonrpc"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name: "ethereum-rpc",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "endpoint"},
			&cli.IntFlag{Name: "batch-max-size", Value: 200},
			&cli.IntFlag{Name: "batch-concurrency-limit", Value: 5},
		},
		Action: func(ctx *cli.Context) error {
			var (
				defaultEndpoint = ctx.String("endpoint")
				// logger          = slogctx.FromCtx(ctx.Context)
				batchOpts []jsonrpc.BatchOptionsFunc
				buf       proto.Buffer

				inputEndpointCol = new(proto.ColStr)
				inputMethodCol   = new(proto.ColStr)
				inputParamsCol   = proto.NewArray(new(proto.ColBytes))
				outputResultCol  = new(proto.ColBytes)

				input = proto.Results{
					{Name: "method", Data: inputMethodCol},
					{Name: "params", Data: inputParamsCol},
					{Name: "endpoint", Data: inputEndpointCol},
				}

				output = proto.Input{
					{Name: "result", Data: outputResultCol},
				}
			)

			batchOpts = append(batchOpts, jsonrpc.WithChunkSize(ctx.Int("batch-max-size")))
			batchOpts = append(batchOpts, jsonrpc.WithConcurrencyLimit(ctx.Int("atch-concurrency-limit")))

			client, err := jsonrpc.NewHTTPClient(ctx.Context)

			if err != nil {
				return err
			}

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

				var (
					requests = make([]*jsonrpc.Message, input.Rows())
					endpoint string
				)

				for i := 0; i < input.Rows(); i++ {
					var (
						method  = inputMethodCol.Row(i)
						params  = inputParamsCol.Row(i)
						req     = jsonrpc.NewMessage()
						jparams = lo.Map(params, prepareParam)
					)

					if edp := inputEndpointCol.Row(i); edp != endpoint {
						if len(endpoint) == 0 {
							endpoint = edp
						}

						return fmt.Errorf("endpoint must be constant for the whole input block")
					}

					js, err := json.Marshal(jparams)

					if err != nil {
						return err
					}

					req.Method = method
					req.Params = js
					requests[i] = req
				}

				if len(endpoint) == 0 {
					endpoint = defaultEndpoint
				}

				responses, err := client.BatchCall(ctx.Context, endpoint, requests, batchOpts...)

				if err != nil {
					return err
				}

				for _, response := range responses {
					if response.Error != nil {
						return response.Error
					}
				}

				for i := 0; i < input.Rows(); i++ {
					outputResultCol.Append(responses[i].Result)
				}

				var outputblock = proto.Block{
					Columns: 1,
					Rows:    input.Rows(),
				}

				if err := outputblock.EncodeRawBlock(&buf, 54451, output); err != nil {
					return err
				}

				if _, err := io.Copy(os.Stdout, buf.Reader()); err != nil {
					return err
				}

				proto.Reset(
					&buf,
					inputEndpointCol,
					inputMethodCol,
					inputParamsCol,
					outputResultCol,
				)

			}
		},
	}
}

func prepareParam(param []byte, _ int) json.RawMessage {
	if len(param) >= 2 && param[0] == '0' && param[1] == 'x' {
		var v = make([]byte, 0, len(param)+2)
		v = append(v, '"')
		v = append(v, param...)
		v = append(v, '"')
		return json.RawMessage(v)
	}

	return json.RawMessage(param)
}
