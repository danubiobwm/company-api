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

func (r *DepartamentoRepository) Create(d *models.Departamento) error {
	return r.db.Create(d).Error
}

func (r *DepartamentoRepository) GetByID(id uuid.UUID) (*models.Departamento, error) {
	var d models.Departamento
	if err := r.db.First(&d, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (r *DepartamentoRepository) Update(d *models.Departamento) error {
	return r.db.Save(d).Error
}

func (r *DepartamentoRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Departamento{}, "id = ?", id).Error
}
