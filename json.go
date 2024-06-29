/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package goutils

import (
	"bytes"

	"github.com/goccy/go-json"
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
