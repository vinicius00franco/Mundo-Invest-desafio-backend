# Padrão de Nomenclatura JSON - camelCase

## Visão Geral

A API do Mundo Invest adota o padrão **camelCase** para nomenclatura de campos JSON, seguindo as melhores práticas de APIs REST modernas e convenções JavaScript/TypeScript.

## Justificativa

### Por que camelCase?

1. **Padrão da Indústria**: A maioria das APIs modernas (Google, Facebook, Twitter, etc.) usa camelCase
2. **Compatibilidade JavaScript**: camelCase é o padrão nativo de JavaScript/TypeScript
3. **Consistência**: Facilita integração com frontends modernos (React, Vue, Angular)
4. **Legibilidade**: Melhor legibilidade para desenvolvedores web
5. **Ferramentas**: Maior suporte em ferramentas de desenvolvimento e documentação

## Padrão Adotado

### Request (Entrada)

Todos os campos de entrada devem usar camelCase:

```json
{
  "nome": "João Silva",
  "email": "joao.silva@example.com",
  "tipoSolicitacao": "Atualização cadastral",
  "valorPatrimonio": 250000
}
```

### Response (Saída)

Todos os campos de resposta devem usar camelCase:

```json
{
  "mensagem": "Cliente criado com sucesso",
  "identificadorInterno": 1,
  "identificadorExterno": "card_123",
  "nome": "João Silva",
  "email": "joao.silva@example.com",
  "valorPatrimonio": 250000,
  "tipoSolicitacao": "Atualização cadastral",
  "status": "Aguardando Análise",
  "dataCriacao": "2026-05-27T10:00:00Z"
}
```

### Webhook Payload

```json
{
  "identificadorEvento": "evt_12345",
  "identificadorCard": "card_67890",
  "clienteEmail": "joao.silva@example.com",
  "dataEvento": "2026-05-27T10:00:00Z"
}
```

## Regras de Conversão

### Go Struct Tags

No código Go, use tags `json` com camelCase:

```go
type CriarClienteRequest struct {
    Nome            string  `json:"nome"`
    Email           string  `json:"email"`
    ValorPatrimonio float64 `json:"valorPatrimonio"`
    TipoSolicitacao string  `json:"tipoSolicitacao"`
}

type CriarClienteResponse struct {
    Mensagem             string    `json:"mensagem"`
    IdentificadorInterno int64     `json:"identificadorInterno"`
    IdentificadorExterno string    `json:"identificadorExterno"`
    Nome                 string    `json:"nome"`
    Email                string    `json:"email"`
    ValorPatrimonio      float64   `json:"valorPatrimonio"`
    TipoSolicitacao      string    `json:"tipoSolicitacao"`
    Status               string    `json:"status"`
    DataCriacao          time.Time `json:"dataCriacao"`
}
```

## Mensagens de Erro

As mensagens de erro também devem usar camelCase para referência aos campos:

```json
{
  "erro": "erro de validação: tipoSolicitacao é obrigatório"
}
```

```json
{
  "erro": "erro de validação: valorPatrimonio deve ser positivo"
}
```

## Exceções

Não há exceções ao padrão camelCase. Todos os campos JSON devem seguir este padrão.

## Migração

### Histórico

- **Antes**: Mix de camelCase e snake_case
- **Depois**: 100% camelCase

### Mudanças Realizadas

1. **Structs Go**: Atualizadas todas as tags `json` para camelCase
2. **Validações**: Mensagens de erro atualizadas para usar camelCase
3. **Testes**: Scripts de integração atualizados para usar camelCase
4. **Documentação**: README e documentação atualizados

### Compatibilidade

Esta mudança é **breaking change** para clientes que consumiam a API com o padrão anterior. Clientes precisam atualizar seu código para usar camelCase.

## Boas Práticas

### Para Desenvolvedores Backend

1. Sempre usar camelCase em tags `json` dos structs
2. Manter consistência com nomenclatura de banco de dados (snake_case) vs JSON (camelCase)
3. Atualizar testes quando adicionar novos campos
4. Documentar novos campos com exemplos em camelCase

### Para Desenvolvedores Frontend

1. Esperar todos os campos em camelCase
2. Não fazer conversão automática de snake_case para camelCase
3. Usar os campos exatamente como retornados pela API
4. Atualizar validações de frontend para usar camelCase

## Referências

- [Google JSON Style Guide](https://google.github.io/styleguide/jsoncstyleguide.xml)
- [Microsoft REST API Guidelines](https://github.com/Microsoft/api-guidelines/blob/master/Guidelines.md)
- [OpenAPI Specification](https://swagger.io/specification/)

## Exemplos Completos

### POST /clientes

**Request:**
```bash
curl -X POST http://localhost:8080/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "João Silva",
    "email": "joao.silva@example.com",
    "tipoSolicitacao": "Atualização cadastral",
    "valorPatrimonio": 250000
  }'
```

**Response (201):**
```json
{
  "mensagem": "Cliente criado com sucesso",
  "identificadorInterno": 1,
  "identificadorExterno": "card_123",
  "nome": "João Silva",
  "email": "joao.silva@example.com",
  "valorPatrimonio": 250000,
  "tipoSolicitacao": "Atualização cadastral",
  "status": "Aguardando Análise",
  "dataCriacao": "2026-05-27T10:00:00Z"
}
```

### POST /webhooks/pipefy/card-updated

**Request:**
```bash
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
  -H "Content-Type: application/json" \
  -d '{
    "identificadorEvento": "evt_12345",
    "identificadorCard": "card_67890",
    "clienteEmail": "joao.silva@example.com",
    "dataEvento": "2026-05-27T10:00:00Z"
  }'
```

**Response (200):**
```json
{
  "mensagem": "Webhook processado com sucesso",
  "identificadorEvento": "evt_12345",
  "identificadorCard": "card_67890",
  "clienteEmail": "joao.silva@example.com"
}
```

## Conclusão

A adoção do padrão camelCase para JSON alinha a API do Mundo Invest com as melhores práticas da indústria, facilitando integração com sistemas modernos e melhorando a experiência dos desenvolvedores.
