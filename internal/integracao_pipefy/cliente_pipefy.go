package integracao_pipefy

import (
	"fmt"
	"time"

	"github.com/MundoInvest/backend/internal/shared/mensagens"
)

// PipefyGraphQLClient define a interface para integração com Pipefy via GraphQL
type PipefyGraphQLClient interface {
	EstruturarMutationCreateCard(pipeID string, fieldsAttributes []FieldAttribute) (string, error)
	EstruturarMutationUpdateCard(cardID string, fieldsAttributes []FieldAttribute) (string, error)
	ExecutarMutation(mutation string) (string, error)
}

// FieldAttribute representa um atributo de campo para mutations do Pipefy
type FieldAttribute struct {
	FieldID string   `json:"fieldId"`
	Values  []string `json:"values"`
}

// pipefyGraphQLClient implementa a interface PipefyGraphQLClient
type pipefyGraphQLClient struct {
	apiToken string
	apiURL   string
}

// NovoPipefyGraphQLClient cria uma nova instância de PipefyGraphQLClient
func NovoPipefyGraphQLClient(apiToken, apiURL string) PipefyGraphQLClient {
	return &pipefyGraphQLClient{
		apiToken: apiToken,
		apiURL:   apiURL,
	}
}

// EstruturarMutationCreateCard estrutura a mutation GraphQL para criar card no Pipefy
// Fonte: https://api-docs.pipefy.com/reference/mutations/#createcard
func (c *pipefyGraphQLClient) EstruturarMutationCreateCard(pipeID string, fieldsAttributes []FieldAttribute) (string, error) {
	catalogo := mensagens.ObterCatalogo()
	// Se pipeID estiver vazio, retorna card ID simulado para ambiente de desenvolvimento
	if pipeID == "" {
		return fmt.Sprintf("card_%d", time.Now().Unix()), nil
	}

	if len(fieldsAttributes) == 0 {
		return "", fmt.Errorf(catalogo.Texto(mensagens.ErrCriarCardPipefy))
	}

	mutation := fmt.Sprintf(`mutation {
		createCard(input: {
			pipe_id: "%s"
			fields_attributes: [
				%s
			]
		}) {
			card {
				id
				title
				url
			}
			success
		}
	}`, pipeID, c.formatarFieldsAttributes(fieldsAttributes))

	return mutation, nil
}

// EstruturarMutationUpdateCard estrutura a mutation GraphQL para atualizar card no Pipefy
// Fonte: https://api-docs.pipefy.com/reference/mutations/#updatecard
func (c *pipefyGraphQLClient) EstruturarMutationUpdateCard(cardID string, fieldsAttributes []FieldAttribute) (string, error) {
	catalogo := mensagens.ObterCatalogo()
	if cardID == "" {
		return "", fmt.Errorf(catalogo.Texto(mensagens.ErrAtualizarCardPipefy))
	}

	if len(fieldsAttributes) == 0 {
		return "", fmt.Errorf(catalogo.Texto(mensagens.ErrAtualizarCardPipefy))
	}

	mutation := fmt.Sprintf(`mutation {
		updateCard(input: {
			card_id: "%s"
			fields_attributes: [
				%s
			]
		}) {
			card {
				id
				title
				url
			}
			success
		}
	}`, cardID, c.formatarFieldsAttributes(fieldsAttributes))

	return mutation, nil
}

// ExecutarMutation executa a mutation no Pipefy (simulado)
func (c *pipefyGraphQLClient) ExecutarMutation(mutation string) (string, error) {
	// Simulação de execução da mutation
	// Em produção, aqui seria feita a requisição HTTP para a API do Pipefy
	// Simular retorno de card ID curto para testes
	return fmt.Sprintf("card_%d", time.Now().Unix()), nil
}

// formatarFieldsAttributes formata os atributos de campos para a mutation
func (c *pipefyGraphQLClient) formatarFieldsAttributes(fields []FieldAttribute) string {
	var formattedFields []string
	for _, field := range fields {
		values := c.formatarValues(field.Values)
		formattedField := fmt.Sprintf(`{ field_id: "%s", values: [%s] }`, field.FieldID, values)
		formattedFields = append(formattedFields, formattedField)
	}
	return joinStrings(formattedFields, ", ")
}

// formatarValues formata os valores para a mutation
func (c *pipefyGraphQLClient) formatarValues(values []string) string {
	var formattedValues []string
	for _, value := range values {
		formattedValue := fmt.Sprintf(`"%s"`, value)
		formattedValues = append(formattedValues, formattedValue)
	}
	return joinStrings(formattedValues, ", ")
}

// joinStrings junta strings com um separador
func joinStrings(strings []string, separator string) string {
	if len(strings) == 0 {
		return ""
	}
	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += separator + strings[i]
	}
	return result
}
