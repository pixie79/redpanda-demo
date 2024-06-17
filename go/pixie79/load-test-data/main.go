package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	pUtils "pixie79/utils"
	pKgo "pixie79/utils/kgo"

	"github.com/joho/godotenv"
	avro "github.com/linkedin/goavro/v2"
	"github.com/twmb/franz-go/pkg/kgo"
)

var (
	destinationTopic    string
	destinationSchemaID string
	schemaURL           string
	seedEnv             string
	seeds               []string
)

func init() {

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Not using .env file")
	}

	pUtils.SetupLogger()

	destinationTopic = os.Getenv("REDPANDA_INPUT_TOPIC")
	if destinationTopic == "" {
		panic("REDPANDA_INPUT_TOPIC environment variable is required")
	}

	destinationSchemaID = os.Getenv("DESTINATION_SCHEMA_ID")
	if destinationSchemaID == "" {
		panic("DESTINATION_SCHEMA_ID environment variable is required")
	}

	schemaURL = os.Getenv("SCHEMA_REGISTRY_URL")
	if schemaURL == "" {
		panic("SCHEMA_REGISTRY_URL environment variable is required")
	}

	seedEnv = os.Getenv("REDPANDA_SEED_URL")
	if seedEnv == "" {
		panic("REDPANDA_SEED_URL environment variable is required")
	}
	seeds = []string{seedEnv}
}

func setupLoader() (*avro.Codec, []byte, string) {
	var (
		err error
	)

	destinationCodec, hdr, err := pKgo.FetchAvroDestinationSchema(schemaURL)
	if err != nil {
		panic(fmt.Sprintf("Error fetching destination schema: %v\n", err))
	}

	return destinationCodec, hdr, destinationTopic
}

func main() {
	var (
		avroRecords      []*kgo.Record
		defaultEventType = "onboardingDataIdentityVerificationSelfieV2"
		fileName         = flag.String("filename", "", "filename")
	)

	eventType := flag.String("t", defaultEventType, "Type of event data to generate (e.g., 'onboardingDataIdentityVerificationSelfieV2', 'alternative')")
	flag.Parse()

	destinationCodec, hdr, destinationTopic := setupLoader()

	// Load JSON data from file
	jsonData, err := os.ReadFile(*fileName)
	if err != nil {
		panic(fmt.Sprintf("Failed to read JSON file: %v", err))
	}

	if eventType == nil {
		slog.Error("Event type pointer is nil")
		return
	}

	eventTypestr := *eventType

	avroRecords, err = pKgo.ConvertToAvroKgoRecords(eventTypestr, jsonData, hdr, destinationCodec, []byte("eventKey"), nil, destinationTopic)
	if err != nil {
		slog.Error("Error converting to Avro records", "Error", err)
		return
	}

	if len(avroRecords) > 0 {
		err := pKgo.SubmitRecords(context.Background(), avroRecords, seeds)
		if err != nil {
			slog.Error("Error submitting records", "Error", err)
		}
	} else {
		slog.Info("No records to submit")
	}

}
