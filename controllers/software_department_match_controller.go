package controllers

import (
	"net/http"
	"software_management/config"
	"software_management/models"
	"software_management/utils"

	"github.com/gin-gonic/gin"
)

// GetSoftwareDepartmentMatches godoc
// @Summary Get all software department matches
// @Description Retrieves all department-level software matches, optionally filtered by software_id or department_id
// @Tags Software Department Matches
// @Produce json
// @Param software_id query int false "Filter by Software ID"
// @Param department_id query int false "Filter by Department ID"
// @Success 200 {array} models.SoftwareDepartmentMatch
// @Failure 500 {object} models.APIResponse
// @Router /api/software-department-matches [get]
func GetSoftwareDepartmentMatches(c *gin.Context) {
	var matches []models.SoftwareDepartmentMatch
	query := config.DB.Preload("Software").Preload("Department")

	if sw := c.Query("software_id"); sw != "" {
		query = query.Where("software_id = ?", sw)
	}
	if dept := c.Query("department_id"); dept != "" {
		query = query.Where("department_id = ?", dept)
	}

	if err := query.Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matches)
}

// GetSoftwareDepartmentMatchByID godoc
// @Summary Get a software department match by ID
// @Description Retrieves a single software department match by ID
// @Tags Software Department Matches
// @Produce json
// @Param id path int true "Match ID"
// @Success 200 {object} models.SoftwareDepartmentMatch
// @Failure 404 {object} models.APIResponse
// @Router /api/software-department-matches/{id} [get]
func GetSoftwareDepartmentMatchByID(c *gin.Context) {
	id := c.Param("id")
	var match models.SoftwareDepartmentMatch
	if err := config.DB.Preload("Software").Preload("Department").First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}
	c.JSON(http.StatusOK, match)
}

// CreateSoftwareDepartmentMatch godoc
// @Summary Create a software department match
// @Description Creates a new department-level software match
// @Tags Software Department Matches
// @Accept json
// @Produce json
// @Param match body models.SoftwareDepartmentMatch true "Software Department Match"
// @Success 201 {object} models.SoftwareDepartmentMatch
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-department-matches [post]
func CreateSoftwareDepartmentMatch(c *gin.Context) {
	var input models.SoftwareDepartmentMatch
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.SoftwareDepartmentMatch
	if err := config.DB.Where("software_id = ? AND department_id = ?", input.SoftwareID, input.DepartmentID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match already exists"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// CreateSoftwareDepartmentMatchWithAutoAssignment godoc
// @Summary Create department match with auto-assignment
// @Description Creates a new match and automatically assigns the software to all staff in the department
// @Tags Software Department Matches
// @Accept json
// @Produce json
// @Param match body models.SoftwareDepartmentMatch true "Software Department Match"
// @Success 201 {object} models.SoftwareDepartmentMatch
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-department-matches/auto-assign [post]
func CreateSoftwareDepartmentMatchWithAutoAssignment(c *gin.Context) {
	var match models.SoftwareDepartmentMatch
	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicates
	var existing models.SoftwareDepartmentMatch
	if err := config.DB.
		Where("software_id = ? AND department_id = ?", match.SoftwareID, match.DepartmentID).
		First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Match already exists"})
		return
	}

	if err := config.DB.Create(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Assign to all staff in department
	var staffList []models.Staff
	if err := config.DB.Where("department_id = ?", match.DepartmentID).Find(&staffList).Error; err == nil {
		utils.AutoAssignSoftwareToStaffByUnit(match.SoftwareID, staffList, utils.SourceDepartment)
	}

	c.JSON(http.StatusCreated, match)
}

// UpdateSoftwareDepartmentMatch godoc
// @Summary Update a software department match
// @Description Updates an existing software department match by ID
// @Tags Software Department Matches
// @Accept json
// @Produce json
// @Param id path int true "Match ID"
// @Param match body models.SoftwareDepartmentMatch true "Updated Software Department Match"
// @Success 200 {object} models.SoftwareDepartmentMatch
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-department-matches/{id} [put]
func UpdateSoftwareDepartmentMatch(c *gin.Context) {
	id := c.Param("id")
	var match models.SoftwareDepartmentMatch
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

// DeleteSoftwareDepartmentMatch godoc
// @Summary Delete a software department match
// @Description Deletes a department-level software match by ID
// @Tags Software Department Matches
// @Produce json
// @Param id path int true "Match ID"
// @Success 200  "No Content - Deletion Success"
// @Failure 500 {object} models.APIResponse
// @Router /api/software-department-matches/{id} [delete]
func DeleteSoftwareDepartmentMatch(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.SoftwareDepartmentMatch{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// DeleteSoftwareDepartmentMatchAndRevokeAssignment godoc
// @Summary Delete match and revoke software from department staff
// @Description Deletes a match and revokes software from all staff in the associated department
// @Tags Software Department Matches
// @Produce json
// @Param id path int true "Match ID"
// @Success 200  "No Content - Deletion Success"
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/software-department-matches/{id}/revoke [delete]
func DeleteSoftwareDepartmentMatchAndRevokeAssignment(c *gin.Context) {
	id := c.Param("id")

	var match models.SoftwareDepartmentMatch
	if err := config.DB.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if err := config.DB.Delete(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Revoke software from staff in department
	var staffList []models.Staff
	if err := config.DB.Where("department_id = ?", match.DepartmentID).Find(&staffList).Error; err == nil {
		utils.AutoRevokeSoftwareFromStaff(match.SoftwareID, staffList, utils.SourceDepartment)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match deleted and software revoked from department staff"})
}
