# Cenários de Validação - Webhook Card Updated

## Cenário: Processar webhook sem event_id

**Descrição**: Tentar processar webhook sem fornecer o event_id.

**Pré-condições**:
- Endpoint POST /webhooks/pipefy/card-updated está disponível

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated sem event_id:
   ```json
   {
     "card_id": "card_456",
     "cliente_email": "joao.silva@example.com",
     "timestamp": "2026-05-18T12:00:00Z"
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica campo obrigatório faltando
- Evento não é salvo no banco

---

## Cenário: Processar webhook sem card_id

**Descrição**: Tentar processar webhook sem fornecer o card_id.

**Pré-condições**:
- Endpoint POST /webhooks/pipefy/card-updated está disponível

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated sem card_id:
   ```json
   {
     "event_id": "evt_123",
     "cliente_email": "joao.silva@example.com",
     "timestamp": "2026-05-18T12:00:00Z"
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica campo obrigatório faltando
- Evento não é salvo no banco

---

## Cenário: Processar webhook sem cliente_email

**Descrição**: Tentar processar webhook sem fornecer o cliente_email.

**Pré-condições**:
- Endpoint POST /webhooks/pipefy/card-updated está disponível

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated sem cliente_email:
   ```json
   {
     "event_id": "evt_123",
     "card_id": "card_456",
     "timestamp": "2026-05-18T12:00:00Z"
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica campo obrigatório faltando
- Evento não é salvo no banco

---

## Cenário: Processar webhook sem timestamp

**Descrição**: Tentar processar webhook sem fornecer o timestamp.

**Pré-condições**:
- Endpoint POST /webhooks/pipefy/card-updated está disponível

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated sem timestamp:
   ```json
   {
     "event_id": "evt_123",
     "card_id": "card_456",
     "cliente_email": "joao.silva@example.com"
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica campo obrigatório faltando
- Evento não é salvo no banco

---

## Cenário: Processar webhook com event_id vazio

**Descrição**: Tentar processar webhook com event_id em branco.

**Pré-condições**:
- Endpoint POST /webhooks/pipefy/card-updated está disponível

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated com event_id vazio:
   ```json
   {
     "event_id": "",
     "card_id": "card_456",
     "cliente_email": "joao.silva@example.com",
     "timestamp": "2026-05-18T12:00:00Z"
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica campo inválido
- Evento não é salvo no banco

---

## Cenário: Processar webhook com cliente_email inválido

**Descrição**: Tentar processar webhook com e-mail em formato inválido.

**Pré-condições**:
- Endpoint POST /webhooks/pipefy/card-updated está disponível

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated com email inválido:
   ```json
   {
     "event_id": "evt_123",
     "card_id": "card_456",
     "cliente_email": "email-invalido",
     "timestamp": "2026-05-18T12:00:00Z"
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica e-mail inválido
- Evento não é salvo no banco

---

## Cenário: Processar webhook com card_id vazio

**Descrição**: Tentar processar webhook com card_id em branco.

**Pré-condições**:
- Endpoint POST /webhooks/pipefy/card-updated está disponível

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated com card_id vazio:
   ```json
   {
     "event_id": "evt_123",
     "card_id": "",
     "cliente_email": "joao.silva@example.com",
     "timestamp": "2026-05-18T12:00:00Z"
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica campo inválido
- Evento não é salvo no banco

---

## Cenário: Processar webhook com cliente_email vazio

**Descrição**: Tentar processar webhook com cliente_email em branco.

**Pré-condições**:
- Endpoint POST /webhooks/pipefy/card-updated está disponível

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated com cliente_email vazio:
   ```json
   {
     "event_id": "evt_123",
     "card_id": "card_456",
     "cliente_email": "",
     "timestamp": "2026-05-18T12:00:00Z"
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica campo inválido
- Evento não é salvo no banco

---

## Cenário: Processar webhook com timestamp inválido

**Descrição**: Tentar processar webhook com timestamp em formato inválido.

**Pré-condições**:
- Endpoint POST /webhooks/pipefy/card-updated está disponível

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated com timestamp inválido:
   ```json
   {
     "event_id": "evt_123",
     "card_id": "card_456",
     "cliente_email": "joao.silva@example.com",
     "timestamp": "invalid-timestamp"
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica timestamp inválido
- Evento não é salvo no banco
