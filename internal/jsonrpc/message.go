package jsonrpc

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type Message struct {
	// common fields
	Version string `json:"jsonrpc"`
	Id      string `json:"id,omitempty"`

	// request fields
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`

	// response fields
	Result json.RawMessage `json:"result,omitempty"`
	Error  *ResponseError  `json:"error,omitempty"`
}

func NewMessage() *Message {
	return &Message{
		Version: "2.0",
		Id:      strconv.FormatUint(rng.Uint64(), 16),
	}
}

type ResponseError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (err ResponseError) Error() string {
	return err.Message
}

type Payload struct {
	Message  *Message
	Messages []*Message
}

func (p *Payload) UnmarshalJSON(input []byte) error {
	switch {
	case input[0] == '{':
		var msg Message

		if err := json.Unmarshal(input, &msg); err != nil {
			return err
		}

		p.Message = &msg

	case input[0] == '[':
		var msgs []*Message

		if err := json.Unmarshal(input, &msgs); err != nil {
			return err
		}

		p.Messages = msgs

	default:
		return fmt.Errorf("cannot recognize rpc request: %s", string(input))
	}

	return nil
}

func (p Payload) MarshalJSON() ([]byte, error) {
	switch {
	case p.Message != nil:
		return json.Marshal(p.Message)
	case len(p.Messages) > 0:
		return json.Marshal(p.Messages)
	default:
		return nil, fmt.Errorf("payload is neither a message nor a batch of messages")
	}
}
