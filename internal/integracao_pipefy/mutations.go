package integracao_pipefy

import (
	"fmt"
)

// Mutations contém as mutations GraphQL do Pipefy
type Mutations struct {
	client PipefyGraphQLClient
}

// NovoMutations cria uma nova instância de Mutations
func NovoMutations(client PipefyGraphQLClient) *Mutations {
	return &Mutations{
		client: client,
	}
}

// CreateCardMutation cria a mutation para criar um card no Pipefy
// Fonte: https://api-docs.pipefy.com/reference/mutations/#createcard
// NOTA: FieldIDs são identificadores de campos do Pipefy e devem ser configurados conforme o setup do pipe
func (m *Mutations) CreateCardMutation(pipeID string, nome, email, tipoSolicitacao string, valorPatrimonio float64) (string, error) {
	fieldsAttributes := []FieldAttribute{
		{
			FieldID: "nome_field_id",
			Values:  []string{nome},
		},
		{
			FieldID: "email_field_id",
			Values:  []string{email},
		},
		{
			FieldID: "tipo_solicitacao_field_id",
			Values:  []string{tipoSolicitacao},
		},
		{
			FieldID: "valor_patrimonio_field_id",
			Values:  []string{fmt.Sprintf("%.2f", valorPatrimonio)},
		},
	}

	return m.client.EstruturarMutationCreateCard(pipeID, fieldsAttributes)
}

// UpdateCardMutation cria a mutation para atualizar um card no Pipefy
// Fonte: https://api-docs.pipefy.com/reference/mutations/#updatecard
// NOTA: FieldIDs são identificadores de campos do Pipefy e devem ser configurados conforme o setup do pipe
func (m *Mutations) UpdateCardMutation(cardID string, nivelPrioridade string) (string, error) {
	fieldsAttributes := []FieldAttribute{
		{
			FieldID: "nivel_prioridade_field_id",
			Values:  []string{nivelPrioridade},
		},
	}

	return m.client.EstruturarMutationUpdateCard(cardID, fieldsAttributes)
}

// CreateCardMutationComCamposPersonalizados cria a mutation para criar um card com campos personalizados
func (m *Mutations) CreateCardMutationComCamposPersonalizados(pipeID string, campos map[string]string) (string, error) {
	var fieldsAttributes []FieldAttribute

	for fieldID, valor := range campos {
		fieldsAttributes = append(fieldsAttributes, FieldAttribute{
			FieldID: fieldID,
			Values:  []string{valor},
		})
	}

	return m.client.EstruturarMutationCreateCard(pipeID, fieldsAttributes)
}

// UpdateCardMutationComCamposPersonalizados cria a mutation para atualizar um card com campos personalizados
func (m *Mutations) UpdateCardMutationComCamposPersonalizados(cardID string, campos map[string]string) (string, error) {
	var fieldsAttributes []FieldAttribute

	for fieldID, valor := range campos {
		fieldsAttributes = append(fieldsAttributes, FieldAttribute{
			FieldID: fieldID,
			Values:  []string{valor},
		})
	}

	return m.client.EstruturarMutationUpdateCard(cardID, fieldsAttributes)
}

// AddCardRelationMutation cria a mutation para adicionar uma relação entre cards
// Fonte: https://api-docs.pipefy.com/reference/mutations/#addcardrelation
func (m *Mutations) AddCardRelationMutation(sourceCardID, destinationCardID, relationTypeID string) (string, error) {
	mutation := fmt.Sprintf(`mutation {
		addCardRelation(input: {
			source_card_id: "%s"
			destination_card_id: "%s"
			relation_type_id: "%s"
		}) {
			success
		}
	}`, sourceCardID, destinationCardID, relationTypeID)

	return mutation, nil
}

// DeleteCardMutation cria a mutation para deletar um card
// Fonte: https://api-docs.pipefy.com/reference/mutations/#deletecard
func (m *Mutations) DeleteCardMutation(cardID string) (string, error) {
	mutation := fmt.Sprintf(`mutation {
		deleteCard(input: {
			card_id: "%s"
		}) {
			success
		}
	}`, cardID)

	return mutation, nil
}

// MoveCardToPhaseMutation cria a mutation para mover um card para outra fase
// Fonte: https://api-docs.pipefy.com/reference/mutations/#movecardtophase
func (m *Mutations) MoveCardToPhaseMutation(cardID, phaseID string) (string, error) {
	mutation := fmt.Sprintf(`mutation {
		moveCardToPhase(input: {
			card_id: "%s"
			phase_id: "%s"
		}) {
			card {
				id
				current_phase {
					id
					name
				}
			}
			success
		}
	}`, cardID, phaseID)

	return mutation, nil
}
