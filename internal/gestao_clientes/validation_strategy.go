package gestao_clientes

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/MundoInvest/backend/internal/shared/config"
	"github.com/MundoInvest/backend/internal/shared/errors"
)

// ValidationStrategy define a interface para estratégias de validação
type ValidationStrategy interface {
	// Validate valida um objeto e retorna erro se inválido
	Validate(interface{}) error
	// StrategyName retorna o nome da estratégia
	StrategyName() string
}

// ClienteValidationStrategy implementa validação de clientes
type ClienteValidationStrategy struct {
	config *config.Config
}

// NewClienteValidationStrategy cria uma nova estratégia de validação de clientes com configuração padrão
func NewClienteValidationStrategy() *ClienteValidationStrategy {
	return &ClienteValidationStrategy{
		config: config.Load(),
	}
}

// NewClienteValidationStrategyComConfig cria uma nova estratégia de validação de clientes com configuração customizada
func NewClienteValidationStrategyComConfig(cfg *config.Config) *ClienteValidationStrategy {
	return &ClienteValidationStrategy{
		config: cfg,
	}
}

// Validate valida um cliente
func (s *ClienteValidationStrategy) Validate(obj interface{}) error {
	request, ok := obj.(CriarClienteRequest)
	if !ok {
		return errors.NewValidationError("", "objeto inválido para validação de cliente")
	}

	var erros []string

	// Validar nome
	if strings.TrimSpace(request.Nome) == "" {
		erros = append(erros, "nome é obrigatório")
	} else if len(strings.TrimSpace(request.Nome)) < 3 {
		erros = append(erros, "nome deve ter pelo menos 3 caracteres")
	} else if len(strings.TrimSpace(request.Nome)) > s.config.Business.MaxNomeLength {
		erros = append(erros, fmt.Sprintf("nome deve ter no máximo %d caracteres", s.config.Business.MaxNomeLength))
	}

	// Validar email
	if strings.TrimSpace(request.Email) == "" {
		erros = append(erros, "email é obrigatório")
	} else {
		if !s.isValidEmail(request.Email) {
			erros = append(erros, "email inválido")
		} else if len(strings.TrimSpace(request.Email)) > s.config.Business.MaxEmailLength {
			erros = append(erros, fmt.Sprintf("email deve ter no máximo %d caracteres", s.config.Business.MaxEmailLength))
		}
	}

	// Validar tipoSolicitacao
	if strings.TrimSpace(request.TipoSolicitacao) == "" {
		erros = append(erros, "tipoSolicitacao é obrigatório")
	} else if len(strings.TrimSpace(request.TipoSolicitacao)) > s.config.Business.MaxTipoSolicitacaoLength {
		erros = append(erros, fmt.Sprintf("tipoSolicitacao deve ter no máximo %d caracteres", s.config.Business.MaxTipoSolicitacaoLength))
	}

	// Validar valorPatrimonio
	if request.ValorPatrimonio < s.config.Business.MinValorPatrimonio {
		erros = append(erros, fmt.Sprintf("valorPatrimonio deve ser maior ou igual a %.2f", s.config.Business.MinValorPatrimonio))
	} else if request.ValorPatrimonio > s.config.Business.MaxValorPatrimonio {
		erros = append(erros, fmt.Sprintf("valorPatrimonio deve ser menor ou igual a %.2f", s.config.Business.MaxValorPatrimonio))
	}

	if len(erros) > 0 {
		return errors.NewValidationError("", strings.Join(erros, "; "))
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
