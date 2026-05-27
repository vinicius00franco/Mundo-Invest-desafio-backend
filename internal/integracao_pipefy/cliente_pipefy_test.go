package integracao_pipefy

import (
	"context"
	"testing"

	"github.com/MundoInvest/backend/internal/shared/config"
)

func TestPipefyGraphQLClient_EstruturarMutationCreateCard(t *testing.T) {
	client := NovoPipefyGraphQLClient("test_token", "https://api.pipefy.com/graphql")

	tests := []struct {
		name         string
		pipeID       string
		fields       []FieldAttribute
		esperadoErro bool
	}{
		{
			name:   "Mutation válida",
			pipeID: "pipe_123",
			fields: []FieldAttribute{
				{FieldID: "nome_field_id", Values: []string{"João Silva"}},
				{FieldID: "email_field_id", Values: []string{"joao@example.com"}},
			},
			esperadoErro: false,
		},
		{
			name:         "Pipe ID vazio",
			pipeID:       "",
			fields:       []FieldAttribute{},
			esperadoErro: true,
		},
		{
			name:         "Fields vazios",
			pipeID:       "pipe_123",
			fields:       []FieldAttribute{},
			esperadoErro: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mutation, err := client.EstruturarMutationCreateCard(tt.pipeID, tt.fields)

			if tt.esperadoErro {
				if err == nil {
					t.Error("Esperado erro, mas não houve erro")
				}
			} else {
				if err != nil {
					t.Errorf("Não esperado erro: %v", err)
				}
				if mutation == "" {
					t.Error("Mutation não deve ser vazia")
				}
			}
		})
	}
}

func TestPipefyGraphQLClient_EstruturarMutationUpdateCard(t *testing.T) {
	client := NovoPipefyGraphQLClient("test_token", "https://api.pipefy.com/graphql")

	tests := []struct {
		name         string
		cardID       string
		fields       []FieldAttribute
		esperadoErro bool
	}{
		{
			name:   "Mutation válida",
			cardID: "card_123",
			fields: []FieldAttribute{
				{FieldID: "nivel_prioridade_field_id", Values: []string{"prioridade_alta"}},
			},
			esperadoErro: false,
			// NOTA: FieldIDs são identificadores de campos do Pipefy (snake_case) e devem ser configurados conforme o setup do pipe
		},
		{
			name:         "Card ID vazio",
			cardID:       "",
			fields:       []FieldAttribute{},
			esperadoErro: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mutation, err := client.EstruturarMutationUpdateCard(tt.cardID, tt.fields)

			if tt.esperadoErro {
				if err == nil {
					t.Error("Esperado erro, mas não houve erro")
				}
			} else {
				if err != nil {
					t.Errorf("Não esperado erro: %v", err)
				}
				if mutation == "" {
					t.Error("Mutation não deve ser vazia")
				}
			}
		})
	}
}

func TestPipefyGraphQLClient_ExecutarMutation(t *testing.T) {
	client := NovoPipefyGraphQLClient("test_token", "https://api.pipefy.com/graphql")

	mutation := "mutation { test }"
	resultado, err := client.ExecutarMutation(mutation)

	if err != nil {
		t.Errorf("Erro ao executar mutation: %v", err)
	}

	if resultado == "" {
		t.Error("Resultado não deve ser vazio")
	}
}

func TestPipefyIntegrationService_CriarCardCliente(t *testing.T) {
	client := NovoPipefyGraphQLClient("test_token", "https://api.pipefy.com/graphql")
	cfg := config.Load()
	service := NewPipefyIntegrationService(client, cfg)

	dados := CardClienteData{
		Nome:            "João Silva",
		Email:           "joao@example.com",
		ValorPatrimonio: 150000.00,
		TipoSolicitacao: "abertura_conta",
	}

	cardID, err := service.CriarCardCliente(context.Background(), "pipe_123", dados)

	if err != nil {
		t.Errorf("Erro ao criar card: %v", err)
	}

	if cardID == "" {
		t.Error("Card ID não deve ser vazio")
	}
}

func TestPipefyIntegrationService_AtualizarCardPrioridade(t *testing.T) {
	client := NovoPipefyGraphQLClient("test_token", "https://api.pipefy.com/graphql")
	cfg := config.Load()
	service := NewPipefyIntegrationService(client, cfg)

	err := service.AtualizarCardPrioridade(context.Background(), "card_123", "prioridade_alta")

	if err != nil {
		t.Errorf("Erro ao atualizar card: %v", err)
	}
}

func TestFieldAttribute(t *testing.T) {
	field := FieldAttribute{
		FieldID: "test_field",
		Values:  []string{"value1", "value2"},
	}

	if field.FieldID != "test_field" {
		t.Errorf("FieldID esperado 'test_field', obtido '%s'", field.FieldID)
	}

	if len(field.Values) != 2 {
		t.Errorf("Values esperado 2, obtido %d", len(field.Values))
	}
}
