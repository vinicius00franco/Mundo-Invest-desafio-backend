package gestao_clientes

import (
	"context"
	"fmt"
	"time"

	"github.com/MundoInvest/backend/internal/dominio"
	"github.com/MundoInvest/backend/internal/integracao_pipefy"
)

// ClienteService define a interface para operações de negócio de clientes
type ClienteService interface {
	CriarCliente(ctx context.Context, request CriarClienteRequest) (*Cliente, error)
}

// clienteService implementa a interface ClienteService
type clienteService struct {
	repository      ClienteRepository
	pipefyService   integracao_pipefy.PipefyIntegrationService
	pipeID          string
	eventDispatcher dominio.EventDispatcher
}

// NovoClienteService cria uma nova instância de ClienteService
func NovoClienteService(
	repository ClienteRepository,
	pipefyService integracao_pipefy.PipefyIntegrationService,
	pipeID string,
	eventDispatcher dominio.EventDispatcher,
) ClienteService {
	return &clienteService{
		repository:      repository,
		pipefyService:   pipefyService,
		pipeID:          pipeID,
		eventDispatcher: eventDispatcher,
	}
}

// CriarCliente cria um novo cliente seguindo as regras de negócio
func (s *clienteService) CriarCliente(ctx context.Context, request CriarClienteRequest) (*Cliente, error) {
	// Validar o payload
	if err := ValidarCriarClienteRequest(request); err != nil {
		return nil, fmt.Errorf("erro de validação: %w", err)
	}

	// Gerar identificador externo simulado (card_id)
	identificadorExterno := gerarIdentificadorCardSimulado()

	// Criar entidade Cliente com status inicial
	cliente := Cliente{
		IdentificadorExterno: identificadorExterno,
		Nome:                 request.Nome,
		Email:                request.Email,
		ValorPatrimonio:      request.ValorPatrimonio,
		TipoSolicitacao:      request.TipoSolicitacao,
		Status:               dominio.StatusAguardandoAnalise,
		NivelPrioridade:      "", // Será definido posteriormente pelo webhook
		DataCriacao:          time.Now(),
		DataAtualizacao:      time.Now(),
	}

	// Salvar cliente no banco de dados
	clienteSalvo, err := s.repository.Salvar(cliente)
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar cliente: %w", err)
	}

	// Integrar com Pipefy usando o serviço dedicado
	cardData := integracao_pipefy.CardClienteData{
		Nome:            clienteSalvo.Nome,
		Email:           clienteSalvo.Email,
		ValorPatrimonio: clienteSalvo.ValorPatrimonio,
		TipoSolicitacao: clienteSalvo.TipoSolicitacao,
	}

	cardID, err := s.pipefyService.CriarCardCliente(ctx, s.pipeID, cardData)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar card no Pipefy: %w", err)
	}

	// Atualizar identificador externo com o card ID retornado
	clienteSalvo.IdentificadorExterno = cardID
	if err := s.repository.Atualizar(*clienteSalvo); err != nil {
		return nil, fmt.Errorf("erro ao atualizar cliente com card ID: %w", err)
	}

	// Emitir evento de domínio
	if s.eventDispatcher != nil {
		evento := dominio.NewClienteCriadoEvent(
			fmt.Sprintf("%d", clienteSalvo.IdentificadorInterno),
			dominio.ClienteCriadoData{
				ClienteID:       fmt.Sprintf("%d", clienteSalvo.IdentificadorInterno),
				Nome:            clienteSalvo.Nome,
				Email:           clienteSalvo.Email,
				ValorPatrimonio: clienteSalvo.ValorPatrimonio,
				TipoSolicitacao: clienteSalvo.TipoSolicitacao,
				DataCriacao:     clienteSalvo.DataCriacao,
			},
		)
		if err := s.eventDispatcher.Dispatch(ctx, evento); err != nil {
			// Log error mas não falhar a operação
			fmt.Printf("Erro ao despachar evento: %v\n", err)
		}
	}

	return clienteSalvo, nil
}

// gerarIdentificadorCardSimulado gera um identificador de card simulado
func gerarIdentificadorCardSimulado() string {
	// Simulação de geração de card_id
	// Em produção, isso seria retornado pelo Pipefy após criar o card
	return fmt.Sprintf("card_%d", time.Now().UnixNano())
}
