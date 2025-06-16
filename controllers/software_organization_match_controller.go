package controllers

import (
	"net/http"
	"software_management/config"
	"software_management/models"
	"software_management/utils"

	"github.com/gin-gonic/gin"
)

// GetSoftwareOrganizationMatches godoc
// @Summary Get all software organization matches
// @Description Retrieves all organization-level software matches, optionally filtered by software_id
// @Tags Software Organization Matches
// @Produce json
// @Param software_id query int false "Filter by Software ID"
// @Success 200 {array} models.SoftwareOrganizationMatch
// @Failure 500 {object} models.APIResponse
// @Router /api/software-organization-matches [get]
func GetSoftwareOrganizationMatches(c *gin.Context) {
	var matches []models.SoftwareOrganizationMatch
	query := config.DB.Preload("Software")

	if sw := c.Query("software_id"); sw != "" {
		query = query.Where("software_id = ?", sw)
	}

	if err := query.Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matches)
}

// GetSoftwareOrganizationMatchByID godoc
// @Summary Get a software organization match by ID
// @Description Retrieves a single software organization match
// @Tags Software Organization Matches
// @Produce json
// @Param id path int true "Match ID"
// @Success 200 {object} models.SoftwareOrganizationMatch
// @Failure 404 {object} models.APIResponse
// @Router /api/software-organization-matches/{id} [get]
func GetSoftwareOrganizationMatchByID(c *gin.Context) {
	id := c.Param("id")
	var match models.SoftwareOrganizationMatch
	if err := config.DB.Preload("Software").First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}
	c.JSON(http.StatusOK, match)
}

// CreateSoftwareOrganizationMatch godoc
// @Summary Create a software organization match
// @Description Creates a new organization-level match for a software
// @Tags Software Organization Matches
// @Accept json
// @Produce json
// @Param match body models.SoftwareOrganizationMatch true "Software Organization Match"
// @Success 201 {object} models.SoftwareOrganizationMatch
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-organization-matches [post]
func CreateSoftwareOrganizationMatch(c *gin.Context) {
	var input models.SoftwareOrganizationMatch
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.SoftwareOrganizationMatch
	if err := config.DB.Where("software_id = ?", input.SoftwareID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match already exists"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// CreateSoftwareOrganizationMatchWithAutoAssignment godoc
// @Summary Create organization match with auto-assignment
// @Description Creates a new match and automatically assigns the software to all staff
// @Tags Software Organization Matches
// @Accept json
// @Produce json
// @Param match body models.SoftwareOrganizationMatch true "Software Organization Match"
// @Success 201 {object} models.SoftwareOrganizationMatch
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-organization-matches/auto-assign [post]
func CreateSoftwareOrganizationMatchWithAutoAssignment(c *gin.Context) {
	var input models.SoftwareOrganizationMatch
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicates
	var existing models.SoftwareOrganizationMatch
	if err := config.DB.
		Where("software_id = ?", input.SoftwareID).
		First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Match already exists"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Assign software to all staff
	var staffList []models.Staff
	if err := config.DB.Find(&staffList).Error; err == nil {
		utils.AutoAssignSoftwareToStaffByUnit(input.SoftwareID, staffList, utils.SourceOrganization)
	}

	c.JSON(http.StatusCreated, input)
}

// UpdateSoftwareOrganizationMatch godoc
// @Summary Update a software organization match
// @Description Updates an existing software organization match by ID
// @Tags Software Organization Matches
// @Accept json
// @Produce json
// @Param id path int true "Match ID"
// @Param match body models.SoftwareOrganizationMatch true "Updated Software Organization Match"
// @Success 200 {object} models.SoftwareOrganizationMatch
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-organization-matches/{id} [put]
func UpdateSoftwareOrganizationMatch(c *gin.Context) {
	id := c.Param("id")
	var match models.SoftwareOrganizationMatch

	if err := config.DB.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}

// DeleteSoftwareOrganizationMatch godoc
// @Summary Delete a software organization match
// @Description Deletes a software organization match by ID
// @Tags Software Organization Matches
// @Produce json
// @Param id path int true "Match ID"
// @Success 200  "No Content - Deletion Success"
// @Failure 500 {object} models.APIResponse
// @Router /api/software-organization-matches/{id} [delete]
func DeleteSoftwareOrganizationMatch(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.SoftwareOrganizationMatch{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// DeleteSoftwareOrganizationMatchAndRevokeAssignment godoc
// @Summary Delete match and revoke software from staff
// @Description Deletes a software organization match and automatically revokes the software from all staff
// @Tags Software Organization Matches
// @Produce json
// @Param id path int true "Match ID"
// @Success 200  "No Content - Deletion Success"
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-organization-matches/{id}/revoke [delete]
func DeleteSoftwareOrganizationMatchAndRevokeAssignment(c *gin.Context) {
	id := c.Param("id")

	var match models.SoftwareOrganizationMatch
	if err := config.DB.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.Delete(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Revoke software from all staff (auto-assigned only)
	var staffList []models.Staff
	if err := config.DB.Find(&staffList).Error; err == nil {
		utils.AutoRevokeSoftwareFromStaff(match.SoftwareID, staffList, utils.SourceOrganization)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match deleted and software revoked from staff"})
}
