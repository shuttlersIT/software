package tests

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"software_management/config"
)

var testRouter *gin.Engine

func TestMain(m *testing.M) {
	// Setup: override DB with test DB
	config.InitTestDB() // ✅ initializes config.DB

	// Run tests
	code := m.Run()
	os.Exit(code)
}
