package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GerenteColaboradoresResponse represents the response structure for gerente's colaboradores
type GerenteColaboradoresResponse struct {
	GerenteID     uuid.UUID               `json:"gerente_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Departamentos []DepartamentoHierarchy `json:"departamentos"`
	Colaboradores []ColaboradorSummary    `json:"colaboradores"`
}

// DepartamentoHierarchy represents a department in the hierarchy
type DepartamentoHierarchy struct {
	ID   uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Nome string    `json:"nome" example:"Technology Department"`
}

// ColaboradorSummary represents a simplified colaborador
type ColaboradorSummary struct {
	ID             uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name           string    `json:"name" example:"John Doe"`
	Email          string    `json:"email" example:"john.doe@company.com"`
	DepartamentoID uuid.UUID `json:"departamento_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// Add other relevant fields as needed
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Bad Request"`
	Message string `json:"message" example:"Invalid input data"`
}

func registerGerentesRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// GET /api/v1/gerentes/:id/colaboradores
	rg.GET("/gerentes/:id/colaboradores", getGerenteColaboradores(db))
}

// GetGerenteColaboradores godoc
// @Summary Get colaboradores under gerente's hierarchy
// @Description Get all colaboradores under a gerente's department hierarchy (including sub-departments)
// @Tags gerentes
// @Accept json
// @Produce json
// @Param id path string true "Gerente ID (UUID)"
// @Success 200 {object} GerenteColaboradoresResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /gerentes/{id}/colaboradores [get]
func getGerenteColaboradores(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		gerenteID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}

		// find department where this gerente_id is set
		var deptID uuid.UUID
		if err := db.Raw("SELECT id FROM departamentos WHERE gerente_id = ? LIMIT 1", gerenteID).Scan(&deptID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if deptID == uuid.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "gerente não vinculado a nenhum departamento"})
			return
		}

		// fetch recursively all sub-departments using CTE
		type DeptRes struct {
			ID   uuid.UUID `json:"id"`
			Nome string    `json:"nome"`
		}
		var depts []DeptRes
		sql := `
		WITH RECURSIVE subdeps AS (
			SELECT id, nome, departamento_superior_id
			FROM departamentos
			WHERE id = ?
			UNION ALL
			SELECT d.id, d.nome, d.departamento_superior_id
			FROM departamentos d
			INNER JOIN subdeps s ON d.departamento_superior_id = s.id
		)
		SELECT id, nome FROM subdeps;
		`
		if err := db.Raw(sql, deptID).Scan(&depts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var ids []uuid.UUID
		for _, d := range depts {
			ids = append(ids, d.ID)
		}
		// fetch collaborators in these departments
		var colabs []ColaboradorSummary
		if len(ids) > 0 {
			if err := db.Table("colaboradores").Where("departamento_id IN ?", ids).Find(&colabs).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// Convert depts to the response format
		var departamentos []DepartamentoHierarchy
		for _, d := range depts {
			departamentos = append(departamentos, DepartamentoHierarchy{
				ID:   d.ID,
				Nome: d.Nome,
			})
		}

		response := GerenteColaboradoresResponse{
			GerenteID:     gerenteID,
			Departamentos: departamentos,
			Colaboradores: colabs,
		}

		c.JSON(http.StatusOK, response)
	}
}
