package hexencoding

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

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
	case reflect.Interface:
		if v.IsNil() {
			return nil
		}
		return ConvertToHexMap(v.Elem())
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return ConvertToHexMap(v.Elem())
	case reflect.Func:
		// Skip function types as they can't be marshaled to JSON
		return nil
	case reflect.Chan:
		// Skip channel types as they can't be marshaled to JSON
		return nil
	case reflect.Struct:
		// Use OrderedMap to maintain field order
		om := &OrderedMap{
			fields: make([]fieldInfo, 0, v.NumField()),
		}
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := t.Field(i)
			fieldName := fieldType.Name

			// Check for JSON tag to see if field should be ignored
			jsonTag := fieldType.Tag.Get("json")
			if jsonTag == "-" {
				// Skip fields with json:"-" tag
				continue
			}

			// Extract JSON field name from tag
			jsonFieldName := fieldName
			if jsonTag != "" {
				// Split by comma to handle options like "name,omitempty"
				parts := strings.Split(jsonTag, ",")
				if parts[0] != "" {
					jsonFieldName = parts[0]
				}
			}

			// Check if field is accessible (exported)
			if !field.CanInterface() {
				// Skip unexported fields
				continue
			}

			// Skip function and channel types
			if field.Kind() == reflect.Func || field.Kind() == reflect.Chan {
				continue
			}

			// Special handling for Version field - convert to string
			if fieldName == "Version" {
				om.fields = append(om.fields, fieldInfo{name: jsonFieldName, value: field.Interface().(spec.DataVersion).String()})
				continue
			}

			// Special handling for data fields
			if fieldName == "Data" || fieldName == "DataSSZ" || fieldName == "DataCd" || fieldName == "DataBlk" {
				if field.Type().Elem().Kind() == reflect.Uint8 {
					bytes := make([]byte, field.Len())
					for j := 0; j < field.Len(); j++ {
						bytes[j] = byte(field.Index(j).Uint())
					}
					om.fields = append(om.fields, fieldInfo{name: jsonFieldName, value: base64.StdEncoding.EncodeToString(bytes)})
					continue
				}
			}

			// Special handling for PubKey and Signature fields
			if fieldName == "PubKey" || fieldName == "Signature" {
				if field.Type().Elem().Kind() == reflect.Uint8 {
					bytes := make([]byte, field.Len())
					for j := 0; j < field.Len(); j++ {
						bytes[j] = byte(field.Index(j).Uint())
					}
					om.fields = append(om.fields, fieldInfo{name: jsonFieldName, value: "0x" + hex.EncodeToString(bytes)})
					continue
				}
			}

			// For all other fields, recursively process
			om.fields = append(om.fields, fieldInfo{name: jsonFieldName, value: ConvertToHexMap(field)})
		}

		// Special handling for qbft.Instance: add forceStop if true
		// Check for both qbft.Instance and *qbft.Instance
		if t.PkgPath() == "github.com/ssvlabs/ssv-spec/qbft" && t.Name() == "Instance" {
			// Look for unexported field "forceStop"
			forceStopField, ok := t.FieldByName("forceStop")
			if ok {
				forceStopVal := v.FieldByIndex(forceStopField.Index)
				if forceStopVal.Kind() == reflect.Bool && forceStopVal.Bool() {
					om.fields = append(om.fields, fieldInfo{name: "forceStop", value: true})
				}
			}
		}

		return om
	case reflect.Slice:
		// Handle nil slices
		if v.IsNil() {
			return nil
		}
		// For arrays/slices, check if they are byte arrays
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// Handle empty byte slices specially - return empty string instead of nil
			if v.Len() == 0 {
				return ""
			}
			// Get the field name from the parent struct if available
			fieldName := ""
			if v.CanAddr() && v.Addr().Type().Name() != "" {
				fieldName = v.Addr().Type().Name()
			}
			// Add 0x prefix for specific fields
			if fieldName == "PubKey" || fieldName == "Signature" {
				return "0x" + hex.EncodeToString(v.Bytes())
			}
			return hex.EncodeToString(v.Bytes())
		}
		// For non-byte arrays/slices, process each element
		l := v.Len()
		// Check SSZ tag sizes
		if v.Type().Elem().Kind() == reflect.Slice && v.Type().Elem().Elem().Kind() == reflect.Uint8 {
			// This is a [][]byte, check if it's Signatures field
			if v.CanAddr() && v.Addr().Type().Name() == "SignedSSVMessage" {
				// Signatures field has ssz-max:"13,256"
				if l > 13 {
					panic("Signatures array length exceeds SSZ max size of 13")
				}
				for i := 0; i < l; i++ {
					if v.Index(i).Len() > 256 {
						panic("Signature length exceeds SSZ max size of 256")
					}
				}
			}
		} else if v.Type().Elem().Kind() == reflect.Uint64 {
			// This is a []uint64, check if it's OperatorIDs field
			if v.CanAddr() && v.Addr().Type().Name() == "SignedSSVMessage" {
				// OperatorIDs field has ssz-max:"13"
				if l > 13 {
					panic("OperatorIDs array length exceeds SSZ max size of 13")
				}
			}
		}
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
			// Get the field name from the parent struct if available
			fieldName := ""
			if v.CanAddr() {
				// Try to get the field name from the parent struct
				parent := v.Addr().Type()
				if parent.Kind() == reflect.Struct {
					for i := 0; i < parent.NumField(); i++ {
						if parent.Field(i).Type == v.Type() {
							fieldName = parent.Field(i).Name
							break
						}
					}
				}
			}
			// Add 0x prefix for specific fields
			if fieldName == "PubKey" || fieldName == "Signature" {
				return "0x" + hex.EncodeToString(bytes)
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
		// For all other types, check if we can access the value
		if v.CanInterface() {
			return v.Interface()
		} else {
			// If we can't access the interface, try to get the value in a different way
			switch v.Kind() {
			case reflect.Bool:
				return v.Bool()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return v.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return v.Uint()
			case reflect.Float32, reflect.Float64:
				return v.Float()
			case reflect.String:
				return v.String()
			case reflect.Complex64, reflect.Complex128:
				return v.Complex()
			default:
				// For other types that we can't access, return a placeholder
				return fmt.Sprintf("<%s>", v.Type().String())
			}
		}
	}
}

func ToHexJSON(v interface{}) ([]byte, error) {
	// Check if the type has a custom MarshalJSON method
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rt := rv.Type().Elem()
		if _, hasCustom := rt.MethodByName("MarshalJSON"); hasCustom {
			// If it has a custom MarshalJSON method, use that instead
			return json.MarshalIndent(v, "", "  ")
		}
	}

	// Fall back to our hex conversion for types without custom marshaling
	return json.MarshalIndent(ConvertToHexMap(reflect.ValueOf(v)), "", "  ")
}
