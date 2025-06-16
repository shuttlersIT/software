package models

import "time"

// Department represents a department in the organization.
// swagger:model
type Department struct {
	ID        uint      `gorm:"primaryKey" json:"id" example:"1"`
	Name      string    `gorm:"unique;not null" json:"name" example:"Engineering"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Department) TableName() string {
	return "departments"
}
