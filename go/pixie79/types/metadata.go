package types

import "time"

type Metadata struct {
	MessageKey          string    `json:"message_key" avro:"messageKey"`                    // AVRO uses camelCase for keys
	CreatedDate         CustpTime `json:"created_date" avro:"createdDate"`                  // Custom AVRO type handling for time
	UpdatedDate         CustpTime `json:"updated_date" avro:"updatedDate"`                  // Same as above
	OutboxPublishedDate CustpTime `json:"outbox_published_date" avro:"outboxPublishedDate"` // More AVRO-specific mapping
	EventType           string    `json:"event_type" avro:"eventType"`                      // Matching the AVRO naming convention
}

func ConvertToAvroMetadataRecord(metadata Metadata) map[string]interface{} {
	return map[string]interface{}{
		"message_key":           metadata.MessageKey,
		"created_date":          metadata.CreatedDate.UnixNano() / int64(time.Millisecond),
		"updated_date":          metadata.UpdatedDate.UnixNano() / int64(time.Millisecond),
		"outbox_published_date": metadata.OutboxPublishedDate.UnixNano() / int64(time.Millisecond),
		"event_type":            metadata.EventType,
	}
}
