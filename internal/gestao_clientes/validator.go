package gestao_clientes

import (
	"fmt"
	"strings"

	"github.com/MundoInvest/backend/internal/shared/config"
	"github.com/MundoInvest/backend/internal/shared/mensagens"
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
	catalogo := mensagens.ObterCatalogo()
	var erros []*ErroValidacao

	if strings.TrimSpace(requisicao.Nome) == "" {
		erros = append(erros, NewErroValidacao("nome", catalogo.Texto(mensagens.ValNomeObrigatorio)))
	} else if len(strings.TrimSpace(requisicao.Nome)) < 3 {
		erros = append(erros, NewErroValidacao("nome", catalogo.Texto(mensagens.ValNomeMinimoCaracteres)))
	}

	if strings.TrimSpace(requisicao.Email) == "" {
		erros = append(erros, NewErroValidacao("email", catalogo.Texto(mensagens.ValEmailObrigatorio)))
	} else {
		strategy := NewClienteValidationStrategy()
		if !strategy.isValidEmail(requisicao.Email) {
			erros = append(erros, NewErroValidacao("email", catalogo.Texto(mensagens.ValEmailInvalido)))
		}
	}

	if strings.TrimSpace(requisicao.TipoSolicitacao) == "" {
		erros = append(erros, NewErroValidacao("tipoSolicitacao", catalogo.Texto(mensagens.ValTipoSolicitacaoObrigatorio)))
	}

	if requisicao.ValorPatrimonio <= 0 {
		erros = append(erros, NewErroValidacao("valorPatrimonio", catalogo.Texto(mensagens.ValValorPatrimonioNegativo)))
	}

	return erros
}
