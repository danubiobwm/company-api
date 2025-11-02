package repositories

import (
	"errors"

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
	var d models.Departamento
	if err := r.db.Preload("Gerente").First(&d, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (r *DepartamentoRepository) GetAll() ([]models.Departamento, error) {
	var list []models.Departamento
	if err := r.db.Preload("Gerente").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *DepartamentoRepository) CountAll() (int64, error) {
	var cnt int64
	if err := r.db.Model(&models.Departamento{}).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}
