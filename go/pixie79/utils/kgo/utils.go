package utils

import (
	"encoding/json"
	"errors"
	"log/slog"

	pTypes "pixie79/types"
	pTransform "pixie79/utils/transforms"

	"github.com/linkedin/goavro/v2"
	"github.com/redpanda-data/redpanda/src/transform-sdk/go/transform"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Interface to convert events to Avro records
type AvroConverter interface {
	ConvertToAvroRecord() (map[string]interface{}, error)
}

func ConvertToAvroKgoRecords(eventType string, jsonData []byte, hdr []byte, codec *goavro.Codec, key []byte, headers []transform.RecordHeader, topic string) ([]*kgo.Record, error) {
	var (
		converter AvroConverter
		records   []*kgo.Record
	)

	switch eventType {

	case "demoEvent":
		var events []pTypes.DemoEvent
		if err := json.Unmarshal(jsonData, &events); err != nil {
			return nil, err
		}
		for _, event := range events {
			converter = &pTypes.DemoEventConverter{Event: event}
			record, err := convertToAvroKgoRecord(converter, hdr, codec, key, headers, topic)
			if err != nil {
				return nil, err
			}
			records = append(records, record)
		}
		return records, nil

	default:
		return nil, errors.New("unsupported event type")
	}
}

func ConvertToAvroKgoRecord(eventType string, jsonData []byte, hdr []byte, codec *goavro.Codec, key []byte, headers []transform.RecordHeader, topic string) (*kgo.Record, error) {
	var converter AvroConverter
	switch eventType {

	case "demoEvent":
		var event pTypes.DemoEvent
		err := json.Unmarshal(jsonData, &event)
		if err != nil {
			return nil, err
		}
		converter = &pTypes.DemoEventConverter{Event: event}
		return convertToAvroKgoRecord(converter, hdr, codec, key, headers, topic)

	default:
		return nil, errors.New("unsupported event type")
	}
}

func convertToAvroKgoRecord(converter AvroConverter, hdr []byte, codec *goavro.Codec, key []byte, headers []transform.RecordHeader, topic string) (*kgo.Record, error) {
	avroRecord, err := converter.ConvertToAvroRecord()
	if err != nil {
		return nil, err
	}

	encodedRecord, err := pTransform.EncodeAvroRecord(avroRecord, codec, hdr, key, headers)
	if err != nil {
		slog.Error("Error encoding Avro record", "Error", err)
	}

	kgoHeaders := make([]kgo.RecordHeader, len(encodedRecord.Headers))

	r := &kgo.Record{
		Key:     encodedRecord.Key,
		Value:   encodedRecord.Value,
		Headers: kgoHeaders,
		Topic:   topic,
	}

	return r, nil
}
