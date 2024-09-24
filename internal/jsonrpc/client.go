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

func (c *HTTPClient) BatchCall(ctx context.Context, endpoint string, reqs []*Message, optFuncs ...CallOptionsFunc) ([]*Message, error) {
	if len(reqs) == 0 {
		return nil, nil
	}

	var opts = NewCallOptions(optFuncs...)

	if opts.disableBatch {
		return c.multiCall(ctx, endpoint, reqs, *opts)
	} else {
		return c.batchCall(ctx, endpoint, reqs, *opts)
	}
}

func (c *HTTPClient) multiCall(ctx context.Context, endpoint string, reqs []*Message, opts CallOptions) ([]*Message, error) {
	var (
		mapper = iter.Mapper[*Message, *Message]{
			MaxGoroutines: opts.maxConcurrentRequests,
		}
	)

	resps, err := mapper.MapErr(reqs, func(req **Message) (*Message, error) {
		return c.doSingleCall(ctx, endpoint, *req)
	})

	if err != nil {
		return nil, err
	}

	return resps, nil
}

func (c *HTTPClient) batchCall(ctx context.Context, endpoint string, reqs []*Message, opts CallOptions) ([]*Message, error) {
	if opts.maxBatchSize >= len(reqs) {
		return c.doBatchCall(ctx, endpoint, reqs)
	}

	var (
		mapper = iter.Mapper[[]*Message, []*Message]{
			MaxGoroutines: opts.maxConcurrentRequests,
		}
		chunks = lo.Chunk(reqs, opts.maxBatchSize)
	)

	chunksRes, err := mapper.MapErr(chunks, func(reqs *[]*Message) ([]*Message, error) {
		return c.doBatchCall(ctx, endpoint, *reqs)
	})

	if err != nil {
		return nil, err
	}

	return lo.Flatten(chunksRes), nil
}

func (c *HTTPClient) doBatchCall(ctx context.Context, endpoint string, reqs []*Message) ([]*Message, error) {
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

func (c *HTTPClient) doSingleCall(ctx context.Context, endpoint string, req *Message) (*Message, error) {
	var (
		buf    bytes.Buffer
		res    Message
		logger = slogctx.FromCtx(ctx)
	)

	if req == nil {
		return nil, nil
	}

	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		return nil, err
	}

	if logger.Enabled(ctx, slog.LevelDebug) {
		logger.Debug("JSON-RPC request", "content", buf.String())
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		endpoint,
		&buf,
	)

	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := c.client.Do(httpReq)

	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	// if httpResp.StatusCode != 200 {
	// 	return nil, fmt.Errorf("bad status code: %d", httpResp.StatusCode)
	// }

	if err := json.NewDecoder(httpResp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if logger.Enabled(ctx, slog.Level(-10)) {
		var js, _ = json.Marshal(res)
		logger.Log(ctx, slog.Level(-10), "JSON-RPC response", "content", js)
	}

	return &res, nil
}

func (c *HTTPClient) Close() error {
	return nil
}
