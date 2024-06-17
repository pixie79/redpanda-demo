package utils

import (
	"log/slog"
	"strings"

	goavro "github.com/linkedin/goavro/v2"
)

// decodeAvro decodes an Avro event using the provided schema and returns a nested map[string]interface{}.
//
// Parameters:
// - schema: The Avro schema used for decoding the event (string).
// - event: The Avro event to be decoded ([]byte).
//
// Returns:
// - nestedMap: The decoded event as a nested map[string]interface{}.
func DecodeAvro(schema string, event []byte) map[string]interface{} {
	sourceCodec, err := goavro.NewCodec(schema)
	if err != nil {
		slog.Error("Error creating Avro codec", "Error", err)
		return nil
	}

	eventStr := strings.Replace(string(event), "\"", "", -1)
	decodedEvent, err := b64DecodeMsg(eventStr, 5)
	if err != nil {
		slog.Error("Error decoding base64", "Error", err)
		return nil
	}

	native, _, err := sourceCodec.NativeFromBinary(decodedEvent)
	if err != nil {
		slog.Error("Error creating native from binary", "Error", err)
		return nil
	}

	nestedMap, ok := native.(map[string]interface{})
	if !ok {
		slog.Error("Unable to convert native to map[string]interface{}")
		return nil
	}

	return nestedMap
}

func WrapUnionSimple(value interface{}, typeName string) interface{} {
	return map[string]interface{}{typeName: value}
}
