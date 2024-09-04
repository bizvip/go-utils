/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package json

import (
	"bytes"

	"github.com/goccy/go-json"
)

func PrettyFormat(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}

// ToJsonWithNoErr 失败回空对象而不是错误
func ToJsonWithNoErr(payload interface{}, pretty bool) string {
	j, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}
	if pretty {
		return PrettyFormat(string(j))
	}
	return string(j)
}
