package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JSONB[T any] struct {
	v *T
}

func (j *JSONB[T]) Scan(src any) error {
	if src == nil {
		j.v = nil
		return nil
	}

	j.v = new(T)
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, j.v)
	case string:
		return json.Unmarshal([]byte(v), j.v)
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (j JSONB[T]) Value() (driver.Value, error) {
	raw, err := json.Marshal(j.v)
	return raw, err
}

func (j *JSONB[T]) Get() *T {
	return j.v
}

func (j *JSONB[T]) Set(v T) {
	j.v = &v
}
