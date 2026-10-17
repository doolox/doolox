package cms

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
)

// Value types supported by the KeyValue model.
const (
	ValueString = "string"
	ValueInt    = "int"
	ValueBool   = "bool"
)

// KeyValue model stores a single keyed value that may be a string, an
// integer, or a boolean. The concrete type is recorded in Type and the value
// itself is persisted as its canonical string representation in Value.
type KeyValue struct {
	gorm.Model
	Key   string `sql:"size:255;unique_index"`
	Type  string `sql:"size:16"`
	Value string `sql:"size:255"`
}

// NewKeyValue builds a KeyValue from any of the supported value types. It
// returns an error if v is not a string, integer, or boolean.
func NewKeyValue(key string, v interface{}) (*KeyValue, error) {
	kv := &KeyValue{Key: key}
	if err := kv.Set(v); err != nil {
		return nil, err
	}
	return kv, nil
}

// Set assigns v to the KeyValue. Supported types are string, int, int64, and bool.
func (kv *KeyValue) Set(v interface{}) error {
	switch value := v.(type) {
	case string:
		kv.Type = ValueString
		kv.Value = value
	case int:
		kv.Type = ValueInt
		kv.Value = strconv.FormatInt(int64(value), 10)
	case int64:
		kv.Type = ValueInt
		kv.Value = strconv.FormatInt(value, 10)
	case bool:
		kv.Type = ValueBool
		kv.Value = strconv.FormatBool(value)
	default:
		return fmt.Errorf("keyvalue: unsupported value type %T", v)
	}
	return nil
}

// Get returns the stored value as the appropriate Go type (string, int64, or bool).
func (kv *KeyValue) Get() (interface{}, error) {
	switch kv.Type {
	case ValueString:
		return kv.Value, nil
	case ValueInt:
		return strconv.ParseInt(kv.Value, 10, 64)
	case ValueBool:
		return strconv.ParseBool(kv.Value)
	default:
		return nil, fmt.Errorf("keyvalue: unknown value type %q", kv.Type)
	}
}

// String returns the stored value as a string.
func (kv *KeyValue) String() string { return kv.Value }

// Int returns the stored value as an int64.
func (kv *KeyValue) Int() (int64, error) {
	return strconv.ParseInt(kv.Value, 10, 64)
}

// Bool returns the stored value as a bool.
func (kv *KeyValue) Bool() (bool, error) {
	return strconv.ParseBool(kv.Value)
}
