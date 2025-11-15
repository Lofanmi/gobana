package gotil

import (
	jsoniter "github.com/json-iterator/go"
)

func JsonAs(src any, dst any) (err error) {
	data, err := jsoniter.Marshal(src)
	if err != nil {
		return
	}
	err = jsoniter.Unmarshal(data, dst)
	return
}
