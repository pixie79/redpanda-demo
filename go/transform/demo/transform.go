package main

import (
	"fmt"
	"log/slog"
	"os"

	pUtils "pixie79/utils"
	pTransforms "pixie79/utils/transforms"

	avro "github.com/linkedin/goavro/v2"
	transform "github.com/redpanda-data/redpanda/src/transform-sdk/go/transform"
)

var (
	destinationCodec    *avro.Codec
	hdr                 []byte
	unmaskedCustomerMap map[string]bool
)

func init() {
	var (
		err               error
		unmaskedCustomers string
	)

	pUtils.SetupLogger()

	unmaskedCustomers = os.Getenv("UNMASKED_CUSTOMERS")
	if unmaskedCustomers == "" {
		slog.Error("UNMASKED_CUSTOMERS environment variable is required")
		panic("UNMASKED_CUSTOMERS environment variable is required")
	}
	slog.Debug("UNMASKED_CUSTOMERS", "unmaskedCustomers", unmaskedCustomers)

	destinationCodec, hdr, err = pTransforms.FetchAvroDestinationSchema()
	if err != nil {
		slog.Error("Error fetching destination schema", "Error", err)
		panic(fmt.Sprintf("Error fetching destination schema: %v\n", err))
	}

	_, unmaskedCustomerMap, err = pUtils.UnmarshalCustomers(unmaskedCustomers)
	if err != nil {
		slog.Error("Error unmarshalling customers", "Error", err)
	}
	slog.Debug("Not Masking Customers with the last_name", "unmaskedCustomerMap", unmaskedCustomerMap)

}

func main() {
	slog.Info("Running transformer")
	transform.OnRecordWritten(toAvro)
}

// func toAvro(e transform.WriteEvent, w transform.RecordWriter) error {
// 	fmt.Println("Returning AVRO", "record", e.Record())
// 	return w.Write(e.Record())
// }

func toAvro(e transform.WriteEvent, w transform.RecordWriter) error {
	var (
		LastName  string
		GivenName string
		err       error
	)
	// Decode the raw event
	nestedMap, err := pTransforms.DecodeAvroRawEvent(e)
	if err != nil {
		slog.Error("Error decoding Avro", "Error", err)
		return err
	}

	Payload := nestedMap["payload"].(map[string]interface{})

	ln, ok := Payload["last_name"].(map[string]interface{})
	if ok {
		LastName = ln["string"].(string)
		if ok {
			// Check if the last name is in the list of customers to not mask
			if !pUtils.StringInMap(LastName, unmaskedCustomerMap) {
				if gn, ok := Payload["given_name"].(map[string]interface{}); ok {
					GivenName = gn["string"].(string)
				}

				Payload["given_name"] = pUtils.WrapUnionSimple(pUtils.MaskString(GivenName, "*", "fixed", 6), "string")
				Payload["last_name"] = pUtils.WrapUnionSimple(pUtils.MaskString(LastName, "*", "fixed", 6), "string")
				slog.Debug("Customer found - masking.")
			} else {
				slog.Info("Unmasked Customer found - not masking.")
			}
		} else {
			slog.Debug("Last name not available.")
		}
	} else {
		slog.Debug("Last name field not found.")
	}

	record, err := pTransforms.EncodeAvroRecord(nestedMap, destinationCodec, hdr, e.Record().Key, e.Record().Headers)
	if err != nil {
		slog.Error("Error encoding Avro", "Error", err)
		return err
	}
	slog.Debug("Returning AVRO", "record", record)
	return w.Write(record)
}
