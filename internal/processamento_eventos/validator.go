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
	IdentificadorEvento string `json:"event_id"`
	IdentificadorCard   string `json:"card_id"`
	EmailCliente        string `json:"cliente_email"`
	DataEvento          string `json:"timestamp"`
}

// ValidarWebhookRequest valida o payload de webhook
func ValidarRequisicaoWebhook(requisicao RequisicaoWebhook) error {
	var erros []string

	// Validar identificadorEvento
	if strings.TrimSpace(requisicao.IdentificadorEvento) == "" {
		erros = append(erros, "event_id é obrigatório")
	}

	// Validar identificadorCard
	if strings.TrimSpace(requisicao.IdentificadorCard) == "" {
		erros = append(erros, "card_id é obrigatório")
	}

	// Validar emailCliente
	if strings.TrimSpace(requisicao.EmailCliente) == "" {
		erros = append(erros, "cliente_email é obrigatório")
	} else {
		if !isValidEmail(requisicao.EmailCliente) {
			erros = append(erros, "cliente_email inválido")
		}
	}

	// Validar dataEvento
	if strings.TrimSpace(requisicao.DataEvento) == "" {
		erros = append(erros, "timestamp é obrigatório")
	} else {
		if !isValidTimestamp(requisicao.DataEvento) {
			erros = append(erros, "timestamp inválido")
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
		return errors.New("event_id é obrigatório")
	}
	return nil
}

// ValidarIdentificadorCard valida o identificador do card
func ValidarIdentificadorCard(identificador string) error {
	if strings.TrimSpace(identificador) == "" {
		return errors.New("card_id é obrigatório")
	}
	return nil
}

// ValidarEmailCliente valida o email do cliente
func ValidarEmailCliente(email string) error {
	if strings.TrimSpace(email) == "" {
		return errors.New("cliente_email é obrigatório")
	}

	if !isValidEmail(email) {
		return errors.New("cliente_email inválido")
	}

	return nil
}

// ValidarDataEvento valida a data do evento
func ValidarDataEvento(data string) error {
	if strings.TrimSpace(data) == "" {
		return errors.New("timestamp é obrigatório")
	}

	if !isValidTimestamp(data) {
		return errors.New("timestamp inválido")
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
		erros = append(erros, NewErroValidacaoWebhook("event_id", "event_id é obrigatório"))
	}

	if strings.TrimSpace(requisicao.IdentificadorCard) == "" {
		erros = append(erros, NewErroValidacaoWebhook("card_id", "card_id é obrigatório"))
	}

	if strings.TrimSpace(requisicao.EmailCliente) == "" {
		erros = append(erros, NewErroValidacaoWebhook("cliente_email", "cliente_email é obrigatório"))
	} else if !isValidEmail(requisicao.EmailCliente) {
		erros = append(erros, NewErroValidacaoWebhook("cliente_email", "cliente_email inválido"))
	}

	if strings.TrimSpace(requisicao.DataEvento) == "" {
		erros = append(erros, NewErroValidacaoWebhook("timestamp", "timestamp é obrigatório"))
	} else if !isValidTimestamp(requisicao.DataEvento) {
		erros = append(erros, NewErroValidacaoWebhook("timestamp", "timestamp inválido"))
	}

	return erros
}
