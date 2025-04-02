package ethereum_rpc_call

import (
	stdjson "encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"reflect"
	"strings"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/agnosticeng/clickhouse-evm/internal/ethereum_rpc"
	"github.com/agnosticeng/clickhouse-evm/internal/jsonrpc"
	"github.com/agnosticeng/clickhouse-evm/internal/jsonrpc_cli"
	"github.com/agnosticeng/clickhouse-evm/internal/types"
	"github.com/agnosticeng/concu/memo"
	"github.com/agnosticeng/evmabi/encoding/json"
	"github.com/agnosticeng/evmabi/fullsig"
	"github.com/agnosticeng/panicsafe"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "ethereum-rpc-call",
		Flags: jsonrpc_cli.Flags(),
		Action: func(ctx *cli.Context) error {
			return panicsafe.Recover(func() error {
				var cache = memo.KeyedErr[string, *abi.Method](
					func(key string) (*abi.Method, error) {
						meth, err := fullsig.ParseMethod(key)

						if err != nil {
							return nil, err
						}

						return &meth, nil
					},
				)

				var (
					defaultEndpoint = ctx.String("endpoint")
					callOpts        = jsonrpc_cli.CallOptionsFromContext(ctx)
					buf             proto.Buffer

					inputToCol          = new(proto.ColStr)
					inputFullSigCol     = new(proto.ColStr)
					inputDataCol        = new(proto.ColBytes)
					inputBlockNumberCol = new(proto.ColInt64)
					inputEndpointCol    = new(proto.ColStr)
					outputResultCol     = new(proto.ColBytes)

					input = proto.Results{
						{Name: "to", Data: inputToCol},
						{Name: "fullsig", Data: inputFullSigCol},
						{Name: "data", Data: inputDataCol},
						{Name: "block_number", Data: inputBlockNumberCol},
						{Name: "endpoint", Data: inputEndpointCol},
					}

					output = proto.Input{
						{Name: "result", Data: outputResultCol},
					}
				)

				callOpts = append(callOpts, jsonrpc.WithRetryableErrorPredicate(ethereum_rpc.RetryableErrorPredicate))
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
							to          = inputToCol.Row(i)
							_fullsig    = inputFullSigCol.Row(i)
							data        = inputDataCol.Row(i)
							blockNumber = inputBlockNumberCol.Row(i)
							inputs      []interface{}
						)

						if edp := inputEndpointCol.Row(i); edp != endpoint {
							if len(endpoint) == 0 {
								endpoint = edp
							} else {
								return fmt.Errorf("endpoint must be constant for the whole input block")
							}
						}

						if len(data) > 0 {
							if err := stdjson.Unmarshal(data, &inputs); err != nil {
								return err
							}
						}

						meth, err := cache(_fullsig)

						if err != nil {
							return err
						}

						inputs, err = prepareParams(meth, inputs)

						if err != nil {
							return fmt.Errorf("failed to prepare params: %v", err)
						}

						var inputData []byte
						inputData = append(inputData, meth.ID...)

						if d, err := meth.Inputs.Pack(inputs...); err != nil {
							return fmt.Errorf("failed to pack inputs values: %v", err)
						} else {
							inputData = append(inputData, d...)
						}

						params, err := stdjson.Marshal([]interface{}{
							TransactionObject{
								To:   to,
								Data: string(hexutil.Encode(inputData)),
							},
							BlockNumberToString(blockNumber),
							map[string]interface{}{},
						})

						if err != nil {
							return err
						}

						batch[i].SetRequest("eth_call", params)
					}

					if !strings.HasPrefix(endpoint, "http") && !strings.HasPrefix(endpoint, "https") {
						endpoint = defaultEndpoint + "#" + endpoint
					}

					if err := client.BatchCall(ctx.Context, endpoint, batch, callOpts...); err != nil {
						return fmt.Errorf("batch call failed: %w", err)
					}

					for i := 0; i < input.Rows(); i++ {
						var res = decodeResult(cache, inputFullSigCol.Row(i), &batch[i])
						js, err := stdjson.Marshal(res)

						if err != nil {
							return fmt.Errorf("failed to decode result: %w", err)
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
						inputToCol,
						inputFullSigCol,
						inputDataCol,
						inputBlockNumberCol,
						inputEndpointCol,
						outputResultCol,
					)
				}
			})
		},
	}
}

func decodeResult(
	cache func(key string) (*abi.Method, error),
	fullsig string,
	resp *jsonrpc.Message,
) types.Result {
	if resp.Error != nil {
		return types.Result{Error: resp.Error.Error()}
	}

	var strData string

	if err := stdjson.Unmarshal(resp.Result, &strData); err != nil {
		return types.Result{Error: err.Error()}
	}

	meth, err := cache(fullsig)

	if err != nil {
		return types.Result{Error: err.Error()}
	}

	data, err := hexutil.Decode(strData)

	if err != nil {
		return types.Result{Error: err.Error()}
	}

	if len(data) == 0 {
		return types.Result{}
	}

	n, err := json.DecodeArguments(data, meth.Outputs)

	if err != nil {
		return types.Result{Error: err.Error()}
	}

	return types.Result{Value: &n}
}

func prepareParams(meth *abi.Method, values []interface{}) ([]interface{}, error) {
	var (
		res = make([]interface{}, len(values))
		err error
	)

	if len(values) != len(meth.Inputs) {
		return nil, fmt.Errorf("invalid values count: %d", len(values))
	}

	for i := 0; i < len(values); i++ {
		switch meth.Inputs[i].Type.T {
		case abi.AddressTy:
			res[i], err = prepareAddress(values[i])
		case abi.IntTy:
			res[i], err = prepareBigInt(values[i])
		case abi.UintTy:
			res[i], err = prepareBigInt(values[i])
		case abi.FixedBytesTy:
			res[i], err = prepareFixedBytes(meth.Inputs[i].Type, values[i])
		case abi.BytesTy:
			res[i], err = prepareBytes(values[i])
		default:
			res[i] = values[i]
		}

		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func prepareFixedBytes(t abi.Type, v interface{}) (interface{}, error) {
	s, ok := v.(string)

	if !ok {
		return nil, fmt.Errorf("invalid string: %v", v)
	}

	var (
		data = common.FromHex(s)
		arr  = reflect.New(t.GetType()).Elem()
	)

	if t.Size != len(data) {
		return nil, fmt.Errorf("invalid %s value: %s (%d)", t.String(), s, len(data))
	}

	for i := 0; i < t.Size; i++ {
		arr.Index(i).Set(reflect.ValueOf(data[i]))
	}

	return arr.Interface(), nil
}

func prepareBytes(v interface{}) ([]byte, error) {
	s, ok := v.(string)

	if !ok {
		return nil, fmt.Errorf("invalid string: %v", v)
	}

	if len(s) >= 2 && s[0:2] == "0x" {
		return hexutil.MustDecode(s), nil
	}

	return []byte(s), nil
}

func prepareAddress(v interface{}) (common.Address, error) {
	s, ok := v.(string)

	if !ok {
		return common.Address{}, fmt.Errorf("invalid address: %v", v)
	}

	return common.HexToAddress(s), nil
}

func prepareBigInt(v interface{}) (*big.Int, error) {
	switch v := v.(type) {
	case string:
		return hexutil.DecodeBig(v)
	case uint64:
		return big.NewInt(0).SetUint64(v), nil
	case float64:
		return big.NewInt(0).SetInt64(int64(v)), nil
	default:
		return nil, fmt.Errorf("invalid big int: %v(%s)", v, reflect.TypeOf(v))
	}
}
