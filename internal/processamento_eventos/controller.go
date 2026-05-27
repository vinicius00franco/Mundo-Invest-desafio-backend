package processamento_eventos

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MundoInvest/backend/internal/shared/errors"
)

// EventoController manipula as requisições HTTP relacionadas a eventos
type EventoController struct {
	service WebhookService
}

// NovoEventoController cria uma nova instância de EventoController
func NovoEventoController(service WebhookService) *EventoController {
	return &EventoController{
		service: service,
	}
}

// ProcessarWebhookHandler manipula a requisição POST /webhooks/pipefy/card-updated
func (c *EventoController) ProcessarWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar se o método é POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear o JSON do corpo da requisição
	var requisicao RequisicaoWebhook
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requisicao); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao parsear JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validar o payload
	if err := ValidarRequisicaoWebhook(requisicao); err != nil {
		http.Error(w, fmt.Sprintf("Erro de validação: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Chamar o serviço para processar o webhook
	if err := c.service.ProcessarWebhook(r.Context(), requisicao); err != nil {
		// Verificar se é erro de cliente não encontrado
		if _, ok := err.(*errors.NotFoundError); ok {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Erro ao processar webhook: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Retornar sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"mensagem":            "Webhook processado com sucesso",
		"identificadorEvento": requisicao.IdentificadorEvento,
		"identificadorCard":   requisicao.IdentificadorCard,
		"emailCliente":        requisicao.EmailCliente,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao gerar resposta: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

// RegistrarRotas registra as rotas do controller no router
func (c *EventoController) RegistrarRotas(mux *http.ServeMux) {
	mux.HandleFunc("/webhooks/pipefy/card-updated", c.ProcessarWebhookHandler)
}
