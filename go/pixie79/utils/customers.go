package utils

import (
	"encoding/json"
	"log/slog"
	pTypes "pixie79/types"
	"strings"
)

// customerExists checks if there is a match in the customer slice for given first and last names.
func CustomerExists(givenName, lastName string, customerIndex map[string]bool) bool {
	// Create the lookup key similar to how we indexed it
	key := strings.ToLower(givenName + lastName)
	if exists := customerIndex[key]; exists {
		// Sensitive Debug output
		// slog.Debug("Customer found:", "GivenName", givenName, "LastName", lastName)
		return true
	}
	return false
}

func IndexCustomers(customers []pTypes.TestCustomer) map[string]bool {
	index := make(map[string]bool)
	for _, customer := range customers {
		// Create a key by concatenating lowercase versions of given name and last name
		key := strings.ToLower(customer.GivenName + customer.LastName)
		index[key] = true
	}
	return index
}

func UnmarshalCustomers(data string) ([]pTypes.TestCustomer, map[string]bool, error) {
	var (
		customers []pTypes.TestCustomer
	)

	err := json.Unmarshal([]byte(data), &customers)
	if err != nil {
		slog.Error("Error unmarshalling JSON", "Error", err)
		return []pTypes.TestCustomer{}, nil, err
	}

	// Index the customers
	customerIndex := IndexCustomers(customers)

	return customers, customerIndex, nil
}
