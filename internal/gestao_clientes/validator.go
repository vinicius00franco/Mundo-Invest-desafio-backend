package gestao_clientes

import (
	"fmt"
	"strings"

	"github.com/MundoInvest/backend/internal/shared/config"
)

// RequisicaoCriarCliente representa o payload para criação de cliente
type RequisicaoCriarCliente struct {
	Nome            string  `json:"nome"`
	Email           string  `json:"email"`
	ValorPatrimonio float64 `json:"valorPatrimonio"`
	TipoSolicitacao string  `json:"tipoSolicitacao"`
}

// ValidarCriarClienteRequest valida o payload de criação de cliente usando Strategy Pattern
func ValidarRequisicaoCriarCliente(requisicao RequisicaoCriarCliente) error {
	strategy := NewClienteValidationStrategy()
	return strategy.Validate(requisicao)
}

// ValidarCriarClienteRequestComConfig valida o payload usando configuração customizada
func ValidarRequisicaoCriarClienteComConfig(requisicao RequisicaoCriarCliente, cfg *config.Config) error {
	strategy := NewClienteValidationStrategyComConfig(cfg)
	return strategy.Validate(requisicao)
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
func ValidarRequisicaoCriarClienteDetalhado(requisicao RequisicaoCriarCliente) []*ErroValidacao {
	var erros []*ErroValidacao

	if strings.TrimSpace(requisicao.Nome) == "" {
		erros = append(erros, NewErroValidacao("nome", "nome é obrigatório"))
	} else if len(strings.TrimSpace(requisicao.Nome)) < 3 {
		erros = append(erros, NewErroValidacao("nome", "nome deve ter pelo menos 3 caracteres"))
	}

	if strings.TrimSpace(requisicao.Email) == "" {
		erros = append(erros, NewErroValidacao("email", "email é obrigatório"))
	} else {
		strategy := NewClienteValidationStrategy()
		if !strategy.isValidEmail(requisicao.Email) {
			erros = append(erros, NewErroValidacao("email", "email inválido"))
		}
	}

	if strings.TrimSpace(requisicao.TipoSolicitacao) == "" {
		erros = append(erros, NewErroValidacao("tipoSolicitacao", "tipoSolicitacao é obrigatório"))
	}

	if requisicao.ValorPatrimonio <= 0 {
		erros = append(erros, NewErroValidacao("valorPatrimonio", "valorPatrimonio deve ser positivo"))
	}

	return erros
}
