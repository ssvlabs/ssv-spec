package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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

			// Special handling for Data and DataSSZ fields - use base64 encoding
			if fieldName == "Data" || fieldName == "DataSSZ" || fieldName == "DataCd" || fieldName == "DataBlk" {
				if field.Type().Elem().Kind() == reflect.Uint8 {
					bytes := make([]byte, field.Len())
					for j := 0; j < field.Len(); j++ {
						bytes[j] = byte(field.Index(j).Uint())
					}
					om.fields = append(om.fields, fieldInfo{name: fieldName, value: base64.StdEncoding.EncodeToString(bytes)})
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
					om.fields = append(om.fields, fieldInfo{name: fieldName, value: "0x" + hex.EncodeToString(bytes)})
					continue
				}
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
			if v.CanAddr() && v.Addr().Type().Name() != "" {
				fieldName = v.Addr().Type().Name()
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
				// Try to decode as hex first
				if bytes, err := hex.DecodeString(vv); err == nil {
					// If the key ends with "Version" or "Root", it's likely a fixed-size array
					if strings.HasSuffix(k, "Version") || strings.HasSuffix(k, "Root") {
						// Convert to fixed-size array if needed
						switch {
						case strings.HasSuffix(k, "Version"):
							var version [4]byte
							copy(version[:], bytes)
							val[k] = version
						case strings.HasSuffix(k, "Root"):
							if k == "ExpectedSigningRoot" {
								val[k] = vv
							} else {
								var root [32]byte
								// Remove 0x prefix if present
								hexStr := vv
								// hexStr = strings.TrimPrefix(hexStr, "0x")
								bytes, err = hex.DecodeString(hexStr)
								if err != nil || len(bytes) != 32 {
									// If still not 32 bytes, raise error
									panic(fmt.Errorf("invalid root: %s", vv))
								}
								copy(root[:], bytes)
								val[k] = root
							}
						default:
							val[k] = bytes
						}
					} else if k == "MsgID" {
						// Special handling for MsgID
						var msgID [56]byte
						hexStr := vv
						if strings.HasPrefix(hexStr, "0x") {
							hexStr = strings.TrimPrefix(hexStr, "0x")
						}
						bytes, err = hex.DecodeString(hexStr)
						if err != nil || len(bytes) != 56 {
							// If still not 56 bytes, raise error
							panic(fmt.Errorf("invalid MsgID: %s", vv))
						}
						copy(msgID[:], bytes)
						val[k] = msgID
					} else if k == "PubKey" {
						// Special handling for BLS public key
						var pubKey [48]byte
						// Remove 0x prefix if present
						hexStr := vv
						hexStr = strings.TrimPrefix(hexStr, "0x")
						bytes, err = hex.DecodeString(hexStr)
						if err != nil || len(bytes) != 48 {
							// If still not 48 bytes, raise error
							panic(fmt.Errorf("invalid BLS public key: %s", vv))
						}
						copy(pubKey[:], bytes)
						val[k] = pubKey
					} else if k == "Signature" {
						// Special handling for BLS signature
						var signature [96]byte
						hexStr := vv
						hexStr = strings.TrimPrefix(hexStr, "0x")
						bytes, err = hex.DecodeString(hexStr)
						if err != nil || len(bytes) != 96 {
							// If still not 96 bytes, raise error
							panic(fmt.Errorf("invalid BLS signature: %s", vv))
						}
						copy(signature[:], bytes)
						val[k] = signature
					} else if k == "CommitteeID" {
						// Special handling for CommitteeID
						var committeeID [32]byte
						hexStr := vv
						hexStr = strings.TrimPrefix(hexStr, "0x")
						bytes, err = hex.DecodeString(hexStr)
						if err != nil || len(bytes) != 32 {
							// If still not 32 bytes, raise error
							panic(fmt.Errorf("invalid CommitteeID: %s", vv))
						}
						copy(committeeID[:], bytes)
						val[k] = committeeID
					} else if k == "DomainType" {
						// Special handling for DomainType
						var domainType [4]byte
						hexStr := vv
						hexStr = strings.TrimPrefix(hexStr, "0x")
						bytes, err = hex.DecodeString(hexStr)
						if err != nil || len(bytes) != 4 {
							// If still not 4 bytes, raise error
							panic(fmt.Errorf("invalid DomainType: %s", vv))
						}
						copy(domainType[:], bytes)
						val[k] = domainType
					} else if k == "WithdrawalCredentials" {
						// Keep WithdrawalCredentials as a string
						val[k] = vv
					} else if k == "Slot" {
						// Special handling for Slot - treat as uint64
						val[k] = vv
					} else {
						val[k] = bytes
					}
				} else {
					// If hex decoding fails, try base64
					if bytes, err := base64.StdEncoding.DecodeString(vv); err == nil {
						val[k] = bytes
					} else {
						// If both hex and base64 decoding fail, keep as is
						val[k] = vv
					}
				}
			case []interface{}:
				// Special handling for OperatorIDs and ValidatorSyncCommitteeIndices
				if k == "OperatorIDs" || k == "ValidatorSyncCommitteeIndices" {
					indices := make([]uint64, len(vv))
					for i, id := range vv {
						switch id := id.(type) {
						case float64:
							indices[i] = uint64(id)
						case string:
							idNum, err := strconv.ParseUint(id, 10, 64)
							if err != nil {
								panic(fmt.Errorf("invalid index: %s", id))
							}
							indices[i] = idNum
						default:
							panic(fmt.Errorf("invalid index type: %T", id))
						}
					}
					val[k] = indices
				} else if k == "MessageIDs" {
					// Special handling for MessageIDs array
					messageIDs := make([][56]byte, len(vv))
					for i, msgID := range vv {
						switch msgID := msgID.(type) {
						case string:
							hexStr := msgID
							hexStr = strings.TrimPrefix(hexStr, "0x")
							bytes, err := hex.DecodeString(hexStr)
							if err != nil || len(bytes) != 56 {
								panic(fmt.Errorf("invalid MessageID: %s", msgID))
							}
							copy(messageIDs[i][:], bytes)
						default:
							panic(fmt.Errorf("invalid MessageID type: %T", msgID))
						}
					}
					val[k] = messageIDs
				} else {
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
				}
			case map[string]interface{}:
				ConvertHexToBytes(vv)
			}
		}
	case []interface{}:
		for i, v := range val {
			switch vv := v.(type) {
			case string:
				// Try to decode as hex first
				if bytes, err := hex.DecodeString(vv); err == nil {
					val[i] = bytes
				} else {
					// If hex decoding fails, try base64
					if bytes, err := base64.StdEncoding.DecodeString(vv); err == nil {
						val[i] = bytes
					} else {
						// If both hex and base64 decoding fail, keep as is
						val[i] = vv
					}
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

// UnmarshalJSONWithHex handles unmarshaling of JSON with hex strings into structs with byte arrays
func UnmarshalJSONWithHex(data []byte, v interface{}) error {
	// Check if the type has a custom UnmarshalJSON method
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rt := rv.Type().Elem()
		if _, hasCustom := rt.MethodByName("UnmarshalJSON"); hasCustom {
			// If it has a custom UnmarshalJSON method, use that instead
			return json.Unmarshal(data, v)
		}
	}

	// First unmarshal into a map
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	// Convert hex strings to bytes
	ConvertHexToBytes(m)

	// Marshal back to JSON
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Unmarshal into the target struct
	return json.Unmarshal(jsonBytes, v)
}
