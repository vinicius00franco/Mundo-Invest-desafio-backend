# Índice de Planos - Mundo Invest

## Visão Geral

Este diretório contém a documentação organizada do plano de implementação do sistema de gestão de clientes e integração com Pipefy para o Mundo Invest. O plano foi dividido em sprints para facilitar o acompanhamento e execução.

## Estrutura de Organização

O plano de implementação original foi reorganizado em pastas separadas por sprint para melhor organização e navegabilidade:

- **[plano-implementacao.md](./plano-implementacao.md)** - Documento original completo (mantido para referência)
- **[sprint-1/](./sprint-1/)** - Sprint 1: Fundamentos e Estrutura
- **[sprint-2/](./sprint-2/)** - Sprint 2: API de Criação de Cliente
- **[sprint-3/](./sprint-3/)** - Sprint 3: API de Webhook
- **[sprint-4/](./sprint-4/)** - Sprint 4: Testes Automatizados
- **[sprint-5/](./sprint-5/)** - Sprint 5: Documentação e Finalização

## Sprints

### Sprint 1: Fundamentos e Estrutura (2 dias)

**Objetivo**: Configurar estrutura do projeto e banco de dados seguindo DDD e normalização

**Conteúdo**:
- Configuração do projeto
- Modelo de dados (3FN + Linguagem Ubíqua + Trigramação)
- Camada de persistência
- Segurança do banco de dados (RBAC)
- Configuração Docker

**Arquivos**:
- [README.md](./sprint-1/README.md)
- [tarefas.md](./sprint-1/tarefas.md)
- [criterios-aceite.md](./sprint-1/criterios-aceite.md)
- [docker-compose.md](./sprint-1/docker-compose.md)

### Sprint 2: API de Criação de Cliente (2 dias)

**Objetivo**: Implementar endpoint POST /clientes seguindo linguagem ubíqua

**Conteúdo**:
- Validação de dados
- Camada de serviço
- Camada de apresentação
- Integração Pipefy (createCard)

**Arquivos**:
- [README.md](./sprint-2/README.md)
- [tarefas.md](./sprint-2/tarefas.md)
- [criterios-aceite.md](./sprint-2/criterios-aceite.md)

### Sprint 3: API de Webhook (2 dias)

**Objetivo**: Implementar endpoint POST /webhooks/pipefy/card-updated seguindo linguagem ubíqua

**Conteúdo**:
- Validação de dados de webhook
- Implementação de idempotência
- Cálculo de prioridade
- Atualização de clientes
- Integração Pipefy (updateCard)

**Arquivos**:
- [README.md](./sprint-3/README.md)
- [tarefas.md](./sprint-3/tarefas.md)
- [criterios-aceite.md](./sprint-3/criterios-aceite.md)

### Sprint 4: Testes Automatizados (2 dias)

**Objetivo**: Implementar testes automatizados cobrindo todos os cenários

**Conteúdo**:
- Testes de criação de cliente (11 cenários)
- Testes de webhook (16 cenários)
- Testes de integração
- Configuração de CI/CD

**Arquivos**:
- [README.md](./sprint-4/README.md)
- [tarefas.md](./sprint-4/tarefas.md)
- [criterios-aceite.md](./sprint-4/criterios-aceite.md)

### Sprint 5: Documentação e Finalização (1 dia)

**Objetivo**: Preparar documentação e README para entrega

**Conteúdo**:
- Criação de README.md
- Documentação técnica
- Preparação para defesa
- Verificação final

**Arquivos**:
- [README.md](./sprint-5/README.md)
- [tarefas.md](./sprint-5/tarefas.md)
- [criterios-aceite.md](./sprint-5/criterios-aceite.md)

## Cronograma Total

**Duração**: 9 dias (2 + 2 + 2 + 2 + 1)

## Contextos Delimitados

O sistema está organizado em três contextos delimitados:

### 1. Contexto: Gestão de Clientes
- **Linguagem Ubíqua:** Cliente, Patrimonio, Solicitacao, StatusCliente, NivelPrioridade
- **Responsabilidade:** Gerenciar o ciclo de vida de clientes e seus patrimônios

### 2. Contexto: Integração Pipefy
- **Linguagem Ubíqua:** Card, IdentificadorCard, Mutation, Query
- **Responsabilidade:** Mapear clientes para cards no Pipefy

### 3. Contexto: Processamento de Eventos
- **Linguagem Ubíqua:** Webhook, Evento, IdentificadorEvento, Idempotencia
- **Responsabilidade:** Processar eventos do Pipefy de forma idempotente

## Cenários Identificados

**Total de Cenários:** 22 cenários cobrindo todos os aspectos de sucesso, falha e validação.

### Criação de Cliente (9 cenários)
- 1 cenário de sucesso
- 2 cenários de falha
- 6 cenários de validação

### Webhook Card Updated (13 cenários)
- 4 cenários de sucesso
- 2 cenários de falha
- 7 cenários de validação

## Como Navegar

1. **Para entender o plano completo**: Comece pelo [plano-implementacao.md](./plano-implementacao.md)
2. **Para trabalhar em uma sprint específica**: Navegue até a pasta da sprint correspondente
3. **Para verificar tarefas**: Abra o arquivo `tarefas.md` na pasta da sprint
4. **Para validar conclusão**: Use o arquivo `criterios-aceite.md` na pasta da sprint

## Princípios Arquiteturais

### Domain-Driven Design (DDD)
- Contextos delimitados bem definidos
- Linguagem ubíqua consistente
- Separação clara de responsabilidades

### Normalização de Dados (3FN)
- Eliminação de redundâncias
- Dependências funcionais bem definidas
- Integridade dos dados

### Feature Folders
- Organização por contexto delimitado
- Coesão de código
- Facilidade de manutenção

### Trigramação
- Nomenclatura consistente de colunas
- Prefixos por contexto
- Clareza e unicidade

## Documentação Relacionada

- [docs/regras-negocio.md](../docs/regras-negocio.md) - Regras de negócio detalhadas
- [docs/requisitos-funcionais.md](../docs/requisitos-funcionais.md) - Requisitos funcionais
- [docs/modelagem-dados.md](../docs/modelagem-dados.md) - Modelagem de dados
- [docs/snapshot-banco-dados.md](../docs/snapshot-banco-dados.md) - Snapshot do banco de dados

## Próximos Passos

1. Revisar este plano com stakeholders
2. Aprovar cronograma e prioridades
3. Iniciar Sprint 1: Fundamentos e Estrutura
4. Daily standups para acompanhar progresso
5. Revisão de sprint ao final de cada sprint
6. Ajustes baseados em aprendizados durante implementação
