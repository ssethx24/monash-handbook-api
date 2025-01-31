package common

import (
	"encoding/json"
	"handbook-scraper/utils"
	"handbook-scraper/utils/log"
)

// LearningOutcomes parses the JSON input into a slice of LearningOutcome.
// It takes a map of string to interface and a path string as input.
// It extracts an array of maps from the given path, marshals it to JSON,
// unmarshals it into a slice of LearningOutcome structs, and removes HTML tags from the descriptions.
// It returns a slice of LearningOutcome structs.
func LearningOutcomes(data map[string]interface{}, path string) []LearningOutcome {
	// Extract the array using GetTypedValue
	arrExtract := utils.GetTypedValue[[]map[string]interface{}](data, path)

	// Marshal the array to a JSON formatted string
	marshalled, err := json.Marshal(arrExtract)
	if err != nil {
		log.Errorf("Failed marshalling array: %v", err)
		return nil
	}

	// Unmarshal directly into a slice of LearningOutcome
	var outcomes []LearningOutcome
	err = json.Unmarshal(marshalled, &outcomes)
	if err != nil {
		log.Errorf("Failed unmarshalling JSON: %v", err)
		return nil
	}

	// Clean up HTML tags in the descriptions
	for i := range outcomes {
		outcomes[i].Description = utils.RemoveHTMLTags(outcomes[i].Description)
	}

	return outcomes
}
