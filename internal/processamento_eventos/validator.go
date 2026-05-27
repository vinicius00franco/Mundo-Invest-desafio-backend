package processamento_eventos

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/MundoInvest/backend/internal/shared/mensagens"
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
	catalogo := mensagens.ObterCatalogo()
	var erros []string

	// Validar identificadorEvento
	if strings.TrimSpace(requisicao.IdentificadorEvento) == "" {
		erros = append(erros, catalogo.Texto(mensagens.ValIdEventoObrigatorio))
	}

	// Validar identificadorCard
	if strings.TrimSpace(requisicao.IdentificadorCard) == "" {
		erros = append(erros, catalogo.Texto(mensagens.ValIdCardObrigatorio))
	}

	// Validar emailCliente
	if strings.TrimSpace(requisicao.EmailCliente) == "" {
		erros = append(erros, catalogo.Texto(mensagens.ValEmailClienteObrigatorio))
	} else {
		if !isValidEmail(requisicao.EmailCliente) {
			erros = append(erros, catalogo.Texto(mensagens.ValEmailClienteInvalido))
		}
	}

	// Validar dataEvento
	if strings.TrimSpace(requisicao.DataEvento) == "" {
		erros = append(erros, catalogo.Texto(mensagens.ValDataEventoObrigatorio))
	} else {
		if !isValidTimestamp(requisicao.DataEvento) {
			erros = append(erros, catalogo.Texto(mensagens.ValDataEventoInvalido))
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
	catalogo := mensagens.ObterCatalogo()
	if strings.TrimSpace(identificador) == "" {
		return errors.New(catalogo.Texto(mensagens.ValIdEventoObrigatorio))
	}
	return nil
}

// ValidarIdentificadorCard valida o identificador do card
func ValidarIdentificadorCard(identificador string) error {
	catalogo := mensagens.ObterCatalogo()
	if strings.TrimSpace(identificador) == "" {
		return errors.New(catalogo.Texto(mensagens.ValIdCardObrigatorio))
	}
	return nil
}

// ValidarEmailCliente valida o email do cliente
func ValidarEmailCliente(email string) error {
	catalogo := mensagens.ObterCatalogo()
	if strings.TrimSpace(email) == "" {
		return errors.New(catalogo.Texto(mensagens.ValEmailClienteObrigatorio))
	}

	if !isValidEmail(email) {
		return errors.New(catalogo.Texto(mensagens.ValEmailClienteInvalido))
	}

	return nil
}

// ValidarDataEvento valida a data do evento
func ValidarDataEvento(data string) error {
	catalogo := mensagens.ObterCatalogo()
	if strings.TrimSpace(data) == "" {
		return errors.New(catalogo.Texto(mensagens.ValDataEventoObrigatorio))
	}

	if !isValidTimestamp(data) {
		return errors.New(catalogo.Texto(mensagens.ValDataEventoInvalido))
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
	catalogo := mensagens.ObterCatalogo()
	var erros []*ErroValidacaoWebhook

	if strings.TrimSpace(requisicao.IdentificadorEvento) == "" {
		erros = append(erros, NewErroValidacaoWebhook("identificadorEvento", catalogo.Texto(mensagens.ValIdEventoObrigatorio)))
	}

	if strings.TrimSpace(requisicao.IdentificadorCard) == "" {
		erros = append(erros, NewErroValidacaoWebhook("identificadorCard", catalogo.Texto(mensagens.ValIdCardObrigatorio)))
	}

	if strings.TrimSpace(requisicao.EmailCliente) == "" {
		erros = append(erros, NewErroValidacaoWebhook("emailCliente", catalogo.Texto(mensagens.ValEmailClienteObrigatorio)))
	} else if !isValidEmail(requisicao.EmailCliente) {
		erros = append(erros, NewErroValidacaoWebhook("emailCliente", catalogo.Texto(mensagens.ValEmailClienteInvalido)))
	}

	if strings.TrimSpace(requisicao.DataEvento) == "" {
		erros = append(erros, NewErroValidacaoWebhook("dataEvento", catalogo.Texto(mensagens.ValDataEventoObrigatorio)))
	} else if !isValidTimestamp(requisicao.DataEvento) {
		erros = append(erros, NewErroValidacaoWebhook("dataEvento", catalogo.Texto(mensagens.ValDataEventoInvalido)))
	}

	return erros
}
