package gestao_clientes

import (
	"errors"
	"fmt"
	"strings"
)

// CriarClienteRequest representa o payload para criação de cliente
type CriarClienteRequest struct {
	Nome            string  `json:"nome"`
	Email           string  `json:"email"`
	ValorPatrimonio float64 `json:"valor_patrimonio"`
	TipoSolicitacao string  `json:"tipo_solicitacao"`
}

// ValidarCriarClienteRequest valida o payload de criação de cliente usando Strategy Pattern
func ValidarCriarClienteRequest(request CriarClienteRequest) error {
	strategy := NewClienteValidationStrategy()
	return strategy.Validate(request)
}

// ValidarEmail valida um email individual
func ValidarEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return errors.New("email é obrigatório")
	}

	if !isValidEmail(email) {
		return errors.New("email inválido")
	}

	return nil
}

// ValidarValorPatrimonio valida o valor do patrimônio
func ValidarValorPatrimonio(valor float64) error {
	if valor <= 0 {
		return errors.New("valor_patrimonio deve ser positivo")
	}

	return nil
}

// ValidarNome valida o nome do cliente
func ValidarNome(nome string) error {
	if strings.TrimSpace(nome) == "" {
		return errors.New("nome é obrigatório")
	}

	return nil
}

// ValidarTipoSolicitacao valida o tipo de solicitação
func ValidarTipoSolicitacao(tipo string) error {
	if strings.TrimSpace(tipo) == "" {
		return errors.New("tipo_solicitacao é obrigatório")
	}

	return nil
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
	}

	if strings.TrimSpace(request.Email) == "" {
		erros = append(erros, NewErroValidacao("email", "email é obrigatório"))
	} else if !isValidEmail(request.Email) {
		erros = append(erros, NewErroValidacao("email", "email inválido"))
	}

	if strings.TrimSpace(request.TipoSolicitacao) == "" {
		erros = append(erros, NewErroValidacao("tipo_solicitacao", "tipo_solicitacao é obrigatório"))
	}

	if request.ValorPatrimonio <= 0 {
		erros = append(erros, NewErroValidacao("valor_patrimonio", "valor_patrimonio deve ser positivo"))
	}

	return erros
}
