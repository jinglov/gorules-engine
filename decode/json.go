package decode

import (
	"encoding/json"
	"math"
	"strconv"
)

const (
	keySeparator = "_"
	maxDepth     = 1
)

func init() {
	RegisterDecoder("json", JsonReport{})
}

type JsonReport struct {
}

func (JsonReport) DecodeReportFromByte(m map[string]string, b []byte) (timestamp int64, err error) {
	data := make(map[string]interface{})
	err = json.Unmarshal(b, &data)
	if err != nil {
		return
	}
	jsonToMap("", data, m, 1)
	return
}

// 解析json: 不支持array, object仅支持2层嵌套
func jsonToMap(prefixKey string, data map[string]interface{}, m map[string]string, depth int) {
	for k, v := range data {
		key := k
		if prefixKey != "" {
			key = prefixKey + keySeparator + k
		}
		switch v.(type) {
		case map[string]interface{}:
			if depth <= maxDepth {
				jsonToMap(key, v.(map[string]interface{}), m, depth+1)
			}
		case float64:
			vInt, vDec := math.Modf(v.(float64))
			if vDec > 0 {
				m[key] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
			} else {
				m[key] = strconv.Itoa(int(vInt))
			}
		case string:
			m[key] = v.(string)
		case bool:
			m[key] = strconv.FormatBool(v.(bool))
		}
	}
	return
}
