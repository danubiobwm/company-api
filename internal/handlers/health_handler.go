package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse representa o formato de resposta do endpoint de health check.
// @Description Estrutura retornada pelo health check
type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

// HealthHandler lida com o endpoint /health
type HealthHandler struct{}

// RegisterRoutes registra as rotas de health check
func (h *HealthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/health")
	r.GET("", h.HealthCheck)
}

// HealthCheck godoc
// @Summary Health check
// @Description Retorna o status de saúde da API
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /api/v1/health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}
