/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package reflects

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-json"
)

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

// StructToMap 将结构体转换为 map
func StructToMap(configStruct interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(configStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct to json: %w", err)
	}

	var result map[string]interface{}
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json to map: %w", err)
	}

	return result, nil
}
