package evm_decode_event

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/abi_provider"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/abi_provider/impl"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/types"
	"github.com/agnosticeng/concu/memo"
	"github.com/agnosticeng/evmabi/encoding/json"
	"github.com/agnosticeng/panicsafe"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name: "evm-decode-event",
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
				var abiProviderCache = memo.KeyedErrTheine[string, abi_provider.ABIProvider](
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
					buf proto.Buffer

					inputTopicsCol  = proto.NewArray(new(proto.ColBytes))
					inputDataCol    = new(proto.ColBytes)
					inputsAbiCol    = proto.NewArray(new(proto.ColStr))
					outputResultCol = new(proto.ColBytes)

					input = proto.Results{
						{Name: "topics", Data: inputTopicsCol},
						{Name: "data", Data: inputDataCol},
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
							topics            = inputTopicsCol.Row(i)
							data              = inputDataCol.Row(i)
							abiProviders      = inputsAbiCol.Row(i)
							decoded           bool
							lastDecodingError error
						)

						if len(topics) == 0 {
							outputResultCol.Append((&types.Result{Error: "cannot decode event with empty topics[0]"}).ToJSON())
							continue
						}

					decodeLoop:
						for _, abiProvider := range abiProviders {
							// slogctx.Info(ctx.Context, "ABIPROVIDER", "str", abiProvider)

							p, err := abiProviderCache(abiProvider)

							if err != nil {
								return fmt.Errorf("failed to parse ABI provider '%s': %w", abiProvider, err)
							}

							evts, err := p.Events(string(topics[0]))

							if err != nil {
								return err
							}

							for _, evt := range evts {
								// slogctx.Info(
								// 	ctx.Context,
								// 	"DECODE",
								// 	"sig", evt.Sig,
								// 	"topics", lo.Map(topics, func(topic []byte, _ int) string { return hexutil.Encode(topic) }),
								// 	"data", hexutil.Encode(data),
								// )

								n, err := json.DecodeLog(topics, data, *evt)

								if err != nil {
									lastDecodingError = err
									continue
								}

								if !n.Exists() {
									continue
								}

								outputResultCol.Append((&types.Result{Value: &n}).ToJSON())
								decoded = true
								break decodeLoop
							}
						}

						if !decoded {
							if lastDecodingError != nil {
								outputResultCol.Append((&types.Result{Error: lastDecodingError.Error()}).ToJSON())
							} else {
								outputResultCol.Append((&types.Result{Error: "cannot decode event"}).ToJSON())
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
						inputTopicsCol,
						inputDataCol,
						inputsAbiCol,
						outputResultCol,
					)

				}
			})
		},
	}
}
