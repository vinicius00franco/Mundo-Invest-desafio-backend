package gestao_clientes

import (
	"context"
	"fmt"
	"time"

	"github.com/MundoInvest/backend/internal/dominio"
	"github.com/MundoInvest/backend/internal/integracao_pipefy"
	"github.com/MundoInvest/backend/internal/shared/config"
	"github.com/MundoInvest/backend/internal/shared/errors"
	"github.com/MundoInvest/backend/internal/shared/mensagens"
)

// ClienteService define a interface para operações de negócio de clientes
type ClienteService interface {
	CriarCliente(ctx context.Context, requisicao RequisicaoCriarCliente) (*Cliente, error)
}

// clienteService implementa a interface ClienteService
type clienteService struct {
	repository      ClienteRepository
	pipefyService   integracao_pipefy.PipefyIntegrationService
	pipeID          string
	eventDispatcher dominio.EventDispatcher
	config          *config.Config
}

// NovoClienteService cria uma nova instância de ClienteService
func NovoClienteService(
	repository ClienteRepository,
	pipefyService integracao_pipefy.PipefyIntegrationService,
	pipeID string,
	eventDispatcher dominio.EventDispatcher,
	cfg *config.Config,
) ClienteService {
	return &clienteService{
		repository:      repository,
		pipefyService:   pipefyService,
		pipeID:          pipeID,
		eventDispatcher: eventDispatcher,
		config:          cfg,
	}
}

// CriarCliente cria um novo cliente seguindo as regras de negócio
func (s *clienteService) CriarCliente(ctx context.Context, requisicao RequisicaoCriarCliente) (*Cliente, error) {
	catalogo := mensagens.ObterCatalogo()
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeouts.Default)
	defer cancel()

	// Validar o payload
	if err := ValidarRequisicaoCriarCliente(requisicao); err != nil {
		return nil, errors.NewValidationError("", err.Error())
	}

	// Gerar identificador externo simulado (card_id)
	identificadorExterno := gerarIdentificadorCardSimulado()

	// Criar entidade Cliente com status inicial
	cliente := Cliente{
		IdentificadorExterno: identificadorExterno,
		Nome:                 requisicao.Nome,
		Email:                requisicao.Email,
		ValorPatrimonio:      requisicao.ValorPatrimonio,
		TipoSolicitacao:      requisicao.TipoSolicitacao,
		Status:               dominio.StatusAguardandoAnalise,
		NivelPrioridade:      "", // Será definido posteriormente pelo webhook
		DataCriacao:          time.Now(),
		DataAtualizacao:      time.Now(),
	}

	// Salvar cliente no banco de dados
	clienteSalvo, err := s.repository.Salvar(ctx, cliente)
	if err != nil {
		return nil, errors.NewServiceError("ClienteService", catalogo.Texto(mensagens.ErrSalvarCliente), err)
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
		return nil, errors.NewIntegrationError("Pipefy", catalogo.Texto(mensagens.ErrCriarCardPipefy), err)
	}

	// Atualizar identificador externo com o card ID retornado
	clienteSalvo.IdentificadorExterno = cardID
	if err := s.repository.Atualizar(ctx, *clienteSalvo); err != nil {
		return nil, errors.NewServiceError("ClienteService", catalogo.Texto(mensagens.ErrAtualizarCliente), err)
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
			catalogo := mensagens.ObterCatalogo()
			fmt.Printf(catalogo.TextoFormatado(mensagens.ErrDespacharEvento), err)
		}
	}

	return clienteSalvo, nil
}

// gerarIdentificadorCardSimulado gera um identificador de card simulado
func gerarIdentificadorCardSimulado() string {
	// Simulação de geração de card_id
	// Em produção, isso seria retornado pelo Pipefy após criar o card
	return fmt.Sprintf("card_%d", time.Now().Unix())
}
