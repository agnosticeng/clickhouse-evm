package main

import (
	"fmt"
	"os"

	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name: "clickhouse-udfs",
		Commands: []*cli.Command{
			function.Command(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
