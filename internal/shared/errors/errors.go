package errors

import (
	"fmt"

	"github.com/MundoInvest/backend/internal/shared/mensagens"
)

// ValidationError representa erros de validação de dados
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	catalogo := mensagens.ObterCatalogo()
	if e.Field != "" {
		return fmt.Sprintf(catalogo.Texto(mensagens.ErrValidacaoCampo), e.Field, e.Message)
	}
	return fmt.Sprintf(catalogo.Texto(mensagens.ErrValidacao), e.Message)
}

// NewValidationError cria um novo erro de validação
func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// RepositoryError representa erros de operações de repositório
type RepositoryError struct {
	Operation string
	Err       error
}

func (e *RepositoryError) Error() string {
	catalogo := mensagens.ObterCatalogo()
	return fmt.Sprintf(catalogo.Texto(mensagens.ErrRepositorio), e.Operation, e.Err)
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// NewRepositoryError cria um novo erro de repositório
func NewRepositoryError(operation string, err error) error {
	return &RepositoryError{
		Operation: operation,
		Err:       err,
	}
}

// ServiceError representa erros de lógica de negócio
type ServiceError struct {
	Service string
	Message string
	Err     error
}

func (e *ServiceError) Error() string {
	catalogo := mensagens.ObterCatalogo()
	if e.Err != nil {
		return fmt.Sprintf(catalogo.Texto(mensagens.ErrServicoComErro), e.Service, e.Message, e.Err)
	}
	return fmt.Sprintf(catalogo.Texto(mensagens.ErrServico), e.Service, e.Message)
}

func (e *ServiceError) Unwrap() error {
	return e.Err
}

// NewServiceError cria um novo erro de serviço
func NewServiceError(service, message string, err error) error {
	return &ServiceError{
		Service: service,
		Message: message,
		Err:     err,
	}
}

// IntegrationError representa erros de integração com sistemas externos
type IntegrationError struct {
	System  string
	Message string
	Err     error
}

func (e *IntegrationError) Error() string {
	catalogo := mensagens.ObterCatalogo()
	if e.Err != nil {
		return fmt.Sprintf(catalogo.Texto(mensagens.ErrIntegracaoComErro), e.System, e.Message, e.Err)
	}
	return fmt.Sprintf(catalogo.Texto(mensagens.ErrIntegracao), e.System, e.Message)
}

func (e *IntegrationError) Unwrap() error {
	return e.Err
}

// NewIntegrationError cria um novo erro de integração
func NewIntegrationError(system, message string, err error) error {
	return &IntegrationError{
		System:  system,
		Message: message,
		Err:     err,
	}
}

// TimeoutError representa erros de timeout
type TimeoutError struct {
	Operation string
	Message   string
}

func (e *TimeoutError) Error() string {
	catalogo := mensagens.ObterCatalogo()
	return fmt.Sprintf(catalogo.Texto(mensagens.ErrTimeout), e.Operation, e.Message)
}

// NewTimeoutError cria um novo erro de timeout
func NewTimeoutError(operation, message string) error {
	return &TimeoutError{
		Operation: operation,
		Message:   message,
	}
}

// NotFoundError representa erros de recurso não encontrado
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	catalogo := mensagens.ObterCatalogo()
	if e.ID != "" {
		return fmt.Sprintf(catalogo.Texto(mensagens.ErrNaoEncontradoComID), e.Resource, e.ID)
	}
	return fmt.Sprintf(catalogo.Texto(mensagens.ErrNaoEncontrado), e.Resource)
}

// NewNotFoundError cria um novo erro de recurso não encontrado
func NewNotFoundError(resource, id string) error {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}
