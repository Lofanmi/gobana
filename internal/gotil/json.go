package gotil

import (
	"encoding/json"
)

func JsonAs(src any, dst any) (err error) {
	data, err := json.Marshal(src)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, dst)
	return
}
