package controllers

import (
	"net/http"
	"software_management/config"
	"software_management/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetSoftwareAssignmentLogs godoc
// @Summary Get all software assignment logs
// @Description Retrieves all software assignment logs with associated staff and software
// @Tags Software Assignment Logs
// @Produce json
// @Success 200 {array} models.SoftwareAssignmentLog
// @Failure 500 {object} models.APIResponse
// @Router /api/logs [get]
func GetSoftwareAssignmentLogs(c *gin.Context) {
	var logs []models.SoftwareAssignmentLog
	if err := config.DB.Preload("Staff").Preload("Software").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// GetAllAssignmentLogsWithDetails godoc
// @Summary Get filtered and paginated logs
// @Description Retrieves logs with optional filtering by date and search query
// @Tags Software Assignment Logs
// @Produce json
// @Param start query string false "Start date (YYYY-MM-DD)"
// @Param end query string false "End date (YYYY-MM-DD)"
// @Param search query string false "Search by staff or software name"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset results"
// @Success 200 {array} models.SoftwareAssignmentLog
// @Failure 500 {object} models.APIResponse
// @Router /api/logs/details [get]
func GetAllAssignmentLogsWithDetails(c *gin.Context) {
	var logs []models.SoftwareAssignmentLog

	startDate := c.Query("start")
	endDate := c.Query("end")
	search := c.Query("search")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	query := config.DB.Preload("Staff").Preload("Software")

	if startDate != "" && endDate != "" {
		start, err1 := time.Parse("2006-01-02", startDate)
		end, err2 := time.Parse("2006-01-02", endDate)
		if err1 == nil && err2 == nil {
			query = query.Where("changed_at BETWEEN ? AND ?", start, end)
		}
	}

	if search != "" {
		query = query.Joins("JOIN staff ON staff.id = software_assignment_logs.staff_id").
			Joins("JOIN software ON software.id = software_assignment_logs.software_id").
			Where("staff.name LIKE ? OR software.name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.
		Order("changed_at DESC").
		Limit(limit).Offset(offset).
		Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetAssignmentLogByID godoc
// @Summary Get a software assignment log by ID
// @Description Retrieves a single software assignment log by its ID
// @Tags Software Assignment Logs
// @Produce json
// @Param id path int true "Log ID"
// @Success 200 {object} models.SoftwareAssignmentLog
// @Failure 404 {object} models.APIResponse
// @Router /api/logs/{id} [get]
func GetAssignmentLogByID(c *gin.Context) {
	id := c.Param("id")
	var log models.SoftwareAssignmentLog

	if err := config.DB.Preload("Staff").Preload("Software").First(&log, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// CreateSoftwareAssignmentLog godoc
// @Summary Create a software assignment log
// @Description Adds a new log entry for a software assignment action
// @Tags Software Assignment Logs
// @Accept json
// @Produce json
// @Param log body models.SoftwareAssignmentLog true "Software Assignment Log"
// @Success 201 {object} models.SoftwareAssignmentLog
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/logs [post]
func CreateSoftwareAssignmentLog(c *gin.Context) {
	var logEntry models.SoftwareAssignmentLog
	if err := c.ShouldBindJSON(&logEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&logEntry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	config.DB.Preload("Staff").Preload("Software").First(&logEntry, logEntry.ID)
	c.JSON(http.StatusCreated, logEntry)
}

// UpdateSoftwareAssignmentLog godoc
// @Summary Update a software assignment log
// @Description Updates an existing log entry by ID
// @Tags Software Assignment Logs
// @Accept json
// @Produce json
// @Param id path int true "Log ID"
// @Param log body models.SoftwareAssignmentLog true "Updated Log"
// @Success 200 {object} models.SoftwareAssignmentLog
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/logs/{id} [put]
func UpdateSoftwareAssignmentLog(c *gin.Context) {
	var logEntry models.SoftwareAssignmentLog
	id := c.Param("id")

	if err := config.DB.First(&logEntry, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	if err := c.ShouldBindJSON(&logEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&logEntry)
	config.DB.Preload("Staff").Preload("Software").First(&logEntry, logEntry.ID)
	c.JSON(http.StatusOK, logEntry)
}

// DeleteSoftwareAssignmentLog godoc
// @Summary Delete a software assignment log
// @Description Deletes a log entry by ID
// @Tags Software Assignment Logs
// @Param id path int true "Log ID"
// @Success 204 "No Content"
// @Failure 500 {object} models.APIResponse
// @Router /api/logs/{id} [delete]
func DeleteSoftwareAssignmentLog(c *gin.Context) {
	var logEntry models.SoftwareAssignmentLog
	id := c.Param("id")

	if err := config.DB.Delete(&logEntry, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetAssignmentLogsForStaff godoc
// @Summary Get logs for specific staff
// @Description Retrieves all logs for a specific staff ID
// @Tags Software Assignment Logs
// @Produce json
// @Param id path int true "Staff ID"
// @Success 200 {array} models.SoftwareAssignmentLog
// @Failure 500 {object} models.APIResponse
// @Router /api/logs/staff/{id} [get]
func GetAssignmentLogsForStaff(c *gin.Context) {
	staffID := c.Param("id")
	var logs []models.SoftwareAssignmentLog

	db := config.DB.
		Preload("Staff").
		Preload("Software").
		Where("staff_id = ?", staffID)

	err := db.Order("changed_at DESC").
		Find(&logs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetAssignmentLogsForStaffWithDetails godoc
// @Summary Get detailed logs for staff
// @Description Paginated and filtered logs for a given staff ID
// @Tags Software Assignment Logs
// @Produce json
// @Param id path int true "Staff ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param search query string false "Search by software name"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} models.SoftwareAssignmentLog
// @Failure 500 {object} models.APIResponse
// @Router /api/logs/staff/{id}/details [get]
func GetAssignmentLogsForStaffWithDetails(c *gin.Context) {
	staffID := c.Param("id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var logs []models.SoftwareAssignmentLog
	query := config.DB.
		Preload("Software").
		Preload("Staff"). // assuming you've added this relation
		Where("staff_id = ?", staffID)

	if startDate != "" && endDate != "" {
		query = query.Where("changed_at BETWEEN ? AND ?", startDate, endDate)
	}
	if search != "" {
		query = query.Joins("JOIN software ON software_assignment_logs.software_id = software.id").
			Where("software.name LIKE ?", "%"+search+"%")
	}

	if err := query.Order("changed_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetAssignmentLogsForStaffWithDetailsRawSQL godoc
// @Summary Get raw logs for staff
// @Description Raw SQL version of staff logs query
// @Tags Software Assignment Logs
// @Produce json
// @Param id path int true "Staff ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param search query string false "Search by software name"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} models.APIResponse
// @Router /api/logs/staff/{id}/details/raw [get]
func GetAssignmentLogsForStaffWithDetailsRawSQL(c *gin.Context) {
	staffID := c.Param("id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	type Result struct {
		LogID         uint      `json:"log_id"`
		SoftwareID    uint      `json:"software_id"`
		Software      string    `json:"software"`
		Action        string    `json:"action"`
		ChangedBy     uint      `json:"changed_by"`
		ChangedByName string    `json:"changed_by_name"`
		ChangedAt     time.Time `json:"changed_at"`
	}

	var results []Result

	query := config.DB.Table("software_assignment_logs").
		Select(`
			software_assignment_logs.id AS log_id,
			software.id AS software_id,
			software.name AS software,
			software_assignment_logs.action,
			software_assignment_logs.changed_by,
			CONCAT(staff.first_name, ' ', staff.last_name) AS changed_by_name,
			software_assignment_logs.changed_at
		`).
		Joins("JOIN software ON software_assignment_logs.software_id = software.id").
		Joins("JOIN staff ON software_assignment_logs.changed_by = staff.id").
		Where("software_assignment_logs.staff_id = ?", staffID)

	if startDate != "" && endDate != "" {
		query = query.Where("software_assignment_logs.changed_at BETWEEN ? AND ?", startDate, endDate)
	}
	if search != "" {
		query = query.Where("software.name LIKE ?", "%"+search+"%")
	}

	err := query.Order("software_assignment_logs.changed_at DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetAssignmentLogsForSoftware godoc
// @Summary Get logs for specific software
// @Description Retrieves all logs for a specific software ID
// @Tags Software Assignment Logs
// @Produce json
// @Param id path int true "Software ID"
// @Success 200 {array} models.SoftwareAssignmentLog
// @Failure 500 {object} models.APIResponse
// @Router /api/logs/software/{id}/plain [get]
func GetAssignmentLogsForSoftware(c *gin.Context) {
	softwareID := c.Param("id")
	var logs []models.SoftwareAssignmentLog

	db := config.DB.
		Preload("Staff").
		Preload("Software").
		Where("software_id = ?", softwareID)

	err := db.Order("changed_at DESC").
		Find(&logs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetAssignmentLogsForSoftwareWithDetails godoc
// @Summary Get detailed logs for software
// @Description Paginated and filtered logs for a specific software ID
// @Tags Software Assignment Logs
// @Produce json
// @Param id path int true "Software ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param search query string false "Search by staff name"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} models.SoftwareAssignmentLog
// @Failure 500 {object} models.APIResponse
// @Router /api/logs/software/{id} [get]
func GetAssignmentLogsForSoftwareWithDetails(c *gin.Context) {
	softwareID := c.Param("id")
	search := c.Query("search")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var logs []models.SoftwareAssignmentLog

	db := config.DB.
		Preload("Staff").
		Preload("Software").
		Where("software_id = ?", softwareID)

	if startDate != "" && endDate != "" {
		db = db.Where("changed_at BETWEEN ? AND ?", startDate, endDate)
	}

	// To filter by staff name, we need a join since GORM can't filter Preloaded fields
	if search != "" {
		db = db.Joins("JOIN staff ON staff.id = software_assignment_logs.staff_id").
			Where("CONCAT(staff.first_name, ' ', staff.last_name) LIKE ?", "%"+search+"%")
	}

	err := db.Order("changed_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&logs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetAssignmentLogsForSoftwareWithDetailsRawSQL godoc
// @Summary Get raw logs for software
// @Description Raw SQL version of software logs query
// @Tags Software Assignment Logs
// @Produce json
// @Param id path int true "Software ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param search query string false "Search by staff name"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} models.APIResponse
// @Router /api/logs/software/{id}/details/raw [get]
func GetAssignmentLogsForSoftwareWithDetailsRawSQL(c *gin.Context) {
	softwareID := c.Param("id")
	search := c.Query("search")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	type Result struct {
		LogID         uint      `json:"log_id"`
		StaffID       uint      `json:"staff_id"`
		StaffName     string    `json:"staff_name"`
		Action        string    `json:"action"`
		ChangedBy     uint      `json:"changed_by"`
		ChangedByName string    `json:"changed_by_name"`
		ChangedAt     time.Time `json:"changed_at"`
	}

	var results []Result

	query := config.DB.Table("software_assignment_logs").
		Select(`
			software_assignment_logs.id AS log_id,
			staff.id AS staff_id,
			CONCAT(staff.first_name, ' ', staff.last_name) AS staff_name,
			software_assignment_logs.action,
			software_assignment_logs.changed_by,
			CONCAT(changer.first_name, ' ', changer.last_name) AS changed_by_name,
			software_assignment_logs.changed_at
		`).
		Joins("JOIN staff ON software_assignment_logs.staff_id = staff.id").
		Joins("JOIN staff AS changer ON software_assignment_logs.changed_by = changer.id").
		Where("software_assignment_logs.software_id = ?", softwareID)

	if search != "" {
		query = query.Where("CONCAT(staff.first_name, ' ', staff.last_name) LIKE ?", "%"+search+"%")
	}
	if startDate != "" && endDate != "" {
		query = query.Where("software_assignment_logs.changed_at BETWEEN ? AND ?", startDate, endDate)
	}

	err := query.Order("software_assignment_logs.changed_at DESC").
		Limit(pageSize).Offset(offset).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
