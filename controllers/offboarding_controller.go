package controllers

import (
	"log"
	"net/http"

	"software_management/config"
	"software_management/models"
	"software_management/utils"

	"github.com/gin-gonic/gin"
)

// OffboardStaff godoc
// @Summary Offboard staff (revoke all software, optionally mark as inactive)
// @Tags Staff
// @Produce json
// @Param id path int true "Staff ID"
// @Success 200 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/staff/{id}/offboard [post]
func OffboardStaff(c *gin.Context) {
	// OffboardStaff revokes all auto-assigned software and optionally marks staff as inactive
	id := c.Param("id")
	var staff models.Staff
	if err := config.DB.First(&staff, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}

	// Revoke software
	if err := utils.RevokeSoftwareAssignmentsForStaff(staff.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke software"})
		return
	}

	// Optionally mark as inactive (if not already)
	if staff.Status != "inactive" {
		staff.Status = "inactive"
		if err := config.DB.Save(&staff).Error; err != nil {
			log.Println("Failed to mark staff inactive during offboarding:", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff offboarded and all software revoked"})
}
