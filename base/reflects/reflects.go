package reflects

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-json"
)

// MergeStruct 使用泛型合并两个相同类型的结构体，零值字段将被跳过
func MergeStruct[T any](existing *T, newData *T) {
	valExisting := reflect.ValueOf(existing).Elem()
	valNewData := reflect.ValueOf(newData).Elem()

	for i := 0; i < valExisting.NumField(); i++ {
		valueFieldExisting := valExisting.Field(i)
		valueFieldNewData := valNewData.Field(i)

		if !valueFieldExisting.CanSet() {
			continue
		}

		// 检查是否为零值
		if !valueFieldNewData.IsZero() {
			valueFieldExisting.Set(valueFieldNewData)
		}
	}
}

// MergeStructData 保持向后兼容的版本
func MergeStructData(existing, newData interface{}) interface{} {
	valExisting := reflect.ValueOf(existing).Elem()
	valNewData := reflect.ValueOf(newData).Elem()

	for i := 0; i < valExisting.NumField(); i++ {
		valueFieldExisting := valExisting.Field(i)
		valueFieldNewData := valNewData.Field(i)

		if !valueFieldExisting.CanSet() {
			continue
		}

		if !valueFieldNewData.IsZero() {
			valueFieldExisting.Set(valueFieldNewData)
		}
	}

	return existing
}

// StructToMap 使用泛型将结构体转换为 map
func StructToMap[T any](configStruct T) (map[string]interface{}, error) {
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
