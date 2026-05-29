# Critérios de Aceite - Sprint 2: API de Criação de Cliente

## Critérios de Aceite

- [ ] Endpoint POST /clientes responde a requisições
- [ ] Cliente válido é salvo no banco com status_cliente correto
- [ ] Payload inválido retorna HTTP 400
- [ ] Mutation createCard está estruturada no código
- [ ] Nomenclatura segue linguagem ubíqua
- [ ] Testes manuais passam

## Validação

Para validar que a sprint foi concluída com sucesso:

### 1. Validação do Endpoint
```bash
# Teste com payload válido
curl -X POST http://localhost:8080/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "cliente_nome": "João Silva",
    "cliente_email": "joao.silva@example.com",
    "tipo_solicitacao": "abertura_conta",
    "valor_patrimonio": 150000.00
  }'

# Esperado: HTTP 201 com cliente criado
```

### 2. Validação de Campos Obrigatórios
```bash
# Teste sem campos obrigatórios
curl -X POST http://localhost:8080/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "cliente_nome": "João Silva"
  }'

# Esperado: HTTP 400 com erro de validação
```

### 3. Validação de E-mail
```bash
# Teste com e-mail inválido
curl -X POST http://localhost:8080/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "cliente_nome": "João Silva",
    "cliente_email": "email-invalido",
    "tipo_solicitacao": "abertura_conta",
    "valor_patrimonio": 150000.00
  }'

# Esperado: HTTP 400 com erro de validação de e-mail
```

### 4. Validação de Patrimônio
```bash
# Teste com patrimônio negativo
curl -X POST http://localhost:8080/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "cliente_nome": "João Silva",
    "cliente_email": "joao.silva@example.com",
    "tipo_solicitacao": "abertura_conta",
    "valor_patrimonio": -1000.00
  }'

# Esperado: HTTP 400 com erro de validação de patrimônio
```

### 5. Validação no Banco de Dados
```bash
# Verificar que cliente foi salvo corretamente
docker exec -it postgres_container psql -U postgres -d mundo_invest -c "
SELECT cliente_nome, cliente_email, status_cliente 
FROM gestao_clientes.vw_cliente 
WHERE cliente_email = 'joao.silva@example.com';
"

# Esperado: Registro com status "Aguardando Análise"
```

### 6. Validação da Mutation GraphQL
Verificar no código que a mutation createCard está estruturada conforme a documentação do Pipefy, com comentários referenciando a fonte oficial.
