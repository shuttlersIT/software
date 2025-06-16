package models

import "time"

// SoftwareAssignment defines a scoped software assignment (e.g. team, department).
// swagger:model
type SoftwareAssignment struct {
	ID         uint      `json:"id" gorm:"primaryKey" example:"1"`
	SoftwareID uint      `json:"software_id" example:"2"`
	ScopeType  string    `json:"scope_type" example:"Team"` // "Department", "Team", "Staff"
	ScopeID    uint      `json:"scope_id" example:"3"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (SoftwareAssignment) TableName() string {
	return "software_assignments"
}
