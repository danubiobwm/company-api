package services

import (
	"strings"

	dderr "github.com/danubiobwm/company-api/internal/errors"
	"github.com/danubiobwm/company-api/internal/models"
	"github.com/danubiobwm/company-api/internal/repositories"
	"github.com/google/uuid"
)

type DepartamentoService struct {
	repo            *repositories.DepartamentoRepository
	colaboradorRepo *repositories.ColaboradorRepository
}

func NewDepartamentoService(r *repositories.DepartamentoRepository, cr *repositories.ColaboradorRepository) *DepartamentoService {
	return &DepartamentoService{repo: r, colaboradorRepo: cr}
}

func (s *DepartamentoService) Create(d *models.Departamento) error {
	if strings.TrimSpace(d.Nome) == "" {
		return dderr.New("nome é obrigatório")
	}

	// count existing departments
	cnt, err := s.repo.CountAll()
	if err != nil {
		return err
	}

	// allow creation of first department without gerente
	if d.GerenteID == nil || *d.GerenteID == uuid.Nil {
		if cnt > 0 {
			return dderr.New("gerente é obrigatório")
		}
	} else {
		// validate gerente exists
		ger, err := s.colaboradorRepo.GetByID(*d.GerenteID)
		if err != nil {
			return err
		}
		if ger == nil {
			return dderr.New("gerente não encontrado")
		}
		// optional: ensure gerente is assigned to same department (business rule)
		// if ger.DepartamentoID != d.ID { ... } // can't check before dept exists
	}

	// validate superior exists and no cycles
	if d.DepartamentoSuperiorID != nil {
		sup, err := s.repo.GetByID(*d.DepartamentoSuperiorID)
		if err != nil {
			return err
		}
		if sup == nil {
			return dderr.New("departamento_superior não encontrado")
		}
		ok, err := s.checkNoCycle(d.ID, *d.DepartamentoSuperiorID)
		if err != nil {
			return err
		}
		if !ok {
			return dderr.New("atribuir esse departamento superior geraria ciclo")
		}
	}

	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}

	return s.repo.Create(d)
}

func (s *DepartamentoService) GetAll() ([]models.Departamento, error) {
	return s.repo.GetAll()
}

func (s *DepartamentoService) GetByID(id uuid.UUID) (*models.Departamento, error) {
	return s.repo.GetByID(id)
}

func (s *DepartamentoService) Update(d *models.Departamento) error {
	// validate exists
	existing, err := s.repo.GetByID(d.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return dderr.New("departamento não encontrado")
	}

	// validate gerente if provided
	if d.GerenteID != nil && *d.GerenteID != uuid.Nil {
		ger, err := s.colaboradorRepo.GetByID(*d.GerenteID)
		if err != nil {
			return err
		}
		if ger == nil {
			return dderr.New("gerente não encontrado")
		}
	}

	// validate superior and cycles
	if d.DepartamentoSuperiorID != nil {
		sup, err := s.repo.GetByID(*d.DepartamentoSuperiorID)
		if err != nil {
			return err
		}
		if sup == nil {
			return dderr.New("departamento_superior não encontrado")
		}
		ok, err := s.checkNoCycle(d.ID, *d.DepartamentoSuperiorID)
		if err != nil {
			return err
		}
		if !ok {
			return dderr.New("atribuir esse departamento superior geraria ciclo")
		}
	}

	// apply update
	return s.repo.Update(d)
}

func (s *DepartamentoService) Delete(id uuid.UUID) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return dderr.New("departamento não encontrado")
	}
	return s.repo.Delete(id)
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
