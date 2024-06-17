package main

import (
	pTypes "pixie79/types"
	pUtils "pixie79/utils"
)

func generateMetadata() pTypes.Metadata {
	return pTypes.Metadata{
		MessageKey:          pUtils.GenerateRandomString("", 10),
		CreatedDate:         pTypes.CustpTime{Time: pUtils.GenerateRandomDate()},
		UpdatedDate:         pTypes.CustpTime{Time: pUtils.GenerateRandomDate()},
		OutboxPublishedDate: pTypes.CustpTime{Time: pUtils.GenerateRandomDate()},
		EventType:           pUtils.GetWeightedRandomDatabaseOperation(),
	}
}
