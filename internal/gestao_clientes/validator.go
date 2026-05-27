package gestao_clientes

import (
	"fmt"
	"strings"

	"github.com/MundoInvest/backend/internal/shared/config"
)

// CriarClienteRequest representa o payload para criação de cliente
type CriarClienteRequest struct {
	Nome            string  `json:"nome"`
	Email           string  `json:"email"`
	ValorPatrimonio float64 `json:"valorPatrimonio"`
	TipoSolicitacao string  `json:"tipoSolicitacao"`
}

// ValidarCriarClienteRequest valida o payload de criação de cliente usando Strategy Pattern
func ValidarCriarClienteRequest(request CriarClienteRequest) error {
	strategy := NewClienteValidationStrategy()
	return strategy.Validate(request)
}

// ValidarCriarClienteRequestComConfig valida o payload usando configuração customizada
func ValidarCriarClienteRequestComConfig(request CriarClienteRequest, cfg *config.Config) error {
	strategy := NewClienteValidationStrategyComConfig(cfg)
	return strategy.Validate(request)
}

// ErroValidacao representa um erro de validação com detalhes
type ErroValidacao struct {
	Campo    string
	Mensagem string
}

// NewErroValidacao cria um novo erro de validação
func NewErroValidacao(campo, mensagem string) *ErroValidacao {
	return &ErroValidacao{
		Campo:    campo,
		Mensagem: mensagem,
	}
}

// Error implementa a interface error
func (e *ErroValidacao) Error() string {
	return fmt.Sprintf("%s: %s", e.Campo, e.Mensagem)
}

// ValidarCriarClienteRequestDetalhado valida o payload e retorna erros detalhados
func ValidarCriarClienteRequestDetalhado(request CriarClienteRequest) []*ErroValidacao {
	var erros []*ErroValidacao

	if strings.TrimSpace(request.Nome) == "" {
		erros = append(erros, NewErroValidacao("nome", "nome é obrigatório"))
	} else if len(strings.TrimSpace(request.Nome)) < 3 {
		erros = append(erros, NewErroValidacao("nome", "nome deve ter pelo menos 3 caracteres"))
	}

	if strings.TrimSpace(request.Email) == "" {
		erros = append(erros, NewErroValidacao("email", "email é obrigatório"))
	} else {
		strategy := NewClienteValidationStrategy()
		if !strategy.isValidEmail(request.Email) {
			erros = append(erros, NewErroValidacao("email", "email inválido"))
		}
	}

	if strings.TrimSpace(request.TipoSolicitacao) == "" {
		erros = append(erros, NewErroValidacao("tipoSolicitacao", "tipoSolicitacao é obrigatório"))
	}

	if request.ValorPatrimonio <= 0 {
		erros = append(erros, NewErroValidacao("valorPatrimonio", "valorPatrimonio deve ser positivo"))
	}

	return erros
}
