package processamento_eventos

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/MundoInvest/backend/internal/dominio"
	"github.com/MundoInvest/backend/internal/gestao_clientes"
	"github.com/MundoInvest/backend/internal/integracao_pipefy"
)

// WebhookController manipula as requisições HTTP relacionadas a webhooks
type WebhookController struct {
	service WebhookService
}

// NovoWebhookController cria uma nova instância de WebhookController
func NovoWebhookController(service WebhookService) *WebhookController {
	return &WebhookController{
		service: service,
	}
}

// ProcessarWebhookHandler manipula a requisição POST /webhooks/pipefy/card-updated
func (c *WebhookController) ProcessarWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar se o método é POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear o JSON do corpo da requisição
	var request WebhookRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao parsear JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validar o payload
	if err := ValidarWebhookRequest(request); err != nil {
		http.Error(w, fmt.Sprintf("Erro de validação: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Chamar o serviço para processar o webhook
	if err := c.service.ProcessarWebhook(r.Context(), request); err != nil {
		// Verificar se é erro de cliente não encontrado
		errMsg := fmt.Sprintf("%v", err)
		if contains(errMsg, "cliente não encontrado") {
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
		"mensagem":             "Webhook processado com sucesso",
		"identificador_evento": request.IdentificadorEvento,
		"identificador_card":   request.IdentificadorCard,
		"cliente_email":        request.ClienteEmail,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao gerar resposta: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

// RegistrarRotas registra as rotas do controller no router
func (c *WebhookController) RegistrarRotas(mux *http.ServeMux) {
	mux.HandleFunc("/webhooks/pipefy/card-updated", c.ProcessarWebhookHandler)
}

// NovoWebhookControllerComDB cria uma nova instância de WebhookController com banco de dados
func NovoWebhookControllerComDB(db *sql.DB) *WebhookController {
	eventoRepository := NovoEventoRepository(db)
	clienteRepository := gestao_clientes.NovoClienteRepository(db)
	prioridadeCalculator := dominio.NovoPrioridadeCalculator()
	pipefyClient := integracao_pipefy.NovoPipefyGraphQLClient("", "")
	service := NovoWebhookService(eventoRepository, clienteRepository, prioridadeCalculator, pipefyClient)
	return NovoWebhookController(service)
}

// contains verifica se uma string contém uma substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
