package integracao_pipefy

import (
	"fmt"
)

// Mutacoes contém as mutations GraphQL do Pipefy
type Mutacoes struct {
	cliente PipefyGraphQLClient
}

// NovoMutacoes cria uma nova instância de Mutacoes
func NovoMutacoes(cliente PipefyGraphQLClient) *Mutacoes {
	return &Mutacoes{
		cliente: cliente,
	}
}

// CreateCardMutation cria a mutation para criar um card no Pipefy
// Fonte: https://api-docs.pipefy.com/reference/mutations/#createcard
// NOTA: FieldIDs são identificadores de campos do Pipefy e devem ser configurados conforme o setup do pipe
func (m *Mutacoes) CreateCardMutation(pipeID string, nome, email, tipoSolicitacao string, valorPatrimonio float64) (string, error) {
	atributosCampo := []FieldAttribute{
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

	return m.cliente.EstruturarMutationCreateCard(pipeID, atributosCampo)
}

// UpdateCardMutation cria a mutation para atualizar um card no Pipefy
// Fonte: https://api-docs.pipefy.com/reference/mutations/#updatecard
// NOTA: FieldIDs são identificadores de campos do Pipefy e devem ser configurados conforme o setup do pipe
func (m *Mutacoes) UpdateCardMutation(cardID string, nivelPrioridade string) (string, error) {
	atributosCampo := []FieldAttribute{
		{
			FieldID: "nivel_prioridade_field_id",
			Values:  []string{nivelPrioridade},
		},
	}

	return m.cliente.EstruturarMutationUpdateCard(cardID, atributosCampo)
}

// CreateCardMutationComCamposPersonalizados cria a mutation para criar um card com campos personalizados
func (m *Mutacoes) CreateCardMutationComCamposPersonalizados(pipeID string, campos map[string]string) (string, error) {
	var atributosCampo []FieldAttribute

	for fieldID, valor := range campos {
		atributosCampo = append(atributosCampo, FieldAttribute{
			FieldID: fieldID,
			Values:  []string{valor},
		})
	}

	return m.cliente.EstruturarMutationCreateCard(pipeID, atributosCampo)
}

// UpdateCardMutationComCamposPersonalizados cria a mutation para atualizar um card com campos personalizados
func (m *Mutacoes) UpdateCardMutationComCamposPersonalizados(cardID string, campos map[string]string) (string, error) {
	var atributosCampo []FieldAttribute

	for fieldID, valor := range campos {
		atributosCampo = append(atributosCampo, FieldAttribute{
			FieldID: fieldID,
			Values:  []string{valor},
		})
	}

	return m.cliente.EstruturarMutationUpdateCard(cardID, atributosCampo)
}

// AddCardRelationMutation cria uma mutation para adicionar uma relação entre cards
// Fonte: https://api-docs.pipefy.com/reference/mutations/#addcardrelation
func (m *Mutacoes) AddCardRelationMutation(sourceCardID, destinationCardID, relationTypeID string) (string, error) {
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

// DeleteCardMutation cria uma mutation para deletar um card
// Fonte: https://api-docs.pipefy.com/reference/mutations/#deletecard
func (m *Mutacoes) DeleteCardMutation(cardID string) (string, error) {
	mutation := fmt.Sprintf(`mutation {
			deleteCard(input: {
				card_id: "%s"
			}) {
				success
			}
		}`, cardID)

	return mutation, nil
}

// MoveCardToPhaseMutation cria uma mutation para mover um card para outra fase
// Fonte: https://api-docs.pipefy.com/reference/mutations/#movecardtophase
func (m *Mutacoes) MoveCardToPhaseMutation(cardID, phaseID string) (string, error) {
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
