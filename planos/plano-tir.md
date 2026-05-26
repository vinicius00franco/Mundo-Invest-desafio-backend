# Plano de TIR (Testes de Integração) - Mundo Invest

## Visão Geral

Este plano define a estratégia de Testes de Integração (TIR) por domínio, compostos por:
- **TIR de Fluxo Completo**: Testes que verificam o fluxo completo de ponta a ponta (Controller → Service → Repository → Banco)
- **TIR de Fluxo Parcial**: Testes que verificam partes menores dos fluxos para facilitar implementação e testes (ex: Service → Repository → Banco, Repository → Banco, Service → PipefyClient)

Os TIR de Fluxo Parcial são quebras estratégicas dos fluxos completos, permitindo:
- Implementação incremental do sistema
- Isolamento de problemas durante desenvolvimento
- Testes mais focados e fáceis de manter
- Validação de integrações específicas entre camadas

**Total de TIR**: 21 testes distribuídos em 3 domínios

## Contextos Delimitados e Domínios

### 1. Domínio: Gestão de Clientes
**Linguagem Ubíqua**: Cliente, Patrimonio, Solicitacao, StatusCliente, NivelPrioridade
**Responsabilidade**: Gerenciar o ciclo de vida de clientes e seus patrimônios

### 2. Domínio: Integração Pipefy
**Linguagem Ubíqua**: Card, IdentificadorCard, Mutation, Query
**Responsabilidade**: Mapear clientes para cards no Pipefy

### 3. Domínio: Processamento de Eventos
**Linguagem Ubíqua**: Webhook, Evento, IdentificadorEvento, Idempotencia
**Responsabilidade**: Processar eventos do Pipefy de forma idempotente

---

## DOMÍNIO 1: Gestão de Clientes

### Fluxo: Criação de Cliente (POST /clientes)

#### TIR de Integração

##### TIR-GC-001: Integração Controller → Service → Repository → Banco
**Objetivo**: Verificar a integração completa do fluxo de criação de cliente
**Cenário BDD**: Criar cliente com dados válidos
**Componentes Envolvidos**:
- ClienteController (gestao_clientes/controller.go)
- ClienteService (gestao_clientes/service.go)
- ClienteRepository (gestao_clientes/repository.go)
- Banco de Dados (SQLite/PostgreSQL)

**Pré-condições**:
- Banco de dados está acessível
- Tabelas estão criadas

**Passos do Teste**:
1. Iniciar transação no banco
2. Enviar requisição POST para /clientes com payload válido
3. Verificar se Controller recebeu a requisição
4. Verificar se Service processou a criação
5. Verificar se Repository persistiu o cliente
6. Verificar se banco contém o registro
7. Rollback transação

**Resultados Esperados**:
- HTTP 201 Created
- Cliente persistido no banco com todos os campos
- Status inicial é "Aguardando Análise"
- Identificador interno e externo gerados

**Dados de Teste**:
```json
{
  "cliente_nome": "João Silva",
  "cliente_email": "joao.silva@example.com",
  "tipo_solicitacao": "Atualização cadastral",
  "valor_patrimonio": 250000
}
```

---

##### TIR-GC-002: Integração Service → PipefyGraphQLClient (Simulação)
**Objetivo**: Verificar a integração entre Service e client GraphQL do Pipefy
**Cenário BDD**: Criar cliente com dados válidos
**Componentes Envolvidos**:
- ClienteService (gestao_clientes/service.go)
- PipefyGraphQLClient (integracao_pipefy/client.go)

**Pré-condições**:
- Cliente criado no banco
- PipefyGraphQLClient configurado em modo simulação

**Passos do Teste**:
1. Mockar ClienteRepository para retornar cliente criado
2. Chamar ClienteService.criarCliente()
3. Verificar se PipefyGraphQLClient foi chamado
4. Verificar se mutation createCard foi estruturada corretamente
5. Verificar se variáveis da mutation estão corretas

**Resultados Esperados**:
- PipefyGraphQLClient recebeu chamada com dados corretos
- Mutation createCard está estruturada conforme especificação Pipefy
- Variáveis contêm: nome, email, patrimonio

**Verificação da Mutation**:
```graphql
mutation createCard($input: CreateCardInput!) {
  createCard(input: $input) {
    card {
      id
    }
  }
}
```

---

#### TIR de Fluxo Parcial

##### TIR-GC-PAR-001: Fluxo Parcial Service → Repository → Banco (Persistência)
**Objetivo**: Verificar integração entre Service, Repository e Banco (sem Controller)
**Cenário BDD**: Criar cliente com dados válidos
**Fluxo**: ClienteService → ClienteRepository → Banco de Dados

**Pré-condições**:
- Banco de dados está acessível
- Tabelas estão criadas

**Passos do Teste**:
1. Iniciar transação no banco
2. Chamar ClienteService.criarCliente() diretamente com payload válido
3. Verificar se Repository persistiu o cliente
4. Verificar se banco contém o registro
5. Buscar cliente por identificador_interno
6. Buscar cliente por cliente_email
7. Rollback transação

**Resultados Esperados**:
- Cliente persistido no banco com todos os campos
- Status inicial é "Aguardando Análise"
- Identificador interno e externo gerados
- Cliente encontrado por identificador_interno
- Cliente encontrado por cliente_email

**Dados de Teste**:
```json
{
  "cliente_nome": "João Silva",
  "cliente_email": "joao.silva@example.com",
  "tipo_solicitacao": "Atualização cadastral",
  "valor_patrimonio": 250000
}
```

---

##### TIR-GC-PAR-002: Fluxo Parcial Repository → Banco (Validação de Campos Obrigatórios)
**Objetivo**: Verificar validação de campos obrigatórios na camada de persistência
**Cenário BDD**: Criar cliente sem campos obrigatórios
**Fluxo**: ClienteRepository → Banco de Dados

**Pré-condições**:
- Banco de dados está acessível
- Tabelas estão criadas

**Passos do Teste**:
1. Iniciar transação no banco
2. Tentar salvar cliente sem cliente_nome
3. Tentar salvar cliente sem cliente_email
4. Tentar salvar cliente sem tipo_solicitacao
5. Tentar salvar cliente sem valor_patrimonio
6. Verificar erros retornados pelo banco
7. Rollback transação

**Resultados Esperados**:
- Banco rejeita inserção sem campos obrigatórios
- Erros indicam quais campos faltam
- Nenhum registro é criado

---

##### TIR-GC-PAR-003: Fluxo Parcial Service → PipefyGraphQLClient (Estruturação Mutation)
**Objetivo**: Verificar integração entre Service e PipefyGraphQLClient (sem banco)
**Cenário BDD**: Criar cliente com dados válidos
**Fluxo**: ClienteService → PipefyGraphQLClient

**Pré-condições**:
- PipefyGraphQLClient configurado em modo simulação

**Passos do Teste**:
1. Mockar ClienteRepository para não chamar banco
2. Chamar ClienteService.criarCliente() com payload válido
3. Verificar se PipefyGraphQLClient foi chamado
4. Verificar se mutation createCard foi estruturada corretamente
5. Verificar se variáveis da mutation estão corretas

**Resultados Esperados**:
- PipefyGraphQLClient recebeu chamada com dados corretos
- Mutation createCard está estruturada conforme especificação Pipefy
- Variáveis contêm: nome, email, patrimonio

---

##### TIR-GC-PAR-004: Fluxo Parcial Controller → Service (Validação de Payload)
**Objetivo**: Verificar integração entre Controller e Service (sem banco)
**Cenário BDD**: Criar cliente sem campos obrigatórios / com e-mail inválido
**Fluxo**: ClienteController → ClienteService

**Pré-condições**:
- Service configurado para validar payload

**Passos do Teste**:
1. Mockar ClienteService para retornar erro de validação
2. Enviar requisição POST para /clientes sem cliente_nome
3. Enviar requisição POST para /clientes sem cliente_email
4. Enviar requisição POST para /clientes com email inválido
5. Verificar se Controller recebeu erros do Service
6. Verificar se Controller retornou HTTP 400

**Resultados Esperados**:
- Controller propaga erros de validação do Service
- HTTP 400 retornado para payloads inválidos
- Mensagens de erro em Português

---

##### TIR-GC-PAR-005: Fluxo Parcial Service → Repository (Busca de Cliente)
**Objetivo**: Verificar integração entre Service e Repository para busca de cliente
**Cenário BDD**: Buscar cliente existente e inexistente
**Fluxo**: ClienteService → ClienteRepository

**Pré-condições**:
- Banco de dados está acessível
- Cliente de teste existe no banco

**Passos do Teste**:
1. Iniciar transação no banco
2. Inserir cliente de teste no banco
3. Chamar ClienteService.buscarClientePorEmail() com email existente
4. Chamar ClienteService.buscarClientePorEmail() com email inexistente
5. Verificar se Repository foi chamado corretamente
6. Verificar resultados retornados
7. Rollback transação

**Resultados Esperados**:
- Cliente encontrado quando email existe
- nil retornado quando email não existe
- Repository chamado com parâmetros corretos

---

### Resumo TIR - Domínio Gestão de Clientes

| TIR ID | Tipo | Fluxo | Cenário BDD | Status |
|--------|------|-------|-------------|--------|
| TIR-GC-001 | Fluxo Completo | Controller→Service→Repository→Banco | Sucesso | Pendente |
| TIR-GC-002 | Fluxo Completo | Service→PipefyGraphQLClient | Sucesso | Pendente |
| TIR-GC-PAR-001 | Fluxo Parcial | Service→Repository→Banco | Sucesso | Pendente |
| TIR-GC-PAR-002 | Fluxo Parcial | Repository→Banco (Validação) | Falha | Pendente |
| TIR-GC-PAR-003 | Fluxo Parcial | Service→PipefyGraphQLClient (Mutation) | Sucesso | Pendente |
| TIR-GC-PAR-004 | Fluxo Parcial | Controller→Service (Validação) | Validação | Pendente |
| TIR-GC-PAR-005 | Fluxo Parcial | Service→Repository (Busca) | Sucesso | Pendente |

---

## DOMÍNIO 2: Integração Pipefy

### Fluxo: Criação de Card (Mutation createCard)

#### TIR de Integração

##### TIR-IP-001: Integração ClienteService → PipefyGraphQLClient → Mutation (createCard)
**Objetivo**: Verificar integração completa na estruturação da mutation createCard
**Cenário BDD**: Criar cliente com dados válidos
**Componentes Envolvidos**:
- ClienteService (gestao_clientes/service.go)
- PipefyGraphQLClient (integracao_pipefy/client.go)

**Pré-condições**:
- Cliente criado com dados válidos

**Passos do Teste**:
1. Mockar ClienteRepository para retornar cliente
2. Chamar ClienteService.criarCliente()
3. Interceptar chamada ao PipefyGraphQLClient
4. Verificar estrutura da mutation
5. Verificar variáveis da mutation
6. Verificar conformidade com documentação Pipefy

**Resultados Esperados**:
- Mutation createCard estruturada corretamente
- Variáveis contêm: nome, email, patrimonio
- Conforme especificação oficial do Pipefy

**Estrutura da Mutation Esperada**:
```graphql
mutation createCard($input: CreateCardInput!) {
  createCard(input: $input) {
    card {
      id
      title
      fields {
        name
        value
      }
    }
  }
}
```

---

##### TIR-IP-002: Integração WebhookService → PipefyGraphQLClient → Mutation (updateCard)
**Objetivo**: Verificar integração completa na estruturação da mutation updateCard
**Cenário BDD**: Processar webhook com prioridade
**Componentes Envolvidos**:
- WebhookService (processamento_eventos/service.go)
- PipefyGraphQLClient (integracao_pipefy/client.go)

**Pré-condições**:
- Cliente existe no banco
- Prioridade calculada

**Passos do Teste**:
1. Mockar ClienteRepository e EventoRepository
2. Chamar WebhookService.processarWebhook()
3. Interceptar chamada ao PipefyGraphQLClient
4. Verificar estrutura da mutation
5. Verificar variáveis da mutation
6. Verificar conformidade com documentação Pipefy

**Resultados Esperados**:
- Mutation updateCard estruturada corretamente
- Variáveis contêm: card_id, status (Processado), prioridade
- Conforme especificação oficial do Pipefy

**Estrutura da Mutation Esperada**:
```graphql
mutation updateCard($input: UpdateCardInput!) {
  updateCard(input: $input) {
    card {
      id
      fields {
        name
        value
      }
    }
  }
}
```

---

#### TIR de Fluxo Parcial

##### TIR-IP-PAR-001: Fluxo Parcial PipefyGraphQLClient (Estruturação createCard)
**Objetivo**: Verificar estruturação correta da mutation createCard
**Cenário BDD**: Criar cliente com dados válidos
**Fluxo**: PipefyGraphQLClient (isolado)

**Pré-condições**:
- PipefyGraphQLClient configurado

**Passos do Teste**:
1. Chamar PipefyGraphQLClient.estruturarMutationCreateCard() com cliente válido
2. Verificar sintaxe GraphQL correta
3. Verificar variáveis definidas corretamente
4. Verificar campos do card mapeados corretamente
5. Verificar conformidade com documentação Pipefy

**Resultados Esperados**:
- Mutation GraphQL sintaticamente correta
- Variáveis: $input do tipo CreateCardInput
- Campos mapeados: nome, email, patrimonio
- Comentário com fonte da especificação

---

##### TIR-IP-PAR-002: Fluxo Parcial PipefyGraphQLClient (Estruturação updateCard)
**Objetivo**: Verificar estruturação correta da mutation updateCard
**Cenário BDD**: Processar webhook com prioridade
**Fluxo**: PipefyGraphQLClient (isolado)

**Pré-condições**:
- PipefyGraphQLClient configurado

**Passos do Teste**:
1. Chamar PipefyGraphQLClient.estruturarMutationUpdateCard() com cliente e prioridade
2. Verificar sintaxe GraphQL correta
3. Verificar variáveis definidas corretamente
4. Verificar campos de status e prioridade mapeados
5. Verificar conformidade com documentação Pipefy

**Resultados Esperados**:
- Mutation GraphQL sintaticamente correta
- Variáveis: card_id, campos a atualizar
- Campos mapeados: status (Processado), prioridade
- Comentário com fonte da especificação

---

### Resumo TIR - Domínio Integração Pipefy

| TIR ID | Tipo | Fluxo | Cenário BDD | Status |
|--------|------|-------|-------------|--------|
| TIR-IP-001 | Fluxo Completo | Service→PipefyGraphQLClient→Mutation | Sucesso | Pendente |
| TIR-IP-002 | Fluxo Completo | Service→PipefyGraphQLClient→Mutation (updateCard) | Sucesso | Pendente |
| TIR-IP-PAR-001 | Fluxo Parcial | PipefyGraphQLClient (createCard) | Sucesso | Pendente |
| TIR-IP-PAR-002 | Fluxo Parcial | PipefyGraphQLClient (updateCard) | Sucesso | Pendente |

---

## DOMÍNIO 3: Processamento de Eventos

### Fluxo: Webhook Card Updated (POST /webhooks/pipefy/card-updated)

#### TIR de Integração

##### TIR-PE-001: Integração Controller → Service → EventoRepository → Banco (Idempotência)
**Objetivo**: Verificar integração completa do fluxo de idempotência
**Cenário BDD**: Processar webhook duplicado
**Componentes Envolvidos**:
- WebhookController (processamento_eventos/controller.go)
- WebhookService (processamento_eventos/service.go)
- EventoRepository (processamento_eventos/repository.go)
- Banco de Dados

**Pré-condições**:
- Evento com event_id "evt_123" já processado no banco
- Cliente existe no banco

**Passos do Teste**:
1. Iniciar transação no banco
2. Inserir evento processado no banco
3. Enviar requisição POST para /webhooks/pipefy/card-updated com event_id duplicado
4. Verificar se Controller recebeu a requisição
5. Verificar se Service detectou duplicidade
6. Verificar se EventoRepository encontrou o evento
7. Verificar se webhook não foi processado novamente
8. Verificar se cliente não foi atualizado
9. Rollback transação

**Resultados Esperados**:
- HTTP 200 OK
- Mensagem indica evento já processado
- Cliente mantém status e prioridade anteriores
- Nenhuma nova atualização no banco

**Dados de Teste**:
```json
{
  "event_id": "evt_123",
  "card_id": "card_456",
  "cliente_email": "joao.silva@example.com",
  "timestamp": "2026-05-18T12:00:00Z"
}
```

---

##### TIR-PE-002: Integração Controller → Service → ClienteRepository → EventoRepository → Banco (Processamento Completo)
**Objetivo**: Verificar integração completa do fluxo de processamento de webhook
**Cenário BDD**: Processar webhook com patrimônio alto (prioridade alta)
**Componentes Envolvidos**:
- WebhookController (processamento_eventos/controller.go)
- WebhookService (processamento_eventos/service.go)
- ClienteRepository (gestao_clientes/repository.go)
- EventoRepository (processamento_eventos/repository.go)
- PrioridadeCalculator (dominio/prioridade_calculator.go)
- PipefyGraphQLClient (integracao_pipefy/client.go)
- Banco de Dados

**Pré-condições**:
- Cliente existe no banco com valor_patrimonio = 250.000
- Evento não foi processado anteriormente

**Passos do Teste**:
1. Iniciar transação no banco
2. Inserir cliente no banco com patrimonio alto
3. Enviar requisição POST para /webhooks/pipefy/card-updated
4. Verificar se Controller recebeu a requisição
5. Verificar se Service processou o webhook
6. Verificar se EventoRepository verificou idempotência
7. Verificar se ClienteRepository buscou o cliente
8. Verificar se PrioridadeCalculator calculou prioridade_alta
9. Verificar se PipefyGraphQLClient estruturou mutation updateCard
10. Verificar se ClienteRepository atualizou status e prioridade
11. Verificar se EventoRepository salvou evento como processado
12. Verificar se banco contém atualizações
13. Rollback transação

**Resultados Esperados**:
- HTTP 200 OK
- Evento marcado como processado (pev_eve_fpr = true)
- Prioridade definida como "prioridade_alta"
- Status atualizado para "Processado"
- Mutation updateCard estruturada
- Response contém dados do processamento

**Dados de Teste**:
```json
{
  "event_id": "evt_123",
  "card_id": "card_456",
  "cliente_email": "joao.silva@example.com",
  "timestamp": "2026-05-18T12:00:00Z"
}
```

---

##### TIR-PE-003: Integração Controller → Service → ClienteRepository (Cliente Não Encontrado)
**Objetivo**: Verificar integração quando cliente não existe
**Cenário BDD**: Processar webhook com cliente não encontrado
**Componentes Envolvidos**:
- WebhookController (processamento_eventos/controller.go)
- WebhookService (processamento_eventos/service.go)
- ClienteRepository (gestao_clientes/repository.go)
- Banco de Dados

**Pré-condições**:
- Cliente com email "inexistente@example.com" não existe no banco

**Passos do Teste**:
1. Iniciar transação no banco
2. Enviar requisição POST para /webhooks/pipefy/card-updated com email inexistente
3. Verificar se Controller recebeu a requisição
4. Verificar se Service processou o webhook
5. Verificar se ClienteRepository não encontrou o cliente
6. Verificar se Service retornou erro
7. Verificar se EventoRepository não salvou evento
8. Rollback transação

**Resultados Esperados**:
- HTTP 404 Not Found
- Mensagem de erro indica cliente não encontrado
- Evento não salvo no banco
- Nenhuma atualização realizada

**Dados de Teste**:
```json
{
  "event_id": "evt_999",
  "card_id": "card_999",
  "cliente_email": "inexistente@example.com",
  "timestamp": "2026-05-18T15:00:00Z"
}
```

---

#### TIR de Fluxo Parcial

##### TIR-PE-PAR-001: Fluxo Parcial Service → EventoRepository → Banco (Idempotência)
**Objetivo**: Verificar integração entre Service, EventoRepository e Banco para idempotência
**Cenário BDD**: Processar webhook duplicado
**Fluxo**: WebhookService → EventoRepository → Banco de Dados

**Pré-condições**:
- Banco de dados está acessível
- Tabelas estão criadas

**Passos do Teste**:
1. Iniciar transação no banco
2. Inserir evento processado no banco (pev_eve_fpr = true)
3. Chamar WebhookService.processarWebhook() com event_id duplicado
4. Verificar se EventoRepository verificou o evento
5. Verificar se Service detectou duplicidade
6. Verificar se processamento não foi realizado
7. Rollback transação

**Resultados Esperados**:
- Evento encontrado com pev_eve_fpr = true
- Service detecta duplicidade
- Processamento não realizado
- Cliente não atualizado

**Dados de Teste**:
```json
{
  "event_id": "evt_123",
  "card_id": "card_456",
  "cliente_email": "joao.silva@example.com",
  "timestamp": "2026-05-18T12:00:00Z"
}
```

---

##### TIR-PE-PAR-002: Fluxo Parcial Service → ClienteRepository → Banco (Busca de Cliente)
**Objetivo**: Verificar integração entre Service, ClienteRepository e Banco para busca de cliente
**Cenário BDD**: Processar webhook com cliente encontrado/não encontrado
**Fluxo**: WebhookService → ClienteRepository → Banco de Dados

**Pré-condições**:
- Banco de dados está acessível
- Cliente de teste existe no banco

**Passos do Teste**:
1. Iniciar transação no banco
2. Inserir cliente de teste no banco
3. Chamar WebhookService.buscarClientePorEmail() com email existente
4. Chamar WebhookService.buscarClientePorEmail() com email inexistente
5. Verificar se ClienteRepository foi chamado corretamente
6. Verificar resultados retornados
7. Rollback transação

**Resultados Esperados**:
- Cliente encontrado quando email existe
- nil retornado quando email não existe
- ClienteRepository chamado com parâmetros corretos

---

##### TIR-PE-PAR-003: Fluxo Parcial Service → PrioridadeCalculator (Cálculo de Prioridade)
**Objetivo**: Verificar integração entre Service e PrioridadeCalculator
**Cenário BDD**: Processar webhook com patrimônio alto/baixo/limite
**Fluxo**: WebhookService → PrioridadeCalculator

**Pré-condições**:
- PrioridadeCalculator configurado

**Passos do Teste**:
1. Chamar WebhookService.calcularPrioridade() com patrimonio 250.000
2. Chamar WebhookService.calcularPrioridade() com patrimonio 150.000
3. Chamar WebhookService.calcularPrioridade() com patrimonio 200.000
4. Chamar WebhookService.calcularPrioridade() com patrimonio 199.999
5. Verificar se PrioridadeCalculator foi chamado corretamente
6. Verificar resultados retornados

**Resultados Esperados**:
- "prioridade_alta" quando patrimonio >= 200.000
- "prioridade_normal" quando patrimonio < 200.000
- Limite exato (200.000) → prioridade_alta

---

##### TIR-PE-PAR-004: Fluxo Parcial Service → PipefyGraphQLClient (Estruturação updateCard)
**Objetivo**: Verificar integração entre Service e PipefyGraphQLClient para updateCard
**Cenário BDD**: Processar webhook com prioridade
**Fluxo**: WebhookService → PipefyGraphQLClient

**Pré-condições**:
- PipefyGraphQLClient configurado em modo simulação

**Passos do Teste**:
1. Mockar ClienteRepository e EventoRepository
2. Chamar WebhookService.processarWebhook() com payload válido
3. Verificar se PipefyGraphQLClient foi chamado
4. Verificar se mutation updateCard foi estruturada corretamente
5. Verificar se variáveis da mutation estão corretas

**Resultados Esperados**:
- PipefyGraphQLClient recebeu chamada com dados corretos
- Mutation updateCard está estruturada conforme especificação Pipefy
- Variáveis contêm: card_id, status (Processado), prioridade

---

##### TIR-PE-PAR-005: Fluxo Parcial Service → ClienteRepository → Banco (Atualização)
**Objetivo**: Verificar integração entre Service, ClienteRepository e Banco para atualização
**Cenário BDD**: Processar webhook com prioridade alta/normal
**Fluxo**: WebhookService → ClienteRepository → Banco de Dados

**Pré-condições**:
- Banco de dados está acessível
- Cliente de teste existe no banco

**Passos do Teste**:
1. Iniciar transação no banco
2. Inserir cliente de teste no banco
3. Chamar WebhookService.atualizarCliente() com prioridade_alta
4. Chamar WebhookService.atualizarCliente() com prioridade_normal
5. Verificar se ClienteRepository atualizou o cliente
6. Verificar se banco contém as atualizações
7. Buscar cliente atualizado
8. Rollback transação

**Resultados Esperados**:
- Status atualizado para "Processado"
- Prioridade salva corretamente (alta/normal)
- Cliente atualizado encontrado no banco

---

##### TIR-PE-PAR-006: Fluxo Parcial Service → EventoRepository → Banco (Persistência de Evento)
**Objetivo**: Verificar integração entre Service, EventoRepository e Banco para persistência de evento
**Cenário BDD**: Processar webhook com sucesso
**Fluxo**: WebhookService → EventoRepository → Banco de Dados

**Pré-condições**:
- Banco de dados está acessível
- Tabelas estão criadas

**Passos do Teste**:
1. Iniciar transação no banco
2. Chamar WebhookService.salvarEvento() com evento válido
3. Verificar se EventoRepository persistiu o evento
4. Verificar se banco contém o registro
5. Verificar se pev_eve_fpr = true
6. Buscar evento por identificador_evento
7. Rollback transação

**Resultados Esperados**:
- Evento persistido no banco
- pev_eve_fpr = true (processado)
- Evento encontrado por identificador_evento

---

##### TIR-PE-PAR-007: Fluxo Parcial Controller → Service (Validação de Payload Webhook)
**Objetivo**: Verificar integração entre Controller e Service para validação de webhook
**Cenário BDD**: Processar webhook sem campos obrigatórios / com campos inválidos
**Fluxo**: WebhookController → WebhookService

**Pré-condições**:
- Service configurado para validar payload

**Passos do Teste**:
1. Mockar WebhookService para retornar erro de validação
2. Enviar requisição POST para /webhooks sem event_id
3. Enviar requisição POST para /webhooks sem card_id
4. Enviar requisição POST para /webhooks sem cliente_email
5. Enviar requisição POST para /webhooks sem timestamp
6. Enviar requisição POST para /webhooks com event_id vazio
7. Enviar requisição POST para /webhooks com email inválido
8. Verificar se Controller recebeu erros do Service
9. Verificar se Controller retornou HTTP 400

**Resultados Esperados**:
- Controller propaga erros de validação do Service
- HTTP 400 retornado para payloads inválidos
- Mensagens de erro em Português

---

### Resumo TIR - Domínio Processamento de Eventos

| TIR ID | Tipo | Fluxo | Cenário BDD | Status |
|--------|------|-------|-------------|--------|
| TIR-PE-001 | Fluxo Completo | Controller→Service→EventoRepo→Banco (Idempotência) | Sucesso (Duplicado) | Pendente |
| TIR-PE-002 | Fluxo Completo | Controller→Service→ClienteRepo→EventoRepo→Banco (Completo) | Sucesso (Prioridade Alta) | Pendente |
| TIR-PE-003 | Fluxo Completo | Controller→Service→ClienteRepo (Não Encontrado) | Falha | Pendente |
| TIR-PE-PAR-001 | Fluxo Parcial | Service→EventoRepo→Banco (Idempotência) | Sucesso (Duplicado) | Pendente |
| TIR-PE-PAR-002 | Fluxo Parcial | Service→ClienteRepo→Banco (Busca) | Sucesso/Falha | Pendente |
| TIR-PE-PAR-003 | Fluxo Parcial | Service→PrioridadeCalculator (Cálculo) | Sucesso (Prioridade Alta/Normal) | Pendente |
| TIR-PE-PAR-004 | Fluxo Parcial | Service→PipefyGraphQLClient (updateCard) | Sucesso | Pendente |
| TIR-PE-PAR-005 | Fluxo Parcial | Service→ClienteRepo→Banco (Atualização) | Sucesso (Prioridade Alta/Normal) | Pendente |
| TIR-PE-PAR-006 | Fluxo Parcial | Service→EventoRepo→Banco (Persistência) | Sucesso | Pendente |
| TIR-PE-PAR-007 | Fluxo Parcial | Controller→Service (Validação) | Validação | Pendente |

---

## Estrutura de Organização dos TIR

### Diretório de Testes

```
backend/
├── tests/
│   ├── integration/
│   │   ├── gestao_clientes/
│   │   │   ├── tir_gc_001_fluxo_completo_controller_service_repository_banco_test.go
│   │   │   └── tir_gc_002_fluxo_completo_service_pipefy_client_test.go
│   │   ├── integracao_pipefy/
│   │   │   ├── tir_ip_001_fluxo_completo_service_pipefy_mutation_create_card_test.go
│   │   │   └── tir_ip_002_fluxo_completo_service_pipefy_mutation_update_card_test.go
│   │   └── processamento_eventos/
│   │       ├── tir_pe_001_fluxo_completo_idempotencia_test.go
│   │       ├── tir_pe_002_fluxo_completo_processamento_completo_test.go
│   │       └── tir_pe_003_fluxo_completo_cliente_nao_encontrado_test.go
│   └── partial/
│       ├── gestao_clientes/
│       │   ├── tir_gc_par_001_fluxo_parcial_service_repository_banco_test.go
│       │   ├── tir_gc_par_002_fluxo_parcial_repository_banco_validacao_test.go
│       │   ├── tir_gc_par_003_fluxo_parcial_service_pipefy_mutation_test.go
│       │   ├── tir_gc_par_004_fluxo_parcial_controller_service_validacao_test.go
│       │   └── tir_gc_par_005_fluxo_parcial_service_repository_busca_test.go
│       ├── integracao_pipefy/
│       │   ├── tir_ip_par_001_fluxo_parcial_pipefy_client_create_card_test.go
│       │   └── tir_ip_par_002_fluxo_parcial_pipefy_client_update_card_test.go
│       └── processamento_eventos/
│           ├── tir_pe_par_001_fluxo_parcial_service_evento_repo_banco_idempotencia_test.go
│           ├── tir_pe_par_002_fluxo_parcial_service_cliente_repo_banco_busca_test.go
│           ├── tir_pe_par_003_fluxo_parcial_service_prioridade_calculator_test.go
│           ├── tir_pe_par_004_fluxo_parcial_service_pipefy_client_update_card_test.go
│           ├── tir_pe_par_005_fluxo_parcial_service_cliente_repo_banco_atualizacao_test.go
│           ├── tir_pe_par_006_fluxo_parcial_service_evento_repo_banco_persistencia_test.go
│           └── tir_pe_par_007_fluxo_parcial_controller_service_validacao_test.go
```

### Convenções de Nomenclatura

- **TIR de Fluxo Completo**: `tir_<sigla_dominio>_<numero>_fluxo_completo_<descricao>_test.go`
- **TIR de Fluxo Parcial**: `tir_<sigla_dominio>_par_<numero>_fluxo_parcial_<descricao>_test.go`

**Siglas de Domínio**:
- GC: Gestão de Clientes
- IP: Integração Pipefy
- PE: Processamento de Eventos

### Estrutura do Arquivo de Teste

```go
package integration

import (
    "testing"
    // imports necessários
)

// TIR-GC-001: Integração Controller → Service → Repository → Banco
// Cenário BDD: Criar cliente com dados válidos
func TestTIR_GC_001_IntegracaoControllerServiceRepositoryBanco(t *testing.T) {
    // Setup
    // Act
    // Assert
    // Teardown
}
```

## Matriz de Rastreabilidade

### Cenários BDD → TIR

| Cenário BDD | TIR de Fluxo Completo | TIR de Fluxo Parcial |
|-------------|----------------------|---------------------|
| Criar cliente com dados válidos | TIR-GC-001, TIR-GC-002 | TIR-GC-PAR-001, TIR-GC-PAR-003, TIR-GC-PAR-005 |
| Criar cliente sem campos obrigatórios | - | TIR-GC-PAR-002, TIR-GC-PAR-004 |
| Criar cliente com e-mail inválido | - | TIR-GC-PAR-004 |
| Criar cliente com patrimônio negativo/zero | - | TIR-GC-PAR-002 |
| Criar cliente com nome vazio | - | TIR-GC-PAR-002 |
| Processar webhook com patrimônio alto | TIR-PE-002, TIR-IP-002 | TIR-PE-PAR-003, TIR-PE-PAR-004, TIR-PE-PAR-005, TIR-IP-PAR-002 |
| Processar webhook com patrimônio baixo | TIR-PE-002, TIR-IP-002 | TIR-PE-PAR-003, TIR-PE-PAR-004, TIR-PE-PAR-005, TIR-IP-PAR-002 |
| Processar webhook com patrimônio no limite | TIR-PE-002, TIR-IP-002 | TIR-PE-PAR-003, TIR-PE-PAR-004, TIR-PE-PAR-005, TIR-IP-PAR-002 |
| Processar webhook duplicado | TIR-PE-001 | TIR-PE-PAR-001, TIR-PE-PAR-006 |
| Processar webhook com cliente não encontrado | TIR-PE-003 | TIR-PE-PAR-002 |
| Processar webhook sem campos obrigatórios | - | TIR-PE-PAR-007 |
| Processar webhook com event_id vazio | - | TIR-PE-PAR-007 |
| Processar webhook com cliente_email inválido | - | TIR-PE-PAR-007 |
| Processar webhook com card_id vazio | - | TIR-PE-PAR-007 |
| Processar webhook com timestamp inválido | - | TIR-PE-PAR-007 |

## Critérios de Sucesso dos TIR

### TIR de Fluxo Completo
- ✅ Todos os componentes envolvidos são testados em conjunto
- ✅ Fluxo completo de ponta a ponta é verificado (Controller → Service → Repository → Banco)
- ✅ Interações entre componentes são validadas
- ✅ Banco de dados é utilizado (não mockado)
- ✅ Transações são usadas para isolamento
- ✅ Cenário BDD é coberto completamente

### TIR de Fluxo Parcial
- ✅ Parte específica do fluxo é testada (ex: Service → Repository → Banco)
- ✅ Integração entre subset de componentes é validada
- ✅ Banco de dados é utilizado quando aplicável
- ✅ Dependências externas ao fluxo podem ser mockadas
- ✅ Facilita implementação incremental
- ✅ Permite isolamento de problemas
- ✅ Comportamentos de borda são testados
- ✅ Mensagens de erro são verificadas

## Ordem de Execução Recomendada

### Fase 1: Fluxos Parciais (Fundação)
1. TIR de Fluxo Parcial - Gestão de Clientes
2. TIR de Fluxo Parcial - Integração Pipefy
3. TIR de Fluxo Parcial - Processamento de Eventos

**Objetivo**: Validar integrações específicas entre camadas antes de implementar fluxos completos

### Fase 2: Fluxos Completos (Integração)
1. TIR de Fluxo Completo - Gestão de Clientes
2. TIR de Fluxo Completo - Integração Pipefy
3. TIR de Fluxo Completo - Processamento de Eventos

**Objetivo**: Validar fluxos de ponta a ponta após fluxos parciais estarem funcionando

### Fase 3: End-to-End (Validação)
1. Executar todos os TIR em conjunto
2. Verificar cobertura de código
3. Validar performance dos testes
4. Validar que todos os cenários BDD estão cobertos

## Ferramentas e Frameworks

### Para Golang
- **testing**: Framework nativo de testes do Go
- **testify/assert**: Biblioteca de assertions
- **testify/mock**: Biblioteca de mocks
- **testcontainers**: Para containers de banco de dados em testes
- **sqlmock**: Para mock de banco de dados (se necessário)

### Comandos Úteis
```bash
# Executar todos os TIR
go test ./tests/...

# Executar TIR de fluxos completos
go test ./tests/integration/...

# Executar TIR de fluxos parciais
go test ./tests/partial/...

# Executar TIR de um domínio específico (fluxos completos)
go test ./tests/integration/gestao_clientes/...

# Executar TIR de um domínio específico (fluxos parciais)
go test ./tests/partial/gestao_clientes/...

# Executar com coverage
go test ./tests/... -cover

# Executar com coverage detalhado
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Conclusão

Este plano define **21 TIR** distribuídos em **3 domínios**:
- **Gestão de Clientes**: 7 TIR (2 fluxo completo + 5 fluxo parcial)
- **Integração Pipefy**: 4 TIR (2 fluxo completo + 2 fluxo parcial)
- **Processamento de Eventos**: 10 TIR (3 fluxo completo + 7 fluxo parcial)

Os TIR cobrem todos os cenários BDD mapeados e garantem que:
1. **Fluxos Parciais** validam integrações específicas entre camadas (Service → Repository → Banco, Service → PipefyClient, etc.)
2. **Fluxos Completos** validam o fluxo de ponta a ponta (Controller → Service → Repository → Banco)
3. Implementação incremental é facilitada pelos fluxos parciais
4. Isolamento de problemas é mais fácil com fluxos parciais
5. Regras de negócio são aplicadas corretamente
6. Validações são realizadas corretamente
7. Idempotência é garantida
8. Mutations GraphQL estão estruturadas corretamente
9. Todos os cenários BDD estão cobertos por TIRs
