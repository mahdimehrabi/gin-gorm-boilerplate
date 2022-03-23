package utils

import (
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func BytesJsonToMap(bytes []byte) (map[string]interface{}, error) {
	mp := make(map[string]interface{})
	err := json.Unmarshal(bytes, &mp)
	return mp, err
}

// JsonStringify convert map to string without error handling
func JsonStringify(input gin.H) string {
	retStr, _ := json.Marshal(input)
	return string(retStr)
}

func MapToJsonBytesBuffer(mp map[string]interface{}) *bytes.Buffer {
	j, err := json.Marshal(mp)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(j)
}
