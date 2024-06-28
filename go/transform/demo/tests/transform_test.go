package main_test

import (
	"log/slog"
	"strconv"
	"testing"

	pUtils "pixie79/utils"
	pKgo "pixie79/utils/kgo"
	pTUtils "pixie79/utils/transforms/tests"

	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"

	"context"
	"os"

	transform "github.com/redpanda-data/redpanda/src/transform-sdk/go/transform"
	"github.com/testcontainers/testcontainers-go/modules/redpanda"
	"github.com/twmb/franz-go/pkg/kadm"
)

var (
	ctx              context.Context                  = context.Background()
	container        *redpanda.Container              = nil
	adminClient      *pTUtils.AdminAPIClient          = nil
	kafkaAdminClient *kadm.Client                     = nil
	schemaClient     *pTUtils.SchemaRegistryAPIClient = nil
	kgoClient        *kgo.Client                      = nil
	stop             func()                           = nil
)

const (
	testDataInput1 = `  {
		"metadata": {
		"message_key": "tnKGDKUndl",
		"created_date": 1296997036167,
		"updated_date": 693745893153,
		"outbox_published_date": 1203458653655,
		"event_type": "INSERT"
		},
		"business_data_payload": {
			"id": "PKs-Is7j",
			"name_prefix": "Mr",
			"preferred_name": "Tom",
			"given_name": "Tom",
			"last_name": "Jones",
			"middle_name": "A",
			"place_of_birth": "Sydney",
			"country_of_residence": "UK"
			}
		}`
	testDataOutput1 = `{
		"metadata": {
		"message_key": "tnKGDKUndl",
		"created_date": 1296997036167,
		"updated_date": 693745893153,
		"outbox_published_date": 1203458653655,
		"event_type": "INSERT"
		},
		"business_data_payload": {
			"id": "PKs-Is7j",
			"name_prefix": "Mr",
			"preferred_name": "Tom",
			"given_name": "******",
			"last_name": "******",
			"middle_name": "A",
			"place_of_birth": "Sydney",
			"country_of_residence": "UK"
			}
		}`
)

func TestMain(m *testing.M) {
	pUtils.SetupLogger()
	stop, kgoClient, kafkaAdminClient, adminClient, schemaClient, container = pTUtils.StartTest(ctx)
	// Run tests
	exitcode := m.Run()
	kgoClient.Close()
	stop()
	os.Exit(exitcode)
}

func TestDemo(t *testing.T) {
	var (
		inputTopic  = "demo"
		outputTopic = "output-demo"
		wasmFile    = "../demo.wasm"
		schemaFile  = "../../../../schemas/" + inputTopic + ".avsc"
		recordType  = "demoEvent"
	)

	t.Parallel()
	binary := pTUtils.LoadWasmFile(t, wasmFile)

	_, _ = pTUtils.DeploySchema(t, inputTopic+"-value", schemaFile, ctx, schemaClient)
	destinationSchemaId, destinationCodec := pTUtils.DeploySchema(t, outputTopic+"-value", schemaFile, ctx, schemaClient)

	metadata := pTUtils.TransformDeployMetadata{
		Name:         outputTopic,
		InputTopic:   inputTopic,
		OutputTopics: []string{outputTopic},
		Environment: []pTUtils.EnvironmentVariable{
			// {Key: "LOG_LEVEL", Value: "DEBUG"},
			{Key: "DESTINATION_SCHEMA_ID", Value: strconv.Itoa(destinationSchemaId)},
			{Key: "UNMASKED_CUSTOMERS", Value: "[{\\\"last_name\\\": \\\"Smith\\\",\\\"first_name\\\": \\\"Jane\\\"}]"},
		},
	}

	slog.Info("Deploying transform", "metadata", metadata)
	pTUtils.DeployTransform(t, metadata, binary, ctx, kafkaAdminClient, adminClient)

	hdr := pUtils.EncodeBuffer(destinationSchemaId)

	inputData1, err := pKgo.ConvertToAvroKgoRecord(recordType, []byte(testDataInput1), hdr, destinationCodec, []byte("eventKey"), []transform.RecordHeader{}, inputTopic)
	if err != nil {
		slog.Error("Error creating record", "Error", err)
	}

	outputData1, err := pKgo.ConvertToAvroKgoRecord(recordType, []byte(testDataOutput1), hdr, destinationCodec, []byte("eventKey"), []transform.RecordHeader{}, inputTopic)
	if err != nil {
		slog.Error("Error creating record", "Error", err)
	}

	slog.Debug("Creating client", "inputTopic", inputTopic, "outputTopic", outputTopic)

	client := pTUtils.MakeClient(t, ctx, container, kgo.DefaultProduceTopic(inputTopic), kgo.ConsumeTopics(outputTopic))

	// Produce records to be transformed
	slog.Debug("Producing record", "record", inputData1)
	defer client.Close()
	err = client.ProduceSync(ctx, inputData1).FirstErr()
	require.NoError(t, err)
	fetches := client.PollFetches(ctx)
	slog.Debug("Transform record produced", "record", fetches)
	pTUtils.RequireRecordsEquals(t, fetches, outputData1)

}
