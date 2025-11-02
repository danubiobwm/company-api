package handlers

import (
	"net/http"

	"github.com/danubiobwm/company-api/internal/models"
	"github.com/danubiobwm/company-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ColaboradorHandler struct {
	service *services.ColaboradorService
}

func NewColaboradorHandler(s *services.ColaboradorService) *ColaboradorHandler {
	return &ColaboradorHandler{service: s}
}

func (h *ColaboradorHandler) RegisterRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/colaboradores")
	r.GET("", h.GetAll)
	r.GET("/:id", h.GetByID)
	r.POST("", h.Create)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}

// GET /colaboradores
func (h *ColaboradorHandler) GetAll(c *gin.Context) {
	filters := make(map[string]interface{})
	colabs, total, err := h.service.List(filters, 1, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total":         total,
		"colaboradores": colabs,
	})
}

// GET /colaboradores/:id
func (h *ColaboradorHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	colab, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if colab == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "colaborador não encontrado"})
		return
	}
	c.JSON(http.StatusOK, colab)
}

// POST /colaboradores
func (h *ColaboradorHandler) Create(c *gin.Context) {
	var colab models.Colaborador
	if err := c.ShouldBindJSON(&colab); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(&colab); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, colab)
}

// PUT /colaboradores/:id
func (h *ColaboradorHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var colab models.Colaborador
	if err := c.ShouldBindJSON(&colab); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	colab.ID = id

	if err := h.service.Update(&colab); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, colab)
}

// DELETE /colaboradores/:id
func (h *ColaboradorHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
