package decoder

import (
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"
)

var (
	supportedFields = []string {
		"nullValue",
		"booleanValue",
		"integerValue",
		"doubleValue",
		"timestampValue",
		"stringValue",
		"bytesValue",
		// "referenceValue",
		// "geoPointValue",
		"arrayValue",
		"mapValue",
	}
)

func generateErrorMessage(path string, typeKey string, reason string) string {
	return fmt.Sprintf(
		`Structure under the path -> %s with the type -> %s is invalid! Reason - %s`, 
		path, typeKey, reason,
	)
}

func handleMapValue(value interface{}, path string) (map[string]interface{}, error) {
	resMap := make(map[string]interface{}) 
	mapStructure, ok := value.(map[string]interface{})
	if !ok {
		return nil, errors.New("can't cast an object under the 'mapValue' to the 'map' type")
	}
	
	fieldsFound := false
	var fieldsMap map[string]interface{}
	for k, v := range mapStructure {
		if k == "fields" {
			fieldsFound = true
			fieldsMap, ok = v.(map[string]interface{})
			if !ok {
				return nil, errors.New(
					`can't cast the 'fields' attribute under the 'mapValue' object to a map type`,
				)
			}
			path += "/fields"
			break
		}
	}

	if !fieldsFound {
		return nil, errors.New("'mapValue' object does not contain an obligatory field - 'fields'")
	}
	
	for k, v := range fieldsMap {
		fieldValMap, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("can't cast the value under path - %s to a map", path+fmt.Sprintf("/%s", k))
		}
		mapVal, err := handleFirestoreType(fieldValMap, path + fmt.Sprintf("/%s", k))
		if err != nil {
			return nil, err
		}
		resMap[k] = mapVal
	}
	
	return resMap, nil
}

func handleArrayValue(value interface{}, path string) ([]interface{}, error) {
	resArr := []interface{} {}
	arrayMap, ok := value.(map[string]interface{})
	if !ok {
		return nil, errors.New("can't cast the value provided for the arrayValue to the map")
	}

	valuesFound := false
	var valuesArray []interface{}
	for k, v := range arrayMap {
		if k == "values" {
			valuesFound = true
			valuesArray, ok = v.([]interface{})
			if !ok {
				return nil, errors.New(
					`can't cast the 'values' attribute under the 'arrayValue' object to an array`,
				)
			}
			path += "/values"
			break
		}
	}

	if !valuesFound {
		return nil, errors.New("'arrayValue' object does not contain an obligatory field - 'values'")
	}

	for i, v := range valuesArray {
		arrValMap, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("can't cast the array val under path - %s to a map", path+fmt.Sprintf("[%d]", i))
		}
		arrVal, err := handleFirestoreType(arrValMap, path + fmt.Sprintf("[%d]", i))
		if err != nil {
			return nil, err
		}
		resArr = append(resArr, arrVal)
	}

	return resArr, nil
}

func handleTimestampValue(value interface{}) (string, error) {
	strTime, ok := value.(string)
	if !ok {
		return "", errors.New("timestamp value is not in a string format.")
	}
	_, err := time.Parse(time.RFC3339, strTime)
	if err != nil {
		return "", err
	}
	return strTime, nil
}

func handleByteValue(value interface{}) (string, error) {
	strValue, ok := value.(string)
	if !ok {
		return "", errors.New("byte value is not in a string format.")
	}

	_, err := base64.StdEncoding.DecodeString(strValue)
	if err != nil {
		return "", err
	}

	return strValue, nil
}

func handleSingularType[T any](value interface{}) (T, error) {
	covertedVal, ok := value.(T)
	if !ok {
		return covertedVal, errors.New("type assertion failed")
	} 
	return covertedVal, nil
}

func handleIntType(value interface{}) (int, error) {
	strInt, ok := value.(string)
	if !ok {
		return 0, errors.New("integer value is supposed to be in a form of a string")
	}

	// Convert to int
	intVal, err := strconv.Atoi(strInt)
	if err != nil {
		return 0, errors.New("string value provided for the integer type is not a number")
	}

	return intVal, nil
}


func handleFirestoreType(childPayload map[string]interface{}, path string) (interface{}, error) {

	if len(childPayload) > 1 {
		return nil, fmt.Errorf("Structure under the path -> %s is invalid! It contains more than one type key.", path)
	}
	
	if len(childPayload) == 0 {
		return nil, fmt.Errorf("Structure under the path -> %s is invalid! It contains no keys.", path)
	}

	var typeKey string
	var typeVal interface{}
	for k, v := range childPayload {
		typeKey = k
		typeVal = v
	}

	if !slices.Contains(supportedFields, typeKey) {
		return nil, fmt.Errorf("Structure under the path -> %s is invalid! It contains an invalid type -> %s", path, typeKey)
	}

	switch typeKey {
	case "nullValue":
		path += "/nullValue"
		if typeVal != nil {
			return nil, errors.New(generateErrorMessage(path, typeKey, "Value is not null"))
		}
		return nil, nil
	case "booleanValue":
		path += "/booleanValue"
		if val, err := handleSingularType[bool](typeVal); err == nil {
			return val, nil
		}
		return nil, errors.New(generateErrorMessage(path, typeKey, "Value is not a boolean type."))
	case "integerValue":
		path += "/integerValue"
		val, err := handleIntType(typeVal)
		if err == nil {
			return val, nil
		}
		return nil, errors.New(generateErrorMessage(path, typeKey, err.Error()))
	case "doubleValue":
		path += "/doubleValue"
		if val, err := handleSingularType[float64](typeVal); err == nil {
			return val, nil
		}
		return nil, errors.New(generateErrorMessage(path, typeKey, "Value is not a double type."))
	case "stringValue":
		path += "/stringValue"
		if val, err := handleSingularType[string](typeVal); err == nil {
			return val, nil
		}
		return nil, errors.New(generateErrorMessage(path, typeKey, "Value is not a string type."))
	case "byteValue":
		path += "/byteValue"
		val, err := handleByteValue(typeVal)
		if err == nil {
			return val, nil
		}
		return nil, errors.New(generateErrorMessage(path, typeKey, err.Error()))
	case "timestampValue":
		path += "/timestampValue"
		val, err := handleTimestampValue(typeVal)
		if err == nil {
			return val, nil
		}
		return nil, errors.New(generateErrorMessage(path, typeKey, err.Error()))
	case "arrayValue":
		path += "/arrayValue"
		val, err := handleArrayValue(typeVal, path)
		if err == nil {
			return val, nil
		}
		return nil, errors.New(generateErrorMessage(path, typeKey, err.Error()))
	case "mapValue":
		path += "/mapValue"
		val, err := handleMapValue(typeVal, path)
		if err == nil {
			return val, nil
		}
		return nil, errors.New(generateErrorMessage(path, typeKey, err.Error()))
	}

	return nil, fmt.Errorf("Unsupported firestore field type - %s", typeKey)
}


func DecodeFromFirestore(payload map[string]interface{}) (map[string]interface{}, error) {
	resPayload := make(map[string]interface{})
	for k, v := range payload {
		valMap, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Can't cast an object under the following key - %s to a map", k)
		}
		val, err := handleFirestoreType(valMap, k)
		if err != nil {
			return nil, err
		}
		resPayload[k] = val
	}
	return resPayload, nil
}