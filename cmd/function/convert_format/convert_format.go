package convert_format

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "convert-format",
		Flags: []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			var (
				buf proto.Buffer

				inputFromFormatCol = new(proto.ColStr)
				inputToFormatCol   = new(proto.ColStr)
				inputStrCol        = new(proto.ColBytes)
				outputResultCol    = new(proto.ColBytes)

				input = proto.Results{
					{Name: "from_format", Data: inputFromFormatCol},
					{Name: "to_format", Data: inputToFormatCol},
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
						from   = inputFromFormatCol.Row(i)
						to     = inputToFormatCol.Row(i)
						input  = inputStrCol.Row(i)
						m      map[string]any
						output []byte
						err    error
					)

					if len(input) == 0 {
						outputResultCol.Append([]byte(""))
						continue
					}

					switch from {
					case "json", "JSON":
						err = json.Unmarshal(input, &m)
					case "toml", "TOML":
						err = toml.Unmarshal(input, &m)
					case "yaml", "YAML":
						err = yaml.Unmarshal(input, &m)
					default:
						return fmt.Errorf("invalid input format: %s", from)
					}

					if err != nil {
						return err
					}

					switch to {
					case "json", "JSON":
						output, err = json.Marshal(m)
					case "toml", "TOML":
						output, err = toml.Marshal(m)
					case "yaml", "YAML":
						output, err = yaml.Marshal(m)
					default:
						return fmt.Errorf("invalid output format: %s", from)
					}

					if err != nil {
						return err
					}

					outputResultCol.Append(output)
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
					inputFromFormatCol,
					inputToFormatCol,
					inputStrCol,
					outputResultCol,
				)

			}
		},
	}
}
