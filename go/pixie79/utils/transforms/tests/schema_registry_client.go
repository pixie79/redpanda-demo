package utils

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
)

// SchemaRegistryAPIClient is a wrapper around Redpanda's SchemaRegistry API
type SchemaRegistryAPIClient struct {
	baseUrl string
	client  *http.Client
}

type SchemaRegistryID struct {
	ID int `json:"id"`
}

type SchemaRegistrySchemaDetails struct {
	Subject    string                          `json:"subject"`
	Version    int                             `json:"version"`
	ID         int                             `json:"id"`
	Schema     string                          `json:"schema"`
	SchemaType string                          `json:"schemaType"`
	References []SchemaRegistrySchemaReference `json:"references"`
}

type SchemaRegistrySchemaReference struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Version int    `json:"version"`
}

// NewSchemaRegistryAPIClient creates a new SchemaRegistry API
func NewSchemaRegistryAPIClient(baseURL string) *SchemaRegistryAPIClient {
	return &SchemaRegistryAPIClient{
		baseUrl: baseURL,
		client:  http.DefaultClient,
	}
}

// Deploy AVRO Schema to Schema Registry
func (cl *SchemaRegistryAPIClient) DeploySchema(ctx context.Context, schemaName string, schema io.Reader) (SchemaRegistryID, error) {
	endpoint, err := url.JoinPath(cl.baseUrl, "/subjects/", schemaName, "/versions")
	if err != nil {
		slog.Error("Failed to join url path", "error", err)
		return SchemaRegistryID{}, err
	}

	slog.Info("Deploying schema to Schema Registry", "endpoint", endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, schema)
	if err != nil {
		slog.Error("Failed to build http request", "error", err)
		return SchemaRegistryID{}, err
	}
	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")

	resp, err := cl.client.Do(req)
	if err != nil {
		slog.Error("Request failed", "error", err)
		return SchemaRegistryID{}, err
	}

	if err := checkResponse(resp); err != nil {
		slog.Error("Failed to check response", "error", err)
		return SchemaRegistryID{}, err
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Unable to read response body", "error", err)
		return SchemaRegistryID{}, err
	}

	metadata := SchemaRegistryID{}
	err = json.Unmarshal(buf, &metadata)
	if err != nil {
		slog.Error("Failed to unmarshal response", "error", err)
		return SchemaRegistryID{}, err
	}
	slog.Info("Schema deployed to Schema Registry", "SchemaName", schemaName, "ID", metadata.ID)

	return metadata, nil
}

// Deploy AVRO Schema to Schema Registry
func (cl *SchemaRegistryAPIClient) GetSchemaDetails(ctx context.Context, schemaName string, version int) (SchemaRegistrySchemaDetails, error) {
	endpoint, err := url.JoinPath(cl.baseUrl, "/subjects/", schemaName, "/versions/", strconv.Itoa(version))
	if err != nil {
		slog.Error("Failed to join url path", "error", err)
		return SchemaRegistrySchemaDetails{}, err
	}

	slog.Debug("Fetching schema details from Schema Registry", "endpoint", endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		slog.Error("Failed to build http request", "error", err)
		return SchemaRegistrySchemaDetails{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := cl.client.Do(req)
	if err != nil {
		slog.Error("Request failed", "error", err)
		return SchemaRegistrySchemaDetails{}, err
	}

	if err := checkResponse(resp); err != nil {
		slog.Error("Failed to check response", "error", err)
		return SchemaRegistrySchemaDetails{}, err
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Unable to read response body", "error", err)
		return SchemaRegistrySchemaDetails{}, err
	}

	metadata := SchemaRegistrySchemaDetails{}
	err = json.Unmarshal(buf, &metadata)
	if err != nil {
		slog.Error("Failed to unmarshal response", "error", err)
		return SchemaRegistrySchemaDetails{}, err
	}
	slog.Debug("Schema deployed to Schema Registry", "response", metadata)

	return metadata, nil
}
