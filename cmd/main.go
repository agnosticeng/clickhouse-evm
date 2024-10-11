package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function"
	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/table_function"
	"github.com/agnosticeng/panicsafe"
	"github.com/agnosticeng/slogcli"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:   "agnostic-clickhouse-udf",
		Flags:  slogcli.SlogFlags(),
		Before: slogcli.SlogBefore,
		After:  slogcli.SlogAfter,
		Commands: []*cli.Command{
			function.Command(),
			table_function.Command(),
		},
	}

	var err = panicsafe.Recover(func() error { return app.Run(os.Args) })

	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
}
