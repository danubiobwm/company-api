package repositories

import (
	"errors"

	"github.com/danubiobwm/company-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ColaboradorRepository struct {
	db *gorm.DB
}

func NewColaboradorRepository(db *gorm.DB) *ColaboradorRepository {
	return &ColaboradorRepository{db: db}
}

func (r *ColaboradorRepository) DB() *gorm.DB {
	return r.db
}

func (r *ColaboradorRepository) Create(c *models.Colaborador) error {
	return r.db.Create(c).Error
}

func (r *ColaboradorRepository) GetByID(id uuid.UUID) (*models.Colaborador, error) {
	var c models.Colaborador
	if err := r.db.First(&c, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *ColaboradorRepository) GetByCPF(cpf string) (*models.Colaborador, error) {
	var c models.Colaborador
	if err := r.db.First(&c, "cpf = ?", cpf).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *ColaboradorRepository) GetByRG(rg string) (*models.Colaborador, error) {
	var c models.Colaborador
	if err := r.db.First(&c, "rg = ?", rg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *ColaboradorRepository) Update(c *models.Colaborador) error {
	return r.db.Save(c).Error
}

func (r *ColaboradorRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Colaborador{}, "id = ?", id).Error
}

func (r *ColaboradorRepository) List(filters map[string]interface{}, page, limit int) ([]models.Colaborador, int64, error) {
	var list []models.Colaborador
	query := r.db.Model(&models.Colaborador{})
	if v, ok := filters["nome"].(string); ok && v != "" {
		query = query.Where("nome ILIKE ?", "%"+v+"%")
	}
	if v, ok := filters["cpf"].(string); ok && v != "" {
		query = query.Where("cpf = ?", v)
	}
	if v, ok := filters["rg"].(string); ok && v != "" {
		query = query.Where("rg = ?", v)
	}
	if v, ok := filters["departamento_id"].(string); ok && v != "" {
		query = query.Where("departamento_id = ?", v)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
