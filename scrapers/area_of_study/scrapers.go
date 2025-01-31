package area_of_study

import (
	"handbook-scraper/scrapers/common"
	"handbook-scraper/utils"
	"handbook-scraper/utils/log"
)

// Scrape extracts the relevant Area of Study data from the raw JSON.
// It parses the JSON and populates the AosScraperData field with the extracted information.
func Scrape(rawJSON map[string]interface{}, baseURL string) (AosData, error) {
	log.Infof("[AREA OF STUDY SCRAPER] Extracting data...")

	curriculum, errCurriculum := common.ParseCurriculum(rawJSON)
	var curriculumError bool
	if errCurriculum != nil {
		log.Errorf("aos scraper: Error parsing curriculum: %v", errCurriculum)
		curriculumError = true
	} else {
		curriculumError = false
	}

	aosScraperData := AosData{
		CommonScraperData: common.CommonScraperData{
			Link:             baseURL,
			Faculty:          utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.school.value"),
			Code:             utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.code"),
			Title:            utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.title"),
			SearchTitle:      utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.search_title"),
			CurrentYear:      utils.StringToInt(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.implementation_year")),
			AcademicItemType: "area_of_study",
		},
		SpecificAosType:      utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.academic_item_type"),
		CreditPoints:         utils.GetTypedValue[int](rawJSON, "props.pageProps.pageContent.credit_points"),
		CurriculumStructure:  curriculum,
		CurriculumError:      curriculumError,
		HandbookDescription:  utils.RemoveHTMLTags(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.handbook_description")),
		InherentRequirements: utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.inherent_requirements"),
		LearningOutcomes:     common.LearningOutcomes(rawJSON, "props.pageProps.pageContent.learning_outcomes"),
		SpecialStatements:    utils.RemoveHTMLTags(utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.special_statements")),
		UndergradPostgrad:    utils.GetTypedValue[string](rawJSON, "props.pageProps.pageContent.undergrad_postgrad.value"),
	}

	log.Success("[AOS SCRAPER] Extraction complete.")
	return aosScraperData, nil
}
