package gestao_clientes

import (
	"fmt"
	"net/mail"
	"strings"
)

// ValidationStrategy define a interface para estratégias de validação
type ValidationStrategy interface {
	// Validate valida um objeto e retorna erro se inválido
	Validate(interface{}) error
	// StrategyName retorna o nome da estratégia
	StrategyName() string
}

// ClienteValidationStrategy implementa validação de clientes
type ClienteValidationStrategy struct{}

// NewClienteValidationStrategy cria uma nova estratégia de validação de clientes
func NewClienteValidationStrategy() *ClienteValidationStrategy {
	return &ClienteValidationStrategy{}
}

// Validate valida um cliente
func (s *ClienteValidationStrategy) Validate(obj interface{}) error {
	request, ok := obj.(CriarClienteRequest)
	if !ok {
		return fmt.Errorf("objeto inválido para validação de cliente")
	}

	var erros []string

	// Validar nome
	if strings.TrimSpace(request.Nome) == "" {
		erros = append(erros, "nome é obrigatório")
	}

	// Validar email
	if strings.TrimSpace(request.Email) == "" {
		erros = append(erros, "email é obrigatório")
	} else {
		if !s.isValidEmail(request.Email) {
			erros = append(erros, "email inválido")
		}
	}

	// Validar tipo_solicitacao
	if strings.TrimSpace(request.TipoSolicitacao) == "" {
		erros = append(erros, "tipo_solicitacao é obrigatório")
	}

	// Validar valor_patrimonio
	if request.ValorPatrimonio <= 0 {
		erros = append(erros, "valor_patrimonio deve ser positivo")
	}

	if len(erros) > 0 {
		return fmt.Errorf("erros de validação: %s", strings.Join(erros, "; "))
	}

	return nil
}

// StrategyName retorna o nome da estratégia
func (s *ClienteValidationStrategy) StrategyName() string {
	return "ClienteValidationStrategy"
}

// isValidEmail valida o formato de email
func (s *ClienteValidationStrategy) isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// WebhookValidationStrategy implementa validação de webhooks
type WebhookValidationStrategy struct{}

// NewWebhookValidationStrategy cria uma nova estratégia de validação de webhooks
func NewWebhookValidationStrategy() *WebhookValidationStrategy {
	return &WebhookValidationStrategy{}
}

// Validate valida um webhook
func (s *WebhookValidationStrategy) Validate(obj interface{}) error {
	// Esta validação seria implementada no contexto de processamento_eventos
	// Por enquanto, retorna nil para manter a estrutura
	return nil
}

// StrategyName retorna o nome da estratégia
func (s *WebhookValidationStrategy) StrategyName() string {
	return "WebhookValidationStrategy"
}

// ValidationContext define o contexto de validação
type ValidationContext struct {
	strategies []ValidationStrategy
}

// NewValidationContext cria um novo contexto de validação
func NewValidationContext() *ValidationContext {
	return &ValidationContext{
		strategies: make([]ValidationStrategy, 0),
	}
}

// AddStrategy adiciona uma estratégia de validação
func (c *ValidationContext) AddStrategy(strategy ValidationStrategy) {
	c.strategies = append(c.strategies, strategy)
}

// Validate executa todas as estratégias de validação
func (c *ValidationContext) Validate(obj interface{}) error {
	for _, strategy := range c.strategies {
		if err := strategy.Validate(obj); err != nil {
			return fmt.Errorf("erro na estratégia %s: %w", strategy.StrategyName(), err)
		}
	}
	return nil
}
