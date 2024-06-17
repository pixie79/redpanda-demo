package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	pUtils "pixie79/utils"
	pTransforms "pixie79/utils/transforms"
	"testing"

	"github.com/linkedin/goavro/v2"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redpanda"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

var (
	adminClient      *AdminAPIClient          = nil
	kafkaClient      *kgo.Client              = nil
	kafkaAdminClient *kadm.Client             = nil
	schemaClient     *SchemaRegistryAPIClient = nil
)

// Load a Wasm file, if it doesn't exist the test is skipped.
func LoadWasmFile(t *testing.T, wasmFile string) []byte {
	contents, err := os.ReadFile(wasmFile)
	if err != nil {
		t.Fatalf("failed to read wasm file %s: %v", wasmFile, err)
	}
	return contents
}

type stdoutLogConsumer struct{}

// Accept prints the log to stdout
func (lc *stdoutLogConsumer) Accept(l testcontainers.Log) {
	slog.Debug(string(l.Content))
}

// Define a wrapper that embeds *slog.Logger and implements testcontainers.Logging
type testcontainersLogger struct {
	logger slog.Logger
}

// Printf method to satisfy the testcontainers.Logging interface
func (l *testcontainersLogger) Printf(format string, v ...interface{}) {
	// You can use fmt.Sprintf to format the string and pass it to any method of slog that accepts a single string
	// For simplicity, assuming you can use Error method for output. Adjust based on your actual logging level needs.
	msg := fmt.Sprintf(format, v...)
	l.logger.Debug(msg)
}

var logger *slog.Logger

// startRedpanda runs the Redpanda binary with a data transforms enabled.
func StartRedpanda(ctx context.Context) (*redpanda.Container, context.CancelFunc) {
	logger = pUtils.SetupLogger()
	// Create an instance of your custom logger
	tcLogger := &testcontainersLogger{logger: *logger}

	redpandaContainer, err := redpanda.RunContainer(ctx,
		testcontainers.WithLogger(tcLogger),
		testcontainers.WithImage("redpandadata/redpanda-nightly:latest"),
		testcontainers.CustomizeRequestOption(func(req *testcontainers.GenericContainerRequest) {
			if req.LogConsumerCfg == nil {
				req.LogConsumerCfg = &testcontainers.LogConsumerConfig{}
			}
			// Comment this out to get broker logs
			req.LogConsumerCfg.Consumers = append(req.LogConsumerCfg.Consumers, &stdoutLogConsumer{})
		}),
		redpanda.WithEnableWasmTransform(),
	)
	if err != nil {
		slog.Error("failed to start container", "Error", err)
	}
	stopFunc := func() {
		if err := redpandaContainer.Terminate(ctx); err != nil {
			slog.Error("failed to terminate container", "Error", err)
		}
	}
	return redpandaContainer, stopFunc
}

func RequireRecordsEquals(t *testing.T, fetches kgo.Fetches, records ...*kgo.Record) {
	require.NoError(t, fetches.Err())
	require.Equal(t, fetches.NumRecords(), len(records))
	for i, got := range fetches.Records() {
		want := records[i]
		require.Equal(t, want.Key, got.Key, "record %d key mismatch", i)
		require.Equal(t, want.Value, got.Value, "record %d value mismatch", i)
		require.Equal(t, want.Headers, got.Headers, "record %d headers mismatch", i)
	}
}

func DeployTransform(t *testing.T, metadata TransformDeployMetadata, binary []byte, ctx context.Context, kafkaAdminClient *kadm.Client, adminClient *AdminAPIClient) {
	topics := []string{metadata.InputTopic}
	topics = append(topics, metadata.OutputTopics...)
	_, err := kafkaAdminClient.CreateTopics(ctx, 1, 1, nil, topics...)
	require.NoError(t, err)
	err = adminClient.DeployTransform(ctx, metadata, bytes.NewReader(binary))
	require.NoError(t, err)
}

func DeploySchema(t *testing.T, schemaName string, filename string, ctx context.Context, schemaClient *SchemaRegistryAPIClient) (int, *goavro.Codec) {
	codec, schemaJson, err := pTransforms.ReadAndValidateAvroSchema(filename)
	require.NoError(t, err)
	response, err := schemaClient.DeploySchema(ctx, schemaName, bytes.NewReader(schemaJson))
	require.NoError(t, err)
	slog.Info("Deployed schema", "schemaId", response.ID, "schemaName", schemaName)
	return response.ID, codec
}

func MakeClient(t *testing.T, ctx context.Context, container *redpanda.Container, opts ...kgo.Opt) *kgo.Client {
	broker, err := container.KafkaSeedBroker(ctx)
	require.NoError(t, err)
	opts = append(opts, kgo.SeedBrokers(broker))
	kgoClient, err := kgo.NewClient(opts...)
	require.NoError(t, err)
	return kgoClient
}

func StartTest(ctx context.Context) (context.CancelFunc, *kgo.Client, *kadm.Client, *AdminAPIClient, *SchemaRegistryAPIClient, *redpanda.Container) {
	slog.Info("starting Redpanda...")
	// Start container, this is shared for all the tests so that they can run in parallel and be faster.
	container, stop := StartRedpanda(ctx)
	slog.Info("Redpanda started!")

	// Setup admin client
	adminURL, err := container.AdminAPIAddress(ctx)
	if err != nil {
		slog.Error("unable to access Admin API Address", "Error", err)
		panic(err)
	}
	adminClient = NewAdminAPIClient(adminURL)
	schemaUrl, err := container.SchemaRegistryAddress(ctx)
	if err != nil {
		slog.Error("unable to access Schema Registry Address", "Error", err)
		panic(err)
	}
	schemaClient = NewSchemaRegistryAPIClient(schemaUrl)

	// Setup broker
	broker, err := container.KafkaSeedBroker(ctx)
	if err != nil {
		slog.Error("unable to access Admin API Address", "Error", err)
		panic(err)
	}
	kgoClient, err := kgo.NewClient(
		kgo.SeedBrokers(broker),
	)
	if err != nil {
		log.Fatalf("unable to create kafka client: %v", err)
	}
	kafkaClient = kgoClient

	kafkaAdminClient = kadm.NewClient(kafkaClient)

	slog.Info("Schema Registry Endpoint", "URL", schemaUrl)
	slog.Info("Admin API Endpoint", "URL", adminURL)
	slog.Info("Kafka Seed Broker", "URL", broker)

	return stop, kgoClient, kafkaAdminClient, adminClient, schemaClient, container
}
