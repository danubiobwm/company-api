package services

import (
	"errors"
	"strings"

	"github.com/danubiobwm/company-api/internal/models"
	"github.com/danubiobwm/company-api/internal/repositories"
	"github.com/google/uuid"
)

type DepartamentoService struct {
	repo            *repositories.DepartamentoRepository
	colaboradorRepo *repositories.ColaboradorRepository
}

func NewDepartamentoService(repo *repositories.DepartamentoRepository, colaboradorRepo *repositories.ColaboradorRepository) *DepartamentoService {
	return &DepartamentoService{repo: repo, colaboradorRepo: colaboradorRepo}
}

func (s *DepartamentoService) Create(d *models.Departamento) error {
	if strings.TrimSpace(d.Nome) == "" {
		return &DomainError{Status: 422, Message: "nome é obrigatório"}
	}

	count, err := s.repo.CountAll()
	if err != nil {
		return err
	}

	// Permitir criar o primeiro departamento sem gerente
	if d.GerenteID == uuid.Nil {
		if count > 0 {
			return &DomainError{Status: 422, Message: "gerente é obrigatório"}
		}
	} else {
		gerente, err := s.colaboradorRepo.GetByID(d.GerenteID)
		if err != nil {
			return err
		}
		if gerente == nil {
			return &DomainError{Status: 404, Message: "gerente não encontrado"}
		}
	}

	if d.DepartamentoSuperiorID != nil {
		ok, err := s.checkNoCycle(d.ID, *d.DepartamentoSuperiorID)
		if err != nil {
			return err
		}
		if !ok {
			return &DomainError{Status: 422, Message: "atribuir esse departamento superior geraria ciclo"}
		}
	}

	if d.ID == uuid.Nil {
		id, _ := uuid.NewV7()
		d.ID = id
	}

	return s.repo.Create(d)
}

func (s *DepartamentoService) checkNoCycle(depID uuid.UUID, superiorID uuid.UUID) (bool, error) {
	superior, err := s.repo.GetByID(superiorID)
	if err != nil {
		return false, err
	}
	if superior == nil {
		return false, errors.New("departamento superior inexistente")
	}
	if superior.DepartamentoSuperiorID != nil {
		if *superior.DepartamentoSuperiorID == depID {
			return false, nil
		}
		return s.checkNoCycle(depID, *superior.DepartamentoSuperiorID)
	}
	return true, nil
}

func (s *DepartamentoService) GetAll() ([]models.Departamento, error) {
	return s.repo.GetAll()
}

func (s *DepartamentoService) GetByID(id uuid.UUID) (*models.Departamento, error) {
	return s.repo.GetByID(id)
}

func (s *DepartamentoService) Update(d *models.Departamento) error {
	return s.repo.Update(d)
}

func (s *DepartamentoService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
