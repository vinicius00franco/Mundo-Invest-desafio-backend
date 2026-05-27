package gestao_clientes

import "time"

// CriarClienteResponse representa a resposta de criação de cliente
type CriarClienteResponse struct {
	Mensagem             string    `json:"mensagem"`
	IdentificadorInterno int64     `json:"identificadorInterno"`
	IdentificadorExterno string    `json:"identificadorExterno"`
	Nome                 string    `json:"nome"`
	Email                string    `json:"email"`
	ValorPatrimonio      float64   `json:"valorPatrimonio"`
	TipoSolicitacao      string    `json:"tipoSolicitacao"`
	Status               string    `json:"status"`
	DataCriacao          time.Time `json:"dataCriacao"`
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
