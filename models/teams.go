package models

import "time"

// Team represents a team under a department.
// swagger:model
type Team struct {
	ID           uint       `json:"id" gorm:"primaryKey" example:"1"`
	Name         string     `json:"name" example:"Frontend"`
	DepartmentID uint       `json:"department_id" example:"2"`
	Department   Department `json:"department" gorm:"foreignKey:ID"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Team) TableName() string {
	return "teams"
}

// Team represents a team under a department.
// swagger:model
type TeamPlain struct {
	ID           uint      `json:"id" gorm:"primaryKey" example:"1"`
	Name         string    `json:"name" example:"Frontend"`
	DepartmentID uint      `json:"department_id" example:"2"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (TeamPlain) TableName() string {
	return "teams"
}
