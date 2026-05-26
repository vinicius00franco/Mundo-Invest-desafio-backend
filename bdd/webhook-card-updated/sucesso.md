# Cenários de Sucesso - Webhook Card Updated

## Cenário: Processar webhook com patrimônio alto (prioridade alta)

**Descrição**: Processar webhook de atualização de card para cliente com patrimônio >= 200.000.

**Pré-condições**:
- Cliente existe no banco com email "joao.silva@example.com"
- Cliente tem valor_patrimonio de 250.000
- Endpoint POST /webhooks/pipefy/card-updated está disponível

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated:
   ```json
   {
     "event_id": "evt_123",
     "card_id": "card_456",
     "cliente_email": "joao.silva@example.com",
     "timestamp": "2026-05-18T12:00:00Z"
   }
   ```

**Resultados Esperados**:
- HTTP 200 OK
- Evento é marcado como processado no banco
- Prioridade do cliente é definida como "prioridade_alta"
- Status do cliente é atualizado para "Processado"
- Mutation GraphQL updateCard é estruturada no código
- Response contém dados do processamento

**Dados de Verificação**:
- Evento com event_id "evt_123" salvo como processado
- Prioridade: "prioridade_alta"
- Status: "Processado"

---

## Cenário: Processar webhook com patrimônio baixo (prioridade normal)

**Descrição**: Processar webhook de atualização de card para cliente com patrimônio < 200.000.

**Pré-condições**:
- Cliente existe no banco com email "maria.santos@example.com"
- Cliente tem valor_patrimonio de 150.000
- Endpoint POST /webhooks/pipefy/card-updated está disponível

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated:
   ```json
   {
     "event_id": "evt_456",
     "card_id": "card_789",
     "cliente_email": "maria.santos@example.com",
     "timestamp": "2026-05-18T13:00:00Z"
   }
   ```

**Resultados Esperados**:
- HTTP 200 OK
- Evento é marcado como processado no banco
- Prioridade do cliente é definida como "prioridade_normal"
- Status do cliente é atualizado para "Processado"

**Dados de Verificação**:
- Evento com event_id "evt_456" salvo como processado
- Prioridade: "prioridade_normal"
- Status: "Processado"

---

## Cenário: Processar webhook com patrimônio no limite (prioridade alta)

**Descrição**: Processar webhook de atualização de card para cliente com patrimônio exatamente 200.000.

**Pré-condições**:
- Cliente existe no banco com email "pedro.oliveira@example.com"
- Cliente tem valor_patrimonio de 200.000
- Endpoint POST /webhooks/pipefy/card-updated está disponível

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated:
   ```json
   {
     "event_id": "evt_789",
     "card_id": "card_101",
     "cliente_email": "pedro.oliveira@example.com",
     "timestamp": "2026-05-18T14:00:00Z"
   }
   ```

**Resultados Esperados**:
- HTTP 200 OK
- Evento é marcado como processado no banco
- Prioridade do cliente é definida como "prioridade_alta"
- Status do cliente é atualizado para "Processado"

**Dados de Verificação**:
- Evento com event_id "evt_789" salvo como processado
- Prioridade: "prioridade_alta"
- Status: "Processado"

---

## Cenário: Processar webhook duplicado (idempotência)

**Descrição**: Processar webhook com event_id que já foi processado anteriormente.

**Pré-condições**:
- Evento com event_id "evt_123" já foi processado
- Endpoint POST /webhooks/pipefy/card-updated está disponível

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated com o mesmo event_id:
   ```json
   {
     "event_id": "evt_123",
     "card_id": "card_456",
     "cliente_email": "joao.silva@example.com",
     "timestamp": "2026-05-18T12:00:00Z"
   }
   ```

**Resultados Esperados**:
- HTTP 200 OK
- Webhook não é processado novamente
- Mensagem indica que o evento já foi processado
- Status e prioridade do cliente permanecem inalterados

**Dados de Verificação**:
- Response contém mensagem de idempotência
- Cliente mantém status e prioridade anteriores
