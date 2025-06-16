package controllers

import (
	"net/http"
	"software_management/config"
	"software_management/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetAssignedSoftware godoc
// @Summary Get all assigned software records
// @Tags Assigned Software
// @Produce json
// @Success 200 {array} models.AssignedSoftware
// @Failure 500 {object} models.APIResponse
// @Router /api/assigned-software [get]
func GetAssignedSoftware(c *gin.Context) {
	var assigned []models.AssignedSoftware
	if err := config.DB.Find(&assigned).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assigned)
}

// GetAssignedSoftwareByID godoc
// @Summary Get a specific assigned software record by ID
// @Tags Assigned Software
// @Param id path int true "Assigned Software ID"
// @Produce json
// @Success 200 {object} models.AssignedSoftware
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/assigned-software/{id} [get]
func GetAssignedSoftwareByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var assigned models.AssignedSoftware
	if err := config.DB.First(&assigned, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assigned software not found"})
		return
	}

	c.JSON(http.StatusOK, assigned)
}

// CreateAssignedSoftware godoc
// @Summary Assign software to a staff member
// @Tags Assigned Software
// @Accept json
// @Produce json
// @Param assignment body models.AssignedSoftware true "Software assignment payload"
// @Success 201 {object} models.AssignedSoftware
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/assign-software [post]
func CreateAssignedSoftware(c *gin.Context) {
	var record models.AssignedSoftware
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, record)
}

// UpdateAssignedSoftware godoc
// @Summary Update an assigned software record
// @Tags Assigned Software
// @Accept json
// @Produce json
// @Param id path int true "Assignment ID"
// @Param assignment body models.AssignedSoftware true "Updated assignment"
// @Success 200 {object} models.AssignedSoftware
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/assigned-software/{id} [put]
func UpdateAssignedSoftware(c *gin.Context) {
	var record models.AssignedSoftware
	id := c.Param("id")

	if err := config.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&record)
	c.JSON(http.StatusOK, record)
}

// DeleteAssignedSoftware godoc
// @Summary Delete an assigned software record
// @Tags Assigned Software
// @Produce json
// @Param id path int true "Assignment ID"
// @Success 204 "No Content - Deletion Success"
// @Failure 500 {object} models.APIResponse
// @Router /api/assigned-software/{id}/force [delete]
func DeleteAssignedSoftware(c *gin.Context) {
	var record models.AssignedSoftware
	id := c.Param("id")

	if err := config.DB.Delete(&record, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// DeleteAssignedSoftwareWithLogging godoc
// @Summary Delete assigned software and log the unassignment
// @Tags Assigned Software
// @Produce json
// @Param id path int true "Assignment ID"
// @Success 200 "No Content - Deletion Success"
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/assigned-software/{id} [delete]
func DeleteAssignedSoftwareWithLogging(c *gin.Context) {
	id := c.Param("id")

	var assignment models.AssignedSoftware
	if err := config.DB.First(&assignment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	if err := config.DB.Delete(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete assignment"})
		return
	}

	// Log unassignment
	log := models.SoftwareAssignmentLog{
		StaffID:    assignment.StaffID,
		SoftwareID: assignment.SoftwareID,
		Action:     "Unassigned",
		ChangedBy:  0, // Set to current user ID if available
		ChangedAt:  time.Now(),
	}
	config.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{"message": "Assignment removed"})
}
