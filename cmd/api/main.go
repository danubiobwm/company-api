package main

import (
	"fmt"
	"log"
	"os"

	"github.com/danubiobwm/company-api/internal/handlers"
	"github.com/danubiobwm/company-api/internal/repositories"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Company API
// @version 1.0
// @description This is a sample company management API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@company.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	cfg := repositories.DBConfig{
		Host:     getenv("DATABASE_HOST", "db"),
		Port:     getenv("DATABASE_PORT", "5432"),
		User:     getenv("DATABASE_USER", "postgres"),
		Password: getenv("DATABASE_PASSWORD", "postgres"),
		DBName:   getenv("DATABASE_NAME", "companydb"),
		SSLMode:  getenv("DATABASE_SSLMODE", "disable"),
	}

	db, err := repositories.NewGormDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	r := gin.Default()
	// HealthCheck godoc
	// @Summary Health check
	// @Description Check if the API is running
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router /health [get]
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	handlers.RegisterRoutes(r, db)

	// swagger (if docs generated)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%s", port)
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
