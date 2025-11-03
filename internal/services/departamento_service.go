package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/danubiobwm/company-api/internal/models"
	"github.com/danubiobwm/company-api/internal/repositories"
	"github.com/google/uuid"
)

// DepartamentoService é responsável pela lógica de negócio dos departamentos
type DepartamentoService struct {
	repo            *repositories.DepartamentoRepository
	colaboradorRepo *repositories.ColaboradorRepository
}

// NewDepartamentoService cria uma nova instância de DepartamentoService
func NewDepartamentoService(
	repo *repositories.DepartamentoRepository,
	colabRepo *repositories.ColaboradorRepository,
) *DepartamentoService {
	return &DepartamentoService{
		repo:            repo,
		colaboradorRepo: colabRepo,
	}
}

// GetAll retorna todos os departamentos
func (s *DepartamentoService) GetAll() ([]models.Departamento, error) {
	return s.repo.FindAll()
}

// GetByID retorna um departamento pelo ID
func (s *DepartamentoService) GetByID(id uuid.UUID) (*models.Departamento, error) {
	return s.repo.GetByID(id)
}

// Create cria um novo departamento
func (s *DepartamentoService) Create(d *models.Departamento) error {
	if strings.TrimSpace(d.Nome) == "" {
		return fmt.Errorf("nome é obrigatório")
	}

	// Se gerente_id foi informado, verifica se existe
	if d.GerenteID != nil && *d.GerenteID != uuid.Nil {
		gerente, err := s.colaboradorRepo.GetByID(*d.GerenteID)
		if err != nil {
			return err
		}
		if gerente == nil {
			return fmt.Errorf("gerente não encontrado")
		}
	}

	// Se departamento superior informado, verifica se existe
	if d.DepartamentoSuperiorID != nil && *d.DepartamentoSuperiorID != uuid.Nil {
		superior, err := s.repo.GetByID(*d.DepartamentoSuperiorID)
		if err != nil {
			return err
		}
		if superior == nil {
			return fmt.Errorf("departamento superior não encontrado")
		}
	}

	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}

	return s.repo.Create(d)
}

// Update atualiza um departamento existente
func (s *DepartamentoService) Update(d *models.Departamento) error {
	existing, err := s.repo.GetByID(d.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("departamento não encontrado")
	}

	// Se gerente informado, valida
	if d.GerenteID != nil && *d.GerenteID != uuid.Nil {
		gerente, err := s.colaboradorRepo.GetByID(*d.GerenteID)
		if err != nil {
			return err
		}
		if gerente == nil {
			return fmt.Errorf("gerente não encontrado")
		}
	}

	existing.Nome = d.Nome
	existing.Descricao = d.Descricao
	existing.GerenteID = d.GerenteID
	existing.DepartamentoSuperiorID = d.DepartamentoSuperiorID

	return s.repo.Update(existing)
}

// Delete remove um departamento pelo ID
func (s *DepartamentoService) Delete(id uuid.UUID) error {
	dept, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if dept == nil {
		return errors.New("departamento não encontrado")
	}

	return s.repo.Delete(id)
}
