package errors

import "fmt"

// DomainError representa um erro de negócio ou validação.
type DomainError struct {
	Message string
	Code    string
}

// Error implementa a interface error.
func (e *DomainError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("[%s] %s", e.Code, e.Message)
	}
	return e.Message
}

// New cria um novo erro de domínio simples.
func New(msg string) *DomainError {
	return &DomainError{Message: msg}
}

// NewWithCode cria um novo erro de domínio com código.
func NewWithCode(code, msg string) *DomainError {
	return &DomainError{Code: code, Message: msg}
}
