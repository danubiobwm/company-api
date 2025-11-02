package handlers

import (
	"net/http"

	"github.com/danubiobwm/company-api/internal/models"
	"github.com/danubiobwm/company-api/internal/repositories"
	"github.com/danubiobwm/company-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func registerColaboradorRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := repositories.NewColaboradorRepository(db)
	deptRepo := repositories.NewDepartamentoRepository(db)
	srv := services.NewColaboradorService(repo, deptRepo)

	rg.POST("/colaboradores", func(c *gin.Context) {
		var payload models.Colaborador
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		if payload.ID == uuid.Nil {
			id, _ := uuid.NewV7()
			payload.ID = id
		}
		if err := srv.Create(&payload); err != nil {
			if de, ok := err.(*services.DomainError); ok {
				c.JSON(de.Status, gin.H{"error": de.Message})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, payload)
	})

}
