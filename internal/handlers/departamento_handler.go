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

// GetAll godoc
// @Summary List all departamentos
// @Description Get a list of all departamentos
// @Tags departamentos
// @Accept json
// @Produce json
// @Success 200 {array} models.Departamento
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /api/v1/departamentos [get]
func (h *DepartamentoHandler) GetAll(c *gin.Context) {
	depts, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, depts)
}

// GetByID godoc
// @Summary Get departamento by ID
// @Description Get detailed information about a specific departamento
// @Tags departamentos
// @Accept json
// @Produce json
// @Param id path string true "Departamento ID (UUID)"
// @Success 200 {object} models.Departamento
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Departamento não encontrado"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /api/v1/departamentos/{id} [get]
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

// Create godoc
// @Summary Create a new departamento
// @Description Create a new departamento with the provided data
// @Tags departamentos
// @Accept json
// @Produce json
// @Param departamento body models.Departamento true "Departamento data"
// @Success 201 {object} models.Departamento
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 422 {object} map[string]string "Entidade não processável"
// @Router /api/v1/departamentos [post]
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

// Update godoc
// @Summary Update a departamento
// @Description Update an existing departamento by ID
// @Tags departamentos
// @Accept json
// @Produce json
// @Param id path string true "Departamento ID (UUID)"
// @Param departamento body models.Departamento true "Departamento data"
// @Success 200 {object} models.Departamento
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 422 {object} map[string]string "Entidade não processável"
// @Router /api/v1/departamentos/{id} [put]
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

// Delete godoc
// @Summary Delete a departamento
// @Description Delete a departamento by ID
// @Tags departamentos
// @Accept json
// @Produce json
// @Param id path string true "Departamento ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /api/v1/departamentos/{id} [delete]
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
