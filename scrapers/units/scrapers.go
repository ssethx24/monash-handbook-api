package units

import (
	"encoding/json"
	"handbook-scraper/scrapers/common"
	"handbook-scraper/utils"
	"handbook-scraper/utils/log"
)

// Scrape extracts the relevant unit data from the raw JSON.
// It parses the JSON and populates the UnitScraperData field with the extracted information.
func Scrape(rawJSON map[string]interface{}, baseURL string) (UnitData, error) {
	log.Infof("[UNIT SCRAPER] Extracting data...")

	unitScraperData := UnitData{
		CommonScraperData: common.CommonScraperData{
			Link:             baseURL,
			Faculty:          utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.academic_org.value"),
			Code:             utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.unit_code"),
			Title:            utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.title"),
			SearchTitle:      utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.search_title"),
			CurrentYear:      utils.StringToInt(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.implementation_year")),
			AcademicItemType: utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.academic_item_type"),
		},
		Synopsis:             utils.RemoveHTMLTags(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.handbook_synopsis")),
		UnitLevel:            utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.level.label"),
		WorkloadRequirements: utils.RemoveHTMLTags(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.workload_requirements")),
		Active:               utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.status.value") == "Active",
		CreditPoints:         utils.GetTypedValue[int](rawJSON, "props.pageProps.pageContent.credit_points"),
		HandbookVersion:      utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.version_name"),
		EFTSL:                utils.GetTypedValue[float32](rawJSON, "props.pageProps.pageContent.eftsl"),
		HighestSCABand:       utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.highest_sca_band"),
		UndergradPostgrad:    utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.undergrad_postgrad_both.value"),
		AreaOfStudy:          utils.StringToArray(utils.RemoveHTMLTags(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.area_of_study_links"))),
		LearningOutcomes:     common.LearningOutcomes(rawJSON, "props.pageProps.pageContent.unit_learning_outcomes"),
		Assessments:          assessments(rawJSON),
		UnitOfferings:        unitOfferings(rawJSON),
		LearningActivities:   learningActivities(rawJSON),
		Requisites:           requisites(rawJSON),
		EnrolmentRules:       enrolmentRules(rawJSON),
	}

	log.Successf("[UNIT SCRAPER] Extraction complete.")

	return unitScraperData, nil
}

// requisites extracts and compresses the requisite data from the raw JSON.
// It navigates to the "requisites" path, extracts the data, and compresses it into a simplified structure.
func requisites(data map[string]interface{}) []CompressedRequisite {

	// Go to the array and get the content
	arrExtract := utils.GetTypedValue[[]map[string]interface{}](data, "props.pageProps.pageContent.requisites")

	// Marshal the array to a JSON formatted string
	marshalled, err := json.Marshal(arrExtract)
	if err != nil {
		log.Errorf("Error marshalling array: %v", err)
		return nil
	}

	// Compress the JSON formatted string
	compressed, err := compressRequisites(marshalled)
	if err != nil {
		log.Errorf("Error compressing requisites: %v", err)
		return nil
	}

	return compressed
}

// assessments parses the JSON input and extracts assessment data into a slice of Assessment structs.
// It navigates to the "assessments" path, extracts the data, and unmarshals it into the Assessment struct.
func assessments(data map[string]interface{}) []Assessment {

	// Extract the array using NavigateToArray
	arrExtract := utils.GetTypedValue[[]map[string]interface{}](data, "props.pageProps.pageContent.assessments")

	// Marshal the array to a JSON formatted string
	marshalled, err := json.Marshal(arrExtract)
	if err != nil {
		log.Errorf("Error marshalling array: %v", err)
		return nil
	}

	// Unmarshal directly into a slice of Assessment
	var assessments []Assessment
	err = json.Unmarshal(marshalled, &assessments)
	if err != nil {
		log.Errorf("Error unmarshalling JSON: %v", err)
		return nil
	}

	// Return the list of assessments
	return assessments
}

// unitOfferings parses the JSON input and extracts unit offering data into a slice of UnitOffering structs.
// It navigates to the "unit_offering" path, extracts the data, and unmarshals it into the UnitOffering struct.
func unitOfferings(data map[string]interface{}) []UnitOffering {
	// Extract the array using NavigateToArray
	arrExtract := utils.GetTypedValue[[]map[string]interface{}](data, "props.pageProps.pageContent.unit_offering")

	// Marshal the array to a JSON formatted string
	marshalled, err := json.Marshal(arrExtract)
	if err != nil {
		log.Errorf("Error marshalling array: %v", err)
		return nil
	}

	// Unmarshal directly into a slice of anonymous structs
	var offerings []struct {
		AttendanceMode struct {
			Value string `json:"value"`
		} `json:"attendance_mode"`
		DisplayName string `json:"display_name"`
		Location    struct {
			Value string `json:"value"`
		} `json:"location"`
		TeachingPeriod struct {
			Value string `json:"value"`
		} `json:"teaching_period"`
	}

	err = json.Unmarshal(marshalled, &offerings)
	if err != nil {
		log.Errorf("Error unmarshalling JSON: %v", err)
		return nil
	}

	// Map to UnitOffering
	var result []UnitOffering
	for _, offering := range offerings {
		result = append(result, UnitOffering{
			AttendanceMode: offering.AttendanceMode.Value,
			DisplayName:    offering.DisplayName,
			Location:       offering.Location.Value,
			Semester:       offering.TeachingPeriod.Value,
		})
	}

	return result
}

// learningActivities parses the JSON input and extracts learning activity data into a slice of LearningActivity structs.
// It navigates to the "learning_activities_grouped" path, extracts the data, and unmarshals it into the LearningActivity struct.
func learningActivities(data map[string]interface{}) []LearningActivity {
	// Extract the array using NavigateToArray
	arrExtract := utils.GetTypedValue[[]map[string]interface{}](data, "props.pageProps.pageContent.learning_activities_grouped")

	// Marshal the array to a JSON formatted string
	marshalled, err := json.Marshal(arrExtract)
	if err != nil {
		log.Errorf("Error marshalling array: %v", err)
		return nil
	}

	// Unmarshal directly into a slice of anonymous structs
	var groupedActivities []struct {
		Activities []struct {
			ActivityType struct {
				Label string `json:"label"`
			} `json:"activity_type"`
			DurationDisplay                      string `json:"duration_display"`
			OfferingsFormattedTeachingActivities string `json:"offerings_formatted_teaching_activities"`
		} `json:"activities"`
	}

	err = json.Unmarshal(marshalled, &groupedActivities)
	if err != nil {
		log.Errorf("Error unmarshalling JSON: %v", err)
		return nil
	}

	// Map to LearningActivity
	var result []LearningActivity
	for _, group := range groupedActivities {
		for _, activity := range group.Activities {
			// Clean up HTML tags in offerings formatted teaching activities
			cleanOfferings := utils.RemoveHTMLTags(activity.OfferingsFormattedTeachingActivities)

			result = append(result, LearningActivity{
				ActivityType:    activity.ActivityType.Label,
				DurationDisplay: activity.DurationDisplay,
				Offerings:       string(cleanOfferings),
			})
		}
	}

	return result
}

// enrolmentRules parses the JSON input and extracts enrolment rule data into a slice of EnrolmentRule structs.
// It navigates to the "enrolment_rules" path, extracts the data, and unmarshals it into the EnrolmentRule struct.
func enrolmentRules(data map[string]interface{}) []EnrolmentRule {
	// Extract the array using NavigateToArray
	arrExtract := utils.GetTypedValue[[]map[string]interface{}](data, "props.pageProps.pageContent.enrolment_rules")

	// Marshal the array to a JSON formatted string
	marshalled, err := json.Marshal(arrExtract)
	if err != nil {
		log.Errorf("Error marshalling array: %v", err)
		return nil
	}

	// Unmarshal directly into a slice of anonymous structs
	var rules []struct {
		Description string `json:"description"`
	}

	err = json.Unmarshal(marshalled, &rules)
	if err != nil {
		log.Errorf("Error unmarshalling JSON: %v", err)
		return nil
	}

	// Map to EnrolmentRule and clean descriptions
	var result []EnrolmentRule
	for _, rule := range rules {
		cleanDescription := utils.RemoveHTMLTags(rule.Description)
		result = append(result, EnrolmentRule{Description: string(cleanDescription)})
	}

	return result
}

// compressRequisites compresses the requisites data into a more manageable format.
// It takes a byte slice of JSON data as input, which represents a list of Requisite structs.
// It returns a slice of CompressedRequisite structs, which is a simplified version of the input data, and an error if any occurs.
func compressRequisites(data []byte) ([]CompressedRequisite, error) {

	// Unmarshal the JSON data into a slice of Requisite structs
	var requisites []Requisite
	err := json.Unmarshal(data, &requisites)
	if err != nil {
		return nil, err
	}

	// Compress the requisites
	var compressedRequisites []CompressedRequisite
	for _, req := range requisites {
		compReq := CompressedRequisite{
			RequisiteType: req.RequisiteType.Label,
			Containers:    []CompressedContainer{},
		}
		for _, container := range req.Containers {
			compContainer := compressContainer(container)
			compReq.Containers = append(compReq.Containers, compContainer)
		}
		compressedRequisites = append(compressedRequisites, compReq)
	}

	return compressedRequisites, nil
}

// compressContainer compresses a Container struct into a CompressedContainer struct.
// It takes a Container struct as input.
// It returns a CompressedContainer struct, which is a simplified version of the input data.
func compressContainer(container Container) CompressedContainer {
	compContainer := CompressedContainer{
		Relationship: container.ParentConnector.Label,
		Units:        []CompressedUnit{},
		Containers:   []CompressedContainer{},
	}

	// Extract units from relationships
	for _, rel := range container.Relationships {
		unit := CompressedUnit{
			UnitCode:   rel.AcademicItemCode,
			UnitNumber: utils.ExtractUnitNumber(rel.AcademicItemCode),
		}
		compContainer.Units = append(compContainer.Units, unit)
	}

	// Recursively compress child containers
	for _, child := range container.Containers {
		childComp := compressContainer(child)
		compContainer.Containers = append(compContainer.Containers, childComp)
	}

	return compContainer
}
