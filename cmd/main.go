package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/agnosticeng/agnostic-clickhouse-udf/cmd/function"
	"github.com/agnosticeng/panicsafe"
	"github.com/urfave/cli/v2"
	slogctx "github.com/veqryn/slog-context"
)

func main() {
	app := cli.App{
		Name: "agnostic-clickhouse-udf",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "log-level",
				Value:   0,
				EnvVars: []string{"LOG_LEVEL"},
			},
			&cli.StringFlag{
				Name:    "log-path",
				EnvVars: []string{"LOG_PATH"},
			},
		},
		Before: func(ctx *cli.Context) error {
			var (
				path = ctx.String("log-path")
				w    io.WriteCloser
				err  error
				lvl  slog.LevelVar
			)

			lvl.Set(slog.Level(ctx.Int("log-level")))

			if len(path) == 0 {
				w = os.Stderr
			} else {
				w, err = os.Create(path)
			}

			if err != nil {
				return err
			}

			slog.NewTextHandler(w, &slog.HandlerOptions{AddSource: true, Level: &lvl})

			var (
				handler = slogctx.NewHandler(slog.NewTextHandler(w, nil), nil)
				logger  = slog.New(handler)
			)

			slog.SetDefault(logger)
			ctx.Context = slogctx.NewCtx(ctx.Context, logger)
			return nil
		},
		Commands: []*cli.Command{
			function.Command(),
		},
	}

	var err = panicsafe.Recover(func() error { return app.Run(os.Args) })

	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
}
