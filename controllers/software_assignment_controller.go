package controllers

import (
	"net/http"
	"strconv"
	"time"

	"software_management/config"
	"software_management/models"

	"github.com/gin-gonic/gin"
)

// GetSoftwareAssignments godoc
// @Summary Get all software assignment rules
// @Tags Software Assignment Rules
// @Produce json
// @Success 200 {array} models.SoftwareAssignment
// @Failure 500 {object} models.APIResponse
// @Router /api/software-assignments/plain [get]
func GetSoftwareAssignments(c *gin.Context) {
	var assignments []models.SoftwareAssignment
	if err := config.DB.Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assignments)
}

// GetSoftwareAssignmentsWithDetail godoc
// @Summary Get software assignment rules with software name and filters
// @Tags Software Assignment Rules
// @Produce json
// @Param search query string false "Search by software name"
// @Param scope_type query string false "Filter by scope type (department/team/organization)"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} models.APIResponse
// @Router /api/software-assignments [get]
func GetSoftwareAssignmentsWithDetail(c *gin.Context) {
	search := c.Query("search")
	scopeType := c.Query("scope_type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	type Result struct {
		ID        uint      `json:"id"`
		Software  string    `json:"software"`
		ScopeType string    `json:"scope_type"`
		ScopeID   uint      `json:"scope_id"`
		CreatedAt time.Time `json:"created_at"`
	}

	var results []Result

	query := config.DB.Table("software_assignments").
		Select(`
			software_assignments.id,
			software.name AS software,
			software_assignments.scope_type,
			software_assignments.scope_id,
			software_assignments.created_at
		`).
		Joins("JOIN software ON software.id = software_assignments.software_id")

	if scopeType != "" {
		query = query.Where("software_assignments.scope_type = ?", scopeType)
	}
	if search != "" {
		query = query.Where("software.name LIKE ?", "%"+search+"%")
	}
	if startDate != "" && endDate != "" {
		query = query.Where("software_assignments.created_at BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.Order("software_assignments.created_at DESC").
		Limit(pageSize).Offset(offset).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetSoftwareAssignmentByID godoc
// @Summary Get a specific software assignment rule by ID
// @Tags Software Assignment Rules
// @Param id path int true "Software Assignment ID"
// @Produce json
// @Success 200 {object} models.SoftwareAssignment
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-assignments/{id} [get]
func GetSoftwareAssignmentByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var assignment models.SoftwareAssignment
	if err := config.DB.First(&assignment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Software assignment not found"})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

// CreateSoftwareAssignment godoc
// @Summary Create a new software assignment rule
// @Tags Software Assignment Rules
// @Accept json
// @Produce json
// @Param assignment body models.SoftwareAssignment true "Software assignment rule"
// @Success 201 {object} models.SoftwareAssignment
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-assignments [post]
func CreateSoftwareAssignment(c *gin.Context) {
	var assignment models.SoftwareAssignment
	if err := c.ShouldBindJSON(&assignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, assignment)
}

// UpdateSoftwareAssignment godoc
// @Summary Update an existing software assignment rule
// @Tags Software Assignment Rules
// @Accept json
// @Produce json
// @Param id path int true "Assignment Rule ID"
// @Param assignment body models.SoftwareAssignment true "Updated software assignment rule"
// @Success 200 {object} models.SoftwareAssignment
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/software-assignments/{id} [put]
func UpdateSoftwareAssignment(c *gin.Context) {
	var assignment models.SoftwareAssignment
	id := c.Param("id")

	if err := config.DB.First(&assignment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	if err := c.ShouldBindJSON(&assignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&assignment)
	c.JSON(http.StatusOK, assignment)
}

// DeleteSoftwareAssignmentWithForce godoc
// @Summary Force delete a software assignment rule (no unassign logging)
// @Tags Software Assignment Rules
// @Produce json
// @Param id path int true "Assignment Rule ID"
// @Success 204 "No Content - Delete Success"
// @Failure 500 {object} models.APIResponse
// @Router /api/software-assignments/{id}/force [delete]
func DeleteSoftwareAssignmentWithForce(c *gin.Context) {
	var assignment models.SoftwareAssignment
	id := c.Param("id")

	if err := config.DB.Delete(&assignment, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// DeleteSoftwareAssignment godoc
// @Summary Delete a software assignment rule and revoke related assignments with logging
// @Tags Software Assignment Rules
// @Produce json
// @Param id path int true "Assignment Rule ID"
// @Success 200 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-assignments/{id} [delete]
func DeleteSoftwareAssignment(c *gin.Context) {
	id := c.Param("id")
	var assignment models.SoftwareAssignment

	if err := config.DB.First(&assignment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	var assigned []models.AssignedSoftware
	if err := config.DB.Where("software_id = ?", assignment.SoftwareID).Find(&assigned).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find assigned software"})
		return
	}

	// Revoke each assignment and log
	for _, a := range assigned {
		log := models.SoftwareAssignmentLog{
			StaffID:    a.StaffID,
			SoftwareID: a.SoftwareID,
			Action:     "Unassigned (Rule Deleted)",
			ChangedBy:  0, // Placeholder for current user
			ChangedAt:  time.Now(),
		}
		config.DB.Create(&log)
		config.DB.Delete(&a)
	}

	// Delete the assignment rule
	if err := config.DB.Delete(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete assignment rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment rule and related assignments deleted"})
}
