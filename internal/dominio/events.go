package dominio

import (
	"context"
	"fmt"
	"time"

	"github.com/MundoInvest/backend/internal/shared/logger"
	"github.com/MundoInvest/backend/internal/shared/mensagens"
)

// Event representa um evento de domínio
type Event interface {
	// AggregateID retorna o ID do agregado que gerou o evento
	AggregateID() string
	// EventType retorna o tipo do evento
	EventType() string
	// OccurredAt retorna quando o evento ocorreu
	OccurredAt() time.Time
	// Data retorna os dados do evento
	Data() interface{}
}

// EventHandler define a interface para handlers de eventos
type EventHandler interface {
	// Handle processa um evento de domínio
	Handle(ctx context.Context, event Event) error
	// EventType retorna o tipo de evento que este handler processa
	EventType() string
}

// EventDispatcher define a interface para despachar eventos
type EventDispatcher interface {
	// Dispatch despacha um evento para os handlers apropriados
	Dispatch(ctx context.Context, event Event) error
	// Register registra um handler para um tipo de evento
	Register(handler EventHandler)
	// Subscribe permite que um handler se inscreva em eventos
	Subscribe(eventType string, handler EventHandler)
}

// BaseEvent fornece uma implementação base para eventos
type BaseEvent struct {
	aggregateID string
	eventType   string
	occurredAt  time.Time
	data        interface{}
}

// NewBaseEvent cria um novo evento base
func NewBaseEvent(aggregateID, eventType string, data interface{}) *BaseEvent {
	return &BaseEvent{
		aggregateID: aggregateID,
		eventType:   eventType,
		occurredAt:  time.Now(),
		data:        data,
	}
}

// AggregateID retorna o ID do agregado
func (e *BaseEvent) AggregateID() string {
	return e.aggregateID
}

// EventType retorna o tipo do evento
func (e *BaseEvent) EventType() string {
	return e.eventType
}

// OccurredAt retorna quando o evento ocorreu
func (e *BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// Data retorna os dados do evento
func (e *BaseEvent) Data() interface{} {
	return e.data
}

// InMemoryEventDispatcher implementa um despachante de eventos em memória
type InMemoryEventDispatcher struct {
	handlers map[string][]EventHandler
}

// NewInMemoryEventDispatcher cria um novo despachante de eventos em memória
func NewInMemoryEventDispatcher() *InMemoryEventDispatcher {
	return &InMemoryEventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

// Register registra um handler para um tipo de evento
func (d *InMemoryEventDispatcher) Register(handler EventHandler) {
	eventType := handler.EventType()
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Subscribe permite que um handler se inscreva em eventos
func (d *InMemoryEventDispatcher) Subscribe(eventType string, handler EventHandler) {
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Dispatch despacha um evento para os handlers apropriados
func (d *InMemoryEventDispatcher) Dispatch(ctx context.Context, event Event) error {
	handlers, exists := d.handlers[event.EventType()]
	if !exists {
		return nil
	}

	for _, handler := range handlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("Panic recuperado no handler de evento",
						"event_type", event.EventType(),
						"handler", handler.EventType(),
						"panic", fmt.Sprintf("%v", r),
					)
				}
			}()

			if err := handler.Handle(ctx, event); err != nil {
				catalogo := mensagens.ObterCatalogo()
				logger.Error(catalogo.Texto(mensagens.ErrProcessarEvento),
					"event_type", event.EventType(),
					"handler", handler.EventType(),
					"error", err.Error(),
				)
			}
		}()
	}

	return nil
}
