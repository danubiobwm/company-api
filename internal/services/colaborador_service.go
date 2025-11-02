package services

import (
	"strings"

	dderr "github.com/danubiobwm/company-api/internal/errors"
	"github.com/danubiobwm/company-api/internal/models"
	"github.com/danubiobwm/company-api/internal/repositories"
	"github.com/google/uuid"
)

type ColaboradorService struct {
	repo     *repositories.ColaboradorRepository
	deptRepo *repositories.DepartamentoRepository
}

func NewColaboradorService(r *repositories.ColaboradorRepository, dr *repositories.DepartamentoRepository) *ColaboradorService {
	return &ColaboradorService{repo: r, deptRepo: dr}
}

// Create cria um novo colaborador com validações (CPF/RG/Depto).
func (s *ColaboradorService) Create(c *models.Colaborador) error {
	if strings.TrimSpace(c.Nome) == "" {
		return dderr.New("nome é obrigatório")
	}

	if !validateCPF(c.CPF) {
		return dderr.New("cpf inválido")
	}

	// CPF único
	existing, err := s.repo.GetByCPF(c.CPF)
	if err != nil {
		return err
	}
	if existing != nil {
		return dderr.New("cpf já cadastrado")
	}

	// RG único (se informado)
	if c.RG != nil && *c.RG != "" {
		existingRG, err := s.repo.GetByRG(*c.RG)
		if err != nil {
			return err
		}
		if existingRG != nil {
			return dderr.New("rg já cadastrado")
		}
	}

	// Departamento existe
	dept, err := s.deptRepo.GetByID(c.DepartamentoID)
	if err != nil {
		return err
	}
	if dept == nil {
		return dderr.New("departamento não existe")
	}

	// gerar ID se ausente
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	return s.repo.Create(c)
}

// GetByID retorna colaborador por UUID
func (s *ColaboradorService) GetByID(id uuid.UUID) (*models.Colaborador, error) {
	return s.repo.GetByID(id)
}

// Update atualiza colaborador (validações básicas)
func (s *ColaboradorService) Update(c *models.Colaborador) error {
	// checar existência
	existing, err := s.repo.GetByID(c.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return dderr.New("colaborador não encontrado")
	}

	// se CPF mudou, validar unicidade e formato
	if c.CPF != existing.CPF {
		if !validateCPF(c.CPF) {
			return dderr.New("cpf inválido")
		}
		if other, err := s.repo.GetByCPF(c.CPF); err != nil {
			return err
		} else if other != nil && other.ID != existing.ID {
			return dderr.New("cpf já cadastrado")
		}
	}

	// se RG mudou, validar unicidade
	if c.RG != nil && (existing.RG == nil || *c.RG != *existing.RG) {
		if other, err := s.repo.GetByRG(*c.RG); err != nil {
			return err
		} else if other != nil && other.ID != existing.ID {
			return dderr.New("rg já cadastrado")
		}
	}

	// departamento existe
	if _, err := s.deptRepo.GetByID(c.DepartamentoID); err != nil {
		return err
	}

	return s.repo.Update(c)
}

// Delete remove colaborador por id
func (s *ColaboradorService) Delete(id uuid.UUID) error {
	// opcional: checar existência
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return dderr.New("colaborador não encontrado")
	}
	return s.repo.Delete(id)
}

// List retorna lista paginada de colaboradores com filtros
func (s *ColaboradorService) List(filters map[string]interface{}, page, limit int) ([]models.Colaborador, int64, error) {
	return s.repo.List(filters, page, limit)
}

// validateCPF - mesma implementação que você já usou (limpa e calcula dígitos)
func validateCPF(c string) bool {
	s := ""
	for _, r := range c {
		if r >= '0' && r <= '9' {
			s += string(r)
		}
	}
	if len(s) != 11 {
		return false
	}
	seq := true
	for i := 1; i < 11; i++ {
		if s[i] != s[0] {
			seq = false
			break
		}
	}
	if seq {
		return false
	}

	calc := func(digs []int) int {
		sum := 0
		for i, v := range digs {
			sum += v * (len(digs) + 1 - i)
		}
		mod := sum % 11
		if mod < 2 {
			return 0
		}
		return 11 - mod
	}

	digs := make([]int, 11)
	for i := 0; i < 11; i++ {
		digs[i] = int(s[i] - '0')
	}
	d1 := calc(digs[:9])
	d2 := calc(append(digs[:9], d1))
	return d1 == digs[9] && d2 == digs[10]
}
