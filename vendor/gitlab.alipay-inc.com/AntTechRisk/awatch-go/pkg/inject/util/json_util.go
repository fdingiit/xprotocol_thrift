package util

import "encoding/json"

func ToJsonString(v interface{}) string {
	if v == nil {
		return ""
	}

	ret, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(ret)
}
