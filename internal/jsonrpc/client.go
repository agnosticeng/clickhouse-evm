package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/samber/lo"
	"github.com/sourcegraph/conc/iter"
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

func (c *HTTPClient) Call(ctx context.Context, endpoint string, msg *Message) (*Message, error) {
	var (
		buf bytes.Buffer
		res Message
	)

	if err := json.NewEncoder(&buf).Encode(msg); err != nil {
		return nil, err
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

	return &res, nil
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
		buf bytes.Buffer
		res Payload
	)

	if len(reqs) == 0 {
		return nil, nil
	}

	if err := json.NewEncoder(&buf).Encode(reqs); err != nil {
		return nil, err
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
