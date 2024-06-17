package utils

import (
	"log/slog"
	"strconv"

	sr "github.com/redpanda-data/redpanda/src/transform-sdk/go/transform/sr"
)

// getSchema retrieves the schema with the given ID from the specified URL.
//
// Parameters:
// - id: The ID of the schema (as a string).
//
// Returns:
// - The retrieved schema (as a string).
func getSchema(id string) (string, error) {
	client := sr.NewClient()
	schemaID, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("SCHEMA_ID not an integer:", "Error", id)
		return "", err
	}
	schema, err := client.LookupSchemaById(schemaID)
	if err != nil {
		slog.Error("Unable to retrieve schema for ID", "Error", id)
		return "", err
	}
	return schema.Schema, err
}
