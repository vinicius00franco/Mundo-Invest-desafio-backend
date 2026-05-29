package integracao_pipefy

import (
	"context"
	"fmt"

	"github.com/MundoInvest/backend/internal/shared/config"
	"github.com/MundoInvest/backend/internal/shared/errors"
	"github.com/MundoInvest/backend/internal/shared/mensagens"
)

// PipefyIntegrationService define a interface para serviços de integração com Pipefy
type PipefyIntegrationService interface {
	CriarCardCliente(ctx context.Context, pipeID string, dados CardClienteData) (string, error)
	AtualizarCardPrioridade(ctx context.Context, cardID string, nivelPrioridade string) error
}

// CardClienteData representa os dados para criar um card de cliente
type CardClienteData struct {
	Nome            string
	Email           string
	ValorPatrimonio float64
	TipoSolicitacao string
}

// pipefyIntegrationService implementa a interface PipefyIntegrationService
type pipefyIntegrationService struct {
	client PipefyGraphQLClient
	config *config.Config
}

// NewPipefyIntegrationService cria uma nova instância de PipefyIntegrationService
func NewPipefyIntegrationService(client PipefyGraphQLClient, cfg *config.Config) PipefyIntegrationService {
	return &pipefyIntegrationService{
		client: client,
		config: cfg,
	}
}

// CriarCardCliente cria um card no Pipefy para um cliente
func (s *pipefyIntegrationService) CriarCardCliente(ctx context.Context, pipeID string, dados CardClienteData) (string, error) {
	catalogo := mensagens.ObterCatalogo()
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeouts.ExternalAPI)
	defer cancel()

	fieldsAttributes := []FieldAttribute{
		{
			FieldID: "nome_field_id",
			Values:  []string{dados.Nome},
		},
		{
			FieldID: "email_field_id",
			Values:  []string{dados.Email},
		},
		{
			FieldID: "patrimonio_field_id",
			Values:  []string{fmt.Sprintf("%.2f", dados.ValorPatrimonio)},
		},
		{
			FieldID: "tipo_solicitacao_field_id",
			Values:  []string{dados.TipoSolicitacao},
		},
	}
	// NOTA: FieldIDs são identificadores de campos do Pipefy (snake_case) e devem ser configurados conforme o setup do pipe

	mutation, err := s.client.EstruturarMutationCreateCard(pipeID, fieldsAttributes)
	if err != nil {
		return "", errors.NewIntegrationError("Pipefy", catalogo.Texto(mensagens.ErrEstruturarMutacaoCriar), err)
	}

	// Executar mutation
	resultado, err := s.client.ExecutarMutation(mutation)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", errors.NewTimeoutError("criar card no Pipefy", "timeout ao executar chamada externa")
		}
		return "", errors.NewIntegrationError("Pipefy", catalogo.Texto(mensagens.ErrExecutarMutacaoCriar), err)
	}

	return resultado, nil
}

// AtualizarCardPrioridade atualiza a prioridade de um card no Pipefy
func (s *pipefyIntegrationService) AtualizarCardPrioridade(ctx context.Context, cardID string, nivelPrioridade string) error {
	catalogo := mensagens.ObterCatalogo()
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeouts.ExternalAPI)
	defer cancel()

	fieldsAttributes := []FieldAttribute{
		{
			FieldID: "nivel_prioridade_field_id",
			Values:  []string{nivelPrioridade},
		},
	}
	// NOTA: FieldIDs são identificadores de campos do Pipefy (snake_case) e devem ser configurados conforme o setup do pipe

	mutation, err := s.client.EstruturarMutationUpdateCard(cardID, fieldsAttributes)
	if err != nil {
		return errors.NewIntegrationError("Pipefy", catalogo.Texto(mensagens.ErrEstruturarMutacaoAtualizar), err)
	}

	// Executar mutation
	_, err = s.client.ExecutarMutation(mutation)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.NewTimeoutError("atualizar card no Pipefy", "timeout ao executar chamada externa")
		}
		return errors.NewIntegrationError("Pipefy", catalogo.Texto(mensagens.ErrExecutarMutacaoAtualizar), err)
	}

	return nil
}
