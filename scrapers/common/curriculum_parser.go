package common

import (
	"fmt"
	"handbook-scraper/utils"
	"handbook-scraper/utils/log"
)

// Curriculum represents the overall curriculum structure.
// It contains the total credit points and a slice of Part structs.
type Curriculum struct {
	TotalCreditPoints int    `json:"total_credit_points"`
	Parts             []Part `json:"parts"`
}

// Part represents a major section of the curriculum (e.g., Part A, Part B).
// It contains the title, description, credit points required, and a slice of Container structs.
type Part struct {
	Title                string         `json:"title"`
	Description          string         `json:"description"`
	CreditPointsRequired int            `json:"credit_points_required"`
	Containers           []Container    `json:"containers"`
	AcademicItems        []AcademicItem `json:"academic_items"`
	Order                int            `json:"order"`
	Connector            string         `json:"connector"` // Represents the connectors between child academicItems OR containers
}

// Container represents a subset of units within a Part (e.g., core units, electives). Containers can be nested
// It contains the title, description, credit points required, a slice of AcademicItem structs, and a connector string (e.g., "AND" or "OR").
// The connector string defines the relationship between the academic items in the container.
// Containers cannot contain both academic items and child containers simultaneously.
type Container struct {
	Title                string         `json:"title"`
	Description          string         `json:"description"`
	CreditPointsRequired int            `json:"credit_points_required"`
	Containers           []Container    `json:"containers"`
	AcademicItems        []AcademicItem `json:"academic_items"`
	Connector            string         `json:"connector"` // Represents the connectors between child academicItems OR containers
}

// AcademicItem represents an academic item (e.g., unit, course, specialization).
// It contains the title, code, description, and credit points.
type AcademicItem struct {
	Type         string `json:"type"` // either units OR area_of_study
	Title        string `json:"title"`
	Code         string `json:"code"`
	Description  string `json:"description"`
	CreditPoints int    `json:"credit_points"`
	URL          string `json:"url"`
}

// ParseCurriculum parses the curriculum JSON into a Curriculum struct.
// It takes a map of string to interface as input, which should contain the curriculum data.
// It extracts the curriculum structure from the given path, parses the total credit points,
// and then iterates through each part of the curriculum, extracting its details and nested containers.
// It returns a Curriculum struct and an error if any parsing fails.
func ParseCurriculum(data map[string]interface{}) (Curriculum, error) {
	data = utils.GetTypedValue[map[string]interface{}](data, "props.pageProps.pageContent.curriculumStructure")

	var curriculum Curriculum
	curriculum.Parts = []Part{} // Initialize as empty slice

	// Extract total credit points.
	totalCreditsStr, ok := data["credit_points"].(string)
	if !ok {
		return curriculum, fmt.Errorf("total credit points not found or not a string")
	}
	curriculum.TotalCreditPoints = utils.StringToInt(totalCreditsStr)

	// Extract the top-level containers (Parts).
	partsData, ok := data["container"].([]interface{})
	if !ok {
		return curriculum, fmt.Errorf("parts container not found or not an array")
	}

	for _, partInterface := range partsData {
		partMap, ok := partInterface.(map[string]interface{})
		if !ok {
			log.Errorf("Skipping invalid part data: %v", partInterface)
			continue
		}

		description, _ := partMap["description"].(string)
		description = utils.RemoveHTMLTags(description)

		// Extract part details.
		title, _ := partMap["title"].(string)
		creditPointsStr, _ := partMap["credit_points"].(string)
		creditPoints := utils.StringToInt(creditPointsStr)
		order := utils.StringToInt(partMap["order"].(string))

		part := Part{
			Title:                title,
			Description:          description,
			CreditPointsRequired: creditPoints,
			Containers:           []Container{}, // Initialize as empty slice
			Order:                order,
			Connector:            "AND", // Default connector
		}

		// Check if the part has nested containers.
		if containersRaw, exists := partMap["container"]; exists {
			containers, _, err := parseContainers(containersRaw)
			if err != nil {
				log.Errorf("Error parsing containers for part '%s': %v", part.Title, err)
			}
			// Append parsed containers.
			part.Containers = append(part.Containers, containers...)

			// CHECK: Check if children container matches the parent connector
			// Logically, if the first child container has the same credit points as the parent container, then the connector should be OR
			if len(part.Containers) > 0 {
				if part.Containers[0].CreditPointsRequired == part.CreditPointsRequired {
					part.Connector = "OR"
				}
			}
		}

		// If no containers were found, check for direct relationships
		if len(part.Containers) == 0 {
			// Extract items from relationships.
			if relationshipsRaw, exists := partMap["relationship"]; exists {
				items, err := parseItems(relationshipsRaw, &partMap)
				if err != nil {
					log.Errorf("Error parsing items for part '%s': %v", part.Title, err)
				}

				if len(items) > 0 {
					part.AcademicItems = items
				}

				// CHECK: Check if children container matches the parent connector
				// Logically, if the first child container has the same credit points as the parent container, then the connector should be OR
				if len(part.AcademicItems) > 0 {
					if part.AcademicItems[0].CreditPoints == part.CreditPointsRequired {
						part.Connector = "OR"
					}
				}
			}
		}

		// Ensure Containers is not nil
		if part.Containers == nil {
			part.Containers = []Container{}
		}

		if part.CreditPointsRequired == curriculum.TotalCreditPoints {
			part.CreditPointsRequired = 0
		}

		// Append the parsed part to the curriculum
		curriculum.Parts = append(curriculum.Parts, part)
	}

	return curriculum, nil
}

// parseContainers recursively parses containers and their nested containers.
// It takes an interface as input, which should be a slice of container data.
// It iterates through each container, extracts its details, and recursively parses nested containers.
// It returns a slice of Container structs and an error if any parsing fails.
// It returns the relationship of the parent container as a string (could be empty string)
func parseContainers(containerData interface{}) ([]Container, string, error) {
	var containers []Container
	var parentConnector string

	// Ensure containerData is a slice.
	containerSlice, ok := containerData.([]interface{})
	if !ok {
		// If containerData is not a slice, return empty slice instead of error.
		return containers, "", nil
	}

	for _, containerInterface := range containerSlice {
		containerMap, ok := containerInterface.(map[string]interface{})
		if !ok {
			log.Errorf("Skipping invalid container data: %v", containerInterface)
			continue
		}

		description, _ := containerMap["description"].(string)
		description = utils.RemoveHTMLTags(description)

		// Extract container details.
		title, _ := containerMap["title"].(string)
		creditPointsStr, _ := containerMap["credit_points"].(string)
		creditPoints := utils.StringToInt(creditPointsStr)

		// Extract parent connector
		conn := containerMap["parent_connector"].(map[string]interface{})
		parentConnector = conn["value"].(string)

		container := Container{
			Title:                title,
			Description:          description,
			CreditPointsRequired: creditPoints,
			AcademicItems:        []AcademicItem{}, // Initialize as empty slice
			Connector:            "AND",            // Default connector
		}

		// Extract items from relationships.
		if relationshipsRaw, exists := containerMap["relationship"]; exists {
			items, err := parseItems(relationshipsRaw, &containerMap)
			if err != nil {
				log.Errorf("Error parsing items for container '%s': %v", container.Title, err)
			}
			container.AcademicItems = items
		}

		// Check for nested containers and parse them recursively.
		if nestedContainersRaw, exists := containerMap["container"]; exists && nestedContainersRaw != nil {
			nestedContainers, relationship, err := parseContainers(nestedContainersRaw)
			if relationship != "" {
				container.Connector = relationship
			}
			if err != nil {
				log.Errorf("Error parsing nested containers for container '%s': %v", container.Title, err)
			}

			// Append nested containers to the current container's containers.
			container.Containers = append(container.Containers, nestedContainers...)
		}

		// CHECK: Check if children container matches the parent connector
		// Logically, if the first child container has the same credit points as the parent container, then the connector should be OR
		// Same goes with academic items
		if len(container.Containers) > 0 {
			if container.Containers[0].CreditPointsRequired == container.CreditPointsRequired {
				container.Connector = "OR"
			}
		}
		if len(container.AcademicItems) > 0 {
			if container.AcademicItems[0].CreditPoints == container.CreditPointsRequired {
				container.Connector = "OR"
			}
		}

		// Append the parsed container to the list
		containers = append(containers, container)
	}

	return containers, parentConnector, nil
}

// parseItems parses the relationship array to extract academic items.
// It takes an interface as input, which should be a slice of item data and the container map.
// It iterates through each item, extracts its details, and returns a slice of AcademicItem structs.
func parseItems(itemsData interface{}, containerMap *map[string]interface{}) ([]AcademicItem, error) {
	var academicItems []AcademicItem

	// Ensure itemsData is a slice.
	itemsSlice, ok := itemsData.([]interface{})
	if !ok {
		// If itemsData is not a slice, return empty slice instead of error.
		return academicItems, nil
	}

	for _, itemInterface := range itemsSlice {
		itemMap, ok := itemInterface.(map[string]interface{})
		if !ok {
			log.Logf("Skipping invalid academic item data: %v", itemInterface)
			continue
		}

		description, _ := itemMap["description"].(string)
		description = utils.RemoveHTMLTags(description)

		code, _ := itemMap["academic_item_code"].(string)
		name, _ := itemMap["academic_item_name"].(string)
		creditPointsStr, _ := itemMap["academic_item_credit_points"].(string)
		creditPoints := utils.StringToInt(creditPointsStr)
		url, _ := itemMap["academic_item_url"].(string)
		itemType, _ := itemMap["academic_item_type"].(map[string]interface{})
		itemTypeValue, _ := itemType["value"].(string)

		academicItem := AcademicItem{
			Type:         itemTypeValue,
			Code:         code,
			Title:        name,
			CreditPoints: creditPoints,
			Description:  description,
			URL:          url,
		}

		academicItems = append(academicItems, academicItem)
	}

	// Extract connector from the relationship.
	if relationship, exists := (*containerMap)["relationship"].([]interface{}); exists && len(relationship) > 0 {
		if relMap, ok := relationship[0].(map[string]interface{}); ok {
			if connectorMap, ok := relMap["parent_connector"].(map[string]interface{}); ok {
				if connectorValue, ok := connectorMap["value"].(string); ok {
					(*containerMap)["connector"] = connectorValue
				}
			}
		}
	}

	return academicItems, nil
}
