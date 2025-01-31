package courses

import (
	"handbook-scraper/scrapers/common"
)

// CourseData holds the extracted data from the handbook.
type CourseData struct {
	common.CommonScraperData  `json:"common"`
	ProfessionalAccreditation string                   `json:"professional_accreditation"` // x.props.pageProps.pageContent.Professional_accreditation
	AbbreviatedName           string                   `json:"abbreviated_name"`           // x.props.pageProps.pageContent.abbreviated_name
	Atar                      string                   `json:"atar"`                       // x.props.pageProps.pageContent.atar
	AwardTitles               []string                 `json:"award_titles"`               // x.props.pageProps.pageContent.award_titles
	CourseDuration            string                   `json:"course_duration"`            // x.props.pageProps.pageContent.course_duration_notes
	CreditPoints              int                      `json:"credit_points"`              // x.props.pageProps.pageContent.credit_points
	CricosCode                string                   `json:"cricos_code"`                // x.props.pageProps.pageContent.cricos_code
	DoubleDegrees             string                   `json:"double_degrees"`             // x.props.pageProps.pageContent.double_degrees
	EnglishLanguage           string                   `json:"english_language"`           // x.props.pageProps.pageContent.english_language
	FullTimeDuration          []string                 `json:"full_time_duration"`         // x.props.pageProps.pageContent.full_time_duration
	IBEnglish                 string                   `json:"ib_english"`                 // x.props.pageProps.pageContent.ib_english
	IBMaths                   string                   `json:"ib_maths"`                   // x.props.pageProps.pageContent.ib_maths
	MaximumDuration           int                      `json:"maximum_duration"`           // x.props.pageProps.pageContent.maximum_duration
	CurriculumStructure       common.Curriculum        `json:"curriculum_structure"`       // x.props.pageProps.pageContent.curriculumStructure (complex)
	CurriculumError           bool                     `json:"curriculum_error"`           // x.props.pageProps.pageContent.curriculumError
	LearningOutcomes          []common.LearningOutcome `json:"learning_outcomes"`          // x.props.pageProps.pageContent.learning_outcomes
}
