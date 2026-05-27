package dominio

import (
	"context"
	"fmt"
	"log"
)

// LoggingEventHandler é um handler que loga eventos
type LoggingEventHandler struct{}

// NewLoggingEventHandler cria um novo handler de logging
func NewLoggingEventHandler() *LoggingEventHandler {
	return &LoggingEventHandler{}
}

// Handle processa o evento logando seus dados
func (h *LoggingEventHandler) Handle(ctx context.Context, event Event) error {
	log.Printf("[EVENT] %s - AggregateID: %s - Data: %+v", 
		event.EventType(), 
		event.AggregateID(), 
		event.Data())
	return nil
}

// EventType retorna o tipo de evento que este handler processa
func (h *LoggingEventHandler) EventType() string {
	return "*" // Processa todos os eventos
}

// ClienteCriadoEventHandler processa eventos de cliente criado
type ClienteCriadoEventHandler struct {
	dispatcher EventDispatcher
}

// NewClienteCriadoEventHandler cria um novo handler para cliente criado
func NewClienteCriadoEventHandler(dispatcher EventDispatcher) *ClienteCriadoEventHandler {
	return &ClienteCriadoEventHandler{
		dispatcher: dispatcher,
	}
}

// Handle processa o evento de cliente criado
func (h *ClienteCriadoEventHandler) Handle(ctx context.Context, event Event) error {
	data, ok := event.Data().(ClienteCriadoData)
	if !ok {
		return fmt.Errorf("dados inválidos para evento %s", event.EventType())
	}

	// Aqui poderia implementar lógica específica para cliente criado
	// Por exemplo: enviar email de boas-vindas, notificar equipe, etc.
	log.Printf("[CLIENTE_CRIADO] Cliente %s (%s) criado com patrimônio %.2f", 
		data.Nome, data.Email, data.ValorPatrimonio)

	return nil
}

// EventType retorna o tipo de evento
func (h *ClienteCriadoEventHandler) EventType() string {
	return ClienteCriadoEvent
}

// ClientePrioridadeCalculadaEventHandler processa eventos de prioridade calculada
type ClientePrioridadeCalculadaEventHandler struct{}

// NewClientePrioridadeCalculadaEventHandler cria um novo handler para prioridade calculada
func NewClientePrioridadeCalculadaEventHandler() *ClientePrioridadeCalculadaEventHandler {
	return &ClientePrioridadeCalculadaEventHandler{}
}

// Handle processa o evento de prioridade calculada
func (h *ClientePrioridadeCalculadaEventHandler) Handle(ctx context.Context, event Event) error {
	data, ok := event.Data().(ClientePrioridadeCalculadaData)
	if !ok {
		return fmt.Errorf("dados inválidos para evento %s", event.EventType())
	}

	log.Printf("[PRIORIDADE_CALCULADA] Cliente %s - Patrimônio: %.2f - Prioridade: %s", 
		data.ClienteID, data.ValorPatrimonio, data.NivelPrioridade)

	return nil
}

// EventType retorna o tipo de evento
func (h *ClientePrioridadeCalculadaEventHandler) EventType() string {
	return ClientePrioridadeCalculadaEvent
}
