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
// @description API para gerenciar colaboradores e departamentos
// @host localhost:8080
// @BasePath /api/v1
func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	dbCfg := repositories.DBConfig{
		Host:     getenv("DATABASE_HOST", "localhost"),
		Port:     getenv("DATABASE_PORT", "5432"),
		User:     getenv("DATABASE_USER", "postgres"),
		Password: getenv("DATABASE_PASSWORD", "postgres"),
		DBName:   getenv("DATABASE_NAME", "companydb"),
		SSLMode:  getenv("DATABASE_SSLMODE", "disable"),
	}

	db, err := repositories.NewGormDB(dbCfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	r := gin.Default()

	// Basic health route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register app routes (handlers)
	handlers.RegisterRoutes(r, db)

	// Swagger
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 starting server on %s", addr)
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
