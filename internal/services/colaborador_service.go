package services

import (
	"strings"

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

func (s *ColaboradorService) Create(c *models.Colaborador) error {
	if strings.TrimSpace(c.Nome) == "" {
		return &DomainError{Status: 422, Message: "nome é obrigatório"}
	}

	if !validateCPF(c.CPF) {
		return &DomainError{Status: 422, Message: "cpf inválido"}
	}

	existing, err := s.repo.GetByCPF(c.CPF)
	if err != nil {
		return err
	}
	if existing != nil {
		return &DomainError{Status: 409, Message: "cpf já cadastrado"}
	}

	if c.RG != nil && *c.RG != "" {
		existingRG, err := s.repo.GetByRG(*c.RG)
		if err != nil {
			return err
		}
		if existingRG != nil {
			return &DomainError{Status: 409, Message: "rg já cadastrado"}
		}
	}

	dept, err := s.deptRepo.GetByID(c.DepartamentoID)
	if err != nil {
		return err
	}
	if dept == nil {
		return &DomainError{Status: 422, Message: "departamento não existe"}
	}

	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	if err := s.repo.Create(c); err != nil {
		return err
	}
	return nil
}

type DomainError struct {
	Status  int
	Message string
}

func (e *DomainError) Error() string { return e.Message }

func validateCPF(c string) bool {
	// remove non-digits
	s := ""
	for _, r := range c {
		if r >= '0' && r <= '9' {
			s += string(r)
		}
	}
	if len(s) != 11 {
		return false
	}
	// reject sequences like 11111111111
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
