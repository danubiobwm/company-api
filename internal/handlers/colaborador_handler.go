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

// GetAll godoc
// @Summary List all colaboradores
// @Description Get a paginated list of colaboradores with optional filtering
// @Tags colaboradores
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(100)
// @Success 200 {object} map[string]interface{} "Lista de colaboradores e total"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /api/v1/colaboradores [get]
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

// GetByID godoc
// @Summary Get colaborador by ID
// @Description Get detailed information about a specific colaborador
// @Tags colaboradores
// @Accept json
// @Produce json
// @Param id path string true "Colaborador ID (UUID)"
// @Success 200 {object} models.Colaborador
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Colaborador não encontrado"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /api/v1/colaboradores/{id} [get]
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

// Create godoc
// @Summary Create a new colaborador
// @Description Create a new colaborador with the provided data
// @Tags colaboradores
// @Accept json
// @Produce json
// @Param colaborador body models.Colaborador true "Colaborador data"
// @Success 201 {object} models.Colaborador
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 422 {object} map[string]string "Entidade não processável"
// @Router /api/v1/colaboradores [post]
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

// Update godoc
// @Summary Update a colaborador
// @Description Update an existing colaborador by ID
// @Tags colaboradores
// @Accept json
// @Produce json
// @Param id path string true "Colaborador ID (UUID)"
// @Param colaborador body models.Colaborador true "Colaborador data"
// @Success 200 {object} models.Colaborador
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 422 {object} map[string]string "Entidade não processável"
// @Router /api/v1/colaboradores/{id} [put]
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

// Delete godoc
// @Summary Delete a colaborador
// @Description Delete a colaborador by ID
// @Tags colaboradores
// @Accept json
// @Produce json
// @Param id path string true "Colaborador ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /api/v1/colaboradores/{id} [delete]
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
