package common

// CommonScraperData represents the common data structure shared by all academic items.
type CommonScraperData struct {
	Link             string `json:"link"`               // https://handbook.monash.edu/current/units/FIT3138
	Faculty          string `json:"faculty"`            // Faculty of Information Technology
	Code             string `json:"code"`               // FIT3138
	Title            string `json:"title"`              // Real time enterprise systems
	SearchTitle      string `json:"search_title"`       // FIT3138 - Real time enterprise systems
	CurrentYear      int    `json:"current_year"`       // 2025
	AcademicItemType string `json:"academic_item_type"` // Unit
}

// LearningOutcome represents the structure of each item in the "unit_learning_outcomes" array.
// It contains the code and description of a learning outcome.
type LearningOutcome struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Unit represents a single academic unit
// It contains the code, name, credit points, description, and URL of the unit.
type Unit struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	CreditPoints int    `json:"credit_points"`
	Description  string `json:"description"`
	URL          string `json:"url"`
}

// StudentProgress represents a student's progress within the curriculum.
// It contains a slice of completed Unit structs and the total credits earned.
type StudentProgress struct {
	CompletedUnits     []Unit `json:"completed_units"`
	TotalCreditsEarned int    `json:"total_credits_earned"`
}
