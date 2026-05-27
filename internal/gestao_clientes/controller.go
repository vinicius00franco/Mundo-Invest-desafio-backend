package gestao_clientes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MundoInvest/backend/internal/integracao_pipefy"
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
	// Verificar se o método é POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear o JSON do corpo da requisição
	var request CriarClienteRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao parsear JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validar o payload
	if err := ValidarCriarClienteRequest(request); err != nil {
		http.Error(w, fmt.Sprintf("Erro de validação: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Chamar o serviço para criar o cliente
	cliente, err := c.service.CriarCliente(request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar cliente: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Retornar o cliente criado com status 201
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"mensagem":              "Cliente criado com sucesso",
		"identificador_interno": cliente.IdentificadorInterno,
		"identificador_externo": cliente.IdentificadorExterno,
		"nome":                  cliente.Nome,
		"email":                 cliente.Email,
		"valor_patrimonio":      cliente.ValorPatrimonio,
		"tipo_solicitacao":      cliente.TipoSolicitacao,
		"status":                cliente.Status,
		"data_criacao":          cliente.DataCriacao,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao gerar resposta: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

// RegistrarRotas registra as rotas do controller no router
func (c *ClienteController) RegistrarRotas(mux *http.ServeMux) {
	mux.HandleFunc("/clientes", c.CriarClienteHandler)
}

// NovoClienteControllerComDB cria uma nova instância de ClienteController com banco de dados
func NovoClienteControllerComDB(db *sql.DB, pipeID string) *ClienteController {
	repository := NovoClienteRepository(db)
	pipefyClient := integracao_pipefy.NovoPipefyGraphQLClient("", "")
	service := NovoClienteService(repository, pipefyClient, pipeID)
	return NovoClienteController(service)
}
