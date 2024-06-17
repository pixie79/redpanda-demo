package types

import (
	"log/slog"
)

// DemoEventPayload represents the primary business data payload in the AVRO schema.
type DemoEventPayload struct {
	Id                 string  `json:"id" avro:"id"`
	NamePrefix         *string `json:"name_prefix,omitempty" avro:"namePrefix"`
	PreferredName      *string `json:"preferred_name,omitempty" avro:"preferredName"`
	GivenName          *string `json:"given_name,omitempty" avro:"givenName"`
	LastName           *string `json:"last_name,omitempty" avro:"lastName"`
	MiddleName         *string `json:"middle_name,omitempty" avro:"middleName"`
	DateOfBirth        *int    `json:"date_of_birth,omitempty" avro:"dateOfBirth"`
	DateOfDeath        *int    `json:"date_of_death,omitempty" avro:"dateOfDeath"`
	Gender             *string `json:"gender,omitempty" avro:"gender"`
	PlaceOfBirth       *string `json:"place_of_birth,omitempty" avro:"placeOfBirth"`
	CountryOfResidence *string `json:"country_of_residence,omitempty" avro:"countryOfResidence"`
}

// DemoEvent represents the top-level event structure in the AVRO schema.
type DemoEvent struct {
	Metadata Metadata         `json:"metadata" avro:"metadata"`
	Payload  DemoEventPayload `json:"business_data_payload" avro:"payload"`
}

// Converter for the current event type
type DemoEventConverter struct {
	Event DemoEvent
}

func (c *DemoEventConverter) ConvertToAvroRecord() (map[string]interface{}, error) {
	return ConvertToAvroRecordDemoEvent(c.Event)
}

func ConvertToAvroRecordDemoEvent(event DemoEvent) (map[string]interface{}, error) {
	// Convert Metadata
	metadataRecord := ConvertToAvroMetadataRecord(event.Metadata)

	// Convert Payload and its nested complex structures
	Payload := map[string]interface{}{
		"id":                   event.Payload.Id,
		"name_prefix":          wrapUnion(event.Payload.NamePrefix, "string"),
		"preferred_name":       wrapUnion(event.Payload.PreferredName, "string"),
		"given_name":           wrapUnion(event.Payload.GivenName, "string"),
		"last_name":            wrapUnion(event.Payload.LastName, "string"),
		"middle_name":          wrapUnion(event.Payload.MiddleName, "string"),
		"date_of_birth":        serializeDate(event.Payload.DateOfBirth),
		"date_of_death":        serializeDate(event.Payload.DateOfDeath),
		"gender":               wrapUnion(event.Payload.Gender, "string"),
		"place_of_birth":       wrapUnion(event.Payload.PlaceOfBirth, "string"),
		"country_of_residence": wrapUnion(event.Payload.CountryOfResidence, "string"),
	}

	slog.Debug("Serialized Payload", "Payload", Payload)

	// Final record
	avroRecord := map[string]interface{}{
		"metadata": metadataRecord,
		"payload":  Payload,
	}

	return avroRecord, nil
}
