package controllers

import (
	"net/http"
	"strconv"

	"software_management/config"
	"software_management/models"

	"github.com/gin-gonic/gin"
)

// GetTeams godoc
// @Summary List all teams
// @Tags Teams
// @Produce json
// @Success 200 {array} models.Team
// @Failure 500 {object} models.APIResponse
// @Router /api/teams [get]
func GetTeams(c *gin.Context) {
	var teams []models.Team
	if err := config.DB.Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teams)
}

// GetTeamsWithDetail godoc
// @Summary List all teams with filters and department preload
// @Tags Teams
// @Produce json
// @Param search query string false "Search by name"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} models.Team
// @Failure 500 {object} models.APIResponse
// @Router /api/teams [get]
func GetTeamsWithDetail(c *gin.Context) {
	search := c.Query("search")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	// Same params
	var teams []models.Team
	query := config.DB.Model(&models.Team{}).Preload("Department")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&teams).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teams)
}

// CreateTeam godoc
// @Summary Create a new team
// @Tags Teams
// @Accept json
// @Produce json
// @Param team body models.Team true "Team object"
// @Success 201 {object} models.Team
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/teams [post]
func CreateTeam(c *gin.Context) {
	var team models.TeamPlain
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, team)
}

// UpdateTeam godoc
// @Summary Update a team by ID
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Param team body models.Team true "Updated team"
// @Success 200 {object} models.Team
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/teams/{id} [put]
func UpdateTeam(c *gin.Context) {
	id := c.Param("id")
	var team models.TeamPlain
	if err := config.DB.First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Save(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, team)
}

// DeleteTeam godoc
// @Summary Delete a team by ID
// @Tags Teams
// @Produce json
// @Param id path int true "Team ID"
// @Success 204 "No Content - Deletion Success"
// @Failure 500 {object} models.APIResponse
// @Router /api/teams/{id} [delete]
func DeleteTeam(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Team{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
