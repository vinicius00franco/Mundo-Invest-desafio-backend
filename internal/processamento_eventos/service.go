package processamento_eventos

import (
	"fmt"
	"time"

	"github.com/MundoInvest/backend/internal/dominio"
	"github.com/MundoInvest/backend/internal/gestao_clientes"
	"github.com/MundoInvest/backend/internal/integracao_pipefy"
)

// WebhookService define a interface para operações de negócio de webhooks
type WebhookService interface {
	ProcessarWebhook(request WebhookRequest) error
}

// webhookService implementa a interface WebhookService
type webhookService struct {
	eventoRepository     EventoRepository
	clienteRepository    gestao_clientes.ClienteRepository
	prioridadeCalculator dominio.PrioridadeCalculator
	pipefyClient         integracao_pipefy.PipefyGraphQLClient
}

// NovoWebhookService cria uma nova instância de WebhookService
func NovoWebhookService(
	eventoRepository EventoRepository,
	clienteRepository gestao_clientes.ClienteRepository,
	prioridadeCalculator dominio.PrioridadeCalculator,
	pipefyClient integracao_pipefy.PipefyGraphQLClient,
) WebhookService {
	return &webhookService{
		eventoRepository:     eventoRepository,
		clienteRepository:    clienteRepository,
		prioridadeCalculator: prioridadeCalculator,
		pipefyClient:         pipefyClient,
	}
}

// ProcessarWebhook processa um webhook do Pipefy de forma idempotente
func (s *webhookService) ProcessarWebhook(request WebhookRequest) error {
	// Validar o payload
	if err := ValidarWebhookRequest(request); err != nil {
		return fmt.Errorf("erro de validação: %w", err)
	}

	// Verificar idempotência: se o evento já foi processado, retornar sucesso
	foiProcessado, err := s.eventoRepository.VerificarFoiProcessado(request.IdentificadorEvento)
	if err != nil {
		return fmt.Errorf("erro ao verificar idempotência: %w", err)
	}

	if foiProcessado {
		// Evento já processado, retornar sucesso (idempotência)
		return nil
	}

	// Buscar cliente por email
	cliente, err := s.clienteRepository.BuscarPorEmail(request.ClienteEmail)
	if err != nil {
		return fmt.Errorf("cliente não encontrado com email %s: %w", request.ClienteEmail, err)
	}

	// Calcular nível de prioridade baseado no patrimônio
	nivelPrioridade := s.prioridadeCalculator.CalcularNivelPrioridade(cliente.ValorPatrimonio)

	// Atualizar cliente com novo status e prioridade
	cliente.Status = "Processado"
	cliente.NivelPrioridade = nivelPrioridade
	cliente.DataAtualizacao = time.Now()

	if err := s.clienteRepository.Atualizar(*cliente); err != nil {
		return fmt.Errorf("erro ao atualizar cliente: %w", err)
	}

	// Estruturar mutation updateCard
	fieldsAttributes := []integracao_pipefy.FieldAttribute{
		{
			FieldID: "nivel_prioridade_field_id",
			Values:  []string{nivelPrioridade},
		},
	}

	mutation, err := s.pipefyClient.EstruturarMutationUpdateCard(request.IdentificadorCard, fieldsAttributes)
	if err != nil {
		return fmt.Errorf("erro ao estruturar mutation updateCard: %w", err)
	}

	// Simular envio da mutation (sem requisição real)
	// Em produção, aqui seria enviada a mutation para o Pipefy
	_, _ = s.pipefyClient.ExecutarMutation(mutation)

	// Salvar evento como processado
	timestampEvento, err := time.Parse(time.RFC3339, request.DataEvento)
	if err != nil {
		return fmt.Errorf("erro ao parsear data do evento: %w", err)
	}

	evento := Evento{
		IdentificadorEvento: request.IdentificadorEvento,
		IdentificadorCard:   request.IdentificadorCard,
		EmailCliente:        request.ClienteEmail,
		TimestampEvento:     timestampEvento,
		FoiProcessado:       true,
		DataCriacao:         time.Now(),
		DataAtualizacao:     time.Now(),
	}

	if _, err := s.eventoRepository.Salvar(evento); err != nil {
		return fmt.Errorf("erro ao salvar evento: %w", err)
	}

	return nil
}
