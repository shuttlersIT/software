package utils

import (
	"time"

	"software_management/config"
	"software_management/models"
)

// Constants for assignment sources
const (
	SourceDepartment   = "department"
	SourceTeam         = "team"
	SourceOrganization = "organization"
)

// AutoAssignSoftwareToStaff assigns software based on department, team, and org matches.
func AutoAssignSoftwareToStaff(staffID, departmentID, teamID uint) error {
	now := time.Now()

	// Fetch already assigned software IDs
	var existingAssignments []uint
	if err := config.DB.Model(&models.AssignedSoftware{}).
		Where("staff_id = ?", staffID).
		Pluck("software_id", &existingAssignments).Error; err != nil {
		return err
	}
	existing := make(map[uint]bool)
	for _, sid := range existingAssignments {
		existing[sid] = true
	}

	// Assign helper
	assignIfNotExists := func(softwareID uint, source string) {
		if !existing[softwareID] {
			config.DB.Create(&models.AssignedSoftware{
				StaffID:    staffID,
				SoftwareID: softwareID,
				AssignedAt: now,
				Source:     source,
			})
			logAssignmentChange(staffID, softwareID, "Assigned")
		}
	}

	// Department matches
	var deptMatches []models.SoftwareDepartmentMatch
	if err := config.DB.Where("department_id = ?", departmentID).Find(&deptMatches).Error; err != nil {
		return err
	}
	for _, match := range deptMatches {
		assignIfNotExists(match.SoftwareID, SourceDepartment)
	}

	// Team matches
	var teamMatches []models.SoftwareTeamMatch
	if err := config.DB.Where("team_id = ?", teamID).Find(&teamMatches).Error; err != nil {
		return err
	}
	for _, match := range teamMatches {
		assignIfNotExists(match.SoftwareID, SourceTeam)
	}

	// Org matches
	var orgMatches []models.SoftwareOrganizationMatch
	if err := config.DB.Find(&orgMatches).Error; err != nil {
		return err
	}
	for _, match := range orgMatches {
		assignIfNotExists(match.SoftwareID, SourceOrganization)
	}

	return nil
}

// SyncSoftwareAssignmentsForStaff is called when staff changes team or department
func SyncSoftwareAssignmentsForStaff(staffID, oldDeptID, oldTeamID, newDeptID, newTeamID uint) error {
	// Revoke old department/team software
	revokeSoftwareFromSource(staffID, oldDeptID, SourceDepartment)
	revokeSoftwareFromSource(staffID, oldTeamID, SourceTeam)

	// Reassign for new department, team, and org
	return AutoAssignSoftwareToStaff(staffID, newDeptID, newTeamID)
}

// RevokeSoftwareAssignmentsForStaff is called during offboarding
func RevokeSoftwareAssignmentsForStaff(staffID uint) error {
	var assignments []models.AssignedSoftware
	if err := config.DB.Where("staff_id = ?", staffID).Find(&assignments).Error; err != nil {
		return err
	}

	for _, assignment := range assignments {
		if isAutoAssigned(assignment.Source) {
			config.DB.Delete(&assignment)
			logAssignmentChange(assignment.StaffID, assignment.SoftwareID, "Unassigned")
		}
	}

	return nil
}

// AutoAssignSoftwareToStaffByUnit assigns a single software to a list of staff
func AutoAssignSoftwareToStaffByUnit(softwareID uint, staffList []models.Staff, source string) {
	now := time.Now()
	for _, staff := range staffList {
		var existing models.AssignedSoftware
		if err := config.DB.
			Where("staff_id = ? AND software_id = ?", staff.ID, softwareID).
			First(&existing).Error; err != nil {
			config.DB.Create(&models.AssignedSoftware{
				StaffID:    staff.ID,
				SoftwareID: softwareID,
				Source:     source,
				AssignedAt: now,
			})
			logAssignmentChange(staff.ID, softwareID, "Assigned")
		}
	}
}

// AutoRevokeSoftwareFromStaff revokes software from a list of staff (only if auto-assigned)
func AutoRevokeSoftwareFromStaff(softwareID uint, staffList []models.Staff, source string) {
	for _, staff := range staffList {
		result := config.DB.Where("staff_id = ? AND software_id = ? AND source = ?", staff.ID, softwareID, source).
			Delete(&models.AssignedSoftware{})

		if result.RowsAffected > 0 {
			logAssignmentChange(staff.ID, softwareID, "Unassigned")
		}
	}
}

// revokeSoftwareFromSource is used internally for department/team match removal
func revokeSoftwareFromSource(staffID, id uint, source string) {
	var softwareIDs []uint

	switch source {
	case SourceDepartment:
		config.DB.Model(&models.SoftwareDepartmentMatch{}).
			Where("department_id = ?", id).
			Pluck("software_id", &softwareIDs)
	case SourceTeam:
		config.DB.Model(&models.SoftwareTeamMatch{}).
			Where("team_id = ?", id).
			Pluck("software_id", &softwareIDs)
	}

	for _, sid := range softwareIDs {
		result := config.DB.Where("staff_id = ? AND software_id = ? AND source = ?", staffID, sid, source).
			Delete(&models.AssignedSoftware{})
		if result.RowsAffected > 0 {
			logAssignmentChange(staffID, sid, "Unassigned")
		}
	}
}

// isAutoAssigned checks if the assignment source is one of the match-based sources
func isAutoAssigned(source string) bool {
	return source == SourceDepartment || source == SourceTeam || source == SourceOrganization
}

// logAssignmentChange writes an assignment or unassignment log
func logAssignmentChange(staffID, softwareID uint, action string) {
	now := time.Now()
	config.DB.Create(&models.SoftwareAssignmentLog{
		StaffID:    staffID,
		SoftwareID: softwareID,
		Action:     action,
		ChangedBy:  0, // System action
		ChangedAt:  now,
		UpdatedAt:  now,
	})
}
