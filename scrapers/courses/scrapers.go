package courses

import (
	"handbook-scraper/scrapers/common"
	"handbook-scraper/utils"
	"handbook-scraper/utils/log"
)

// Scrape extracts the relevant course data from the raw JSON.
// It parses the JSON and populates the CourseScraperData field with the extracted information.
func Scrape(rawJSON map[string]interface{}, baseURL string) (CourseData, error) {
	log.Infof("[COURSE SCRAPER] Extracting data...")

	curriculum, errCurriculum := common.ParseCurriculum(rawJSON)
	var curriculumError bool
	if errCurriculum != nil {
		log.Errorf("[COURSE SCRAPER]: Error parsing curriculum: %v", errCurriculum)
		curriculumError = true
	} else {
		curriculumError = false
	}

	courseScraperData := CourseData{
		CommonScraperData: common.CommonScraperData{
			Link:             baseURL,
			Faculty:          utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.school.value"),
			Code:             utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.course_code"),
			Title:            utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.title"),
			SearchTitle:      utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.search_title"),
			CurrentYear:      utils.StringToInt(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.implementation_year")),
			AcademicItemType: utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.academic_item_type"),
		},
		ProfessionalAccreditation: utils.RemoveHTMLTags(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.professional_accreditation")),
		AbbreviatedName:           utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.abbreviated_name"),
		Atar:                      utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.atar"),
		AwardTitles:               extractAwardTitles(rawJSON),
		CourseDuration:            utils.RemoveHTMLTags(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.course_duration_notes")),
		CreditPoints:              utils.StringToInt(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.credit_points")),
		CricosCode:                utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.cricos_code"),
		DoubleDegrees:             utils.RemoveHTMLTags(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.double_degrees")),
		EnglishLanguage:           utils.RemoveHTMLTags(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.english_language")),
		FullTimeDuration:          extractFullTimeDurations(rawJSON),
		IBEnglish:                 utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.ib_english"),
		IBMaths:                   utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.ib_maths"),
		MaximumDuration:           utils.StringToInt(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.maximum_duration")),
		LearningOutcomes:          common.LearningOutcomes(rawJSON, "props.pageProps.pageContent.learning_outcomes"),
		CurriculumStructure:       curriculum,
		CurriculumError:           curriculumError,
	}

	log.Success("[COURSE SCRAPER] Extraction complete.")
	return courseScraperData, nil
}

// ExtractAwardTitles parses the JSON input and extracts award titles into a slice of strings.
// It navigates to the specified path in the JSON and extracts the "award_title" values.
func extractAwardTitles(data map[string]interface{}) []string {
	path := "props.pageProps.pageContent.award_titles"
	// Extract the array using NavigateToArray
	arrExtract := utils.GetTypedValue[[]map[string]interface{}](data, path)
	if len(arrExtract) == 0 {
		log.Errorf("No data found at path: %s", path)
		return nil
	}

	var awardTitles []string
	// Iterate through the array
	for _, item := range arrExtract {
		if title, ok := item["award_title"].(string); ok {
			awardTitles = append(awardTitles, title)
		} else {
			log.Errorf("award_title is missing or not a string in item: %v", item)
		}
	}

	return awardTitles
}

// ExtractFullTimeDurations parses the JSON input and extracts full-time durations into a slice of strings.
// It navigates to the specified path in the JSON and extracts the "duration_display" values for "Full time" entries.
func extractFullTimeDurations(data map[string]interface{}) []string {
	path := "props.pageProps.pageContent.full_time_duration"
	// Extract the array using NavigateToArray
	arrExtract := utils.GetTypedValue[[]map[string]interface{}](data, path)
	if len(arrExtract) == 0 {
		log.Errorf("No data found at path: %s", path)
		return nil
	}

	var durations []string
	// Iterate through the array
	for _, item := range arrExtract {
		if typeInfo, ok := item["type"].(map[string]interface{}); ok {
			if typeValue, ok := typeInfo["value"].(string); ok && typeValue == "Full time" {
				if durationDisplay, ok := item["duration_display"].(string); ok {
					durations = append(durations, durationDisplay)
				} else {
					log.Errorf("duration_display is missing or not a string in item: %v", item)
				}
			}
		}
	}

	return durations
}
