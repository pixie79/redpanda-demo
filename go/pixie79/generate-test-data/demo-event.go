package main

import (
	"math/rand"
	pTypes "pixie79/types"
	pUtils "pixie79/utils"
)

func generateTestEventDemoEvent(customers []pTypes.TestCustomer) pTypes.DemoEvent {
	customer := customers[rand.Intn(len(customers))]

	// Generate the business data payload with random values
	Payload := pTypes.DemoEventPayload{
		Id:                 pUtils.GenerateRandomString("PK", 6),
		NamePrefix:         pUtils.IfEmptyReturnNilString(pUtils.RandomChoiceString([]string{"Mr", "Mrs", "Ms", ""})),
		PreferredName:      pUtils.IfEmptyReturnNilString(&customer.GivenName),
		GivenName:          pUtils.IfEmptyReturnNilString(&customer.GivenName),
		LastName:           pUtils.IfEmptyReturnNilString(&customer.LastName),
		MiddleName:         pUtils.IfEmptyReturnNilString(pUtils.RandomChoiceString([]string{"A", "B", "C", ""})),
		DateOfBirth:        pUtils.DateOrNil(),
		DateOfDeath:        nil, // Not likely to be populated in a normal case
		Gender:             pUtils.IfEmptyReturnNilString(pUtils.RandomChoiceString([]string{"Male", "Female", ""})),
		PlaceOfBirth:       pUtils.IfEmptyReturnNilString(pUtils.RandomChoiceString([]string{"London", "New York", "Sydney", ""})),
		CountryOfResidence: pUtils.IfEmptyReturnNilString(pUtils.RandomChoiceString([]string{"UK", "USA", "Australia", ""})),
	}

	// Return the complete test event
	return pTypes.DemoEvent{
		Metadata: generateMetadata(), // Generate the metadata
		Payload:  Payload,            // Assign the business data payload
	}

}
