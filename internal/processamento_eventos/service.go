package processamento_eventos

import (
	"context"
	"time"

	"github.com/MundoInvest/backend/internal/dominio"
	"github.com/MundoInvest/backend/internal/gestao_clientes"
	"github.com/MundoInvest/backend/internal/integracao_pipefy"
	"github.com/MundoInvest/backend/internal/shared/config"
	"github.com/MundoInvest/backend/internal/shared/errors"
)

// WebhookService define a interface para operações de negócio de webhooks
type WebhookService interface {
	ProcessarWebhook(ctx context.Context, requisicao RequisicaoWebhook) error
}

// webhookService implementa a interface WebhookService
type webhookService struct {
	eventoRepository      EventoRepository
	clienteRepository     gestao_clientes.ClienteRepository
	calculadoraPrioridade dominio.CalculadoraPrioridade
	pipefyClient          integracao_pipefy.PipefyGraphQLClient
	config                *config.Config
}

// NovoWebhookService cria uma nova instância de WebhookService
func NovoWebhookService(
	eventoRepository EventoRepository,
	clienteRepository gestao_clientes.ClienteRepository,
	calculadoraPrioridade dominio.CalculadoraPrioridade,
	pipefyClient integracao_pipefy.PipefyGraphQLClient,
	cfg *config.Config,
) WebhookService {
	return &webhookService{
		eventoRepository:      eventoRepository,
		clienteRepository:     clienteRepository,
		calculadoraPrioridade: calculadoraPrioridade,
		pipefyClient:          pipefyClient,
		config:                cfg,
	}
}

// ProcessarWebhook processa um webhook do Pipefy de forma idempotente
func (s *webhookService) ProcessarWebhook(ctx context.Context, requisicao RequisicaoWebhook) error {
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeouts.Default)
	defer cancel()

	// Validar o payload
	if err := ValidarRequisicaoWebhook(requisicao); err != nil {
		return errors.NewValidationError("", err.Error())
	}

	// Verificar idempotência: se o evento já foi processado, retornar sucesso
	foiProcessado, err := s.eventoRepository.VerificarFoiProcessado(ctx, requisicao.IdentificadorEvento)
	if err != nil {
		return errors.NewServiceError("WebhookService", "erro ao verificar idempotência", err)
	}

	if foiProcessado {
		// Evento já processado, retornar sucesso (idempotência)
		return nil
	}

	// Buscar cliente por email
	cliente, err := s.clienteRepository.BuscarPorEmail(ctx, requisicao.EmailCliente)
	if err != nil {
		return errors.NewServiceError("WebhookService", "cliente não encontrado com email "+requisicao.EmailCliente, err)
	}

	// Calcular nível de prioridade baseado no patrimônio
	nivelPrioridade := s.calculadoraPrioridade.CalcularNivelPrioridade(cliente.ValorPatrimonio)

	// Atualizar cliente com novo status e prioridade
	cliente.Status = dominio.StatusProcessado
	cliente.NivelPrioridade = nivelPrioridade
	cliente.DataAtualizacao = time.Now()

	if err := s.clienteRepository.Atualizar(ctx, *cliente); err != nil {
		return errors.NewServiceError("WebhookService", "erro ao atualizar cliente", err)
	}

	// Estruturar mutation updateCard
	fieldsAttributes := []integracao_pipefy.FieldAttribute{
		{
			FieldID: "nivel_prioridade_field_id",
			Values:  []string{nivelPrioridade},
		},
	}
	// NOTA: FieldIDs são identificadores de campos do Pipefy (snake_case) e devem ser configurados conforme o setup do pipe

	mutation, err := s.pipefyClient.EstruturarMutationUpdateCard(requisicao.IdentificadorCard, fieldsAttributes)
	if err != nil {
		return errors.NewIntegrationError("Pipefy", "erro ao estruturar mutation updateCard", err)
	}

	// Simular envio da mutation (sem requisição real)
	// Em produção, aqui seria enviada a mutation para o Pipefy
	_, _ = s.pipefyClient.ExecutarMutation(mutation)

	// Salvar evento como processado
	timestampEvento, err := time.Parse(time.RFC3339, requisicao.DataEvento)
	if err != nil {
		return errors.NewServiceError("WebhookService", "erro ao parsear data do evento", err)
	}

	evento := Evento{
		IdentificadorEvento: requisicao.IdentificadorEvento,
		IdentificadorCard:   requisicao.IdentificadorCard,
		EmailCliente:        requisicao.EmailCliente,
		TimestampEvento:     timestampEvento,
		FoiProcessado:       true,
		DataCriacao:         time.Now(),
		DataAtualizacao:     nil, // Field is nullable
	}

	if _, err := s.eventoRepository.Salvar(ctx, evento); err != nil {
		return errors.NewServiceError("WebhookService", "erro ao salvar evento", err)
	}

	return nil
}
