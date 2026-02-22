package engine

import (
	"fmt"
	"log/slog"
)

type Processor struct {
	payload map[string]interface{}
}

func (prc *Processor) Convert() (map[string]interface{}, error) {
	decodedPayload, decodeErr := DecodeFromFirestore(prc.payload)
	if decodeErr == nil {
		return decodedPayload, nil
	}
	
	slog.Warn(
		fmt.Sprintf(
			"Payload provided by the user can't be decoded from Firestore format, due to the following reason - %s.",
			decodeErr.Error(),
		),
	)
	
	fmt.Println("Proceeding with checking, whether payload is suitable encoding into the Firestore format.")

	encodedPayload, encodeErr := EncodeToFirestore(prc.payload)
	if encodeErr == nil {
		return encodedPayload, nil
	}
	
	slog.Warn(
		fmt.Sprintf(
			`Payload provided by the user can't be encoded from Firestore format, due to the following reason - %s.`,
			encodeErr.Error(),
		),
	)

	return nil, fmt.Errorf(
		`Can't decode the payload from Firestore type due to the following error -> %s.
		At the same time, can't encode the payload to the Firestore type due to the following error -> %s`, 
		decodeErr, encodeErr,
	)
}

func (prc *Processor) ConvertToFirestore() {

}

func (prc *Processor) ConvertFromFirestore() {

}

func NewProcessor(payload map[string]interface{}) *Processor {
	return &Processor{
		payload: payload,
	}
}