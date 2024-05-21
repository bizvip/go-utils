/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package goutils

import (
	"bytes"
	"encoding/json"
)

type JsonUtils struct{}

func NewJsonUtils() *JsonUtils {
	return &JsonUtils{}
}

func (ju *JsonUtils) JSONPrettyFormat(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}

// ToJSONUnsafe returns "{}" on failure case
func (ju *JsonUtils) ToJSONUnsafe(payload interface{}, pretty bool) string {
	j, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}
	if pretty {
		return ju.JSONPrettyFormat(string(j))
	}
	return string(j)
}
