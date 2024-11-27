package jsonrpc_cli

import (
	"github.com/agnosticeng/agnostic-clickhouse-udf/internal/jsonrpc"
	"github.com/urfave/cli/v2"
)

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "endpoint",
			EnvVars: []string{"ETHEREUM_RPC_ENDPOINT"},
		},
		&cli.IntFlag{
			Name:    "max-batch-size",
			Value:   200,
			EnvVars: []string{"ETHEREUM_RPC_MAX_BATCH_SIZE"},
		},
		&cli.IntFlag{
			Name:    "max-concurrent-requests",
			Value:   5,
			EnvVars: []string{"ETHEREUM_RPC_MAX_CONCURRENT_REQUESTS"},
		},
		&cli.BoolFlag{
			Name:    "disable-batch",
			Value:   false,
			EnvVars: []string{"ETHEREUM_RPC_DISABLE_BATCH"},
		},
		&cli.BoolFlag{
			Name:    "fail-on-error",
			Value:   false,
			EnvVars: []string{"ETHEREUM_RPC_FAIL_ON_ERROR"},
		},
		&cli.BoolFlag{
			Name:    "fail-on-null",
			Value:   false,
			EnvVars: []string{"ETHEREUM_RPC_FAIL_ON_NULL"},
		},
	}
}

func CallOptionsFromContext(ctx *cli.Context) []jsonrpc.CallOptionsFunc {
	var callOpts []jsonrpc.CallOptionsFunc
	callOpts = append(callOpts, jsonrpc.WithMatchBatchSize(ctx.Int("max-batch-size")))
	callOpts = append(callOpts, jsonrpc.WithMaxConcurrentRequests(ctx.Int("max-concurrent-requests")))
	callOpts = append(callOpts, jsonrpc.WithDisableBatch(ctx.Bool("disable-batch")))
	callOpts = append(callOpts, jsonrpc.WithFailOnError(ctx.Bool("fail-on-error")))
	callOpts = append(callOpts, jsonrpc.WithFailOnNull(ctx.Bool("fail-on-null")))
	return callOpts
}
