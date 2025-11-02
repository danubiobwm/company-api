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

// DB returns the underlying gorm DB (read-only access)
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
