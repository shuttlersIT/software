package models

import "time"

// SoftwareOrganizationMatch represents software assigned to all staff in the organization.
// swagger:model
type SoftwareOrganizationMatch struct {
	// ID is the unique identifier of the organization match
	// example: 1
	ID uint `json:"id" gorm:"primaryKey"`

	// SoftwareID is the ID of the software assigned to the entire organization
	// example: 2
	SoftwareID uint `json:"software_id"`

	// Software holds detailed software information
	Software Software `json:"software" gorm:"foreignKey:SoftwareID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
