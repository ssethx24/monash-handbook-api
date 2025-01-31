package units

import (
	"fmt"
	"handbook-scraper/scrapers/common"
)

// CheckRequisites checks if a student meets the prerequisites and prohibitions for a given unit.
// It takes a UnitData struct and a slice of completed units as input.
// It returns true if all prerequisites are met and no prohibitions are violated, false otherwise,
// a message explaining why the prereqs are not met or prohibitions are violated, and an error if any occurs.
func CheckRequisites(unitData UnitData, completedUnits []common.Unit) (bool, []string, error) {
	if len(unitData.Requisites) == 0 {
		// If there are no requisites, the student automatically meets the requirements
		return true, []string{}, nil
	}

	var unmetRequisites []string

	// Iterate through each requisite
	for _, requisite := range unitData.Requisites {
		if requisite.RequisiteType == "Prerequisite" {
			met, messages, err := checkContainer(requisite.Containers, completedUnits, false)
			if err != nil {
				return false, []string{}, fmt.Errorf("error checking prerequisite container: %w", err)
			}
			if !met {
				unmetRequisites = append(unmetRequisites, messages...)
			}
		} else if requisite.RequisiteType == "Prohibition" {
			met, messages, err := checkContainer(requisite.Containers, completedUnits, true)
			if err != nil {
				return false, []string{}, fmt.Errorf("error checking prohibition container: %w", err)
			}
			if !met {
				unmetRequisites = append(unmetRequisites, messages...)
			}
		}
	}

	if len(unmetRequisites) > 0 {
		return false, unmetRequisites, nil
	}

	// If all prerequisites are met and no prohibitions are violated, return true
	return true, []string{}, nil
}

// checkContainer recursively checks if a container's requirements are met.
// It takes a slice of CompressedContainer, a slice of completed units, and a boolean indicating whether to check for prohibitions as input.
// It returns true if the requirements of all containers are met (or no prohibitions are violated), false otherwise,
// a message explaining why the prereqs are not met or prohibitions are violated, and an error if any occurs.
func checkContainer(containers []CompressedContainer, completedUnits []common.Unit, isProhibition bool) (bool, []string, error) {
	if len(containers) == 0 {
		return true, []string{}, nil // No containers, consider it met
	}

	var unmetRequisites []string

	for _, container := range containers {
		met, messages, err := checkContainerLogic(container, completedUnits, isProhibition)
		if err != nil {
			return false, []string{}, fmt.Errorf("error checking container logic: %w", err)
		}
		if !met {
			unmetRequisites = append(unmetRequisites, messages...)
		}
	}

	if len(unmetRequisites) > 0 {
		return false, unmetRequisites, nil
	}
	return true, []string{}, nil // All containers met or no prohibitions violated
}

// checkContainerLogic checks if a single container's logic is met.
// It takes a CompressedContainer, a slice of completed units, and a boolean indicating whether to check for prohibitions as input.
// It returns true if the container's logic is met (or no prohibitions are violated), false otherwise,
// a message explaining why the prereqs are not met or prohibitions are violated, and an error if any occurs.
func checkContainerLogic(container CompressedContainer, completedUnits []common.Unit, isProhibition bool) (bool, []string, error) {
	if container.Relationship == "AND" {
		return checkAndLogic(container, completedUnits, isProhibition)
	} else if container.Relationship == "OR" {
		return checkOrLogic(container, completedUnits, isProhibition)
	} else {
		return false, []string{}, fmt.Errorf("unknown relationship type: %s", container.Relationship)
	}
}

// checkAndLogic checks if all units in an AND container are met or no prohibitions are violated.
// It takes a CompressedContainer, a slice of completed units, and a boolean indicating whether to check for prohibitions as input.
// It returns true if all units and subcontainers are met (or no prohibitions are violated), false otherwise,
// a message explaining why the prereqs are not met or prohibitions are violated, and an error if any occurs.
func checkAndLogic(container CompressedContainer, completedUnits []common.Unit, isProhibition bool) (bool, []string, error) {
	if len(container.Units) == 0 && len(container.Containers) == 0 {
		return true, []string{}, nil // No units or containers, consider it met
	}

	var unmetRequisites []string
	var mentionedUnits []string

	if !isProhibition {
		for _, unit := range container.Units {
			if !isUnitCompleted(unit, completedUnits) {
				mentionedUnits = append(mentionedUnits, unit.UnitCode)
			}
		}

		for _, subContainer := range container.Containers {
			met, messages, err := checkContainer([]CompressedContainer{subContainer}, completedUnits, isProhibition)
			if err != nil {
				return false, []string{}, fmt.Errorf("error checking subcontainer: %w", err)
			}
			if !met {
				unmetRequisites = append(unmetRequisites, messages...)
			}
		}

		if len(mentionedUnits) > 0 {
			var message string
			message += "Requires: "
			for i, unit := range mentionedUnits {
				message += unit
				if i < len(mentionedUnits)-1 {
					message += " and "
				}
			}
			unmetRequisites = append(unmetRequisites, message)
		}
	} else {
		for _, unit := range container.Units {
			if isUnitCompleted(unit, completedUnits) {
				mentionedUnits = append(mentionedUnits, unit.UnitCode)
			}
		}

		for _, subContainer := range container.Containers {
			met, messages, err := checkContainer([]CompressedContainer{subContainer}, completedUnits, isProhibition)
			if err != nil {
				return false, []string{}, fmt.Errorf("error checking subcontainer: %w", err)
			}
			if !met {
				unmetRequisites = append(unmetRequisites, messages...)
			}
		}

		if len(mentionedUnits) > 0 {
			var message string
			message += "Prohibited by: "
			for i, unit := range mentionedUnits {
				message += unit
				if i < len(mentionedUnits)-1 {
					message += " and "
				}
			}
			unmetRequisites = append(unmetRequisites, message)
		}
	}

	if len(unmetRequisites) > 0 {
		return false, unmetRequisites, nil
	}

	return true, []string{}, nil // All units and subcontainers met or no prohibitions violated
}

// checkOrLogic checks if at least one unit in an OR container is met or no prohibitions are violated.
// It takes a CompressedContainer, a slice of completed units, and a boolean indicating whether to check for prohibitions as input.
// It returns true if at least one unit or subcontainer is met (or no prohibitions are violated), false otherwise,
// a message explaining why the prereqs are not met or prohibitions are violated, and an error if any occurs.
func checkOrLogic(container CompressedContainer, completedUnits []common.Unit, isProhibition bool) (bool, []string, error) {
	if len(container.Units) == 0 && len(container.Containers) == 0 {
		return true, []string{}, nil // No units or containers, consider it met
	}

	var mentionedUnits []string

	if !isProhibition {
		for _, unit := range container.Units {
			if isUnitCompleted(unit, completedUnits) {
				return true, []string{}, nil // If any unit is completed, return true
			}
		}

		for _, subContainer := range container.Containers {
			met, _, err := checkContainer([]CompressedContainer{subContainer}, completedUnits, isProhibition)
			if err != nil {
				return false, []string{}, fmt.Errorf("error checking subcontainer: %w", err)
			}
			if met {
				return true, []string{}, nil // If any subcontainer is met, return true
			}
		}

		// Collect unmet units
		for _, unit := range container.Units {
			if !isUnitCompleted(unit, completedUnits) {
				mentionedUnits = append(mentionedUnits, unit.UnitCode)
			}
		}
	} else {
		for _, unit := range container.Units {
			if isUnitCompleted(unit, completedUnits) {
				mentionedUnits = append(mentionedUnits, unit.UnitCode)
			}
		}

		for _, subContainer := range container.Containers {
			met, _, err := checkContainer([]CompressedContainer{subContainer}, completedUnits, isProhibition)
			if err != nil {
				return false, []string{}, fmt.Errorf("error checking subcontainer: %w", err)
			}
			if met {
				mentionedUnits = append(mentionedUnits, "sub-container")
			}
		}
	}

	// Collect unmet units or prohibitions
	var unmetRequisites []string
	if !isProhibition {
		if len(mentionedUnits) > 0 {
			var message string
			message += "Requires one of: "
			for i, unit := range mentionedUnits {
				message += unit
				if i < len(mentionedUnits)-1 {
					message += " or "
				}
			}
			unmetRequisites = append(unmetRequisites, message)
		}
	} else {
		if len(mentionedUnits) > 0 {
			var message string
			message += "Prohibited by one of: "
			for i, unit := range mentionedUnits {
				message += unit
				if i < len(mentionedUnits)-1 {
					message += " or "
				}
			}
			unmetRequisites = append(unmetRequisites, message)
		}
	}

	for _, subContainer := range container.Containers {
		_, messages, _ := checkContainer([]CompressedContainer{subContainer}, completedUnits, isProhibition)
		unmetRequisites = append(unmetRequisites, messages...)
	}

	if len(unmetRequisites) > 0 {
		return false, unmetRequisites, nil
	}

	return isProhibition, []string{}, nil // No units or subcontainers met or all prohibitions violated
}

// isUnitCompleted checks if a unit is in the list of completed units.
// It takes a CompressedUnit and a slice of completed units as input.
// It returns true if the unit is in the list of completed units, false otherwise.
func isUnitCompleted(unit CompressedUnit, completedUnits []common.Unit) bool {
	for _, completed := range completedUnits {
		if completed.Code == unit.UnitCode {
			return true
		}
	}
	return false
}
