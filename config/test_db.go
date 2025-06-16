package config

import (
	"log"
	"software_management/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

func InitTestDB() {
	var err error
	TestDB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test DB: %v", err)
	}

	// Migrate all models
	err = TestDB.AutoMigrate(
		&models.Department{}, &models.Team{}, &models.Staff{}, &models.Software{},
		&models.AssignedSoftware{}, &models.SoftwareAssignmentLog{},
		&models.SoftwareDepartmentMatch{}, &models.SoftwareTeamMatch{}, &models.SoftwareOrganizationMatch{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate test DB: %v", err)
	}
}
