package prop_test

import (
	"fmt"
	"general-agent/extension/prop"
	"sort"
	"testing"
)

func TestMapToJSON(t *testing.T) {
	data := map[string]string{
		"key1":                  "value1",
		"key2.subkey":           "value2",
		"key3.subkey":           "value3",
		"key4.subkey.subsubkey": "value4",
	}

	jsonData, err := prop.KVToJSON(data)
	if err != nil {
		fmt.Println("Failed to convert map to JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
func TestJSONToMap(t *testing.T) {
	jsonData := []byte(`
{
  "protect_param": [
    {
      "mode": 2,
      "url": "/aaa",
      "srcip": "any",
      "set_periodic": 0,
      "week_day": "",
      "day_enable_time": ""
    },
    {
      "mode": 2,
      "url": "/bbbb",
      "srcip": "any",
      "set_periodic": 0,
      "week_day": "",
      "day_enable_time": ""
    }
  ],
  "profile": "",
  "log": 1,
  "log_level": 6,
  "match_style": 0,
  "enable": 0,
  "action": 1,
  "priority": 0,
  "max_count": 256,
  "resp": 4224,
  "ha_backup":{
      "mode": 2,
      "url": "/bbbb",
      "srcip": "any",
      "set_periodic": 0,
      "week_day": "",
      "day_enable_time": ""
    },
  "sessionid": ""
}
`)

	result, err := prop.JSONToKV(jsonData)
	if err != nil {
		fmt.Println("Failed to convert JSON to map:", err)
		return
	}

	keys := make([]string, 0, len(result))
	for key := range result {
		keys = append(keys, key)
	}

	// 对切片进行排序
	sort.Strings(keys)

	// 循环排序后的切片，访问map中的值
	for _, key := range keys {
		fmt.Printf("%s\t%s\n", key, result[key])
	}

}
