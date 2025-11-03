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

func (r *DepartamentoRepository) FindAll() ([]models.Departamento, error) {
	var departamentos []models.Departamento
	if err := r.db.Preload("Gerente").Find(&departamentos).Error; err != nil {
		return nil, err
	}
	return departamentos, nil
}

func (r *DepartamentoRepository) GetByID(id uuid.UUID) (*models.Departamento, error) {
	var dept models.Departamento
	if err := r.db.Preload("Gerente").First(&dept, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dept, nil
}

func (r *DepartamentoRepository) Create(d *models.Departamento) error {
	return r.db.Create(d).Error
}

func (r *DepartamentoRepository) Update(d *models.Departamento) error {
	return r.db.Save(d).Error
}

func (r *DepartamentoRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Departamento{}, "id = ?", id).Error
}
