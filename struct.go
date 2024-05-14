/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package goutils

import (
	"reflect"
)

type StructUtils struct{}

func NewStructUtils() *StructUtils {
	return &StructUtils{}
}

// MergeStructData 使用反射来合并两个struct 反射影响高性能
func MergeStructData(existing, newData interface{}) interface{} {
	valExisting := reflect.ValueOf(existing).Elem()
	valNewData := reflect.ValueOf(newData).Elem()

	for i := 0; i < valExisting.NumField(); i++ {
		valueFieldExisting := valExisting.Field(i)
		valueFieldNewData := valNewData.Field(i)

		if !valueFieldExisting.CanSet() {
			continue
		}

		if !reflect.DeepEqual(valueFieldNewData.Interface(), reflect.Zero(valueFieldNewData.Type()).Interface()) {
			valueFieldExisting.Set(valueFieldNewData)
		}
	}

	return existing
}
