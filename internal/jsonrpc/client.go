package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

func (c *HTTPClient) BatchCall(ctx context.Context, endpoint string, msgs []*Message) ([]*Message, error) {
	var (
		buf bytes.Buffer
		res Payload
		m   = make(map[string]int)
	)

	if len(msgs) == 0 {
		return nil, nil
	}

	if err := json.NewEncoder(&buf).Encode(msgs); err != nil {
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

	for i, msg := range msgs {
		if len(msg.Id) > 0 {
			m[msg.Id] = i
		}
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

	if len(msgs) != len(res.Messages) {
		return nil, fmt.Errorf("sent %d request but received %d responses", len(msgs), len(res.Messages))
	}

	return res.Messages, nil
}

func (c *HTTPClient) Close() error {
	return nil
}
