package models

import "time"

// AssignedSoftware represents the relationship between a staff and an assigned software.
// swagger:model
type AssignedSoftware struct {
	ID         uint      `json:"id" gorm:"primaryKey" example:"1"`
	StaffID    uint      `json:"staff_id" gorm:"index;not null" example:"2"`
	SoftwareID uint      `json:"software_id" gorm:"index;not null" example:"3"`
	Source     string    `json:"source" gorm:"type:enum('manual','organization','department','team');default:'manual'" example:"department"`
	AssignedAt time.Time `json:"assigned_at" gorm:"column:assigned_at;autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (AssignedSoftware) TableName() string {
	return "assigned_software"
}

// AssignedSoftwareDetail represents an enriched view of software assignments.
// swagger:model
type AssignedSoftwareDetail struct {
	ID         uint      `json:"id" example:"1"`
	StaffID    uint      `json:"staff_id" example:"2"`
	SoftwareID uint      `json:"software_id" example:"3"`
	Software   string    `json:"software" example:"ClickUp"`                                                                                 // <-- this is software.name
	Source     string    `json:"source" gorm:"type:enum('manual','organization','department','team');default:'manual'" example:"department"` // "manual", "organization", "department", "team"
	AssignedAt time.Time `json:"assigned_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

/**
type AssignedSoftware struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	StaffID    uint      `json:"staff_id"` // ✅ DO NOT USE `gorm:"-"` here
	SoftwareID uint      `json:"software_id"`
	AssignedAt time.Time `json:"assigned_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (AssignedSoftware) TableName() string {
	return "assigned_software"
}

**/
