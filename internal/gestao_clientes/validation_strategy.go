package gestao_clientes

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/MundoInvest/backend/internal/shared/config"
	"github.com/MundoInvest/backend/internal/shared/errors"
	"github.com/MundoInvest/backend/internal/shared/mensagens"
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
	catalogo := mensagens.ObterCatalogo()
	requisicao, ok := obj.(RequisicaoCriarCliente)
	if !ok {
		return errors.NewValidationError("", catalogo.Texto(mensagens.ErrObjetoInvalido))
	}

	var erros []string

	// Validar nome
	if strings.TrimSpace(requisicao.Nome) == "" {
		erros = append(erros, catalogo.Texto(mensagens.ValNomeObrigatorio))
	} else if len(strings.TrimSpace(requisicao.Nome)) < 3 {
		erros = append(erros, catalogo.Texto(mensagens.ValNomeMinimoCaracteres))
	} else if len(strings.TrimSpace(requisicao.Nome)) > s.config.Business.MaxNomeLength {
		erros = append(erros, catalogo.TextoFormatado(mensagens.ValNomeMaximoCaracteres, s.config.Business.MaxNomeLength))
	}

	// Validar email
	if strings.TrimSpace(requisicao.Email) == "" {
		erros = append(erros, catalogo.Texto(mensagens.ValEmailObrigatorio))
	} else {
		if !s.isValidEmail(requisicao.Email) {
			erros = append(erros, catalogo.Texto(mensagens.ValEmailInvalido))
		} else if len(strings.TrimSpace(requisicao.Email)) > s.config.Business.MaxEmailLength {
			erros = append(erros, catalogo.TextoFormatado(mensagens.ValEmailMaximoCaracteres, s.config.Business.MaxEmailLength))
		}
	}

	// Validar tipoSolicitacao
	if strings.TrimSpace(requisicao.TipoSolicitacao) == "" {
		erros = append(erros, catalogo.Texto(mensagens.ValTipoSolicitacaoObrigatorio))
	} else if len(strings.TrimSpace(requisicao.TipoSolicitacao)) > s.config.Business.MaxTipoSolicitacaoLength {
		erros = append(erros, catalogo.TextoFormatado(mensagens.ValTipoSolicitacaoMaximo, s.config.Business.MaxTipoSolicitacaoLength))
	}

	// Validar valorPatrimonio
	if requisicao.ValorPatrimonio < 0 {
		erros = append(erros, catalogo.Texto(mensagens.ValValorPatrimonioNegativo))
	} else if requisicao.ValorPatrimonio > s.config.Business.MaxValorPatrimonio {
		erros = append(erros, catalogo.TextoFormatado(mensagens.ValorPatrimonioMaximo, s.config.Business.MaxValorPatrimonio))
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
	catalogo := mensagens.ObterCatalogo()
	for _, strategy := range c.strategies {
		if err := strategy.Validate(obj); err != nil {
			return fmt.Errorf(catalogo.Texto(mensagens.ErrEstrategiaValidacao), strategy.StrategyName(), err)
		}
	}
	return nil
}
