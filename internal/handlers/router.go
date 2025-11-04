package handlers

import (
	_ "github.com/danubiobwm/company-api/docs"
	"github.com/danubiobwm/company-api/internal/repositories"
	"github.com/danubiobwm/company-api/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api/v1")

	// Health check
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

	// Registrar rotas do Gerente
	RegisterGerenteRoutes(api, db)

	// Registrar rotas do Swagger
	RegisterSwaggerRoutes(r)
}
