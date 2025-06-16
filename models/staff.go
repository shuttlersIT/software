package models

import "time"

// Staff represents a simple view employee in the organization.
// swagger:model
type StaffPlain struct {
	ID           uint      `gorm:"primaryKey" json:"id" example:"1"`
	FirstName    string    `json:"first_name" example:"John"`
	LastName     string    `json:"last_name" example:"Doe"`
	Email        string    `gorm:"unique" json:"email" example:"john.doe@shuttlers.co"`
	DepartmentID uint      `json:"department_id" example:"2"`
	TeamID       uint      `json:"team_id" example:"5"`
	Status       string    `json:"status"  example:"Active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (StaffPlain) TableName() string {
	return "staff"
}

// StaffWithDetail represents an enriched view of the employee in the organization.
// swagger:model
type Staff struct {
	ID           uint       `gorm:"primaryKey" json:"id" example:"1"`
	FirstName    string     `json:"first_name" example:"John"`
	LastName     string     `json:"last_name" example:"Doe"`
	Email        string     `gorm:"unique" json:"email" example:"john.doe@shuttlers.co"`
	DepartmentID uint       `json:"department_id" example:"2"`
	Department   Department `gorm:"foreignKey:DepartmentID" json:"department"`
	TeamID       uint       `json:"team_id" example:"5"`
	Team         Team       `gorm:"foreignKey:TeamID" json:"team"`
	Status       string     `json:"status"  example:"Active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Staff) TableName() string {
	return "staff"
}
