package xun

import "encoding/json"

// R alias map[string]interface{}, R is the first letter of "Row"
type R map[string]interface{}

// AnyToR cast any inteface to R type
func AnyToR(v interface{}) (R, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var r R
	err = json.Unmarshal(bytes, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// AnyToRs cast any inteface to R Slice
func AnyToRs(v interface{}) ([]R, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var rs []R
	err = json.Unmarshal(bytes, &rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
