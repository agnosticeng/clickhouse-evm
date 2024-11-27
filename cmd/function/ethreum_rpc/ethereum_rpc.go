package ethereum_rpc

import (
	"encoding/json"
	stdjson "encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/jsonrpc"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/jsonrpc_cli"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/types"
	"github.com/agnosticeng/panicsafe"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "ethereum-rpc",
		Flags: jsonrpc_cli.Flags(),
		Action: func(ctx *cli.Context) error {
			return panicsafe.Recover(func() error {

				var (
					defaultEndpoint = ctx.String("endpoint")
					callOpts        = jsonrpc_cli.CallOptionsFromContext(ctx)
					buf             proto.Buffer

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
						batch    = make(jsonrpc.Batch, input.Rows())
						endpoint string
					)

					for i := 0; i < input.Rows(); i++ {
						var (
							method  = inputMethodCol.Row(i)
							params  = inputParamsCol.Row(i)
							jparams = lo.Map(params, prepareParam)
						)

						if edp := inputEndpointCol.Row(i); edp != endpoint {
							if len(endpoint) == 0 {
								endpoint = edp
							} else {
								return fmt.Errorf("endpoint must be constant for the whole input block")
							}
						}

						js, err := json.Marshal(jparams)

						if err != nil {
							return err
						}

						batch[i].SetRequest(method, js)
					}

					if !strings.HasPrefix(endpoint, "http") && !strings.HasPrefix(endpoint, "https") {
						endpoint = defaultEndpoint + "#" + endpoint
					}

					if err := client.BatchCall(ctx.Context, endpoint, batch, callOpts...); err != nil {
						return err
					}

					for i := 0; i < input.Rows(); i++ {
						if resp := batch[i]; resp.Error != nil {
							outputResultCol.Append(lo.Must(stdjson.Marshal(types.Result{Error: resp.Error.Error()})))
						} else {
							outputResultCol.Append(lo.Must(stdjson.Marshal(types.Result{Value: resp.Result})))
						}
					}

					var outputblock = proto.Block{
						Columns: 1,
						Rows:    input.Rows(),
					}

					if err := outputblock.EncodeRawBlock(&buf, 54451, output); err != nil {
						return err
					}

					if _, err := os.Stdout.Write(buf.Buf); err != nil {
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
			})
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
