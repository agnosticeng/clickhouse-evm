package evm_decode_call

import (
	"errors"
	"io"
	"os"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/abi_provider"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/abi_provider/impl"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/memo"
	"github.com/agnosticeng/evmabi/json"
	"github.com/urfave/cli/v2"
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
			var cache = memo.KeyedErr[string, abi_provider.ABIProvider](
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
				inputsAbiCol    = new(proto.ColStr)
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
						input  = inputDataCol.Row(i)
						output = outputDataCol.Row(1)
						key    = inputsAbiCol.Row(i)
					)

					p, err := cache(key)

					if err != nil {
						return err
					}

					meth, err := p.Method(string(input[:4]))

					if err != nil {
						return err
					}

					if meth == nil {
						outputResultCol.Append([]byte("{}"))
						continue
					}

					n, err := json.DecodeTrace(input, output, *meth)

					if err != nil {
						return err
					}

					if !n.Exists() {
						outputResultCol.Append([]byte("{}"))
						continue
					}

					js, err := n.MarshalJSON()

					if err != nil {
						return err
					}

					outputResultCol.Append(js)
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
		},
	}
}
