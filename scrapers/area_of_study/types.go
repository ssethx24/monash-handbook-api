package area_of_study

import (
	"handbook-scraper/scrapers/common"
)

// AosData holds the extracted data from the handbook.
type AosData struct {
	common.CommonScraperData `json:"common"`
	SpecificAosType          string                   `json:"specific_aos_type"`     // x.props.pageProps.pageContent.academic_item_type (e.g. Major)
	CreditPoints             int                      `json:"credit_points"`         // x.props.pageProps.pageContent.credit_points
	CurriculumStructure      common.Curriculum        `json:"curriculum_structure"`  // x.props.pageProps.pageContent.curriculumStructure
	CurriculumError          bool                     `json:"curriculum_error"`      // x.props.pageProps.pageContent.curriculumError
	HandbookDescription      string                   `json:"handbook_description"`  // x.props.pageProps.pageContent.handbook_description
	InherentRequirements     string                   `json:"inherent_requirements"` // x.props.pageProps.pageContent.inherent_requirements
	LearningOutcomes         []common.LearningOutcome `json:"learning_outcomes"`     // x.props.pageProps.pageContent.learning_outcomes
	SpecialStatements        string                   `json:"special_statements"`    // x.props.pageProps.pageContent.special_statements
	UndergradPostgrad        string                   `json:"undergrad_postgrad"`    // x.props.pageProps.pageContent.undergrad_postgrad.value
}
