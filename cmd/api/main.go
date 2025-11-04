package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/danubiobwm/company-api/docs"
	"github.com/danubiobwm/company-api/internal/handlers"
	"github.com/danubiobwm/company-api/internal/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @title Company API
// @version 1.0
// @description API para gerenciar colaboradores e departamentos
// @host localhost:8080
// @BasePath /api/v1

type Config struct {
	AppPort string
	DB      repositories.DBConfig
}

func loadConfig() Config {
	return Config{
		AppPort: getenv("APP_PORT", "8080"),
		DB: repositories.DBConfig{
			Host:     getenv("DATABASE_HOST", "db"),
			Port:     getenv("DATABASE_PORT", "5432"),
			User:     getenv("DATABASE_USER", "postgres"),
			Password: getenv("DATABASE_PASSWORD", "postgres"),
			DBName:   getenv("DATABASE_NAME", "companydb"),
			SSLMode:  getenv("DATABASE_SSLMODE", "disable"),
		},
	}
}

func main() {
	config := loadConfig()

	var db *gorm.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = repositories.NewGormDB(config.DB)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	r := setupRouter(db)

	addr := fmt.Sprintf(":%s", config.AppPort)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func setupRouter(db *gorm.DB) *gin.Engine {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	handlers.RegisterRoutes(r, db)

	return r
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
