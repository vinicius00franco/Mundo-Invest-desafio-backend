package processamento_eventos

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

// WebhookRequest representa o payload para processamento de webhook
type WebhookRequest struct {
	IdentificadorEvento string `json:"identificador_evento"`
	IdentificadorCard   string `json:"identificador_card"`
	ClienteEmail        string `json:"cliente_email"`
	DataEvento          string `json:"data_evento"`
}

// ValidarWebhookRequest valida o payload de webhook
func ValidarWebhookRequest(request WebhookRequest) error {
	var erros []string

	// Validar identificador_evento
	if strings.TrimSpace(request.IdentificadorEvento) == "" {
		erros = append(erros, "identificador_evento é obrigatório")
	}

	// Validar identificador_card
	if strings.TrimSpace(request.IdentificadorCard) == "" {
		erros = append(erros, "identificador_card é obrigatório")
	}

	// Validar cliente_email
	if strings.TrimSpace(request.ClienteEmail) == "" {
		erros = append(erros, "cliente_email é obrigatório")
	} else {
		if !isValidEmail(request.ClienteEmail) {
			erros = append(erros, "cliente_email inválido")
		}
	}

	// Validar data_evento
	if strings.TrimSpace(request.DataEvento) == "" {
		erros = append(erros, "data_evento é obrigatório")
	} else {
		if !isValidTimestamp(request.DataEvento) {
			erros = append(erros, "data_evento inválido")
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
		return errors.New("identificador_evento é obrigatório")
	}
	return nil
}

// ValidarIdentificadorCard valida o identificador do card
func ValidarIdentificadorCard(identificador string) error {
	if strings.TrimSpace(identificador) == "" {
		return errors.New("identificador_card é obrigatório")
	}
	return nil
}

// ValidarClienteEmail valida o email do cliente
func ValidarClienteEmail(email string) error {
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
		return errors.New("data_evento é obrigatório")
	}

	if !isValidTimestamp(data) {
		return errors.New("data_evento inválido")
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
func ValidarWebhookRequestDetalhado(request WebhookRequest) []*ErroValidacaoWebhook {
	var erros []*ErroValidacaoWebhook

	if strings.TrimSpace(request.IdentificadorEvento) == "" {
		erros = append(erros, NewErroValidacaoWebhook("identificador_evento", "identificador_evento é obrigatório"))
	}

	if strings.TrimSpace(request.IdentificadorCard) == "" {
		erros = append(erros, NewErroValidacaoWebhook("identificador_card", "identificador_card é obrigatório"))
	}

	if strings.TrimSpace(request.ClienteEmail) == "" {
		erros = append(erros, NewErroValidacaoWebhook("cliente_email", "cliente_email é obrigatório"))
	} else if !isValidEmail(request.ClienteEmail) {
		erros = append(erros, NewErroValidacaoWebhook("cliente_email", "cliente_email inválido"))
	}

	if strings.TrimSpace(request.DataEvento) == "" {
		erros = append(erros, NewErroValidacaoWebhook("data_evento", "data_evento é obrigatório"))
	} else if !isValidTimestamp(request.DataEvento) {
		erros = append(erros, NewErroValidacaoWebhook("data_evento", "data_evento inválido"))
	}

	return erros
}
