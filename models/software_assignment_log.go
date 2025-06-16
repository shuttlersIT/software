package models

import "time"

// SoftwareAssignmentLog represents an audit trail for software assignment and unassignment actions.
// swagger:model
type SoftwareAssignmentLog struct {
	ID         uint       `json:"id" example:"1"`
	StaffID    uint       `json:"staff_id" example:"3"`
	Staff      StaffPlain `json:"staff" gorm:"foreignKey:StaffID"`
	SoftwareID uint       `json:"software_id" example:"7"`
	Software   Software   `json:"software" gorm:"foreignKey:SoftwareID"`
	Action     string     `json:"action" example:"Assigned"` // Assigned | Unassigned
	ChangedBy  uint       `json:"changed_by" example:"2"`
	ChangedAt  time.Time  `json:"changed_at" example:"2025-06-11T15:04:05Z"`
	UpdatedAt  time.Time  `json:"updated_at" example:"2025-06-11T15:05:00Z"`
}

func (SoftwareAssignmentLog) TableName() string {
	return "software_assignment_logs"
}

type NewSoftwareAssignmentLog struct {
	ID         uint      `json:"id"`
	StaffID    uint      `json:"staff_id"`
	SoftwareID uint      `json:"software_id"`
	Action     string    `json:"action"`     // e.g., "Assigned", "Unassigned"
	ChangedBy  uint      `json:"changed_by"` // User ID who performed the action
	ChangedAt  time.Time `json:"changed_at"` // Optional: will be set automatically if omitted
	UpdatedAt  time.Time `json:"updated_at"` // Automatically set
}

func (NewSoftwareAssignmentLog) TableName() string {
	return "software_assignment_logs"
}
