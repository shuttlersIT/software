package main

import (
	"log"
	"os"
	"software_management/config"
	_ "software_management/docs" // 👈 Required for Swagger docs
	"software_management/models"
	"software_management/routes"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/joho/godotenv"
)

func main() {
	config.InitDB()
	config.DB.AutoMigrate(&models.Department{})

	r := routes.RegisterRoutes()

	// ✅ Register this route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on http://localhost:%s", port)
	_ = r.Run(":" + port)
}

func init() {
	_ = godotenv.Load()
}
