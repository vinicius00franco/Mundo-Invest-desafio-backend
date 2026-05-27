package processamento_eventos

import (
	"encoding/json"
	"net/http"

	"github.com/MundoInvest/backend/internal/shared/errors"
	"github.com/MundoInvest/backend/internal/shared/mensagens"
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
	catalogo := mensagens.ObterCatalogo()
	mapeadorHTTP := mensagens.NovoMapeadorHTTP()

	// Verificar se o método é POST
	if r.Method != http.MethodPost {
		http.Error(w, catalogo.Texto(mensagens.ErrMetodoNaoPermitido), http.StatusMethodNotAllowed)
		return
	}

	// Parsear o JSON do corpo da requisição
	var requisicao RequisicaoWebhook
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requisicao); err != nil {
		http.Error(w, catalogo.Texto(mensagens.ErrProcessarJSON), mapeadorHTTP.StatusPara(mensagens.TipoHTTP))
		return
	}
	defer r.Body.Close()

	// Validar o payload
	if err := ValidarRequisicaoWebhook(requisicao); err != nil {
		http.Error(w, err.Error(), mapeadorHTTP.StatusPara(mensagens.TipoValidacao))
		return
	}

	// Chamar o serviço para processar o webhook
	if err := c.service.ProcessarWebhook(r.Context(), requisicao); err != nil {
		// Verificar se é erro de cliente não encontrado
		if _, ok := err.(*errors.NotFoundError); ok {
			http.Error(w, catalogo.Texto(mensagens.ErrClienteNaoEncontrado), http.StatusNotFound)
			return
		}
		http.Error(w, catalogo.Texto(mensagens.ErrProcessarWebhook), mapeadorHTTP.StatusPara(mensagens.TipoServico))
		return
	}

	// Retornar sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"mensagem":            catalogo.Texto(mensagens.WebhookProcessado),
		"identificadorEvento": requisicao.IdentificadorEvento,
		"identificadorCard":   requisicao.IdentificadorCard,
		"emailCliente":        requisicao.EmailCliente,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, catalogo.Texto(mensagens.ErrGerarResposta), mapeadorHTTP.StatusPara(mensagens.TipoHTTP))
		return
	}
}

// RegistrarRotas registra as rotas do controller no router
func (c *EventoController) RegistrarRotas(mux *http.ServeMux) {
	mux.HandleFunc("/webhooks/pipefy/card-updated", c.ProcessarWebhookHandler)
}
