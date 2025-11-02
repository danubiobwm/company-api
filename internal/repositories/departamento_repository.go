package repositories

import (
	"github.com/danubiobwm/company-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DepartamentoRepository struct {
	db *gorm.DB
}

func NewDepartamentoRepository(db *gorm.DB) *DepartamentoRepository {
	return &DepartamentoRepository{db: db}
}

func (r *DepartamentoRepository) Create(dep *models.Departamento) error {
	return r.db.Create(dep).Error
}

func (r *DepartamentoRepository) Update(dep *models.Departamento) error {
	return r.db.Save(dep).Error
}

func (r *DepartamentoRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Departamento{}, "id = ?", id).Error
}

func (r *DepartamentoRepository) GetByID(id uuid.UUID) (*models.Departamento, error) {
	var dep models.Departamento
	if err := r.db.Preload("Gerente").First(&dep, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dep, nil
}

func (r *DepartamentoRepository) GetAll() ([]models.Departamento, error) {
	var deps []models.Departamento
	if err := r.db.Preload("Gerente").Find(&deps).Error; err != nil {
		return nil, err
	}
	return deps, nil
}

func (r *DepartamentoRepository) CountAll() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Departamento{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
