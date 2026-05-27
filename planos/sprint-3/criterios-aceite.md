# Critérios de Aceite - Sprint 3: API de Webhook

## Critérios de Aceite

- [ ] Endpoint POST /webhooks/pipefy/card-updated responde a requisições
- [ ] Webhook válido processa cliente corretamente
- [ ] NivelPrioridade é calculado corretamente
- [ ] Idempotencia funciona (identificador_evento duplicado não reprocessa)
- [ ] Cliente não encontrado retorna HTTP 404
- [ ] Mutation updateCard está estruturada no código
- [ ] Nomenclatura segue linguagem ubíqua
- [ ] Testes manuais passam

## Validação

Para validar que a sprint foi concluída com sucesso:

### 1. Validação do Endpoint
```bash
# Teste com webhook válido
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
  -H "Content-Type: application/json" \
  -d '{
    "identificador_evento": "evt_12345",
    "identificador_card": "card_67890",
    "cliente_email": "joao.silva@example.com",
    "data_evento": "2026-05-27T10:00:00Z"
  }'

# Esperado: HTTP 200 com processamento concluído
```

### 2. Validação de Cálculo de Prioridade (Patrimônio Alto)
```bash
# Teste com patrimônio >= 200.000
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
  -H "Content-Type: application/json" \
  -d '{
    "identificador_evento": "evt_prioridade_alta",
    "identificador_card": "card_123",
    "cliente_email": "cliente.alto@example.com",
    "data_evento": "2026-05-27T10:00:00Z"
  }'

# Esperado: Cliente atualizado com nivel_prioridade = "prioridade_alta"
```

### 3. Validação de Cálculo de Prioridade (Patrimônio Normal)
```bash
# Teste com patrimônio < 200.000
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
  -H "Content-Type: application/json" \
  -d '{
    "identificador_evento": "evt_prioridade_normal",
    "identificador_card": "card_456",
    "cliente_email": "cliente.normal@example.com",
    "data_evento": "2026-05-27T10:00:00Z"
  }'

# Esperado: Cliente atualizado com nivel_prioridade = "prioridade_normal"
```

### 4. Validação de Idempotência
```bash
# Primeira requisição
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
  -H "Content-Type: application/json" \
  -d '{
    "identificador_evento": "evt_idempotente",
    "identificador_card": "card_789",
    "cliente_email": "cliente.teste@example.com",
    "data_evento": "2026-05-27T10:00:00Z"
  }'

# Segunda requisição com mesmo identificador_evento
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
  -H "Content-Type: application/json" \
  -d '{
    "identificador_evento": "evt_idempotente",
    "identificador_card": "card_789",
    "cliente_email": "cliente.teste@example.com",
    "data_evento": "2026-05-27T10:00:00Z"
  }'

# Esperado: Segunda requisição retorna HTTP 200 sem reprocessar
```

### 5. Validação de Cliente Não Encontrado
```bash
# Teste com cliente que não existe
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
  -H "Content-Type: application/json" \
  -d '{
    "identificador_evento": "evt_inexistente",
    "identificador_card": "card_999",
    "cliente_email": "nao.existe@example.com",
    "data_evento": "2026-05-27T10:00:00Z"
  }'

# Esperado: HTTP 404 com erro de cliente não encontrado
```

### 6. Validação de Campos Obrigatórios
```bash
# Teste sem campos obrigatórios
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
  -H "Content-Type: application/json" \
  -d '{
    "identificador_evento": "evt_sem_campos"
  }'

# Esperado: HTTP 400 com erro de validação
```

### 7. Validação no Banco de Dados
```bash
# Verificar que cliente foi atualizado corretamente
docker exec -it postgres_container psql -U postgres -d mundo_invest -c "
SELECT cliente_email, status_cliente, nivel_prioridade 
FROM gestao_clientes.vw_cliente 
WHERE cliente_email = 'cliente.alto@example.com';
"

# Esperado: status_cliente = "Processado", nivel_prioridade calculado corretamente
```

### 8. Validação da Mutation GraphQL
Verificar no código que a mutation updateCard está estruturada conforme a documentação do Pipefy, com comentários referenciando a fonte oficial.
