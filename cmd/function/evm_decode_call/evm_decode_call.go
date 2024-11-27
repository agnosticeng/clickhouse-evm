package evm_decode_call

import (
	stdjson "encoding/json"
	"errors"
	"io"
	"os"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/abi_provider"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/abi_provider/impl"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/memo"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/types"
	"github.com/agnosticeng/evmabi/json"
	"github.com/agnosticeng/panicsafe"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	slogctx "github.com/veqryn/slog-context"
)

func Command() *cli.Command {
	return &cli.Command{
		Name: "evm-decode-call",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "abi-provider",
				EnvVars: []string{"EVM_ABI_PROVIDER"},
			},
		},
		Action: func(ctx *cli.Context) error {
			return panicsafe.Recover(func() error {
				var abiProviderCache = memo.KeyedErr[string, abi_provider.ABIProvider](
					func(key string) (abi_provider.ABIProvider, error) {
						if len(key) == 0 {
							return impl.NewABIProvider(ctx.String("abi-provider"))
						} else {
							return impl.NewABIProvider(key)
						}
					},
				)

				var (
					buf proto.Buffer

					inputDataCol    = new(proto.ColBytes)
					outputDataCol   = new(proto.ColBytes)
					inputsAbiCol    = proto.NewArray(new(proto.ColStr))
					outputResultCol = new(proto.ColBytes)

					input = proto.Results{
						{Name: "input", Data: inputDataCol},
						{Name: "output", Data: outputDataCol},
						{Name: "abi", Data: inputsAbiCol},
					}

					output = proto.Input{
						{Name: "result", Data: outputResultCol},
					}
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

					for i := 0; i < input.Rows(); i++ {
						var (
							input             = inputDataCol.Row(i)
							output            = outputDataCol.Row(i)
							abiProviders      = inputsAbiCol.Row(i)
							decoded           bool
							lastDecodingError error
						)

					decodeLoop:
						for _, abiProvider := range abiProviders {
							p, err := abiProviderCache(abiProvider)

							if err != nil {
								slogctx.FromCtx(ctx.Context).Info(err.Error())
								return err
							}

							meths, err := p.Methods(string(input[:4]))

							if err != nil {
								return err
							}

							for _, meth := range meths {
								n, err := json.DecodeTrace(input, output, *meth)

								if err != nil {
									lastDecodingError = err
									continue
								}

								if !n.Exists() {
									continue
								}

								outputResultCol.Append(lo.Must(stdjson.Marshal(types.Result{Value: lo.Must(n.MarshalJSON())})))
								decoded = true
								break decodeLoop
							}
						}

						if !decoded {
							if lastDecodingError != nil {
								outputResultCol.Append(lo.Must(stdjson.Marshal(types.Result{Error: lastDecodingError.Error()})))
							} else {
								outputResultCol.Append(lo.Must(stdjson.Marshal(types.Result{})))
							}
						}
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
						inputDataCol,
						outputDataCol,
						inputsAbiCol,
						outputResultCol,
					)

				}
			})
		},
	}
}
