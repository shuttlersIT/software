package controllers

import (
	"fmt"
	"net/http"
	"software_management/config"
	"software_management/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetSoftware godoc
// @Summary Get all software
// @Tags Software
// @Produce json
// @Success 200 {array} models.Software
// @Failure 500 {object} models.APIResponse
// @Router /api/software/plain [get]
func GetAllSoftware(c *gin.Context) {
	var softwares []models.Software
	if err := config.DB.Find(&softwares).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, softwares)
}

// GetSoftwareWithDetail godoc
// @Summary Get paginated software with filtering and search
// @Tags Software
// @Produce json
// @Param search query string false "Search by software name"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} models.Software
// @Failure 500 {object} models.APIResponse
// @Router /api/software [get]
func GetAllSoftwareWithDetail(c *gin.Context) {
	search := c.Query("search")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var software []models.Software
	query := config.DB.Model(&models.Software{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&software).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, software)
}

// GetSoftwareByID godoc
// @Summary Get a software item by ID
// @Tags Software
// @Param id path int true "Software ID"
// @Produce json
// @Success 200 {object} models.Software
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software/{id} [get]
func GetSoftwareByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid software ID"})
		return
	}

	var software models.Software
	if err := config.DB.First(&software, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Software not found"})
		return
	}

	c.JSON(http.StatusOK, software)
}

// CreateSoftware godoc
// @Summary Create new software
// @Tags Software
// @Accept json
// @Produce json
// @Param software body models.Software true "Software to create"
// @Success 201 {object} models.Software
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software [post]
func CreateSoftware(c *gin.Context) {
	var software models.Software
	if err := c.ShouldBindJSON(&software); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate software name
	var existing models.Software
	if err := config.DB.Where("name = ?", software.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Software with this name already exists"})
		return
	}

	if err := config.DB.Create(&software).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, software)
}

// UpdateSoftware godoc
// @Summary Update existing software
// @Tags Software
// @Accept json
// @Produce json
// @Param id path int true "Software ID"
// @Param software body models.Software true "Software fields to update"
// @Success 200 {object} models.Software
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software/{id} [put]
func UpdateSoftware(c *gin.Context) {
	id := c.Param("id")
	var software models.Software
	if err := config.DB.First(&software, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Software not found"})
		return
	}

	var input models.Software
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate name (excluding self)
	var duplicate models.Software
	if err := config.DB.Where("name = ? AND id != ?", input.Name, id).First(&duplicate).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Another software with this name already exists"})
		return
	}

	software.Name = input.Name
	software.Description = input.Description // Add any other fields you use

	if err := config.DB.Save(&software).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, software)
}

// DeleteSoftware godoc
// @Summary Delete software and revoke all its assignments
// @Tags Software
// @Produce json
// @Param id path int true "Software ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software/{id} [delete]
func DeleteSoftware(c *gin.Context) {
	idParam := c.Param("id")
	softwareID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid software ID"})
		return
	}

	var software models.Software
	if err := config.DB.First(&software, softwareID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Software not found"})
		return
	}

	// Step 1: Revoke all assignments
	var assignments []models.AssignedSoftware
	if err := config.DB.Where("software_id = ?", softwareID).Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
		return
	}

	now := time.Now()
	for _, assignment := range assignments {
		// Delete assignment
		if err := config.DB.Delete(&assignment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke assignment"})
			return
		}

		// Log unassignment
		log := models.SoftwareAssignmentLog{
			StaffID:    assignment.StaffID,
			SoftwareID: assignment.SoftwareID,
			Action:     "Unassigned",
			ChangedBy:  0, // system
			ChangedAt:  now,
			UpdatedAt:  now,
		}
		config.DB.Create(&log)
	}

	// Step 2: Delete the software record itself
	if err := config.DB.Delete(&software).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete software"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Software deleted and all assignments revoked"})
}

// GetStaffAssignedToSoftware godoc
// @Summary Get staff assigned to a specific software by ID
// @Tags Software
// @Produce json
// @Param id path int true "Software ID"
// @Success 200 {array} models.Staff
// @Failure 500 {object} models.APIResponse
// @Router /api/software/{software_id}/assigned-staff [get]
func GetStaffAssignedToSoftware(c *gin.Context) {
	softwareID := c.Param("id")
	var assignments []models.AssignedSoftware

	if err := config.DB.Where("software_id = ?", softwareID).Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assignments)
}

// GetStaffAssignedToSoftware godoc
// @Summary Get list of staff assigned to a specific software
// @Tags Software
// @Produce json
// @Param id path int true "Software ID"
// @Param search query string false "Search by staff name"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} models.Staff
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software/{software_id}/assigned-staff [get]
func GetStaffAssignedToSoftwareWithDetails(c *gin.Context) {
	softwareIDParam := c.Param("id")
	softwareID, err := strconv.Atoi(softwareIDParam)
	if err != nil {
		e := fmt.Sprintln("Invalid staff ID = " + softwareIDParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": e})
		return
	}

	search := c.Query("search")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// Step 1: Get staff IDs from assigned_software table
	var staffIDs []uint
	err = config.DB.
		Model(&models.AssignedSoftware{}).
		Where("software_id = ?", softwareID).
		Pluck("staff_id", &staffIDs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(staffIDs) == 0 {
		c.JSON(http.StatusOK, []models.Staff{}) // Return empty list if none found
		return
	}

	// Step 2: Retrieve staff details by IDs
	var staff []models.Staff
	query := config.DB.Preload("Department").Preload("Team").
		Where("id IN ?", staffIDs)

	if search != "" {
		query = query.Where("CONCAT(first_name, ' ', last_name) LIKE ?", "%"+search+"%")
	}
	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	err = query.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&staff).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// GetAllSoftwareNames godoc
// @Summary Get a list of all software names
// @Tags Software
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} models.APIResponse
// @Router /api/software/names [get]
func GetAllSoftwareNames(c *gin.Context) {
	var names []string
	if err := config.DB.Model(&models.Software{}).Pluck("name", &names).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, names)
}

// GetAllSoftwareSummaries godoc
// @Summary Get a list of all software (id and name only)
// @Tags Software
// @Produce json
// @Success 200 {array} models.SoftwareSummary
// @Failure 500 {object} models.APIResponse
// @Router /api/software/summaries [get]
func GetAllSoftwareSummaries(c *gin.Context) {
	type SoftwareSummary struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	var summaries []SoftwareSummary
	if err := config.DB.Model(&models.Software{}).Select("id", "name").Scan(&summaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summaries)
}
