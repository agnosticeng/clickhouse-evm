package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
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

func (c *HTTPClient) BatchCall(ctx context.Context, endpoint string, batch Batch, optFuncs ...CallOptionsFunc) error {
	if len(batch) == 0 {
		return nil
	}

	var opts = NewCallOptions(optFuncs...)

	if err := opts.ParseFromEndpoint(endpoint); err != nil {
		return err
	}

	if opts.disableBatch {
		return c.multiCall(ctx, endpoint, batch, *opts)
	} else {
		return c.batchCall(ctx, endpoint, batch, *opts)
	}
}

func (c *HTTPClient) multiCall(ctx context.Context, endpoint string, batch Batch, opts CallOptions) error {
	var pool = pool.New().
		WithContext(ctx).
		WithCancelOnError().
		WithFirstError().
		WithMaxGoroutines(opts.maxConcurrentRequests)

	for i := 0; i < len(batch); i++ {
		pool.Go(func(ctx context.Context) error {
			return c.doCall(ctx, endpoint, &batch[i], opts)
		})
	}

	return pool.Wait()
}

func (c *HTTPClient) batchCall(ctx context.Context, endpoint string, batch Batch, opts CallOptions) error {
	if opts.maxBatchSize >= len(batch) {
		return c.doBatchCall(ctx, endpoint, batch, opts)
	}

	var pool = pool.New().
		WithContext(ctx).
		WithCancelOnError().
		WithFirstError().
		WithMaxGoroutines(opts.maxConcurrentRequests)

	for _, chunk := range lo.Chunk(batch, opts.maxBatchSize) {
		pool.Go(func(ctx context.Context) error {
			return c.doBatchCall(ctx, endpoint, chunk, opts)
		})
	}

	return pool.Wait()
}

func (c *HTTPClient) doBatchCall(ctx context.Context, endpoint string, batch Batch, opts CallOptions) error {
	var (
		buf    bytes.Buffer
		res    = MessageOrBatch{Batch: batch}
		logger = slogctx.FromCtx(ctx)
	)

	if len(batch) == 0 {
		return nil
	}

	if err := json.NewEncoder(&buf).Encode(batch); err != nil {
		return fmt.Errorf("failed to encode JSON request: %w", err)
	}

	if logger.Enabled(ctx, slog.LevelDebug) {
		logger.Debug("JSON-RPC request", "endpoint", endpoint, "content", buf.String())
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, &buf)

	if err != nil {
		return fmt.Errorf("failed to create new HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)

	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if logger.Enabled(ctx, slog.Level(-10)) {
		logger.Log(ctx, slog.Level(-10), "JSON-RPC response", "endpoint", endpoint, "content", lo.Must(json.Marshal(res)))
	}

	if len(res.Batch) == 0 && res.Message != nil {
		return fmt.Errorf("JSONRPC error response: %s", res.Message.Error)
	}

	if opts.failOnError || opts.failOnNull || opts.failOnRetryableError {
		for _, msg := range res.Batch {
			if err := processMessageError(&msg, opts); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *HTTPClient) Call(ctx context.Context, endpoint string, msg *Message, optFuncs ...CallOptionsFunc) error {
	var opts = NewCallOptions(optFuncs...)

	if err := opts.ParseFromEndpoint(endpoint); err != nil {
		return err
	}

	return c.doCall(ctx, endpoint, msg, *opts)
}

func (c *HTTPClient) doCall(ctx context.Context, endpoint string, msg *Message, opts CallOptions) error {
	var (
		buf    bytes.Buffer
		logger = slogctx.FromCtx(ctx)
	)

	if msg == nil {
		return nil
	}

	if err := json.NewEncoder(&buf).Encode(msg); err != nil {
		return nil
	}

	if logger.Enabled(ctx, slog.LevelDebug) {
		logger.Debug("JSON-RPC request", "endpoint", endpoint, "content", buf.String())
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, &buf)

	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := c.client.Do(httpReq)

	if err != nil {
		return err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return fmt.Errorf("bad status code: %d", httpResp.StatusCode)
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&msg); err != nil {
		return err
	}

	if logger.Enabled(ctx, slog.Level(-10)) {
		logger.Log(ctx, slog.Level(-10), "JSON-RPC response", "endpoint", endpoint, "content", lo.Must(json.Marshal(msg)))
	}

	if opts.failOnError || opts.failOnNull || opts.failOnRetryableError {
		if err := processMessageError(msg, opts); err != nil {
			return err
		}
	}

	return nil
}

func (c *HTTPClient) Close() error {
	return nil
}

func processMessageError(msg *Message, opts CallOptions) error {
	switch {
	case opts.failOnError && msg.Error != nil:
		return msg.Error

	case opts.failOnRetryableError && msg.Error != nil && opts.retryableErrorPredicate != nil:
		if opts.retryableErrorPredicate(msg.Error) {
			return msg.Error
		}

	case opts.failOnNull && (msg.Result == nil || bytes.Equal(msg.Result, []byte(`null`))):
		return fmt.Errorf("null result")
	}

	return nil
}
