package keccak256

import (
	"errors"
	"io"
	"os"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "keccak256",
		Flags: []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			var (
				buf proto.Buffer

				inputStrCol     = new(proto.ColBytes)
				outputResultCol = new(proto.ColBytes)

				input = proto.Results{
					{Name: "str", Data: inputStrCol},
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
						str = inputStrCol.Row(i)
						res = crypto.Keccak256Hash(str).Bytes()
					)

					outputResultCol.Append(res)
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
					inputStrCol,
					outputResultCol,
				)

			}
		},
	}
}
