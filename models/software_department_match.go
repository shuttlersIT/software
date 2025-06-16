package models

import "time"

// SoftwareDepartmentMatch represents software assigned to a specific department.
// swagger:model
type SoftwareDepartmentMatch struct {
	// ID is the unique identifier of the department match
	// example: 1
	ID uint `json:"id" gorm:"primaryKey"`

	// SoftwareID is the ID of the software
	// example: 2
	SoftwareID uint `json:"software_id"`

	// DepartmentID is the ID of the department receiving the software
	// example: 3
	DepartmentID uint `json:"department_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SoftwareDepartmentMatch) TableName() string {
	return "software_department_matches"
}
