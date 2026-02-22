package engine

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
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
		// "referenceValue", // TODO
		// "geoPointValue",  // TODO
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
				return nil, fmt.Errorf("can't cast the 'fields' attribute under the 'mapValue' object to a map type. Path - %s", path)
			}
			path += "/fields"
			break
		}
	}

	if !fieldsFound {
		return nil, fmt.Errorf("'mapValue' object does not contain an obligatory field - 'fields'. Path - %s", path)
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

func handleGoMap(payloadVal interface{}, path string) (map[string]map[string]interface{} , error) {
	firestoreMapObject := map[string]map[string]interface{}{"mapValue": {"fields": map[string]interface{}{}}}

	mapVal, ok := payloadVal.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't cast the value under the following path - %s to a map.", path)
	}

	for k, v := range mapVal {
		if slices.Contains(supportedFields, k) {
			return nil, fmt.Errorf("Object under the path -> %s, contains the key -> %s, which is the Firestore type", path, k)
		}
		processedVal, err := handleGoType(v, path + "/" + k)
		if err != nil {
			return  nil, err
		}
		mapVal[k] = processedVal
	}
	
	firestoreMapObject["mapValue"]["fields"] = mapVal
	return firestoreMapObject, nil
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

func handleGoArray(payloadVal interface{}, path string) (interface{}, error) {

	firestoreArrayObject := map[string]map[string]interface{} {"arrayValue":{"values": []interface{}{}}}

	payloadArr, ok := payloadVal.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Can't cast a value under the following path in the payload -> %s to an array.", path)
	}

	for i, elem := range payloadArr {
		processedElem, err := handleGoType(elem, path +  fmt.Sprintf("[%d]", i))
		if err != nil {
			return nil, err
		}
		payloadArr[i] = processedElem
	}

	firestoreArrayObject["arrayValue"]["values"] = payloadArr
	return firestoreArrayObject, nil	
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

	if !strings.HasSuffix(strValue, "=") {
		return "", fmt.Errorf(
			"The following string -> %s should be padded, in order to be considered as a valid base64 encoded byte value",
			strValue,
		)
	}

	_, err := base64.StdEncoding.Strict().DecodeString(strValue)
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

func handleGoSingularType(val interface{}, firestoreType string) map[string]interface{} {
	firestoreObject := make(map[string]interface{})
	if firestoreType == "integerValue" || firestoreType == "doubleValue" {
		firestoreObject = map[string]interface{}{firestoreType:strconv.FormatFloat(val.(float64), 'f', -1, 64)}
	} else {
		firestoreObject = map[string]interface{}{firestoreType:val}
	}
	return firestoreObject
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
		// Check error handling
		path += "/arrayValue"
		val, err := handleArrayValue(typeVal, path)
		if err == nil {
			return val, nil
		}
		// return nil, errors.New(generateErrorMessage(path, typeKey, err.Error()))
		return nil, err
	case "mapValue":
		// Check error handling
		path += "/mapValue"
		val, err := handleMapValue(typeVal, path)
		if err == nil {
			return val, nil
		}
		// return nil, errors.New(generateErrorMessage(path, typeKey, err.Error()))
		return nil, err
	}

	return nil, fmt.Errorf("Unsupported firestore field type - %s. Path -> %s", typeKey, path)
}

func handleGoType(payloadVal interface{}, path string) (interface{}, error) {
	var generalErr error = nil
	switch t := payloadVal.(type) {
		case string:
			// Check if byte
			if _, err := handleByteValue(payloadVal); err == nil {
				return handleGoSingularType(payloadVal, "byteValue"), nil
			}
			// Check if timestamp
			if _, err := handleTimestampValue(payloadVal); err == nil {
				return handleGoSingularType(payloadVal, "timestampValue"), nil
			}
			return handleGoSingularType(payloadVal, "stringValue"), nil
		case nil:
			return handleGoSingularType(payloadVal, "nullValue"), nil
		case bool:
			return handleGoSingularType(payloadVal, "booleanValue"), nil
		case float64:
			if math.Mod(payloadVal.(float64), 1) == 0 {
				return handleGoSingularType(payloadVal, "integerValue"), nil
			}
			return handleGoSingularType(payloadVal, "doubleValue"), nil
		case []interface{}:
			firestoreArrayObject, err := handleGoArray(payloadVal, path)
			if err != nil {
				return nil, err
			}
			return firestoreArrayObject, nil 
		case map[string]interface{}:
			firestoreMapObject, err := handleGoMap(payloadVal, path)
			if err != nil {
				return nil, err
			}
			return firestoreMapObject, nil
		default:
			generalErr = fmt.Errorf("the following type - %s, which was found under the path - %s is not supported!", t, path)
	}

	return nil, generalErr
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

func EncodeToFirestore(payload map[string]interface{}) (map[string]interface{}, error) {
	encodedPayload := make(map[string]interface{})
	for k, v := range payload {
		encodedVal, err := handleGoType(v, k)
		if err != nil {
			return nil, err
		}
		encodedPayload[k] = encodedVal
	}
	return encodedPayload, nil
}