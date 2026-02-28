package test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"runtime"
	"testing"

	"github.com/mvksxm/firestore-json-convert/engine"
)


type EngineFunction func() (interface{}, error)

func executeEngineTesting(
	t *testing.T,
	isEncode bool,
	resPayload map[string]interface{}, 
	testId string,
	functions []EngineFunction,
) {
	operation := "decoding"
	if isEncode {
		operation = "encoding"
	}

	resPayloadByte, err := json.Marshal(resPayload)
	if err != nil {
		t.Errorf(
			"Error occured, when marshalling res payload. Err: %s. (Test Id #%s)", 
			err.Error(),
			testId, 
		)
		return
	}

	for _, funcToTest := range functions {
		funcName :=  runtime.FuncForPC(reflect.ValueOf(funcToTest).Pointer()).Name()
		encodedPl, err := funcToTest()
		if err != nil {
			t.Errorf(
				"Error occured, when %s payload by using the method -> '%s'. Err: %s. (Test Id #%s)", 
				operation,
				funcName, 
				err.Error(),
				testId,
			)
			continue
		}

		if encodedPlByte, err := json.Marshal(encodedPl); err == nil {

			if !bytes.Equal(resPayloadByte, encodedPlByte) {
				t.Errorf(
					"Payload received, when %s the data by using the method -> '%s' is not equal to the intended result (Test Id #%s)", 
					operation,
					funcName, 
					testId,
				)
			}

		} else { 
			t.Errorf(
				"Error occured when marshaling the result returned from the method -> '%s'. Err -> %s. (Test Id #%s)", 
				funcName,
				err.Error(), 
				testId,
			)
		}
	}
}

func TestEncodeFirestore(t *testing.T) {

	isEncode := true
	testPayloadsMap := getPayloads(isEncode)
	
	for k, v := range testPayloadsMap {
		testPayload := v[0]
		resPayload := v[1]

		// Testing encoding methods
		processor := engine.NewProcessor(testPayload)
		executeEngineTesting(
			t,
			isEncode,
			resPayload,
			k,
			[]EngineFunction{processor.ConvertToFirestore},
		)
	}
}


func TestDecodeFirestore(t *testing.T) {
	isEncode := false
	testPayloadsMap := getPayloads(isEncode)
	
	for k, v := range testPayloadsMap {
		testPayload := v[0]
		resPayload := v[1]

		// Testing decoding methods
		processor := engine.NewProcessor(testPayload)
		executeEngineTesting(
			t,
			isEncode,
			resPayload,
			k,
			[]EngineFunction{processor.ConvertFromFirestore},
		)
	}

}