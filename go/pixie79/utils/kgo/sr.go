package utils

import (
	"fmt"
	"strconv"

	sr "github.com/landoop/schema-registry"
)

// getSchema retrieves the schema for a given schema ID from the Schema Registry.
//
// Parameters:
// - schemaID: The ID of the schema to retrieve.
// - registryURL: The URL of the Schema Registry.
//
// Returns:
// - string: The retrieved schema.
// - error: An error if the schema retrieval fails.
func getSchema(schemaID string, registryURL string) (string, error) {
	registry, err := sr.NewClient(registryURL)
	if err != nil {
		return "", fmt.Errorf("cannot connect to Schema Registry: %w", err)
	}

	schemaIDInt, err := strconv.Atoi(schemaID)
	if err != nil {
		return "", fmt.Errorf("schema ID not an integer: %w", err)
	}

	schema, err := registry.GetSchemaByID(schemaIDInt)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve schema for ID %s: %w", schemaID, err)
	}

	return schema, nil
}
