package e2e

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

type simpleLogConsumer struct {
	t *testing.T
}

func (lg *simpleLogConsumer) Accept(log testcontainers.Log) {
	lg.t.Log(string(log.Content))
}

//go:embed *.sql
var sqlFiles embed.FS

func TestE2E(t *testing.T) {
	var ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// setup network
	net, err := network.New(ctx)
	require.NoError(t, err)
	testcontainers.CleanupNetwork(t, net)

	// start clickhouse-server container and mount UDF bundle
	configPath, err := filepath.Abs(filepath.Join(".", "config.xml"))
	require.NoError(t, err)
	installBundleScriptPath, err := filepath.Abs(filepath.Join(".", "install-bundle.sh"))
	require.NoError(t, err)
	bundlePath, err := filepath.Abs(filepath.Join("..", "tmp", "bundle.tar.gz"))
	require.NoError(t, err)

	clickhouseContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:    "clickhouse/clickhouse-server:25.3",
			Networks: []string{net.Name},
			Env: map[string]string{
				"CLICKHOUSE_PASSWORD": "test",
			},
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      configPath,
					ContainerFilePath: "/etc/clickhouse-server/config.d/config.xml",
					FileMode:          700,
				},
				{
					HostFilePath:      installBundleScriptPath,
					ContainerFilePath: "/docker-entrypoint-initdb.d/install-bundle.sh",
					FileMode:          755,
				},
				{
					HostFilePath:      bundlePath,
					ContainerFilePath: "/bundle.tar.gz",
					FileMode:          700,
				},
			},
			// LogConsumerCfg: &testcontainers.LogConsumerConfig{
			// 	Consumers: []testcontainers.LogConsumer{&simpleLogConsumer{t}},
			// },
		},
		Started: true,
	})
	require.NoError(t, err)
	defer testcontainers.TerminateContainer(clickhouseContainer)

	// create clickhouse client
	clickhousPort, err := clickhouseContainer.MappedPort(ctx, "9000")
	require.NoError(t, err)
	opts, err := clickhouse.ParseDSN(fmt.Sprintf("tcp://default:test@localhost:%d/default", clickhousPort.Int()))
	require.NoError(t, err)
	conn, err := clickhouse.Open(opts)
	require.NoError(t, err)
	defer conn.Close()

	pingCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

pingLoop:
	for {
		select {
		case <-pingCtx.Done():
			t.FailNow()
		case <-time.After(time.Second):
			t.Log("Trying to ping ClickHouse server...")
			err = conn.Ping(ctx)
			if err == nil {
				break pingLoop
			}
			t.Logf("Failed to ping ClickHouse server: %s", err.Error())
		}
	}

	ctx = clickhouse.Context(
		ctx,
		clickhouse.WithSettings(clickhouse.Settings{
			"send_logs_level": "debug",
		}),
		clickhouse.WithLogs(func(l *clickhouse.Log) {
			if strings.Contains(l.Text, "Executable generates stderr:") {
				t.Log(l.Text)
			}
		}),
	)

	// reload UDFs
	conn.Exec(ctx, "system reload functions")

	// run SQL file queries
	entries, err := sqlFiles.ReadDir(".")
	require.NoError(t, err)

	for _, entry := range entries {
		t.Run(entry.Name(), func(t *testing.T) {
			content, err := os.ReadFile(entry.Name())
			require.NoError(t, err)
			queries := lo.Compact(strings.Split(string(content), ";;"))

			for i, query := range queries {
				var queryName = strconv.FormatInt(int64(i), 10)

				if match := regexp.MustCompile(`\s*--\s(.+)`).FindStringSubmatch(query); len(match) >= 2 {
					queryName = match[1]
				}

				t.Run(queryName, func(t *testing.T) {
					err = conn.Exec(ctx, query)
					require.NoError(t, err)
				})
			}
		})
	}
}
