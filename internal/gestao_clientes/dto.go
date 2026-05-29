package gestao_clientes

import "time"

// CriarClienteResponse representa a resposta de criação de cliente
type CriarClienteResponse struct {
	Mensagem             string    `json:"mensagem"`
	IdentificadorInterno int64     `json:"identificador_interno"`
	IdentificadorExterno string    `json:"identificador_externo"`
	Nome                 string    `json:"cliente_nome"`
	Email                string    `json:"cliente_email"`
	ValorPatrimonio      float64   `json:"valor_patrimonio"`
	TipoSolicitacao      string    `json:"tipo_solicitacao"`
	Status               string    `json:"status"`
	DataCriacao          time.Time `json:"data_criacao"`
}

// ErrorResponse representa uma resposta de erro padronizada
type ErrorResponse struct {
	Erro string `json:"erro"`
}

// NewCriarClienteResponse cria uma nova resposta de criação de cliente
func NewCriarClienteResponse(cliente *Cliente) CriarClienteResponse {
	return CriarClienteResponse{
		Mensagem:             "", // Será definido pelo controller
		IdentificadorInterno: cliente.IdentificadorInterno,
		IdentificadorExterno: cliente.IdentificadorExterno,
		Nome:                 cliente.Nome,
		Email:                cliente.Email,
		ValorPatrimonio:      cliente.ValorPatrimonio,
		TipoSolicitacao:      cliente.TipoSolicitacao,
		Status:               cliente.Status,
		DataCriacao:          cliente.DataCriacao,
	}
}

// NewErrorResponse cria uma nova resposta de erro
func NewErrorResponse(erro string) ErrorResponse {
	return ErrorResponse{
		Erro: erro,
	}
}
