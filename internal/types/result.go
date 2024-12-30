package types

import (
	"encoding/json"

	"github.com/samber/lo"
)

type Result struct {
	Error string `json:"error,omitempty"`
	Value any    `json:"value,omitempty"`
}

func (res *Result) ToJSON() []byte {
	return lo.Must(json.Marshal(res))
}
