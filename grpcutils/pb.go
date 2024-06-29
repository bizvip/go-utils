/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package grpcutils

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// ToAny 封装了将任意 proto.Message 转换为 anypb.Any 类型的通用函数
func ToAny(pb proto.Message) (*anypb.Any, error) {
	anyData, err := anypb.New(pb)
	if err != nil {
		return nil, err
	}
	return anyData, nil
}
