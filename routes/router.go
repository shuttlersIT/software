package routes

import (
	"software_management/controllers"

	"github.com/gin-gonic/gin"
)

// @title Software Management API
// @version 1.0
// @description API for managing software assignments in an organization.
// @host localhost:8080
// @BasePath /api

func RegisterRoutes() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{

		// ===== Department Routes =====
		api.GET("/departments/plain", controllers.GetDepartments)
		api.GET("/departments", controllers.GetDepartmentsWithDetail)
		api.POST("/departments", controllers.CreateDepartment)
		api.PUT("/departments/:id", controllers.UpdateDepartment)
		api.DELETE("/departments/:id", controllers.DeleteDepartment)

		// Nested: Teams under a Department
		api.GET("/departments/:id/teams", controllers.GetTeamsByDepartment)

		// ===== Team Routes =====
		api.GET("/teams", controllers.GetTeams)
		api.POST("/teams", controllers.CreateTeam)
		api.PUT("/teams/:id", controllers.UpdateTeam)
		api.DELETE("/teams/:id", controllers.DeleteTeam)

		// ===== Staff Routes =====
		api.GET("/staff/plain", controllers.GetAllStaff)
		api.GET("/staff", controllers.GetStaffWithDetail)
		api.GET("/staff/:id", controllers.GetStaffByID)
		api.POST("/staff", controllers.CreateStaff)
		api.POST("/staff/with-software", controllers.CreateStaffWithSoftwareMatch)
		api.PUT("/staff/:id", controllers.UpdateStaff)
		api.PUT("/staff/:id/with-software", controllers.UpdateStaffWithSoftwareMatch)
		api.DELETE("/staff/:id", controllers.DeleteStaff)
		api.PUT("/staff/:id/offboard", controllers.OffboardStaff) // ✅ Offboarding

		// Nested: Staff-related software & logs
		api.GET("/staff/:id/assigned-software", controllers.GetSoftwareAssignedToStaff)
		api.GET("/staff/:id/assigned-software/detail", controllers.GetSoftwareAssignedToStaffWithDetails)
		api.GET("/staff/:id/logs", controllers.GetAssignmentLogsForStaff)

		// ===== Software Routes =====
		api.GET("/software/plain", controllers.GetAllSoftware)
		api.GET("/software", controllers.GetAllSoftwareWithDetail)
		api.GET("/software/names", controllers.GetAllSoftwareNames)
		api.GET("/software/:id", controllers.GetSoftwareByID)
		api.POST("/software", controllers.CreateSoftware)
		api.PUT("/software/:id", controllers.UpdateSoftware)
		api.DELETE("/software/:id", controllers.DeleteSoftware)

		// Nested: Software-related staff & logs
		api.GET("/software/:id/assigned-staff", controllers.GetStaffAssignedToSoftware)
		api.GET("/software/:id/assigned-staff/detail", controllers.GetStaffAssignedToSoftwareWithDetails)

		// ===== Manual Software Assignment Routes =====
		api.GET("/software-assignments/plain", controllers.GetSoftwareAssignments)
		api.GET("/software-assignments", controllers.GetSoftwareAssignmentsWithDetail)
		api.GET("/api/software-assignments/:id", controllers.GetSoftwareAssignmentByID)
		api.POST("/software-assignments", controllers.CreateSoftwareAssignment)
		api.PUT("/software-assignments/:id", controllers.UpdateSoftwareAssignment)
		api.DELETE("/software-assignments/:id", controllers.DeleteSoftwareAssignment)
		api.DELETE("/software-assignments/:id/force", controllers.DeleteSoftwareAssignmentWithForce)

		// ===== Assigned Software Routes (Actual assignments) =====
		api.GET("/assigned-software", controllers.GetAssignedSoftware)
		api.POST("/assign-software", controllers.CreateAssignedSoftware)
		api.PUT("/assigned-software/:id", controllers.UpdateAssignedSoftware)
		api.DELETE("/assigned-software/:id/force", controllers.DeleteAssignedSoftware)
		api.DELETE("/assigned-software/:id", controllers.DeleteAssignedSoftwareWithLogging)

		// ===== Software Assignment Logs =====
		api.GET("/logs", controllers.GetAllAssignmentLogsWithDetails)
		api.GET("/logs/details", controllers.GetAllAssignmentLogsWithDetails)
		api.GET("/logs/:id", controllers.GetAssignmentLogByID)
		api.POST("/logs", controllers.CreateSoftwareAssignmentLog)
		api.PUT("/logs/:id", controllers.UpdateSoftwareAssignmentLog)
		api.DELETE("/logs/:id", controllers.DeleteSoftwareAssignmentLog)
		api.GET("/logs/staff/:id", controllers.GetAssignmentLogsForStaff)
		api.GET("/logs/staff/:id/details", controllers.GetAssignmentLogsForStaffWithDetails)
		api.GET("/logs/software/:id/plain", controllers.GetAssignmentLogsForSoftware)
		api.GET("/logs/software/:id", controllers.GetAssignmentLogsForSoftwareWithDetails)

		// ===== Auto-assignment Match Controllers =====

		// Organization-level match routes
		orgMatch := api.Group("/software-organization-matches")
		{
			orgMatch.GET("", controllers.GetSoftwareOrganizationMatches)
			orgMatch.POST("", controllers.CreateSoftwareOrganizationMatch)
			orgMatch.PUT("/:id", controllers.UpdateSoftwareOrganizationMatch)
			orgMatch.DELETE("/:id", controllers.DeleteSoftwareOrganizationMatch)
		}

		// Department-level match routes
		deptMatch := api.Group("/software-department-matches")
		{
			deptMatch.GET("", controllers.GetSoftwareDepartmentMatches)
			deptMatch.POST("", controllers.CreateSoftwareDepartmentMatch)
			deptMatch.PUT("/:id", controllers.UpdateSoftwareDepartmentMatch)
			deptMatch.DELETE("/:id", controllers.DeleteSoftwareDepartmentMatch)
		}

		// Team-level match routes
		teamMatch := api.Group("/software-team-matches")
		{
			teamMatch.GET("", controllers.GetSoftwareTeamMatches)
			teamMatch.POST("", controllers.CreateSoftwareTeamMatch)
			teamMatch.PUT("/:id", controllers.UpdateSoftwareTeamMatch)
			teamMatch.DELETE("/:id", controllers.DeleteSoftwareTeamMatch)
		}
	}

	return r
}
