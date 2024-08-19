/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package pb

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// PbToJson protobuf 消息转换为 json，保留零值
func PbToJson(data proto.Message) ([]byte, error) {
	opts := protojson.MarshalOptions{
		EmitUnpopulated: true, // 默认值不忽略
		UseProtoNames:   true, // 使用proto name返回http字段
		UseEnumNumbers:  true, // 将枚举值作为数字发出，默认为枚举值的字符串
	}
	respBytes, err := opts.Marshal(data)
	if err != nil {
		return nil, err
	}
	return respBytes, nil
}
