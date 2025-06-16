package controllers

import (
	"net/http"
	"software_management/config"
	"software_management/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetDepartments godoc
// @Summary List all departments
// @Tags Departments
// @Produce json
// @Success 200 {array} models.Department
// @Router /api/departments/plain [get]
func GetDepartments(c *gin.Context) {
	var departments []models.Department
	config.DB.Find(&departments)
	c.JSON(http.StatusOK, departments)
}

// GetDepartmentsWithDetail godoc
// @Summary List all departments with filters
// @Tags Departments
// @Produce json
// @Param search query string false "Search by name"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} models.Department
// @Router /api/departments [get]
func GetDepartmentsWithDetail(c *gin.Context) {
	search := c.Query("search")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var departments []models.Department
	query := config.DB.Model(&models.Department{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&departments).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, departments)
}

// CreateDepartment godoc
// @Summary Create a new department
// @Tags Departments
// @Accept json
// @Produce json
// @Param department body models.Department true "Department to create"
// @Success 201 {object} models.Department
// @Failure 400 {object} models.APIResponse
// @Router /api/departments [post]
func CreateDepartment(c *gin.Context) {
	var department models.Department
	if err := c.ShouldBindJSON(&department); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&department)
	c.JSON(http.StatusCreated, department)
}

// UpdateDepartment godoc
// @Summary Update a department by ID
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param department body models.Department true "Updated department"
// @Success 200 {object} models.Department
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/departments/{id} [put]
func UpdateDepartment(c *gin.Context) {
	var department models.Department
	id := c.Param("id")

	if err := config.DB.First(&department, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	if err := c.ShouldBindJSON(&department); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&department)
	c.JSON(http.StatusOK, department)
}

// DeleteDepartment godoc
// @Summary Delete a department by ID
// @Tags Departments
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 "No Content - Deletion Success"
// @Failure 500 {object} models.APIResponse
// @Router /api/departments/{id} [delete]
func DeleteDepartment(c *gin.Context) {
	var department models.Department
	id := c.Param("id")

	if err := config.DB.Delete(&department, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// GetTeamsByDepartment godoc
// @Summary List teams under a specific department
// @Tags Departments
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {array} models.Team
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/departments/{id}/teams [get]
func GetTeamsByDepartment(c *gin.Context) {
	deptID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	var teams []models.Team
	if err := config.DB.Where("department_id = ?", deptID).Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve teams"})
		return
	}

	c.JSON(http.StatusOK, teams)
}
