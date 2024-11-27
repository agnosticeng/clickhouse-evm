package main

import (
	"os"

	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function"
	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/table_function"
	"github.com/agnosticeng/panicsafe"
	"github.com/agnosticeng/slogcli"
	"github.com/urfave/cli/v2"
	slogctx "github.com/veqryn/slog-context"
)

func main() {
	app := cli.App{
		Name:   "agnostic-clickhouse-udf",
		Flags:  slogcli.SlogFlags(),
		Before: slogcli.SlogBefore,
		After:  slogcli.SlogAfter,
		ExitErrHandler: func(ctx *cli.Context, err error) {
			slogctx.FromCtx(ctx.Context).Error(err.Error())
		},
		Commands: []*cli.Command{
			function.Command(),
			table_function.Command(),
		},
	}

	var err = panicsafe.Recover(func() error { return app.Run(os.Args) })

	if err != nil {
		os.Exit(1)
	}
}
