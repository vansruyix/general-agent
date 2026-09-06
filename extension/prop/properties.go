package prop

import (
	"encoding/json"
	"fmt"
	"strings"
)

// KVToMap 将键值对转换为嵌套的Map
func KVToMap(keyValuePairs map[string]string) (map[string]interface{}, error) {
	jsonData := make(map[string]interface{})
	for key, value := range keyValuePairs {
		if err := setValueInJSONData(jsonData, key, value); err != nil {
			return nil, err
		}
	}
	return jsonData, nil
}

// KVToJSON 将键值对转换为JSON
func KVToJSON(keyValuePairs map[string]string) ([]byte, error) {
	data, err := KVToMap(keyValuePairs)
	if err != nil {
		return nil, fmt.Errorf("convert to map failed: %w", err)
	}
	return json.Marshal(data)
}

// 在JSON数据中设置键值对
func setValueInJSONData(data map[string]interface{}, key string, value string) error {
	segments := splitEscapedDot(key)
	for i, segment := range segments {
		isArray := strings.Contains(segment, "[")
		if isArray {
			arrayKey := segment[:strings.Index(segment, "[")]
			arrayIndex := segment[strings.Index(segment, "[")+1 : strings.Index(segment, "]")]

			if data[arrayKey] == nil {
				data[arrayKey] = make([]interface{}, 0)
			}

			index, err := getIndex(arrayIndex)
			if err != nil {
				return err
			}
			if i == len(segments)-1 {
				// 判断索引是否在切片范围内
				if index < len(data[arrayKey].([]interface{})) {
					data[arrayKey].([]interface{})[index] = value
				} else {
					// 如果索引超过切片长度，使用 append 进行扩容
					data[arrayKey] = append(data[arrayKey].([]interface{}), make([]interface{}, index-len(data[arrayKey].([]interface{}))+1)...)
					data[arrayKey].([]interface{})[index] = value
				}
			} else {
				if len(data[arrayKey].([]interface{})) <= index {
					data[arrayKey] = append(data[arrayKey].([]interface{}), make(map[string]interface{}))
				}
				data = data[arrayKey].([]interface{})[index].(map[string]interface{})
			}
		} else {
			if i == len(segments)-1 {
				data[segment] = value
			} else {
				if data[segment] == nil {
					data[segment] = make(map[string]interface{})
				}
				data = data[segment].(map[string]interface{})
			}
		}
	}
	return nil
}
func splitEscapedDot(input string) []string {
	var segments []string
	var currentSegment string

	for i := 0; i < len(input); i++ {
		if input[i] == '\\' && i+1 < len(input) && input[i+1] == '.' {
			// Escaped dot, skip the backslash and add dot to current segment
			currentSegment += "."
			i++ // Skip the next character (dot)
		} else if input[i] == '.' {
			// Unescaped dot, start a new segment
			segments = append(segments, currentSegment)
			currentSegment = ""
		} else {
			// Other characters, add to current segment
			currentSegment += string(input[i])
		}
	}

	// Add the last segment
	segments = append(segments, currentSegment)

	return segments
}

// 获取数组索引
func getIndex(indexStr string) (int, error) {
	var index int
	if _, err := fmt.Sscanf(indexStr, "%d", &index); err != nil {
		return 0, err
	}
	return index, nil
}

// JSONToKV 将Json转为扁平化的kv。类似：user.name=xxx
func JSONToKV(jsonData []byte) (map[string]string, error) {
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}
	return MapToKV(data)
}

// MapToKV 将嵌套的map转为扁平化的map kv
func MapToKV(data map[string]interface{}) (map[string]string, error) {
	// 将嵌套数组展平
	keyValuePairs := make(map[string]string)
	flattenJSON(data, "", keyValuePairs)
	return keyValuePairs, nil
}

// 扁平化JSON
func flattenJSON(data map[string]interface{}, prefix string, keyValuePairs map[string]string) {
	for key, value := range data {
		fullKey := prefix + strings.ReplaceAll(key, ".", "\\.")
		switch value := value.(type) {
		case map[string]interface{}:
			flattenJSON(value, fullKey+".", keyValuePairs)
		case []interface{}:
			flattenJSONArray(value, fullKey, keyValuePairs)
		default:
			keyValuePairs[fullKey] = fmt.Sprintf("%v", value)
		}
	}
}

// 将嵌套数组扁平化为键值对
func flattenJSONArray(array []interface{}, prefix string, keyValuePairs map[string]string) {
	for i, item := range array {
		switch item := item.(type) {
		case map[string]interface{}:
			flattenJSON(item, fmt.Sprintf("%s[%d].", prefix, i), keyValuePairs)
		default:
			keyValuePairs[fmt.Sprintf("%s[%d]", prefix, i)] = fmt.Sprintf("%v", item)
		}
	}
}
