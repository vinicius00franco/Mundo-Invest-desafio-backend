# Critérios de Aceite - Sprint 4: Testes Automatizados

## Critérios de Aceite

- [ ] Todos os 22 cenários de teste passam
- [ ] Cobertura de testes > 80%
- [ ] Testes rodam em pipeline CI/CD

## Validação

Para validar que a sprint foi concluída com sucesso:

### 1. Executar Todos os Testes
```bash
# Executar todos os testes
go test ./... -v

# Esperado: Todos os testes passam (22 cenários)
```

### 2. Verificar Cobertura de Testes
```bash
# Gerar relatório de cobertura
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Esperado: Cobertura > 80%
```

### 3. Executar Testes por Feature

#### Testes de Criação de Cliente
```bash
# Executar testes do contexto de Gestão de Clientes
go test ./internal/gestao_clientes/... -v

# Esperado: 11 testes passam (1 sucesso + 2 falha + 8 validação)
```

#### Testes de Webhook
```bash
# Executar testes do contexto de Processamento de Eventos
go test ./internal/processamento_eventos/... -v

# Esperado: 16 testes passam (4 sucesso + 2 falha + 10 validação)
```

#### Testes de Integração
```bash
# Executar testes de integração
go test ./internal/shared/database/... -v

# Esperado: 3 testes passam
```

### 4. Validação dos 22 Cenários

#### Cenários de Criação de Cliente (9 cenários)
1. ✅ Criação com payload válido
2. ❌ Validação de campos obrigatórios (nome)
3. ❌ Validação de campos obrigatórios (email)
4. ❌ Validação de campos obrigatórios (tipo_solicitacao)
5. ❌ Validação de campos obrigatórios (valor_patrimonio)
6. ⚠️ Validação de e-mail inválido
7. ⚠️ Validação de patrimônio negativo
8. ⚠️ Validação de patrimônio zero
9. ⚠️ Validação de nome vazio

#### Cenários de Webhook (13 cenários)
1. ✅ Processamento com prioridade alta (patrimonio >= 200.000)
2. ✅ Processamento com prioridade normal (patrimonio < 200.000)
3. ✅ Processamento no limite (patrimonio == 200.000)
4. ✅ Idempotência (event_id duplicado)
5. ❌ Cliente não encontrado
6. ⚠️ Validação de campos obrigatórios (event_id)
7. ⚠️ Validação de campos obrigatórios (card_id)
8. ⚠️ Validação de campos obrigatórios (cliente_email)
9. ⚠️ Validação de campos obrigatórios (timestamp)
10. ⚠️ Validação de event_id vazio
11. ⚠️ Validação de card_id vazio
12. ⚠️ Validação de cliente_email vazio
13. ⚠️ Validação de timestamp inválido

### 5. Validação de Pipeline CI/CD
Verificar que:
- Pipeline CI/CD está configurado
- Testes executam automaticamente em cada commit/PR
- Relatórios de cobertura são gerados
- Notificações de falha estão configuradas

### 6. Testes de Integração
```bash
# Executar testes de integração com banco de dados real
docker-compose up -d postgres
go test ./internal/shared/database/... -v -tags=integration

# Esperado: Testes de integração passam
```

### 7. Testes de Performance (Opcional)
```bash
# Executar testes de performance
go test ./... -bench=. -benchmem

# Esperado: Performance dentro de limites aceitáveis
```

## Metodologia de Testes

### Testes Unitários
- Testam funções e métodos isoladamente
- Usam mocks para dependências externas
- São rápidos de executar
- Focam em lógica de negócio

### Testes de Integração
- Testam interações entre componentes
- Usam banco de dados real (via Docker)
- São mais lentos que testes unitários
- Focam em integração de sistemas

### Testes de Aceite
- Testam o sistema como um todo
- Simulam cenários reais de uso
- São os mais lentos
- Focam em comportamento esperado

## Cobertura Esperada por Módulo

- `internal/gestao_clientes/`: > 85%
- `internal/processamento_eventos/`: > 85%
- `internal/integracao_pipefy/`: > 70% (código simulado)
- `internal/dominio/`: > 90%
- `internal/shared/`: > 80%
