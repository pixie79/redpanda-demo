package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	pUtils "pixie79/utils"
	"strconv"

	avro "github.com/linkedin/goavro/v2"
	transform "github.com/redpanda-data/redpanda/src/transform-sdk/go/transform"
	sr "github.com/redpanda-data/redpanda/src/transform-sdk/go/transform/sr"
)

func DecodeAvroRawEvent(e transform.WriteEvent) (map[string]interface{}, error) {
	rawEvent, err := json.Marshal(e.Record().Value)
	if err != nil {
		slog.Error("Unable to marshal raw event", "Error", e.Record().Value)
		return nil, err
	}

	// Sensitive Debug output
	// slog.Debug("Raw Event", "Event", rawEvent)

	// Extract the source schema ID from the event
	sourceSchemaID, err := sr.ExtractID(e.Record().Value)
	if err != nil {
		slog.Error("SCHEMA_ID not an integer", "Error", strconv.Itoa(sourceSchemaID))
		return nil, err
	}

	sourceSchema, err := getSchema(fmt.Sprintf("%d", sourceSchemaID))
	if err != nil {
		slog.Error("Error retrieving source schema", "Error", err)
		return nil, err
	}
	nestedMap := pUtils.DecodeAvro(sourceSchema, rawEvent)
	if nestedMap == nil {
		slog.Error("Unable to decode Avro event")
		return nil, err
	}
	return nestedMap, nil
}

func FetchAvroDestinationSchema() (*avro.Codec, []byte, error) {
	var (
		destinationSchemaID    string
		destinationSchemaIDInt int
		destinationCodec       *avro.Codec
		err                    error
	)

	destinationSchemaID = os.Getenv("DESTINATION_SCHEMA_ID")
	if destinationSchemaID == "" {
		panic("DESTINATION_SCHEMA_ID environment variable is required")
	}

	destinationSchemaIDInt, err = strconv.Atoi(destinationSchemaID)
	if err != nil {
		panic(fmt.Sprintf("DESTINATION_SCHEMA_ID not an integer: %s", destinationSchemaID))
	}

	destinationSchema, err := getSchema(destinationSchemaID)
	if err != nil {
		panic(fmt.Sprintf("Error retrieving destination schema: %v\n", err))
	}

	destinationCodec, err = avro.NewCodec(destinationSchema)
	if err != nil {
		panic(fmt.Sprintf("Error creating Avro codec: %v\n", err))
	}

	hdr := pUtils.EncodeBuffer(destinationSchemaIDInt)

	return destinationCodec, hdr, nil
}

func EncodeAvroRecord(nestedMap map[string]interface{}, destinationCodec *avro.Codec, hdr []byte, key []byte, headers []transform.RecordHeader) (transform.Record, error) {

	encoded, err := destinationCodec.BinaryFromNative(hdr, nestedMap)
	if err != nil {
		slog.Error("Error encoding Avro", "Error", err)
		return transform.Record{}, err
	}

	record := transform.Record{
		Key:     key,
		Value:   encoded,
		Headers: headers,
	}

	return record, nil
}

func ReadAndValidateAvroSchema(filename string) (*avro.Codec, []byte, error) {
	schemaFile, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("Failed to read Avro schema file", "Error", err)
		return nil, nil, err
	}

	// Convert the schema file content into a string
	schemaString := string(schemaFile)

	// Create a new Codec using the schema
	codec, err := avro.NewCodec(schemaString)
	if err != nil {
		slog.Error("Failed to create Avro codec", "Error", err)
		return nil, nil, err
	}

	slog.Debug("Successfully read and parsed Avro schema file")

	schemaForRegistry := map[string]string{
		"schema": codec.Schema(),
	}

	slog.Debug("Schema for registry", "Schema", schemaForRegistry)
	// Convert the map to a JSON string
	schemaJSON, err := json.Marshal(schemaForRegistry)
	if err != nil {
		slog.Error("Failed to marshal schema for registry", "Error", err)
	}

	return codec, schemaJSON, err
}
