package evm_signature_from_descriptor

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/types"
	"github.com/agnosticeng/evmabi/abi"
	"github.com/agnosticeng/evmabi/fullsig"
	"github.com/bytedance/sonic"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/urfave/cli/v2"
)

type Signature struct {
	Selector  string `json:"selector"`
	Signature string `json:"signature"`
	FullSig   string `json:"fullsig"`
}

func Command() *cli.Command {
	return &cli.Command{
		Name:  "evm-signature-from-descriptor",
		Flags: []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			var (
				buf proto.Buffer

				inptDescriptorCol = new(proto.ColBytes)
				outputResultCol   = new(proto.ColBytes)

				input = proto.Results{
					{Name: "event_descriptor", Data: inptDescriptorCol},
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
					var content = inptDescriptorCol.Row(i)

					node, err := sonic.Get(content, "type")

					if err != nil {
						outputResultCol.Append((&types.Result{Error: err.Error()}).ToJSON())
						continue
					}

					_type, err := node.String()

					if err != nil {
						outputResultCol.Append((&types.Result{Error: err.Error()}).ToJSON())
						continue
					}

					switch _type {
					case "event":
						evt, err := abi.JSONEvent(content)

						if err != nil {
							outputResultCol.Append((&types.Result{Error: err.Error()}).ToJSON())
							continue
						}

						outputResultCol.Append((&types.Result{
							Value: &Signature{
								Selector:  evt.ID.Hex(),
								Signature: evt.Sig,
								FullSig:   fullsig.StringifyEvent(evt),
							},
						}).ToJSON())

					case "function":
						meth, err := abi.JSONMethod(content)

						if err != nil {
							outputResultCol.Append((&types.Result{Error: err.Error()}).ToJSON())
							continue
						}

						outputResultCol.Append((&types.Result{
							Value: &Signature{
								Selector:  hexutil.Encode(meth.ID),
								Signature: meth.Sig,
								FullSig:   fullsig.StringifyMethod(meth),
							},
						}).ToJSON())

					default:
						outputResultCol.Append((&types.Result{Error: fmt.Sprintf("unknown type: %s", _type)}).ToJSON())
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
					inptDescriptorCol,
					outputResultCol,
				)

			}
		},
	}
}
