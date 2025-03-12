/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package ex

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
)

type Error struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// Error 实现 error 接口
func (e *Error) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// String 实现 fmt.Stringer 接口
func (e *Error) String() string {
	if len(e.Meta) == 0 {
		return e.Error()
	}

	metaJSON, err := json.Marshal(e.Meta)
	if err != nil {
		return fmt.Sprintf("Error %d: %s (meta marshal failed: %v)", e.Code, e.Message, err)
	}

	return fmt.Sprintf("Error %d: %s | Meta: %s", e.Code, e.Message, string(metaJSON))
}

// MarshalZerologObject 实现 zerolog.LogObjectMarshaler 接口
func (e *Error) MarshalZerologObject(event *zerolog.Event) {
	event.Int("code", e.Code).
		Str("message", e.Message)

	if len(e.Meta) > 0 {
		metaEvent := zerolog.Dict()
		for k, v := range e.Meta {
			switch vt := v.(type) {
			case string:
				metaEvent.Str(k, vt)
			case int:
				metaEvent.Int(k, vt)
			case int64:
				metaEvent.Int64(k, vt)
			case float64:
				metaEvent.Float64(k, vt)
			case bool:
				metaEvent.Bool(k, vt)
			default:
				if bs, err := json.Marshal(v); err == nil {
					metaEvent.RawJSON(k, bs)
				}
			}
		}
		event.Dict("meta", metaEvent)
	}
}

func (e *Error) SetMessage(message string) *Error {
	e.Message = message
	return e
}

func (e *Error) SetMeta(key string, value interface{}) *Error {
	if e.Meta == nil {
		e.Meta = make(map[string]interface{})
	}
	e.Meta[key] = value
	return e
}
