package services

import (
	"strings"

	"github.com/danubiobwm/company-api/internal/models"
	"github.com/danubiobwm/company-api/internal/repositories"
	"github.com/google/uuid"
)

type DepartamentoService struct {
	repo            *repositories.DepartamentoRepository
	colaboradorRepo *repositories.ColaboradorRepository
}

func NewDepartamentoService(r *repositories.DepartamentoRepository, cr *repositories.ColaboradorRepository) *DepartamentoService {
	return &DepartamentoService{
		repo:            r,
		colaboradorRepo: cr,
	}
}

func (s *DepartamentoService) Create(d *models.Departamento) error {
	if strings.TrimSpace(d.Nome) == "" {
		return &DomainError{Status: 422, Message: "nome é obrigatório"}
	}

	if d.GerenteID == uuid.Nil {
		return &DomainError{Status: 422, Message: "gerente_id é obrigatório"}
	}

	gerente, err := s.colaboradorRepo.GetByID(d.GerenteID)
	if err != nil {
		return err
	}
	if gerente == nil {
		return &DomainError{Status: 422, Message: "gerente não encontrado"}
	}

	if d.DepartamentoSuperiorID != nil {
		superior, err := s.repo.GetByID(*d.DepartamentoSuperiorID)
		if err != nil {
			return err
		}
		if superior == nil {
			return &DomainError{Status: 422, Message: "departamento_superior não encontrado"}
		}

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

func (s *DepartamentoService) Update(d *models.Departamento) error {
	if strings.TrimSpace(d.Nome) == "" {
		return &DomainError{Status: 422, Message: "nome é obrigatório"}
	}

	if d.GerenteID == uuid.Nil {
		return &DomainError{Status: 422, Message: "gerente_id é obrigatório"}
	}

	gerente, err := s.colaboradorRepo.GetByID(d.GerenteID)
	if err != nil {
		return err
	}
	if gerente == nil {
		return &DomainError{Status: 422, Message: "gerente não encontrado"}
	}

	if d.DepartamentoSuperiorID != nil {
		superior, err := s.repo.GetByID(*d.DepartamentoSuperiorID)
		if err != nil {
			return err
		}
		if superior == nil {
			return &DomainError{Status: 422, Message: "departamento_superior não encontrado"}
		}

		ok, err := s.checkNoCycle(d.ID, *d.DepartamentoSuperiorID)
		if err != nil {
			return err
		}
		if !ok {
			return &DomainError{Status: 422, Message: "atribuir esse departamento superior geraria ciclo"}
		}
	}

	return s.repo.Update(d)
}

func (s *DepartamentoService) checkNoCycle(child, parent uuid.UUID) (bool, error) {
	cur := parent
	for {
		dept, err := s.repo.GetByID(cur)
		if err != nil {
			return false, err
		}
		if dept == nil {
			break
		}
		if dept.ID == child {
			return false, nil
		}
		if dept.DepartamentoSuperiorID == nil {
			break
		}
		cur = *dept.DepartamentoSuperiorID
	}
	return true, nil
}

func (s *DepartamentoService) Delete(id uuid.UUID) error {
	dept, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if dept == nil {
		return &DomainError{Status: 404, Message: "departamento não encontrado"}
	}
	return s.repo.Delete(id)
}
