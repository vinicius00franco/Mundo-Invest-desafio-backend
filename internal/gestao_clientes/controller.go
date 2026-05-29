package gestao_clientes

import (
	"encoding/json"
	"net/http"

	"github.com/MundoInvest/backend/internal/shared/mensagens"
)

// ClienteController manipula as requisições HTTP relacionadas a clientes
type ClienteController struct {
	service ClienteService
}

// NovoClienteController cria uma nova instância de ClienteController
func NovoClienteController(service ClienteService) *ClienteController {
	return &ClienteController{
		service: service,
	}
}

// CriarClienteHandler manipula a requisição POST /clientes
func (c *ClienteController) CriarClienteHandler(w http.ResponseWriter, r *http.Request) {
	catalogo := mensagens.ObterCatalogo()
	mapeadorHTTP := mensagens.NovoMapeadorHTTP()

	// Verificar se o método é POST
	if r.Method != http.MethodPost {
		http.Error(w, catalogo.Texto(mensagens.ErrMetodoNaoPermitido), http.StatusMethodNotAllowed)
		return
	}

	// Parsear o JSON do corpo da requisição
	var requisicao RequisicaoCriarCliente
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requisicao); err != nil {
		http.Error(w, catalogo.Texto(mensagens.ErrProcessarJSON), mapeadorHTTP.StatusPara(mensagens.TipoHTTP))
		return
	}
	defer r.Body.Close()

	// Validar o payload
	if err := ValidarRequisicaoCriarCliente(requisicao); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(mapeadorHTTP.StatusPara(mensagens.TipoValidacao))
		json.NewEncoder(w).Encode(NewErrorResponse(err.Error()))
		return
	}

	// Chamar o serviço para criar o cliente
	cliente, err := c.service.CriarCliente(r.Context(), requisicao)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(mapeadorHTTP.StatusPara(mensagens.TipoServico))
		json.NewEncoder(w).Encode(NewErrorResponse(err.Error()))
		return
	}

	// Retornar o cliente criado com status 201
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := NewCriarClienteResponse(cliente)
	response.Mensagem = catalogo.Texto(mensagens.CliCriadoSucesso)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, catalogo.Texto(mensagens.ErrGerarResposta), mapeadorHTTP.StatusPara(mensagens.TipoHTTP))
		return
	}
}

// RegistrarRotas registra as rotas do controller no router
func (c *ClienteController) RegistrarRotas(mux *http.ServeMux) {
	mux.HandleFunc("/clientes", c.CriarClienteHandler)
}
