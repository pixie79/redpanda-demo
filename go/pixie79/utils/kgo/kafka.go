package utils

import (
	"context"
	"fmt"

	"log/slog"
	"pixie79/utils"

	"github.com/twmb/franz-go/pkg/kgo"
)

// createKafkaConnection creates a Kafka connection.
//
// ctx: the context.Context to use for the connection.
//
// Returns:
// - *kgo.Client: the Kafka client.
// - error: an error if the connection could not be established.
func createKafkaConnection(ctx context.Context, seeds []string) (*kgo.Client, error) {
	transactionID := utils.RandomString(20)

	opts := []kgo.Opt{
		kgo.SeedBrokers(seeds...),
		kgo.TransactionalID(transactionID),
		kgo.RecordPartitioner(kgo.RoundRobinPartitioner()),
		kgo.RecordRetries(4),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.AllowAutoTopicCreation(),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),
	}
	// Initialize public CAs for TLS
	// opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))

	//// Initializes SASL/SCRAM 256
	// Get Credentials from context
	// opts = append(opts, kgo.SASL(scram.Auth{
	// 	User: credentials.Username,
	// 	Pass: credentials.Password,
	// }.AsSha256Mechanism()))

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("could not connect to Kafka: %w", err)
	}

	return client, nil
}

// submitRecords submits Kafka records to a Kafka broker and commits them in a transaction.
//
// ctx - The context.Context object for cancellation signals and deadlines.
// kafkaRecords - A slice of *kgo.Record objects representing the Kafka records to be submitted.
// Returns an error if any step in the process fails.
func SubmitRecords(ctx context.Context, kafkaRecords []*kgo.Record, seeds []string) error {
	client, err := createKafkaConnection(context.Background(), seeds)
	if err != nil {
		return fmt.Errorf("failed to create Kafka client: %v", err)
	}

	if err := client.BeginTransaction(); err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	if err := produceMessages(ctx, client, kafkaRecords); err != nil {
		if rollbackErr := rollbackTransaction(client); rollbackErr != nil {
			return rollbackErr
		}
		return fmt.Errorf("failed to produce messages: %v", err)
	}

	if err := client.EndTransaction(ctx, kgo.TryCommit); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	slog.Info("produced Kafka records", "Count", len(kafkaRecords))
	return nil
}

// produceMessages is a function that produces messages using a Kafka client.
//
// It takes the following parameters:
// - ctx: the context.Context object used for cancellation and timeouts.
// - client: a pointer to a kgo.Client object representing the Kafka client.
// - records: a slice of pointers to kgo.Record objects representing the records to be produced.
//
// It returns an error if any error occurred while producing the messages.
func produceMessages(ctx context.Context, client *kgo.Client, records []*kgo.Record) error {
	var errPromise kgo.FirstErrPromise

	for _, record := range records {
		client.Produce(ctx, record, errPromise.Promise())
	}

	return errPromise.Err()
}

// rollbackTransaction is a function that rolls back a transaction.
//
// It takes a *kgo.Client as a parameter.
// It returns an error.
func rollbackTransaction(client *kgo.Client) error {
	ctx := context.Background()

	if err := client.AbortBufferedRecords(ctx); err != nil {
		return err
	}

	if err := client.EndTransaction(ctx, kgo.TryAbort); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
