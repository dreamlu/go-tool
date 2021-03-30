package json

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type CJSON []byte

func (j CJSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return string(j), nil
}

func (j *CJSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("不合法的JSON数据")
	}
	*j = append((*j)[0:0], s...)
	return nil
}

func (CJSON) GormDataType() string {
	return "json"
}

func (j CJSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		// use "" replace null
		return []byte("\"\""), nil
	}
	return j, nil
}

func (j *CJSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("CJSON nil error")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

func (j CJSON) Equals(j1 CJSON) bool {
	return bytes.Equal(j, j1)
}

func (j CJSON) String() string {
	return string(j)
}

// Deprecated
// use Unmarshal replace
func (j CJSON) Struct(value interface{}) error {
	return j.Unmarshal(value)
}

// Unmarshal support Struct/Array
func (j CJSON) Unmarshal(value interface{}) error {
	return json.Unmarshal(j, value)
}
