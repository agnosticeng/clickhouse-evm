package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/samber/lo"
	"github.com/sourcegraph/conc/iter"
	slogctx "github.com/veqryn/slog-context"
)

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient(ctx context.Context) (*HTTPClient, error) {
	var client = &http.Client{}
	return &HTTPClient{
		client: client,
	}, nil
}

func (c *HTTPClient) BatchCall(ctx context.Context, endpoint string, reqs []*Message, optFuncs ...BatchOptionsFunc) ([]*Message, error) {
	if len(reqs) == 0 {
		return nil, nil
	}

	var opts = NewBatchOptions(optFuncs...)

	if opts.chunkSize <= 0 || opts.chunkSize >= len(reqs) {
		return c.batchCall(ctx, endpoint, reqs)
	}

	var (
		mapper = iter.Mapper[[]*Message, []*Message]{
			MaxGoroutines: opts.concurrencyLimit,
		}
		chunks = lo.Chunk(reqs, opts.chunkSize)
	)

	chunksRes, err := mapper.MapErr(chunks, func(reqs *[]*Message) ([]*Message, error) {
		return c.batchCall(ctx, endpoint, *reqs)
	})

	if err != nil {
		return nil, err
	}

	return lo.Flatten(chunksRes), nil
}

func (c *HTTPClient) batchCall(ctx context.Context, endpoint string, reqs []*Message) ([]*Message, error) {
	var (
		buf    bytes.Buffer
		res    Payload
		logger = slogctx.FromCtx(ctx)
	)

	if len(reqs) == 0 {
		return nil, nil
	}

	if err := json.NewEncoder(&buf).Encode(reqs); err != nil {
		return nil, err
	}

	if logger.Enabled(ctx, slog.LevelDebug) {
		logger.Debug("JSON-RPC request", "content", buf.String())
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		endpoint,
		&buf,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if logger.Enabled(ctx, slog.Level(-10)) {
		var js, _ = json.Marshal(res)
		logger.Log(ctx, slog.Level(-10), "JSON-RPC response", "content", js)
	}

	if len(res.Messages) == 0 && res.Message != nil {
		return nil, res.Message.Error
	}

	if len(reqs) != len(res.Messages) {
		return nil, fmt.Errorf("sent %d request but received %d responses", len(reqs), len(res.Messages))
	}

	return res.Messages, nil
}

func (c *HTTPClient) Close() error {
	return nil
}
