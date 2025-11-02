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

func registerDepartamentoRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := repositories.NewDepartamentoRepository(db)
	colabRepo := repositories.NewColaboradorRepository(db)
	srv := services.NewDepartamentoService(repo, colabRepo)

	rg.POST("/departamentos", func(c *gin.Context) {
		var payload models.Departamento
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

	rg.GET("/departamentos/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}

		dept, err := repo.GetByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if dept == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "departamento não encontrado"})
			return
		}

		type Node struct {
			ID                     uuid.UUID  `json:"id"`
			Nome                   string     `json:"nome"`
			GerenteID              uuid.UUID  `json:"gerente_id"`
			DepartamentoSuperiorID *uuid.UUID `json:"departamento_superior_id"`
		}

		var nodes []Node
		sql := `
		WITH RECURSIVE tree AS (
			SELECT id, nome, gerente_id, departamento_superior_id
			FROM departamentos
			WHERE id = ?
			UNION ALL
			SELECT d.id, d.nome, d.gerente_id, d.departamento_superior_id
			FROM departamentos d
			INNER JOIN tree t ON d.departamento_superior_id = t.id
		)
		SELECT * FROM tree;`
		if err := db.Raw(sql, id).Scan(&nodes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"departamento": dept,
			"hierarquia":   nodes,
		})
	})

	rg.PUT("/departamentos/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}

		existing, err := repo.GetByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if existing == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "departamento não encontrado"})
			return
		}

		var payload models.Departamento
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		existing.Nome = payload.Nome
		existing.GerenteID = payload.GerenteID
		existing.DepartamentoSuperiorID = payload.DepartamentoSuperiorID

		if err := srv.Update(existing); err != nil {
			if de, ok := err.(*services.DomainError); ok {
				c.JSON(de.Status, gin.H{"error": de.Message})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, existing)
	})

	rg.DELETE("/departamentos/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
			return
		}

		if err := repo.Delete(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	})

	rg.POST("/departamentos/listar", func(c *gin.Context) {
		var body struct {
			Nome                   string     `json:"nome"`
			GerenteNome            string     `json:"gerente_nome"`
			DepartamentoSuperiorID *uuid.UUID `json:"departamento_superior_id"`
			Page, Limit            int        `json:"page"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
			return
		}

		page := body.Page
		if page < 1 {
			page = 1
		}
		limit := body.Limit
		if limit < 1 || limit > 100 {
			limit = 10
		}

		query := db.Model(&models.Departamento{})
		if body.Nome != "" {
			query = query.Where("nome ILIKE ?", "%"+body.Nome+"%")
		}
		if body.DepartamentoSuperiorID != nil {
			query = query.Where("departamento_superior_id = ?", body.DepartamentoSuperiorID)
		}

		var results []models.Departamento
		if err := query.Offset((page - 1) * limit).Limit(limit).Find(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"page":          page,
			"limit":         limit,
			"departamentos": results,
		})
	})
}
