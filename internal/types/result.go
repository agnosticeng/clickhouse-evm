package types

import "encoding/json"

type Result struct {
	Error string          `json:"error,omitempty"`
	Value json.RawMessage `json:"value,omitempty"`
}
