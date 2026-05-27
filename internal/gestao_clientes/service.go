package gestao_clientes

import (
	"fmt"
	"time"

	"github.com/MundoInvest/backend/internal/integracao_pipefy"
)

// ClienteService define a interface para operações de negócio de clientes
type ClienteService interface {
	CriarCliente(request CriarClienteRequest) (*Cliente, error)
}

// clienteService implementa a interface ClienteService
type clienteService struct {
	repository   ClienteRepository
	pipefyClient integracao_pipefy.PipefyGraphQLClient
	pipeID       string
}

// NovoClienteService cria uma nova instância de ClienteService
func NovoClienteService(repository ClienteRepository, pipefyClient integracao_pipefy.PipefyGraphQLClient, pipeID string) ClienteService {
	return &clienteService{
		repository:   repository,
		pipefyClient: pipefyClient,
		pipeID:       pipeID,
	}
}

// CriarCliente cria um novo cliente seguindo as regras de negócio
func (s *clienteService) CriarCliente(request CriarClienteRequest) (*Cliente, error) {
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
		Status:               "Aguardando Análise",
		NivelPrioridade:      "", // Será definido posteriormente pelo webhook
		DataCriacao:          time.Now(),
		DataAtualizacao:      time.Now(),
	}

	// Salvar cliente no banco de dados
	clienteSalvo, err := s.repository.Salvar(cliente)
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar cliente: %w", err)
	}

	// Estruturar mutation createCard usando o cliente Pipefy
	fieldsAttributes := []integracao_pipefy.FieldAttribute{
		{
			FieldID: "nome_field_id",
			Values:  []string{clienteSalvo.Nome},
		},
		{
			FieldID: "email_field_id",
			Values:  []string{clienteSalvo.Email},
		},
		{
			FieldID: "patrimonio_field_id",
			Values:  []string{fmt.Sprintf("%.2f", clienteSalvo.ValorPatrimonio)},
		},
		{
			FieldID: "tipo_solicitacao_field_id",
			Values:  []string{clienteSalvo.TipoSolicitacao},
		},
	}

	mutation, err := s.pipefyClient.EstruturarMutationCreateCard(s.pipeID, fieldsAttributes)
	if err != nil {
		return nil, fmt.Errorf("erro ao estruturar mutation createCard: %w", err)
	}

	// Simular envio da mutation (sem requisição real)
	// Em produção, aqui seria enviada a mutation para o Pipefy
	_, _ = s.pipefyClient.ExecutarMutation(mutation)

	return clienteSalvo, nil
}

// gerarIdentificadorCardSimulado gera um identificador de card simulado
func gerarIdentificadorCardSimulado() string {
	// Simulação de geração de card_id
	// Em produção, isso seria retornado pelo Pipefy após criar o card
	return fmt.Sprintf("card_%d", time.Now().UnixNano())
}
