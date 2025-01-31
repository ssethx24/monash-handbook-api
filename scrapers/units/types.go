package units

import (
	"handbook-scraper/scrapers/common"
)

// UnitData holds the extracted data from the handbook.
type UnitData struct {
	common.CommonScraperData `json:"common"`
	Synopsis                 string                   `json:"synopsis"`              //
	UnitLevel                string                   `json:"unit_level"`            //
	WorkloadRequirements     string                   `json:"workload_requirements"` //
	Active                   bool                     `json:"active"`                //
	CreditPoints             int                      `json:"credit_points"`         //
	HandbookVersion          string                   `json:"handbook_version"`      //
	EFTSL                    float32                  `json:"eftsl"`                 //
	HighestSCABand           string                   `json:"highest_sca_band"`      //
	UndergradPostgrad        string                   `json:"undergrad_postgrad"`    //
	AreaOfStudy              []string                 `json:"area_of_study"`         //
	LearningOutcomes         []common.LearningOutcome `json:"learning_outcomes"`     //
	Assessments              []Assessment             `json:"assessments"`           //
	UnitOfferings            []UnitOffering           `json:"unit_offerings"`        //
	LearningActivities       []LearningActivity       `json:"learning_activities"`   //
	Requisites               []CompressedRequisite    `json:"requisites"`            //
	EnrolmentRules           []EnrolmentRule          `json:"enrolment_rules"`       //
}

// Assessment represents a single assessment with relevant fields
type Assessment struct {
	AssessmentName string `json:"assessment_name"`
	AssessmentType struct {
		Label string `json:"label"`
		Value string `json:"value"`
	} `json:"assessment_type"`
	Number      string `json:"number"`
	Weight      string `json:"weight"`
	Description string `json:"description,omitempty"`
}

// UnitOffering represents the structured data for each unit offering
type UnitOffering struct {
	AttendanceMode string `json:"attendance_mode"`
	DisplayName    string `json:"display_name"`
	Location       string `json:"location"`
	Semester       string `json:"semester"`
}

// LearningActivity represents a single learning activity with relevant fields
type LearningActivity struct {
	ActivityType    string `json:"activity_type"`
	DurationDisplay string `json:"duration_display"`
	Offerings       string `json:"offerings_formatted_teaching_activities"`
}

// EnrolmentRule represents a single enrolment rule with the description
type EnrolmentRule struct {
	Description string `json:"description"`
}

// REQUISITE TYPES START HERE

type Requisite struct {
	AcademicItemCode    string        `json:"academic_item_code"`
	AcademicItemVersion string        `json:"academic_item_version_number"`
	Active              string        `json:"active"`
	ClID                ClID          `json:"cl_id"`
	Containers          []Container   `json:"container"`
	Description         string        `json:"description"`
	EndDate             string        `json:"end_date"`
	Order               string        `json:"order"`
	RequisiteClID       string        `json:"requisite_cl_id"`
	RequisiteType       RequisiteType `json:"requisite_type"`
	StartDate           string        `json:"start_date"`
}

type ClID struct {
	ClID  string `json:"cl_id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RequisiteType struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Container struct {
	ClID                 string         `json:"cl_id"`
	Containers           []Container    `json:"containers"`
	ParentConnector      Connector      `json:"parent_connector"`
	ParentContainerTable string         `json:"parent_container_table"`
	ParentRecord         ParentRecord   `json:"parent_record"`
	Relationships        []Relationship `json:"relationships"`
	Title                string         `json:"title"`
}

type Connector struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type ParentRecord struct {
	ClID  string `json:"cl_id"`
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Relationship struct {
	AbbrName                 string           `json:"abbr_name"`
	AbbreviatedNameAndMajor  interface{}      `json:"abbreviated_name_and_major"`
	AcademicItem             AcademicItem     `json:"academic_item"`
	AcademicItemCode         string           `json:"academic_item_code"`
	AcademicItemCreditPoints string           `json:"academic_item_credit_points"`
	AcademicItemName         string           `json:"academic_item_name"`
	AcademicItemType         AcademicItemType `json:"academic_item_type"`
	AcademicItemURL          string           `json:"academic_item_url"`
	AcademicItemVersionName  string           `json:"academic_item_version_name"`
	ClID                     string           `json:"cl_id"`
	Order                    string           `json:"order"`
	ParentRecord             ParentRecord     `json:"parent_record"`
}

type AcademicItem struct {
	ClID  string `json:"cl_id"`
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type AcademicItemType struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Compressed structures
type CompressedRequisite struct {
	RequisiteType string                `json:"requisite_type"` // "Prerequisite" or "Prohibition"
	Containers    []CompressedContainer `json:"containers"`
}

type CompressedContainer struct {
	Relationship string                `json:"relationship"` // "AND" or "OR"
	Units        []CompressedUnit      `json:"units"`
	Containers   []CompressedContainer `json:"containers,omitempty"`
}

type CompressedUnit struct {
	UnitCode   string `json:"unit_code"`
	UnitNumber string `json:"unit_number"`
}
