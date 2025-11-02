package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func registerGerentesRoutes(rg *gin.RouterGroup, db *gorm.DB) {

	rg.GET("/gerentes/:id/colaboradores", func(c *gin.Context) {
		gerenteID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}

		var deptID uuid.UUID
		err = db.Raw(`
			SELECT d.id
			FROM departamentos d
			WHERE d.gerente_id = ?
			LIMIT 1;
		`, gerenteID).Scan(&deptID).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if deptID == uuid.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "gerente não vinculado a nenhum departamento"})
			return
		}

		type Result struct {
			ID   uuid.UUID `json:"id"`
			Nome string    `json:"nome"`
		}
		var depts []Result
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

		var deptIDs []uuid.UUID
		for _, d := range depts {
			deptIDs = append(deptIDs, d.ID)
		}

		type Colab struct {
			ID             uuid.UUID `json:"id"`
			Nome           string    `json:"nome"`
			CPF            string    `json:"cpf"`
			DepartamentoID uuid.UUID `json:"departamento_id"`
		}
		var colaboradores []Colab
		if err := db.Table("colaboradores").Where("departamento_id IN ?", deptIDs).Find(&colaboradores).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"gerente_id":    gerenteID,
			"departamentos": depts,
			"colaboradores": colaboradores,
			"total_colabs":  len(colaboradores),
		})
	})
}
