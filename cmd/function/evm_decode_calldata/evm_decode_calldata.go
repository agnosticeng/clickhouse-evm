package evm_decode_calldata

import (
	"errors"
	"io"
	"os"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/abi_provider"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/abi_provider/impl"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/types"
	"github.com/agnosticeng/concu/memo"
	"github.com/agnosticeng/evmabi/encoding/json"
	"github.com/agnosticeng/evmabi/fullsig"
	"github.com/agnosticeng/panicsafe"
	"github.com/bytedance/sonic/ast"
	"github.com/urfave/cli/v2"
	slogctx "github.com/veqryn/slog-context"
)

func Command() *cli.Command {
	return &cli.Command{
		Name: "evm-decode-calldata",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "abi-provider",
				EnvVars: []string{"EVM_ABI_PROVIDER"},
			},
			&cli.Uint64Flag{
				Name:    "abi-provider-cache-size",
				Value:   2000,
				EnvVars: []string{"EVM_ABI_PROVIDER_CACHE_SIZE"},
			},
		},
		Action: func(ctx *cli.Context) error {
			return panicsafe.Recover(func() error {
				var abiProviderCache = memo.KeyedErrTheine(
					func(key string) (abi_provider.ABIProvider, error) {
						if len(key) == 0 {
							return impl.NewABIProvider(ctx.String("abi-provider"))
						} else {
							return impl.NewABIProvider(key)
						}
					},
					int64(ctx.Uint64("abi-provider-cache-size")),
				)

				var (
					logger = slogctx.FromCtx(ctx.Context)
					buf    proto.Buffer

					inputInputDataCol  = new(proto.ColBytes)
					inputOutputDataCol = new(proto.ColBytes)
					inputsAbiCol       = proto.NewArray(new(proto.ColStr))
					outputResultCol    = new(proto.ColBytes)

					input = proto.Results{
						{Name: "input", Data: inputInputDataCol},
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
							input             = inputInputDataCol.Row(i)
							abiProviders      = inputsAbiCol.Row(i)
							decoded           bool
							lastDecodingError error
						)

						if len(input) < 4 {
							outputResultCol.Append((&types.Result{Error: "cannot decode calldata with length < 4"}).ToJSON())
							continue
						}

					decodeLoop:
						for _, abiProvider := range abiProviders {
							p, err := abiProviderCache(abiProvider)

							if err != nil {
								logger.Error(err.Error(), "abi_provider", abiProvider)
								return err
							}

							meths, err := p.Methods(string(input[:4]))

							if err != nil {
								return err
							}

							for _, meth := range meths {
								n, err := json.DecodeCallData(input, *meth)

								if err != nil {
									lastDecodingError = err
									continue
								}

								if !n.Exists() {
									continue
								}

								n.Set("fullsig", ast.NewString(fullsig.StringifyMethod(meth)))
								outputResultCol.Append((&types.Result{Value: &n}).ToJSON())
								decoded = true
								break decodeLoop
							}
						}

						if !decoded {
							if lastDecodingError != nil {
								outputResultCol.Append((&types.Result{Error: lastDecodingError.Error()}).ToJSON())
							} else {
								outputResultCol.Append((&types.Result{Error: "cannot decode call"}).ToJSON())
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
						inputInputDataCol,
						inputOutputDataCol,
						inputsAbiCol,
						outputResultCol,
					)

				}
			})
		},
	}
}
