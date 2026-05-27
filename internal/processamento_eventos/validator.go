package processamento_eventos

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

// WebhookRequest representa o payload para processamento de webhook
type RequisicaoWebhook struct {
	IdentificadorEvento string `json:"identificadorEvento"`
	IdentificadorCard   string `json:"identificadorCard"`
	EmailCliente        string `json:"emailCliente"`
	DataEvento          string `json:"dataEvento"`
}

// ValidarWebhookRequest valida o payload de webhook
func ValidarRequisicaoWebhook(requisicao RequisicaoWebhook) error {
	var erros []string

	// Validar identificadorEvento
	if strings.TrimSpace(requisicao.IdentificadorEvento) == "" {
		erros = append(erros, "identificadorEvento é obrigatório")
	}

	// Validar identificadorCard
	if strings.TrimSpace(requisicao.IdentificadorCard) == "" {
		erros = append(erros, "identificadorCard é obrigatório")
	}

	// Validar emailCliente
	if strings.TrimSpace(requisicao.EmailCliente) == "" {
		erros = append(erros, "emailCliente é obrigatório")
	} else {
		if !isValidEmail(requisicao.EmailCliente) {
			erros = append(erros, "emailCliente inválido")
		}
	}

	// Validar dataEvento
	if strings.TrimSpace(requisicao.DataEvento) == "" {
		erros = append(erros, "dataEvento é obrigatório")
	} else {
		if !isValidTimestamp(requisicao.DataEvento) {
			erros = append(erros, "dataEvento inválido")
		}
	}

	if len(erros) > 0 {
		return errors.New(strings.Join(erros, "; "))
	}

	return nil
}

// isValidEmail valida o formato de email
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// isValidTimestamp valida o formato de timestamp
func isValidTimestamp(timestamp string) bool {
	_, err := time.Parse(time.RFC3339, timestamp)
	return err == nil
}

// ValidarIdentificadorEvento valida o identificador do evento
func ValidarIdentificadorEvento(identificador string) error {
	if strings.TrimSpace(identificador) == "" {
		return errors.New("identificadorEvento é obrigatório")
	}
	return nil
}

// ValidarIdentificadorCard valida o identificador do card
func ValidarIdentificadorCard(identificador string) error {
	if strings.TrimSpace(identificador) == "" {
		return errors.New("identificadorCard é obrigatório")
	}
	return nil
}

// ValidarEmailCliente valida o email do cliente
func ValidarEmailCliente(email string) error {
	if strings.TrimSpace(email) == "" {
		return errors.New("emailCliente é obrigatório")
	}

	if !isValidEmail(email) {
		return errors.New("emailCliente inválido")
	}

	return nil
}

// ValidarDataEvento valida a data do evento
func ValidarDataEvento(data string) error {
	if strings.TrimSpace(data) == "" {
		return errors.New("dataEvento é obrigatório")
	}

	if !isValidTimestamp(data) {
		return errors.New("dataEvento inválido")
	}

	return nil
}

// ErroValidacaoWebhook representa um erro de validação de webhook com detalhes
type ErroValidacaoWebhook struct {
	Campo    string
	Mensagem string
}

// NewErroValidacaoWebhook cria um novo erro de validação de webhook
func NewErroValidacaoWebhook(campo, mensagem string) *ErroValidacaoWebhook {
	return &ErroValidacaoWebhook{
		Campo:    campo,
		Mensagem: mensagem,
	}
}

// Error implementa a interface error
func (e *ErroValidacaoWebhook) Error() string {
	return fmt.Sprintf("%s: %s", e.Campo, e.Mensagem)
}

// ValidarWebhookRequestDetalhado valida o payload e retorna erros detalhados
func ValidarRequisicaoWebhookDetalhado(requisicao RequisicaoWebhook) []*ErroValidacaoWebhook {
	var erros []*ErroValidacaoWebhook

	if strings.TrimSpace(requisicao.IdentificadorEvento) == "" {
		erros = append(erros, NewErroValidacaoWebhook("identificadorEvento", "identificadorEvento é obrigatório"))
	}

	if strings.TrimSpace(requisicao.IdentificadorCard) == "" {
		erros = append(erros, NewErroValidacaoWebhook("identificadorCard", "identificadorCard é obrigatório"))
	}

	if strings.TrimSpace(requisicao.EmailCliente) == "" {
		erros = append(erros, NewErroValidacaoWebhook("emailCliente", "emailCliente é obrigatório"))
	} else if !isValidEmail(requisicao.EmailCliente) {
		erros = append(erros, NewErroValidacaoWebhook("emailCliente", "emailCliente inválido"))
	}

	if strings.TrimSpace(requisicao.DataEvento) == "" {
		erros = append(erros, NewErroValidacaoWebhook("dataEvento", "dataEvento é obrigatório"))
	} else if !isValidTimestamp(requisicao.DataEvento) {
		erros = append(erros, NewErroValidacaoWebhook("dataEvento", "dataEvento inválido"))
	}

	return erros
}
