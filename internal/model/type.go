package model

import (
	"database/sql/driver"
	"encoding/json"
)

type StringList []string

func (o *StringList) Scan(value interface{}) error {
	bytes, _ := value.([]byte)
	return json.Unmarshal(bytes, o)
}

func (o StringList) Value() (driver.Value, error) {
	return json.Marshal(o)
}
