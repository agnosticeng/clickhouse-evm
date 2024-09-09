package evm_decode_event

import (
	"errors"
	"io"
	"os"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/agnosticeng/evmabi/json"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name: "evm-decode-event",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "abi-file", Required: true},
		},
		Action: func(ctx *cli.Context) error {
			var (
				buf proto.Buffer

				inputTopicsCol  = proto.NewArray(new(proto.ColBytes))
				inputDataCol    = new(proto.ColBytes)
				outputResultCol = new(proto.ColBytes)

				input = proto.Results{
					{Name: "topics", Data: inputTopicsCol},
					{Name: "data", Data: inputDataCol},
				}

				output = proto.Input{
					{Name: "result", Data: outputResultCol},
				}
			)

			f, err := os.Open(ctx.String("abi-file"))

			if err != nil {
				return err
			}

			defer f.Close()

			_abi, err := ParseIndexedABI(f)

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

				for i := 0; i < input.Rows(); i++ {
					var js, err = decodeEvent(
						inputTopicsCol.Row(i),
						inputDataCol.Row(i),
						_abi,
					)

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
					inputTopicsCol,
					inputDataCol,
					outputResultCol,
				)

			}
		},
	}
}

func decodeEvent(topics [][]byte, data []byte, _abi IndexedABI) ([]byte, error) {
	var eventDesc = _abi.EventsSigHashIndex[string(topics[0])]

	if eventDesc == nil {
		return []byte("{}"), nil
	}

	evt, err := json.DecodeLog(data, topics, *eventDesc)

	if err != nil {
		return nil, err
	}

	js, err := evt.MarshalJSON()

	if err != nil {
		return nil, err
	}

	return js, nil
}
