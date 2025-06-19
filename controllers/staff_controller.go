package controllers

import (
	"fmt"
	"log"
	"net/http"
	"software_management/config"
	"software_management/models"
	"software_management/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetStaff godoc
// @Summary Get all staff
// @Tags Staff
// @Produce json
// @Success 200 {array} models.StaffPlain
// @Failure 500 {object} models.APIResponse
// @Router /api/staff/plain [get]
func GetAllStaff(c *gin.Context) {
	var staff []models.StaffPlain
	if err := config.DB.Find(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, staff)
}

// GetStaffWithDetail godoc
// @Summary Get staff with filters and pagination
// @Tags Staff
// @Produce json
// @Param search query string false "Search keyword (first name, last name, or email)"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} models.Staff
// @Failure 500 {object} models.APIResponse
// @Router /api/staff [get]
func GetStaffWithDetail(c *gin.Context) {
	search := c.Query("search")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var staff []models.Staff
	query := config.DB.Preload("Department").Preload("Team")

	if search != "" {
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// GetStaffByID godoc
// @Summary Get a staff member by ID
// @Description Retrieve a single staff member by their unique ID
// @Tags Staff
// @Param id path int true "Staff ID"
// @Produce json
// @Success 200 {object} models.Staff
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/staff/{id} [get]
func GetStaffByID(c *gin.Context) {
	idParam := c.Param("id")
	fmt.Println(idParam)

	// Parse ID param to uint
	id, err := strconv.Atoi(idParam)
	if err != nil {
		e := fmt.Sprintln("Invalid staff ID = " + idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": e})
		return
	}

	var staff models.Staff
	if err := config.DB.First(&staff, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// CreateStaff godoc
// @Summary Create a new staff
// @Tags Staff
// @Accept json
// @Produce json
// @Param staff body models.Staff true "Staff object"
// @Success 201 {object} models.Staff
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/staff [post]
func CreateStaff(c *gin.Context) {
	var staff models.Staff
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate email
	var existing models.Staff
	if err := config.DB.Where("email = ?", staff.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A staff with this email already exists"})
		return
	}

	if err := config.DB.Create(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, staff)
}

// CreateStaffWithSoftwareMatch godoc
// @Summary Create staff and auto-assign software
// @Tags Staff
// @Accept json
// @Produce json
// @Param staff body models.Staff true "Staff object"
// @Success 201 {object} models.Staff
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/staff/with-software [post]
func CreateStaffWithSoftwareMatch(c *gin.Context) {
	var staff models.Staff
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate email
	var existing models.Staff
	if err := config.DB.Where("email = ?", staff.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A staff with this email already exists"})
		return
	}

	if err := config.DB.Create(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auto assign software after creation
	if err := utils.AutoAssignSoftwareToStaff(staff.ID, staff.DepartmentID, staff.TeamID); err != nil {
		log.Println("Auto-assignment error:", err)
	}

	c.JSON(http.StatusCreated, staff)
}

// UpdateStaff godoc
// @Summary Update staff by ID
// @Tags Staff
// @Accept json
// @Produce json
// @Param id path int true "Staff ID"
// @Param staff body models.Staff true "Updated staff object"
// @Success 200 {object} models.Staff
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/staff/{id} [put]
func UpdateStaff(c *gin.Context) {
	id := c.Param("id")
	var staff models.Staff
	if err := config.DB.First(&staff, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, staff)
}

// UpdateStaffWithSoftwareMatch godoc
// @Summary Update staff and sync software assignments
// @Tags Staff
// @Accept json
// @Produce json
// @Param id path int true "Staff ID"
// @Param staff body models.Staff true "Updated staff object"
// @Success 200 {object} models.Staff
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/staff/{id}/with-software [put]
func UpdateStaffWithSoftwareMatch(c *gin.Context) {
	id := c.Param("id")

	var existing models.Staff
	if err := config.DB.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}

	var input models.Staff
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate email (excluding current staff)
	var duplicate models.Staff
	if err := config.DB.Where("email = ? AND id != ?", input.Email, id).First(&duplicate).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Another staff with this email already exists"})
		return
	}

	// Track changes
	oldDeptID := existing.DepartmentID
	oldTeamID := existing.TeamID
	deptChanged := input.DepartmentID != oldDeptID
	teamChanged := input.TeamID != oldTeamID
	statusChanged := input.Status != existing.Status

	// Apply updates
	existing.FirstName = input.FirstName
	existing.LastName = input.LastName
	existing.Email = input.Email
	existing.DepartmentID = input.DepartmentID
	existing.TeamID = input.TeamID
	existing.Status = input.Status

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Revoke software if staff is now inactive
	if statusChanged && input.Status == "inactive" {
		if err := utils.RevokeSoftwareAssignmentsForStaff(existing.ID); err != nil {
			log.Println("Failed to revoke software for inactive staff:", err)
		}
	}

	// Sync if dept or team changed
	if deptChanged || teamChanged {
		if err := utils.SyncSoftwareAssignmentsForStaff(
			existing.ID,
			oldDeptID,
			oldTeamID,
			input.DepartmentID,
			input.TeamID,
		); err != nil {
			log.Println("Software sync error:", err)
		}
	}

	c.JSON(http.StatusOK, existing)
}

// DeleteStaff godoc
// @Summary Delete staff by ID (also revokes assigned software)
// @Tags Staff
// @Produce json
// @Param id path int true "Staff ID"
// @Success 200 "No Content - Deletion Success"
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/staff/{id} [delete]
func DeleteStaff(c *gin.Context) {
	id := c.Param("id")
	var staff models.Staff
	if err := config.DB.First(&staff, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}

	// Revoke software before deleting
	if err := utils.RevokeSoftwareAssignmentsForStaff(staff.ID); err != nil {
		log.Println("Failed to revoke software before delete:", err)
	}

	// Delete staff
	if err := config.DB.Delete(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff deleted and software unassigned"})
}

// GetSoftwareAssignedToStaff godoc
// @Summary Get software assigned to a specific staff by ID
// @Tags Staff
// @Produce json
// @Param id path int true "Staff ID"
// @Success 200 {array} models.AssignedSoftware
// @Failure 500 {object} models.APIResponse
// @Router /api/staff/{staff_id}/assigned-software [get]
func GetSoftwareAssignedToStaff(c *gin.Context) {
	staffID := c.Param("id")
	var assignments []models.AssignedSoftware

	if err := config.DB.Where("staff_id = ?", staffID).Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assignments)
}

// GetSoftwareAssignedToStaffWithDetails godoc
// @Summary Get detailed software assigned to a specific staff using Software model
// @Tags Staff
// @Produce json
// @Param id path int true "Staff ID"
// @Param search query string false "Search by software name"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} models.Software
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/staff/{staff_id}/assigned-software/detail [get]
func GetSoftwareAssignedToStaffWithDetails(c *gin.Context) {
	staffIDParam := c.Param("id")
	staffID, err := strconv.Atoi(staffIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid staff ID"})
		return
	}

	search := c.Query("search")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var softwareList []models.Software

	query := config.DB.Table("software").
		Select("DISTINCT software.*").
		Joins("JOIN assigned_software ON assigned_software.software_id = software.id").
		Where("assigned_software.staff_id = ?", staffID)

	if search != "" {
		query = query.Where("software.name LIKE ?", "%"+search+"%")
	}
	if startDate != "" && endDate != "" {
		query = query.Where("assigned_software.assigned_at BETWEEN ? AND ?", startDate, endDate)
	}

	err = query.Order("software.name ASC").
		Limit(pageSize).Offset(offset).
		Scan(&softwareList).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, softwareList)
}
