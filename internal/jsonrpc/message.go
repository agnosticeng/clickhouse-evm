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

func (m *Message) SetRequest(method string, params json.RawMessage) {
	m.Version = "2.0"
	m.Id = strconv.FormatUint(rng.Uint64(), 16)
	m.Method = method
	m.Params = params
}

func (m *Message) Clear() {
	m.Version = ""
	m.Id = ""
	m.Method = ""
	m.Params = nil
	m.Result = nil
	m.Error = nil
}

func NewRequest(method string, params json.RawMessage) *Message {
	var m Message
	m.SetRequest(method, params)
	return &m
}

type ResponseError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (err ResponseError) Error() string {
	return err.Message
}

type Batch []Message

func (b Batch) Clear() {
	for i := 0; i < len(b); i++ {
		(&b[i]).Clear()
	}
}

type MessageOrBatch struct {
	Message *Message
	Batch   Batch
}

func (mob *MessageOrBatch) UnmarshalJSON(input []byte) error {
	switch {
	case input[0] == '{':
		if err := json.Unmarshal(input, &mob.Message); err != nil {
			return err
		}

	case input[0] == '[':
		if err := json.Unmarshal(input, &mob.Batch); err != nil {
			return err
		}

	default:
		return fmt.Errorf("neither a message nor a batch of messages: %s", string(input))
	}

	return nil
}

func (mob *MessageOrBatch) MarshalJSON() ([]byte, error) {
	switch {
	case mob.Message != nil:
		return json.Marshal(mob.Message)
	case len(mob.Batch) > 0:
		return json.Marshal(mob.Batch)
	default:
		return nil, fmt.Errorf("neither a message nor a batch of messages")
	}
}
