# Cenários de Falha - Webhook Card Updated

## Cenário: Processar webhook com cliente não encontrado

**Descrição**: Tentar processar webhook para cliente que não existe no banco.

**Pré-condições**:
- Endpoint POST /webhooks/pipefy/card-updated está disponível
- Cliente com email "inexistente@example.com" não existe no banco

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated com email inexistente:
   ```json
   {
     "event_id": "evt_999",
     "card_id": "card_999",
     "cliente_email": "inexistente@example.com",
     "timestamp": "2026-05-18T15:00:00Z"
   }
   ```

**Resultados Esperados**:
- HTTP 404 Not Found
- Mensagem de erro indica cliente não encontrado
- Evento não é salvo no banco
- Nenhuma atualização é realizada

**Dados de Verificação**:
- Response contém erro "CLIENTE_NAO_ENCONTRADO"
- Banco não contém registro do evento

---

## Cenário: Processar webhook com erro de banco de dados

**Descrição**: Tentar processar webhook quando banco de dados está indisponível.

**Pré-condições**:
- Endpoint POST /webhooks/pipefy/card-updated está disponível
- Banco de dados está indisponível ou em erro

**Passos**:
1. Enviar requisição POST para /webhooks/pipefy/card-updated:
   ```json
   {
     "event_id": "evt_888",
     "card_id": "card_888",
     "cliente_email": "joao.silva@example.com",
     "timestamp": "2026-05-18T16:00:00Z"
   }
   ```

**Resultados Esperados**:
- HTTP 500 Internal Server Error
- Mensagem de erro indica problema no banco de dados
- Evento não é salvo no banco
- Nenhuma atualização é realizada

**Dados de Verificação**:
- Response contém erro "ERRO_BANCO_DADOS"
- Sistema permanece em estado consistente
