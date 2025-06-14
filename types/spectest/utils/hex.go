package utils

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/attestantio/go-eth2-client/spec"
)

// OrderedMap is a map that preserves field order
type OrderedMap struct {
	fields []fieldInfo
}

type fieldInfo struct {
	name  string
	value interface{}
}

func (om *OrderedMap) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteString("{")
	for i, f := range om.fields {
		if i > 0 {
			buf.WriteString(",")
		}
		// Marshal the field name
		name, err := json.Marshal(f.name)
		if err != nil {
			return nil, err
		}
		buf.Write(name)
		buf.WriteString(":")
		// Marshal the field value
		val, err := json.Marshal(f.value)
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

// ConvertToHexMap recursively converts byte arrays to hex strings while preserving field order
func ConvertToHexMap(v reflect.Value) interface{} {
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return ConvertToHexMap(v.Elem())
	case reflect.Struct:
		// Use OrderedMap to maintain field order
		om := &OrderedMap{
			fields: make([]fieldInfo, 0, v.NumField()),
		}
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldName := t.Field(i).Name

			// Special handling for Version field - convert to string
			if fieldName == "Version" {
				om.fields = append(om.fields, fieldInfo{name: fieldName, value: field.Interface().(spec.DataVersion).String()})
				continue
			}

			// For all other fields, recursively process
			om.fields = append(om.fields, fieldInfo{name: fieldName, value: ConvertToHexMap(field)})
		}
		return om
	case reflect.Slice:
		// Handle nil slices
		if v.IsNil() {
			return nil
		}
		// For arrays/slices, check if they are byte arrays
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return hex.EncodeToString(v.Bytes())
		}
		// For non-byte arrays/slices, process each element
		l := v.Len()
		arr := make([]interface{}, l)
		for i := 0; i < l; i++ {
			arr[i] = ConvertToHexMap(v.Index(i))
		}
		return arr
	case reflect.Array:
		// For arrays, check if they are byte arrays
		if v.Type().Elem().Kind() == reflect.Uint8 {
			bytes := make([]byte, v.Len())
			for i := 0; i < v.Len(); i++ {
				bytes[i] = byte(v.Index(i).Uint())
			}
			return hex.EncodeToString(bytes)
		}
		// For non-byte arrays, process each element
		l := v.Len()
		arr := make([]interface{}, l)
		for i := 0; i < l; i++ {
			arr[i] = ConvertToHexMap(v.Index(i))
		}
		return arr
	case reflect.Map:
		// Handle nil maps
		if v.IsNil() {
			return nil
		}
		m := make(map[string]interface{})
		for _, key := range v.MapKeys() {
			m[fmt.Sprintf("%v", key.Interface())] = ConvertToHexMap(v.MapIndex(key))
		}
		return m
	default:
		// For all other types, return the value as is
		return v.Interface()
	}
}

// ConvertHexToBytes recursively converts hex strings and integer arrays to byte arrays in a map
func ConvertHexToBytes(v interface{}) {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, v := range val {
			switch vv := v.(type) {
			case string:
				// Try to decode as hex
				if bytes, err := hex.DecodeString(vv); err == nil {
					val[k] = bytes
				} else {
					ConvertHexToBytes(v)
				}
			case []interface{}:
				// Check if it's an array of integers (byte array)
				if len(vv) > 0 {
					if _, ok := vv[0].(float64); ok {
						// Convert []interface{} of float64 to []byte
						bytes := make([]byte, len(vv))
						for i, f := range vv {
							bytes[i] = byte(f.(float64))
						}
						val[k] = bytes
					} else {
						ConvertHexToBytes(vv)
					}
				}
			case map[string]interface{}:
				ConvertHexToBytes(vv)
			}
		}
	case []interface{}:
		for i, v := range val {
			switch vv := v.(type) {
			case string:
				// Try to decode as hex
				if bytes, err := hex.DecodeString(vv); err == nil {
					val[i] = bytes
				} else {
					ConvertHexToBytes(v)
				}
			case []interface{}:
				// Check if it's an array of integers (byte array)
				if len(vv) > 0 {
					if _, ok := vv[0].(float64); ok {
						// Convert []interface{} of float64 to []byte
						bytes := make([]byte, len(vv))
						for i, f := range vv {
							bytes[i] = byte(f.(float64))
						}
						val[i] = bytes
					} else {
						ConvertHexToBytes(vv)
					}
				}
			case map[string]interface{}:
				ConvertHexToBytes(vv)
			}
		}
	}
}
