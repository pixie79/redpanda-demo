package utils

import (
	"fmt"
	"log/slog"
	"os"
	"pixie79/utils"
	"strconv"

	avro "github.com/linkedin/goavro/v2"
	"github.com/twmb/franz-go/pkg/kgo"
)

func FetchAvroDestinationSchema(schemaURL string) (*avro.Codec, []byte, error) {
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

	remoteSchema, err := getSchema(destinationSchemaID, schemaURL)
	if err != nil {
		panic(fmt.Sprintf("Error retrieving destination schema: %v\n", err))
	}

	destinationCodec, err = avro.NewCodec(remoteSchema)
	if err != nil {
		panic(fmt.Sprintf("Error creating Avro codec: %v\n", err))
	}

	hdr := utils.EncodeBuffer(destinationSchemaIDInt)

	return destinationCodec, hdr, nil
}

func EncodeAvroRecord(nestedMap map[string]interface{}, destinationCodec *avro.Codec, hdr []byte, key []byte, destinationTopic string) (*kgo.Record, error) {

	encoded, err := destinationCodec.BinaryFromNative(hdr, nestedMap)
	if err != nil {
		slog.Error("Error encoding Avro", "Error", err)
		return &kgo.Record{}, err
	}

	record := &kgo.Record{
		Key:   key,
		Value: encoded,
		Topic: destinationTopic,
	}

	return record, nil
}
