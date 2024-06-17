package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	pTypes "pixie79/types"
	pUtils "pixie79/utils"
	"strconv"
)

const (
	defaultCustomers = `
	[
		{
			"given_name": "John",
			"last_name": "Doe",
			"national_identity_numbers": ["123456789"]
		},
		{
			"given_name": "Jane",
			"last_name": "Smith",
			"national_identity_numbers": ["7707077777087"]
		},
		{
			"given_name": "Fred",
			"last_name": "Blogs",
			"national_identity_numbers": ["987654321"]
		},
		{
			"given_name": "Amy",
			"last_name": "Dune",
			"national_identity_numbers": ["7707077777087"]
		},
		{
			"given_name": "Michael",
			"last_name": "Johnson",
			"national_identity_numbers": ["987654321"]
		},
		{
			"given_name": "Tom",
			"last_name": "Jones",
			"national_identity_numbers": ["987654321", "7707077777087"]
		},
		{
			"given_name": "",
			"last_name": "",
			"national_identity_numbers": ["7707077777087"]
		},
		{
			"given_name": "",
			"last_name": "",
			"national_identity_numbers": [""]
		}
	]`
)

func main() {
	pUtils.SetupLogger()
	// Default values for flags
	defaultNumEvents := 2000
	defaultFilename := "demo_event_data.json"
	defaultEventType := "demoEvent" // Default event type to generate

	var customers []pTypes.TestCustomer
	if err := json.Unmarshal([]byte(defaultCustomers), &customers); err != nil {
		slog.Error("Error unmarshalling JSON", "Error", err)
		panic(err)
	}

	// Flags for custom input
	numEvents := flag.Int("n", defaultNumEvents, "Number of events to generate")
	outputFilename := flag.String("o", defaultFilename, "Output filename")
	eventType := flag.String("t", defaultEventType, "Type of event data to generate (e.g., 'demoEvent', 'alternative')")
	flag.Parse()

	var events []interface{}

	// Determine the type of data to generate based on the CLI argument
	switch *eventType {
	case "demoEvent":
		events = make([]interface{}, *numEvents)
		for i := 0; i < *numEvents; i++ {
			events[i] = generateTestEventDemoEvent(customers)
		}

	default:
		slog.Error("Unknown event type", "Error", *eventType)
		return
	}

	// Serialize to JSON and save to file
	file, err := os.Create(*outputFilename)
	if err != nil {
		slog.Error("Error creating file", "Error", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print
	if err := encoder.Encode(events); err != nil {
		slog.Error("Error encoding JSON", "Error", err)
		return
	}

	slog.Info("Generating", "events", strconv.Itoa(*numEvents), "Type", *eventType, "Output File", *outputFilename)
}
