package hexencoding

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/types"
)

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

	// Check if the JSON is an array or object
	var firstChar byte
	for i := 0; i < len(data); i++ {
		if data[i] != ' ' && data[i] != '\n' && data[i] != '\r' && data[i] != '\t' {
			firstChar = data[i]
			break
		}
	}

	if firstChar == '[' {
		// Handle array - unmarshal into generic slice first
		var arr []interface{}
		if err := json.Unmarshal(data, &arr); err != nil {
			return err
		}

		// Convert hex strings to bytes in the array
		ConvertHexToBytes(arr)

		// Marshal back to JSON
		jsonBytes, err := json.Marshal(arr)
		if err != nil {
			return err
		}

		// Unmarshal into the target struct
		return json.Unmarshal(jsonBytes, v)
	} else {
		// Handle object - unmarshal into generic map first
		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}

		// Handle Root fields in Source and Target maps
		if source, ok := m["Source"].(map[string]interface{}); ok {
			if root, ok := source["Root"].(string); ok {
				if !strings.HasPrefix(root, "0x") {
					source["Root"] = "0x" + root
				}
			}
		}
		if target, ok := m["Target"].(map[string]interface{}); ok {
			if root, ok := target["Root"].(string); ok {
				if !strings.HasPrefix(root, "0x") {
					target["Root"] = "0x" + root
				}
			}
		}

		// Convert hex strings to bytes
		ConvertHexToBytes(m)

		// marshal back to json
		jsonBytes, err := json.Marshal(m)
		if err != nil {
			return err
		}

		return json.Unmarshal(jsonBytes, v)
	}
}

// UnmarshalJSONWithHex handles unmarshaling of JSON with hex strings into structs with byte arrays
func SSVUnmarshalJSONWithHex(data []byte, v interface{}) error {
	// Check if the type has a custom UnmarshalJSON method
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rt := rv.Type().Elem()
		if _, hasCustom := rt.MethodByName("UnmarshalJSON"); hasCustom {
			// If it has a custom UnmarshalJSON method, use that instead
			return json.Unmarshal(data, v)
		}
	}

	// Check if the JSON is an array or object
	var firstChar byte
	for i := 0; i < len(data); i++ {
		if data[i] != ' ' && data[i] != '\n' && data[i] != '\r' && data[i] != '\t' {
			firstChar = data[i]
			break
		}
	}

	if firstChar == '[' {
		// Handle array - unmarshal into generic slice first
		var arr []interface{}
		if err := json.Unmarshal(data, &arr); err != nil {
			return err
		}

		// Convert hex strings to bytes in the array
		ConvertHexToBytes(arr)

		// Marshal back to JSON
		jsonBytes, err := json.Marshal(arr)
		if err != nil {
			return err
		}

		// Unmarshal into the target struct
		return json.Unmarshal(jsonBytes, v)
	} else {
		// Handle object - unmarshal into generic map first
		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}

		// Handle Root fields in Source and Target maps
		if source, ok := m["Source"].(map[string]interface{}); ok {
			if root, ok := source["Root"].(string); ok {
				if !strings.HasPrefix(root, "0x") {
					source["Root"] = "0x" + root
				}
			}
		}
		if target, ok := m["Target"].(map[string]interface{}); ok {
			if root, ok := target["Root"].(string); ok {
				if !strings.HasPrefix(root, "0x") {
					target["Root"] = "0x" + root
				}
			}
		}

		// Convert hex strings to bytes
		ConvertHexToBytes(m)

		// Use direct reflection assignment instead of marshaling back to JSON
		return assignMapToStruct(m, v)
	}
}

// assignMapToStruct uses reflection to directly assign values from a map to a struct
func assignMapToStruct(m map[string]interface{}, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a struct")
	}

	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)

		// Get the JSON field name
		jsonTag := fieldType.Tag.Get("json")
		fieldName := fieldType.Name
		if jsonTag != "" && jsonTag != "-" {
			// Split by comma to handle options like "name,omitempty"
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				fieldName = parts[0]
			}
		}

		// Skip if field is not in the map
		mapValue, exists := m[fieldName]
		if !exists {
			continue
		}

		// Assign the value from map to struct field
		if err := assignValueToField(field, mapValue); err != nil {
			// Special handling for qbft.Instance.forceStop (unexported)
			if rt.PkgPath() == "github.com/ssvlabs/ssv-spec/qbft" && rt.Name() == "Instance" && fieldType.Name == "forceStop" {
				if b, ok := mapValue.(bool); ok {
					setUnexportedBoolField(rv, "forceStop", b)
					continue
				}
			}
			return fmt.Errorf("failed to assign field %s: %v", fieldName, err)
		}
	}

	// Also handle forceStop if present as ForceStop or forceStop in the map (for custom JSON)
	if rt.PkgPath() == "github.com/ssvlabs/ssv-spec/qbft" && rt.Name() == "Instance" {
		if vMap, ok := m["forceStop"]; ok {
			if b, ok := vMap.(bool); ok {
				setUnexportedBoolField(rv, "forceStop", b)
			}
		} else if vMap, ok := m["ForceStop"]; ok {
			if b, ok := vMap.(bool); ok {
				setUnexportedBoolField(rv, "forceStop", b)
			}
		}
	}

	return nil
}

// assignValueToField assigns a value from the map to a struct field using reflection
func assignValueToField(field reflect.Value, value interface{}) error {
	if !field.CanSet() {
		// Instead of returning error, just skip for unexported fields (except for special handling above)
		return nil
	}

	// Handle nil values
	if value == nil {
		field.Set(reflect.Zero(field.Type()))
		return nil
	}

	valueType := reflect.TypeOf(value)
	fieldType := field.Type()

	// If types match exactly, assign directly
	if valueType == fieldType {
		field.Set(reflect.ValueOf(value))
		return nil
	}

	// Handle special cases for byte arrays
	switch fieldType.Kind() {
	case reflect.Array:
		if fieldType.Elem().Kind() == reflect.Uint8 {
			// Handle [N]byte arrays
			switch v := value.(type) {
			case []byte:
				if len(v) != fieldType.Len() {
					return fmt.Errorf("byte array length mismatch: expected %d, got %d", fieldType.Len(), len(v))
				}
				reflect.Copy(field, reflect.ValueOf(v))
				return nil
			case [32]byte:
				if fieldType.Len() == 32 {
					field.Set(reflect.ValueOf(v))
					return nil
				}
			case [48]byte:
				if fieldType.Len() == 48 {
					field.Set(reflect.ValueOf(v))
					return nil
				}
			case [56]byte:
				if fieldType.Len() == 56 {
					field.Set(reflect.ValueOf(v))
					return nil
				}
			case [96]byte:
				if fieldType.Len() == 96 {
					field.Set(reflect.ValueOf(v))
					return nil
				}
			case [4]byte:
				if fieldType.Len() == 4 {
					field.Set(reflect.ValueOf(v))
					return nil
				}
			case [20]byte:
				if fieldType.Len() == 20 {
					field.Set(reflect.ValueOf(v))
					return nil
				}
			}
		}
	case reflect.Slice:
		if fieldType.Elem().Kind() == reflect.Uint8 {
			// Handle []byte slices
			if v, ok := value.([]byte); ok {
				field.Set(reflect.ValueOf(v))
				return nil
			}
			// Handle fixed-size byte arrays converting to slices
			if v, ok := value.([48]byte); ok {
				field.Set(reflect.ValueOf(v[:]))
				return nil
			}
			if v, ok := value.([32]byte); ok {
				field.Set(reflect.ValueOf(v[:]))
				return nil
			}
			if v, ok := value.([96]byte); ok {
				field.Set(reflect.ValueOf(v[:]))
				return nil
			}
			if v, ok := value.([4]byte); ok {
				field.Set(reflect.ValueOf(v[:]))
				return nil
			}
			if v, ok := value.([20]byte); ok {
				field.Set(reflect.ValueOf(v[:]))
				return nil
			}
			if v, ok := value.([56]byte); ok {
				field.Set(reflect.ValueOf(v[:]))
				return nil
			}
		}
	case reflect.Ptr:
		// Handle pointer types
		if value == nil {
			field.Set(reflect.Zero(fieldType))
			return nil
		}

		// Create a new instance of the pointed-to type
		if field.IsNil() {
			field.Set(reflect.New(fieldType.Elem()))
		}

		// Recursively assign to the pointed-to value
		return assignValueToField(field.Elem(), value)
	case reflect.Struct:
		// Handle nested structs
		if v, ok := value.(map[string]interface{}); ok {
			return assignMapToStruct(v, field.Addr().Interface())
		}
	case reflect.Map:
		// Handle map types recursively
		// Special handling for SlashableSlots map
		if fieldType.String() == "map[string][]phase0.Slot" {
			if v, ok := value.(map[string][]uint64); ok {
				slashableSlots := make(map[string][]phase0.Slot)
				for key, slotValues := range v {
					slots := make([]phase0.Slot, len(slotValues))
					for i, slotValue := range slotValues {
						slots[i] = phase0.Slot(slotValue)
					}
					slashableSlots[key] = slots
				}
				field.Set(reflect.ValueOf(slashableSlots))
				return nil
			}
		}

		if v, ok := value.(map[string]interface{}); ok {
			// Create a new map of the correct type
			mapType := field.Type()
			newMap := reflect.MakeMap(mapType)

			// Recursively process each key-value pair
			for key, val := range v {
				// Convert key to the correct type
				keyType := mapType.Key()
				newKey := reflect.New(keyType).Elem()

				// Use the same recursive logic for key conversion
				if err := assignValueToField(newKey, key); err != nil {
					return fmt.Errorf("failed to convert map key %v: %v", key, err)
				}

				// Create a new value of the correct type
				valueType := mapType.Elem()
				newValue := reflect.New(valueType).Elem()

				// Recursively assign the value
				if err := assignValueToField(newValue, val); err != nil {
					return fmt.Errorf("failed to assign map value for key %s: %v", key, err)
				}

				// Set the key-value pair in the map
				newMap.SetMapIndex(newKey, newValue)
			}

			field.Set(newMap)
			return nil
		}
	case reflect.Interface:
		// Handle interface types
		if v, ok := value.(map[string]interface{}); ok {
			// Special handling for types.Duty interface
			if fieldType.String() == "types.Duty" {
				// Check if it's a ValidatorDuty by looking for Type field
				if _, hasType := v["Type"]; hasType {
					// Create a new ValidatorDuty
					validatorDuty := &types.ValidatorDuty{}
					if err := assignMapToStruct(v, validatorDuty); err != nil {
						return fmt.Errorf("failed to convert to ValidatorDuty: %v", err)
					}
					field.Set(reflect.ValueOf(validatorDuty))
					return nil
				} else {
					// Check if it's a CommitteeDuty by looking for ValidatorDuties field
					if _, hasValidatorDuties := v["ValidatorDuties"]; hasValidatorDuties {
						// Create a new CommitteeDuty
						committeeDuty := &types.CommitteeDuty{}
						if err := assignMapToStruct(v, committeeDuty); err != nil {
							return fmt.Errorf("failed to convert to CommitteeDuty: %v", err)
						}
						field.Set(reflect.ValueOf(committeeDuty))
						return nil
					}
				}
			}
		}
	case reflect.String:
		// Handle string types
		if v, ok := value.(string); ok {
			field.SetString(v)
			return nil
		}
		// Handle fixed-size byte arrays converting to hex strings
		if v, ok := value.([32]byte); ok {
			field.SetString(hex.EncodeToString(v[:]))
			return nil
		}
		if v, ok := value.([48]byte); ok {
			field.SetString("0x" + hex.EncodeToString(v[:]))
			return nil
		}
		if v, ok := value.([96]byte); ok {
			field.SetString("0x" + hex.EncodeToString(v[:]))
			return nil
		}
		if v, ok := value.([4]byte); ok {
			field.SetString(hex.EncodeToString(v[:]))
			return nil
		}
		if v, ok := value.([20]byte); ok {
			field.SetString("0x" + hex.EncodeToString(v[:]))
			return nil
		}
		if v, ok := value.([56]byte); ok {
			field.SetString(hex.EncodeToString(v[:]))
			return nil
		}
	}

	// For other types, try to convert
	valueReflect := reflect.ValueOf(value)
	if valueReflect.Type().ConvertibleTo(fieldType) {
		field.Set(valueReflect.Convert(fieldType))
		return nil
	}

	// Handle special type conversions
	switch fieldType.String() {
	case "phase0.ValidatorIndex", "phase0.Slot", "qbft.Round", "qbft.Height":
		// Convert string to ValidatorIndex
		if str, ok := value.(string); ok {
			if idx, err := strconv.ParseUint(str, 10, 64); err == nil {
				// Create a new ValidatorIndex value
				validatorIndex := reflect.New(fieldType).Elem()
				validatorIndex.SetUint(idx)
				field.Set(validatorIndex)
				return nil
			}
		}
		// Try direct uint64 conversion
		if idx, ok := value.(uint64); ok {
			validatorIndex := reflect.New(fieldType).Elem()
			validatorIndex.SetUint(idx)
			field.Set(validatorIndex)
			return nil
		}
		// Try float64 conversion (from JSON numbers)
		if f, ok := value.(float64); ok {
			validatorIndex := reflect.New(fieldType).Elem()
			validatorIndex.SetUint(uint64(f))
			field.Set(validatorIndex)
			return nil
		}
		// Handle []uint8 (byte array) conversion
		if bytes, ok := value.([]uint8); ok {
			if len(bytes) <= 8 {
				// Convert 1-8 bytes to uint64 (little-endian, as SSZ uses)
				var idx uint64
				for i := 0; i < len(bytes); i++ {
					idx |= uint64(bytes[i]) << (i * 8)
				}
				validatorIndex := reflect.New(fieldType).Elem()
				validatorIndex.SetUint(idx)
				field.Set(validatorIndex)
				return nil
			}
		}
		// Handle []byte conversion
		if bytes, ok := value.([]byte); ok {
			if len(bytes) <= 8 {
				// Convert 1-8 bytes to uint64 (little-endian, as SSZ uses)
				var idx uint64
				for i := 0; i < len(bytes); i++ {
					idx |= uint64(bytes[i]) << (i * 8)
				}
				validatorIndex := reflect.New(fieldType).Elem()
				validatorIndex.SetUint(idx)
				field.Set(validatorIndex)
				return nil
			}
		}
	case "uint64":
		// Handle string to uint64 conversion for map keys
		if str, ok := value.(string); ok {
			if idx, err := strconv.ParseUint(str, 10, 64); err == nil {
				field.SetUint(idx)
				return nil
			}
		}
		// Try direct uint64 conversion
		if idx, ok := value.(uint64); ok {
			field.SetUint(idx)
			return nil
		}
		// Try float64 conversion (from JSON numbers)
		if f, ok := value.(float64); ok {
			field.SetUint(uint64(f))
			return nil
		}
		// Handle []uint8 (byte array) conversion
		if bytes, ok := value.([]uint8); ok {
			if len(bytes) <= 8 {
				// Convert 1-8 bytes to uint64 (little-endian, as SSZ uses)
				var idx uint64
				for i := 0; i < len(bytes); i++ {
					idx |= uint64(bytes[i]) << (i * 8)
				}
				field.SetUint(idx)
				return nil
			}
		}
		// Handle []byte conversion
		if bytes, ok := value.([]byte); ok {
			if len(bytes) <= 8 {
				// Convert 1-8 bytes to uint64 (little-endian, as SSZ uses)
				var idx uint64
				for i := 0; i < len(bytes); i++ {
					idx |= uint64(bytes[i]) << (i * 8)
				}
				field.SetUint(idx)
				return nil
			}
		}
	case "phase0.BLSPubKey":
		// Handle string to BLSPubKey conversion
		if str, ok := value.(string); ok {
			// Remove 0x prefix if present
			hexStr := strings.TrimPrefix(str, "0x")
			bytes, err := hex.DecodeString(hexStr)
			if err != nil || len(bytes) != 48 {
				return fmt.Errorf("invalid BLS public key: %s", str)
			}
			var pubKey phase0.BLSPubKey
			copy(pubKey[:], bytes)
			field.Set(reflect.ValueOf(pubKey))
			return nil
		}
		// Try direct [48]byte conversion
		if pubKey, ok := value.([48]byte); ok {
			field.Set(reflect.ValueOf(phase0.BLSPubKey(pubKey)))
			return nil
		}
	case "spec.DataVersion":
		// Handle string to DataVersion conversion
		if str, ok := value.(string); ok {
			if idx, err := strconv.ParseUint(str, 10, 64); err == nil {
				field.SetUint(idx)
				return nil
			}
			return nil
		}
		if num, ok := value.(int); ok {
			field.Set(reflect.ValueOf(spec.DataVersion(num)))
			return nil
		}
	case "types.RunnerRole":
		// Handle string to RunnerRole conversion
		if str, ok := value.(string); ok {
			// Try to parse as string representation first
			if role, found := stringToRunnerRole(str); found {
				runnerRole := reflect.New(fieldType).Elem()
				runnerRole.SetInt(int64(role))
				field.Set(runnerRole)
				return nil
			}
			// Try to parse as numeric string
			if idx, err := strconv.ParseInt(str, 10, 32); err == nil {
				runnerRole := reflect.New(fieldType).Elem()
				runnerRole.SetInt(idx)
				field.Set(runnerRole)
				return nil
			}
		}
		// Try direct int32 conversion
		if idx, ok := value.(int32); ok {
			runnerRole := reflect.New(fieldType).Elem()
			runnerRole.SetInt(int64(idx))
			field.Set(runnerRole)
			return nil
		}
		// Try float64 conversion (from JSON numbers)
		if f, ok := value.(float64); ok {
			runnerRole := reflect.New(fieldType).Elem()
			runnerRole.SetInt(int64(f))
			field.Set(runnerRole)
			return nil
		}
	}

	// Handle other slices
	if v, ok := value.([]interface{}); ok {
		slice := reflect.MakeSlice(fieldType, len(v), len(v))
		for i, elem := range v {
			if err := assignValueToField(slice.Index(i), elem); err != nil {
				return fmt.Errorf("failed to assign slice element %d: %v", i, err)
			}
		}
		field.Set(slice)
		return nil
	}

	return fmt.Errorf("cannot assign %T to %s", value, fieldType)
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
					// If the key ends with "Version", it's a fixed-size array
					if strings.HasSuffix(k, "Version") {
						var version [4]byte
						copy(version[:], bytes)
						val[k] = version
					} else if strings.HasSuffix(k, "Root") {
						if k == "ExpectedSigningRoot" || k == "ControllerPostRoot" || k == "PostRoot" {
							val[k] = vv
						} else if k == "BlockRoot" {
							// Add 0x prefix for BlockRoot
							if !strings.HasPrefix(vv, "0x") {
								vv = "0x" + vv
							}
							val[k] = vv
						} else {
							// For other roots, remove 0x prefix if present
							hexStr := vv
							if strings.HasPrefix(hexStr, "0x") {
								hexStr = strings.TrimPrefix(hexStr, "0x")
							}
							bytes, err = hex.DecodeString(hexStr)
							if err != nil {
								// If still not 32 bytes, raise error
								panic(fmt.Errorf("invalid root: %s", vv))
							}
							if len(bytes) == 32 {
								var root [32]byte
								copy(root[:], bytes)
								val[k] = root
							} else {
								val[k] = bytes
							}
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
					} else if k == "Value" {
						// Special handling for Value field (32-byte array)
						var value [32]byte
						hexStr := vv
						hexStr = strings.TrimPrefix(hexStr, "0x")
						bytes, err = hex.DecodeString(hexStr)
						if err != nil || len(bytes) != 32 {
							// If still not 32 bytes, raise error
							panic(fmt.Errorf("invalid Value: %s", vv))
						}
						copy(value[:], bytes)
						val[k] = value
					} else if k == "FeeRecipientAddress" {
						// Special handling for FeeRecipient
						var feeRecipient [20]byte
						hexStr := vv
						hexStr = strings.TrimPrefix(hexStr, "0x")
						bytes, err = hex.DecodeString(hexStr)
						if err != nil || len(bytes) != 20 {
							panic(fmt.Errorf("invalid FeeRecipient: %s", vv))
						}
						copy(feeRecipient[:], bytes)
						val[k] = feeRecipient
					} else if k == "ValidatorPubKey" {
						// Special handling for ValidatorPubKey (48-byte BLS public key)
						var validatorPubKey [48]byte
						hexStr := vv
						hexStr = strings.TrimPrefix(hexStr, "0x")
						bytes, err = hex.DecodeString(hexStr)
						if err != nil || len(bytes) != 48 {
							panic(fmt.Errorf("invalid ValidatorPubKey: %s", vv))
						}
						copy(validatorPubKey[:], bytes)
						val[k] = validatorPubKey
					} else if k == "WithdrawalCredentials" {
						// Keep WithdrawalCredentials as a string
						val[k] = vv
					} else if k == "Slot" || k == "DutySlot" || k == "ValidatorPK" || k == "ValidatorIndex" {
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
			case float64:
				val[k] = vv
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
				} else if k == "Heights" {
					// Special handling for Heights and Rounds arrays
					values := make([]uint64, len(vv))
					for i, value := range vv {
						switch value := value.(type) {
						case float64:
							values[i] = uint64(value)
						case string:
							valueNum, err := strconv.ParseUint(value, 10, 64)
							if err != nil {
								panic(fmt.Errorf("invalid height/round: %s", value))
							}
							values[i] = valueNum
						default:
							panic(fmt.Errorf("invalid height/round type: %T", value))
						}
					}
					val[k] = values
				} else if k == "Rounds" {
					// Special handling for Rounds array
					rounds := make([]uint64, len(vv))
					for i, round := range vv {
						switch round := round.(type) {
						case float64:
							rounds[i] = uint64(round)
						case string:
							roundNum, err := strconv.ParseUint(round, 10, 64)
							if err != nil {
								panic(fmt.Errorf("invalid round: %s", round))
							}
							rounds[i] = roundNum
						default:
							panic(fmt.Errorf("invalid round type: %T", round))
						}
					}
					val[k] = rounds
				} else if k == "Proposers" {
					// Special handling for Proposers array
					proposers := make([]uint64, len(vv))
					for i, proposer := range vv {
						switch proposer := proposer.(type) {
						case float64:
							proposers[i] = uint64(proposer)
						case string:
							proposerNum, err := strconv.ParseUint(proposer, 10, 64)
							if err != nil {
								panic(fmt.Errorf("invalid proposer: %s", proposer))
							}
							proposers[i] = proposerNum
						default:
							panic(fmt.Errorf("invalid proposer type: %T", proposer))
						}
					}
					val[k] = proposers
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
				} else if k == "ExpectedRoots" {
					// Special handling for ExpectedRoots array
					expectedRoots := make([][32]byte, len(vv))
					for i, root := range vv {
						switch root := root.(type) {
						case string:
							hexStr := root
							hexStr = strings.TrimPrefix(hexStr, "0x")
							bytes, err := hex.DecodeString(hexStr)
							if err != nil || len(bytes) != 32 {
								panic(fmt.Errorf("invalid ExpectedRoot: %s", root))
							}
							copy(expectedRoots[i][:], bytes)
						default:
							panic(fmt.Errorf("invalid ExpectedRoot type: %T", root))
						}
					}
					val[k] = expectedRoots
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
				// Special handling for SlashableSlots map
				if k == "SlashableSlots" {
					slashableSlots := make(map[string][]uint64)
					for slotKey, slotValues := range vv {
						switch slotValues := slotValues.(type) {
						case []interface{}:
							slots := make([]uint64, len(slotValues))
							for i, slotValue := range slotValues {
								switch slotValue := slotValue.(type) {
								case string:
									// Convert string to uint64 (treat as decimal, not hex)
									slotNum, err := strconv.ParseUint(slotValue, 10, 64)
									if err != nil {
										panic(fmt.Errorf("invalid slot value: %s", slotValue))
									}
									slots[i] = slotNum
								case float64:
									slots[i] = uint64(slotValue)
								default:
									panic(fmt.Errorf("invalid slot value type: %T", slotValue))
								}
							}
							slashableSlots[slotKey] = slots
						default:
							panic(fmt.Errorf("invalid slot values type: %T", slotValues))
						}
					}
					val[k] = slashableSlots
				} else {
					ConvertHexToBytes(vv)
				}
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

// Helper to convert RunnerRole string to enum value
func stringToRunnerRole(s string) (types.RunnerRole, bool) {
	switch s {
	case "COMMITTEE_RUNNER":
		return types.RoleCommittee, true
	case "AGGREGATOR_RUNNER":
		return types.RoleAggregator, true
	case "PROPOSER_RUNNER":
		return types.RoleProposer, true
	case "SYNC_COMMITTEE_CONTRIBUTION_RUNNER":
		return types.RoleSyncCommitteeContribution, true
	case "VALIDATOR_REGISTRATION_RUNNER":
		return types.RoleValidatorRegistration, true
	case "VOLUNTARY_EXIT_RUNNER":
		return types.RoleVoluntaryExit, true
	default:
		return types.RoleUnknown, false
	}
}

// Helper to set unexported bool field using unsafe
func setUnexportedBoolField(structVal reflect.Value, fieldName string, value bool) {
	field := structVal.FieldByName(fieldName)
	if !field.IsValid() {
		return
	}
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().SetBool(value)
}
