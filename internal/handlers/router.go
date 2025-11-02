package handlers

import (
	"github.com/danubiobwm/company-api/internal/repositories"
	"github.com/danubiobwm/company-api/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api/v1")

	// Health check simples
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Instâncias de repositórios
	deptRepo := repositories.NewDepartamentoRepository(db)
	colabRepo := repositories.NewColaboradorRepository(db)

	// Services
	deptService := services.NewDepartamentoService(deptRepo, colabRepo)
	colabService := services.NewColaboradorService(colabRepo, deptRepo)

	// Handlers
	deptHandler := NewDepartamentoHandler(deptService)
	colabHandler := NewColaboradorHandler(colabService)

	// Registrar rotas
	deptHandler.RegisterRoutes(api)
	colabHandler.RegisterRoutes(api)
	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
