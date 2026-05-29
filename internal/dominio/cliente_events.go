package dominio

import (
	"time"
)

// Tipos de eventos de cliente
const (
	ClienteCriadoEvent       = "cliente.criado"
	ClienteAtualizadoEvent    = "cliente.atualizado"
	ClientePrioridadeCalculadaEvent = "cliente.prioridade_calculada"
	ClienteProcessadoEvent   = "cliente.processado"
)

// ClienteCriadoData representa os dados do evento de cliente criado
type ClienteCriadoData struct {
	ClienteID      string
	Nome           string
	Email          string
	ValorPatrimonio float64
	TipoSolicitacao string
	DataCriacao     time.Time
}

// ClienteAtualizadoData representa os dados do evento de cliente atualizado
type ClienteAtualizadoData struct {
	ClienteID       string
	CamposAlterados []string
	DataAtualizacao  time.Time
}

// ClientePrioridadeCalculadaData representa os dados do evento de prioridade calculada
type ClientePrioridadeCalculadaData struct {
	ClienteID        string
	ValorPatrimonio  float64
	NivelPrioridade  string
	DataCalculo      time.Time
}

// ClienteProcessadoData representa os dados do evento de cliente processado
type ClienteProcessadoData struct {
	ClienteID        string
	NivelPrioridade  string
	DataProcessamento time.Time
}

// NewClienteCriadoEvent cria um novo evento de cliente criado
func NewClienteCriadoEvent(clienteID string, data ClienteCriadoData) Event {
	return NewBaseEvent(clienteID, ClienteCriadoEvent, data)
}

// NewClienteAtualizadoEvent cria um novo evento de cliente atualizado
func NewClienteAtualizadoEvent(clienteID string, data ClienteAtualizadoData) Event {
	return NewBaseEvent(clienteID, ClienteAtualizadoEvent, data)
}

// NewClientePrioridadeCalculadaEvent cria um novo evento de prioridade calculada
func NewClientePrioridadeCalculadaEvent(clienteID string, data ClientePrioridadeCalculadaData) Event {
	return NewBaseEvent(clienteID, ClientePrioridadeCalculadaEvent, data)
}

// NewClienteProcessadoEvent cria um novo evento de cliente processado
func NewClienteProcessadoEvent(clienteID string, data ClienteProcessadoData) Event {
	return NewBaseEvent(clienteID, ClienteProcessadoEvent, data)
}
