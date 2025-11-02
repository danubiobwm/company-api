package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func registerGerentesRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// GET /api/v1/gerentes/:id/colaboradores
	rg.GET("/gerentes/:id/colaboradores", func(c *gin.Context) {
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
		var colabs []map[string]interface{}
		if len(ids) > 0 {
			if err := db.Table("colaboradores").Where("departamento_id IN ?", ids).Find(&colabs).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"gerente_id": gerenteID, "departamentos": depts, "colaboradores": colabs})
	})
}
