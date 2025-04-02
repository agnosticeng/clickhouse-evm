package evm_descriptor_from_fullsig

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/agnosticeng/clickhouse-evm/internal/types"
	"github.com/agnosticeng/evmabi/abi"
	"github.com/agnosticeng/evmabi/fullsig"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "evm-descriptor-from-fullsig",
		Flags: []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			var (
				buf proto.Buffer

				inputFullSigCol = new(proto.ColStr)
				outputResultCol = new(proto.ColBytes)

				input = proto.Results{
					{Name: "fullsig", Data: inputFullSigCol},
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
					var s = inputFullSigCol.Row(i)

					switch {
					case strings.HasPrefix(s, "event"):
						evt, err := fullsig.ParseEvent(s)

						if err != nil {
							outputResultCol.Append((&types.Result{Error: err.Error()}).ToJSON())
							continue
						}

						m, err := abi.EventToFieldMarshaling(&evt)

						if err != nil {
							outputResultCol.Append((&types.Result{Error: err.Error()}).ToJSON())
							continue
						}

						outputResultCol.Append((&types.Result{Value: m}).ToJSON())

					case strings.HasPrefix(s, "function"):
						meth, err := fullsig.ParseMethod(s)

						if err != nil {
							outputResultCol.Append((&types.Result{Error: err.Error()}).ToJSON())
							continue
						}

						m, err := abi.MethodToFieldMarshaling(&meth)

						if err != nil {
							outputResultCol.Append((&types.Result{Error: err.Error()}).ToJSON())
							continue
						}

						outputResultCol.Append((&types.Result{Value: m}).ToJSON())

					default:
						outputResultCol.Append((&types.Result{Error: fmt.Sprintf("unknown type: %s", s)}).ToJSON())
						continue
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
					inputFullSigCol,
					outputResultCol,
				)

			}
		},
	}
}
