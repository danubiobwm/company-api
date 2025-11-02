package handlers

import (
	"net/http"

	"github.com/danubiobwm/company-api/internal/models"
	"github.com/danubiobwm/company-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DepartamentoHandler struct {
	service *services.DepartamentoService
}

func NewDepartamentoHandler(s *services.DepartamentoService) *DepartamentoHandler {
	return &DepartamentoHandler{service: s}
}

func (h *DepartamentoHandler) RegisterRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/departamentos")
	r.GET("", h.GetAll)
	r.GET("/:id", h.GetByID)
	r.POST("", h.Create)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}

// GET /departamentos
func (h *DepartamentoHandler) GetAll(c *gin.Context) {
	depts, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, depts)
}

// GET /departamentos/:id
func (h *DepartamentoHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	dept, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if dept == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "departamento não encontrado"})
		return
	}
	c.JSON(http.StatusOK, dept)
}

// POST /departamentos
func (h *DepartamentoHandler) Create(c *gin.Context) {
	var dept models.Departamento
	if err := c.ShouldBindJSON(&dept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(&dept); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dept)
}

// PUT /departamentos/:id
func (h *DepartamentoHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var dept models.Departamento
	if err := c.ShouldBindJSON(&dept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dept.ID = id

	if err := h.service.Update(&dept); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dept)
}

// DELETE /departamentos/:id
func (h *DepartamentoHandler) Delete(c *gin.Context) {
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
