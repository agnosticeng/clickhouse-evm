package jsonrpc_cli

import (
	"time"

	"github.com/agnosticeng/clickhouse-evm/internal/jsonrpc"
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
			Name:    "fail-on-retryable-error",
			Value:   false,
			EnvVars: []string{"ETHEREUM_RPC_FAIL_ON_RETRYABLE_ERROR"},
		},
		&cli.BoolFlag{
			Name:    "fail-on-null",
			Value:   false,
			EnvVars: []string{"ETHEREUM_RPC_FAIL_ON_NULL"},
		},
		&cli.IntSliceFlag{
			Name:    "retryable-status-codes",
			Value:   cli.NewIntSlice(429, 502, 503, 504),
			EnvVars: []string{"ETHEREUM_RPC_RETRYABLE_STATUS_CODES"},
		},
		&cli.DurationFlag{
			Name:    "retry-initial-interval",
			Value:   500 * time.Millisecond,
			EnvVars: []string{"ETHEREUM_RPC_RETRY_INITIAL_INTERVAL"},
		},
		&cli.Float64Flag{
			Name:    "retry-randomization-factor",
			Value:   0.5,
			EnvVars: []string{"ETHEREUM_RPC_RETRY_RANDOMIZATION_FACTOR"},
		},
		&cli.Float64Flag{
			Name:    "retry-multiplier",
			Value:   1.5,
			EnvVars: []string{"ETHEREUM_RPC_RETRY_MULTIPLIER"},
		},
		&cli.DurationFlag{
			Name:    "retry-max-interval",
			Value:   60 * time.Second,
			EnvVars: []string{"ETHEREUM_RPC_RETRY_MAX_INTERVAL"},
		},
		&cli.DurationFlag{
			Name:    "retry-max-elapsed-time",
			Value:   300 * time.Second,
			EnvVars: []string{"ETHEREUM_RPC_RETRY_MAX_ELAPSED_TIME"},
		},
		&cli.UintFlag{
			Name:    "retry-max-tries",
			Value:   20,
			EnvVars: []string{"ETHEREUM_RPC_RETRY_MAX_TRIES"},
		},
	}
}

func CallOptionsFromContext(ctx *cli.Context) []jsonrpc.CallOptionsFunc {
	var callOpts []jsonrpc.CallOptionsFunc
	callOpts = append(callOpts, jsonrpc.WithMatchBatchSize(ctx.Int("max-batch-size")))
	callOpts = append(callOpts, jsonrpc.WithMaxConcurrentRequests(ctx.Int("max-concurrent-requests")))
	callOpts = append(callOpts, jsonrpc.WithDisableBatch(ctx.Bool("disable-batch")))
	callOpts = append(callOpts, jsonrpc.WithFailOnError(ctx.Bool("fail-on-error")))
	callOpts = append(callOpts, jsonrpc.WithFailOnError(ctx.Bool("fail-on-retryable-error")))
	callOpts = append(callOpts, jsonrpc.WithFailOnNull(ctx.Bool("fail-on-null")))
	callOpts = append(callOpts, jsonrpc.WithRetryableStatusCodes(ctx.IntSlice("retryable-status-codes")))
	callOpts = append(callOpts, jsonrpc.WithRetryInitialInterval(ctx.Duration("retry-initial-interval")))
	callOpts = append(callOpts, jsonrpc.WithRetryRandomizationFactor(ctx.Float64("retry-randomization-factor")))
	callOpts = append(callOpts, jsonrpc.WithRetryMultiplier(ctx.Float64("retry-multiplier")))
	callOpts = append(callOpts, jsonrpc.WithRetryMaxInterval(ctx.Duration("retry-max-interval")))
	callOpts = append(callOpts, jsonrpc.WithRetryMaxElapsedTime(ctx.Duration("retry-max-elapsed-time")))
	callOpts = append(callOpts, jsonrpc.WithRetryMaxTries(ctx.Uint("retry-max-tries")))

	return callOpts
}
