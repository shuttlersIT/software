package models

import "time"

// Software represents a piece of software managed by the system.
// swagger:model
type Software struct {
	ID          uint      `gorm:"primaryKey" json:"id" example:"1"`
	Name        string    `gorm:"unique" json:"name" example:"ClickUp"`
	Description string    `json:"description" example:"Communication Suite"`
	Type        string    `json:"type" example:"SaaS"` // e.g., License, SaaS, etc.
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Software) TableName() string {
	return "software"
}
