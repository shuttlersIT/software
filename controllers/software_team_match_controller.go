package controllers

import (
	"net/http"
	"software_management/config"
	"software_management/models"
	"software_management/utils"

	"github.com/gin-gonic/gin"
)

// GetSoftwareTeamMatches godoc
// @Summary Get all software team matches
// @Description Retrieves all team-level software matches, optionally filtered by software_id or team_id
// @Tags Software Team Matches
// @Produce json
// @Param software_id query int false "Filter by Software ID"
// @Param team_id query int false "Filter by Team ID"
// @Success 200 {array} models.SoftwareTeamMatch
// @Failure 500 {object} models.APIResponse
// @Router /api/software-team-matches [get]
func GetSoftwareTeamMatches(c *gin.Context) {
	var matches []models.SoftwareTeamMatch
	query := config.DB.Preload("Software").Preload("Team")

	if sw := c.Query("software_id"); sw != "" {
		query = query.Where("software_id = ?", sw)
	}
	if tm := c.Query("team_id"); tm != "" {
		query = query.Where("team_id = ?", tm)
	}

	if err := query.Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matches)
}

// GetSoftwareTeamMatchByID godoc
// @Summary Get a software team match by ID
// @Description Retrieves a specific software team match by its ID
// @Tags Software Team Matches
// @Produce json
// @Param id path int true "Match ID"
// @Success 200 {object} models.SoftwareTeamMatch
// @Failure 404 {object} models.APIResponse
// @Router /api/software-team-matches/{id} [get]
func GetSoftwareTeamMatchByID(c *gin.Context) {
	id := c.Param("id")
	var match models.SoftwareTeamMatch
	if err := config.DB.Preload("Software").Preload("Team").First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}
	c.JSON(http.StatusOK, match)
}

// CreateSoftwareTeamMatch godoc
// @Summary Create a software team match
// @Description Creates a new software-team association
// @Tags Software Team Matches
// @Accept json
// @Produce json
// @Param match body models.SoftwareTeamMatch true "Software Team Match"
// @Success 201 {object} models.SoftwareTeamMatch
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-team-matches [post]
func CreateSoftwareTeamMatch(c *gin.Context) {
	var input models.SoftwareTeamMatch
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.SoftwareTeamMatch
	if err := config.DB.Where("software_id = ? AND team_id = ?", input.SoftwareID, input.TeamID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match already exists"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// CreateSoftwareTeamMatchWithAutoAssignment godoc
// @Summary Create team match with auto-assignment
// @Description Creates a match and assigns the software to all team members
// @Tags Software Team Matches
// @Accept json
// @Produce json
// @Param match body models.SoftwareTeamMatch true "Software Team Match"
// @Success 201 {object} models.SoftwareTeamMatch
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-team-matches/auto-assign [post]
func CreateSoftwareTeamMatchWithAutoAssignment(c *gin.Context) {
	var match models.SoftwareTeamMatch
	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicates
	var existing models.SoftwareTeamMatch
	if err := config.DB.
		Where("software_id = ? AND team_id = ?", match.SoftwareID, match.TeamID).
		First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Match already exists"})
		return
	}

	if err := config.DB.Create(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Assign to all staff in team
	var staffList []models.Staff
	if err := config.DB.Where("team_id = ?", match.TeamID).Find(&staffList).Error; err == nil {
		utils.AutoAssignSoftwareToStaffByUnit(match.SoftwareID, staffList, utils.SourceTeam)
	}

	c.JSON(http.StatusCreated, match)
}

// UpdateSoftwareTeamMatch godoc
// @Summary Update a software team match
// @Description Updates an existing software team match by ID
// @Tags Software Team Matches
// @Accept json
// @Produce json
// @Param id path int true "Match ID"
// @Param match body models.SoftwareTeamMatch true "Updated Software Team Match"
// @Success 200 {object} models.SoftwareTeamMatch
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-team-matches/{id} [put]
func UpdateSoftwareTeamMatch(c *gin.Context) {
	id := c.Param("id")
	var match models.SoftwareTeamMatch

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

// DeleteSoftwareTeamMatch godoc
// @Summary Delete a software team match
// @Description Deletes a software-team association by ID
// @Tags Software Team Matches
// @Produce json
// @Param id path int true "Match ID"
// @Success 200  "No Content - Deletion Success"
// @Failure 500 {object} models.APIResponse
// @Router /api/software-team-matches/{id} [delete]
func DeleteSoftwareTeamMatch(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.SoftwareTeamMatch{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// DeleteSoftwareTeamMatchAndRevokeAssignments godoc
// @Summary Delete match and revoke software from team members
// @Description Deletes the match and removes the software from all staff in the team
// @Tags Software Team Matches
// @Produce json
// @Param id path int true "Match ID"
// @Success 200  "No Content - Deletion Success"
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-team-matches/{id}/revoke [delete]
func DeleteSoftwareTeamMatchAndRevokeAssignments(c *gin.Context) {
	id := c.Param("id")

	var match models.SoftwareTeamMatch
	if err := config.DB.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.Delete(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Revoke software from staff in team
	var staffList []models.Staff
	if err := config.DB.Where("team_id = ?", match.TeamID).Find(&staffList).Error; err == nil {
		utils.AutoRevokeSoftwareFromStaff(match.SoftwareID, staffList, utils.SourceTeam)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match deleted and software revoked from team staff"})
}
