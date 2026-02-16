package engine

import (
	"errors"
	"fmt"
	"log/slog"
	"github.com/mvksxm/firestore-json-convert/engine/decoder"
)

type Processor struct {
	payload map[string]interface{}
}

func (prc *Processor) Convert() (map[string]interface{}, error) {
	decodedPayload, decodeErr := decoder.DecodeFromFirestore(prc.payload)
	if decodeErr != nil {
		slog.Warn(
			fmt.Sprintf(
				`Payload provided by the user can't be decoded from Firestore format, due to the following reason - %s.
				Proceeding with checking, whether payload is suitable encoding into the Firestore format.`,
			 	decodeErr.Error(),
			),
		)
	}
	if decodeErr == nil {
		return decodedPayload, nil
	}
	return nil, errors.New("Not implemented yet...")
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