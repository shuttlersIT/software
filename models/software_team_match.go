package models

import "time"

// SoftwareTeamMatch represents software assigned to a specific team.
// swagger:model
type SoftwareTeamMatch struct {
	// ID is the unique identifier of the team match
	// example: 1
	ID uint `json:"id" gorm:"primaryKey" example:"1"`

	// SoftwareID is the ID of the software
	// example: 2
	SoftwareID uint `json:"software_id" example:"2"`

	Software Software `json:"software" gorm:"foreignKey:SoftwareID"`

	// TeamID is the ID of the team receiving the software
	// example: 3
	TeamID uint `json:"team_id" example:"3"`

	Team Team `json:"team" gorm:"foreignKey:TeamID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SoftwareTeamMatch) TableName() string {
	return "software_team_matches"
}
