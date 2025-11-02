package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	v1 := r.Group("/api/v1")
	{
		registerColaboradorRoutes(v1, db)
		registerDepartamentoRoutes(v1, db)
		registerGerentesRoutes(v1, db)
		registerHealthRoutes(r)
	}
}
