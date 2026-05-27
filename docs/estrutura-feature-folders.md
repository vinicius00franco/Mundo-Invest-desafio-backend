# Estrutura de Feature Folders - Contextos Delimitados

## Visão Geral

O projeto Mundo Invest utiliza a arquitetura de **Feature Folders** organizada por **Contextos Delimitados** (Bounded Contexts) seguindo os princípios de **Domain-Driven Design (DDD)**. Cada contexto representa um domínio específico do negócio com sua própria linguagem ubíqua.

## Contextos Delimitados

### 1. Gestão de Clientes (`gestao_clientes`)

**Prefixo de Trigramação**: `gcl` (Gestão CLientes)

**Responsabilidade**: Gerenciar o ciclo de vida de clientes, desde o cadastro até a atualização de informações.

**Arquivos**:
- `internal/gestao_clientes/controller.go` - Camada de apresentação HTTP
- `internal/gestao_clientes/service.go` - Camada de serviço e lógica de negócio
- `internal/gestao_clientes/repository.go` - Camada de persistência de dados
- `internal/gestao_clientes/validator.go` - Validação de dados de entrada

**Entidades**:
- `Cliente` - Entidade principal do contexto
- `CriarClienteRequest` - DTO para criação de cliente

**Linguagem Ubíqua**:
- `IdentificadorInterno` (gcl_cli_int) - ID interno gerado pelo sistema
- `IdentificadorExterno` (gcl_cli_ide) - ID externo do Pipefy
- `Nome` (gcl_cli_nom) - Nome do cliente
- `Email` (gcl_cli_ema) - Email do cliente
- `TipoSolicitacao` (gcl_cli_tso) - Tipo de solicitação
- `ValorPatrimonio` (gcl_cli_vpa) - Valor do patrimônio
- `Status` (gcl_cli_stc) - Status do cliente

**Schema do Banco de Dados**: `gestao_clientes`

### 2. Processamento de Eventos (`processamento_eventos`)

**Prefixo de Trigramação**: `pev` (Processamento de EVentos)

**Responsabilidade**: Processar eventos do Pipefy de forma idempotente, calcular prioridade e atualizar clientes.

**Arquivos**:
- `internal/processamento_eventos/controller.go` - Camada de apresentação HTTP
- `internal/processamento_eventos/service.go` - Camada de serviço e lógica de negócio
- `internal/processamento_eventos/repository.go` - Camada de persistência de dados
- `internal/processamento_eventos/validator.go` - Validação de dados de entrada

**Entidades**:
- `Evento` - Entidade principal do contexto
- `WebhookRequest` - DTO para recebimento de webhooks

**Linguagem Ubíqua**:
- `IdentificadorInterno` (pev_eve_int) - ID interno gerado pelo sistema
- `IdentificadorEvento` (pev_eve_ide) - ID do evento do Pipefy
- `IdentificadorCard` (pev_eve_idc) - ID do card do Pipefy
- `EmailCliente` (pev_eve_ema) - Email do cliente relacionado
- `TimestampEvento` (pev_eve_tms) - Timestamp do evento
- `FoiProcessado` (pev_eve_fpr) - Flag de idempotência

**Schema do Banco de Dados**: `processamento_eventos`

### 3. Domínio Compartilhado (`domino`)

**Responsabilidade**: Contém lógica de domínio compartilhada entre contextos.

**Arquivos**:
- `internal/dominio/prioridade_calculator.go` - Cálculo de nível de prioridade

**Linguagem Ubíqua**:
- `PrioridadeCalculator` - Calculadora de prioridade
- `NivelPrioridade` - Nível de prioridade (prioridade_alta, prioridade_normal)

**Regras de Negócio**:
- Patrimônio >= 200.000 → prioridade_alta
- Patrimônio < 200.000 → prioridade_normal

### 4. Integração Pipefy (`integracao_pipefy`)

**Responsabilidade**: Gerenciar a integração com a API do Pipefy via GraphQL.

**Arquivos**:
- `internal/integracao_pipefy/client.go` - Cliente GraphQL
- `internal/integracao_pipefy/mutations.go` - Mutations GraphQL do Pipefy

**Linguagem Ubíqua**:
- `PipefyGraphQLClient` - Cliente GraphQL do Pipefy
- `FieldAttribute` - Atributos de campo para mutations
- `Mutations` - Container de mutations GraphQL

**Mutations Implementadas**:
- `createCard` - Criar card no Pipefy
- `updateCard` - Atualizar card no Pipefy
- `addCardRelation` - Adicionar relação entre cards
- `deleteCard` - Deletar card
- `moveCardToPhase` - Mover card para outra fase

**Referências**:
- [Pipefy API Documentation - createCard](https://api-docs.pipefy.com/reference/mutations/#createcard)
- [Pipefy API Documentation - updateCard](https://api-docs.pipefy.com/reference/mutations/#updatecard)

### 5. Compartilhado (`shared`)

**Responsabilidade**: Funcionalidades compartilhadas entre contextos.

**Arquivos**:
- `internal/shared/database/connection.go` - Conexão com banco de dados
- `internal/shared/database/transacao.go` - Gerenciamento de transações

**Linguagem Ubíqua**:
- `ConfiguracaoBancoDados` - Configuração de conexão
- `BancoDeTransacao` - Interface para transações

## Mapeamento de Arquivos por Contexto

### Contexto: Gestão de Clientes
```
internal/gestao_clientes/
├── controller.go          # HTTP: POST /clientes
├── service.go             # Lógica: CriarCliente, BuscarPorEmail, Atualizar
├── repository.go          # Persistência: Salvar, BuscarPorEmail, Atualizar
└── validator.go          # Validação: ValidarCriarClienteRequest
```

### Contexto: Processamento de Eventos
```
internal/processamento_eventos/
├── controller.go          # HTTP: POST /webhooks/pipefy/card-updated
├── service.go             # Lógica: ProcessarWebhook, VerificarIdempotencia
├── repository.go          # Persistência: Salvar, VerificarFoiProcessado
└── validator.go          # Validação: ValidarWebhookRequest
```

### Contexto: Domínio
```
internal/dominio/
└── prioridade_calculator.go  # Lógica: CalcularNivelPrioridade
```

### Contexto: Integração Pipefy
```
internal/integracao_pipefy/
├── client.go              # GraphQL: Estruturar mutations, Executar mutations
└── mutations.go           # Mutations: createCard, updateCard, etc.
```

### Contexto: Compartilhado
```
internal/shared/database/
├── connection.go          # Conexão: NovaConexao, NovoBancoDados
└── transacao.go          # Transações: BancoDeTransacao
```

## Trigramação

A trigramação é aplicada consistentemente em toda a codebase:

### Padrão de Nomenclatura

**Tabelas**: `{prefixo_contexto}_{prefixo_entidade}`
- `gestao_clientes.cliente` → `gcl_cli`
- `processamento_eventos.evento` → `pev_eve`

**Colunas**: `{prefixo_tabela}_{prefixo_campo}`
- `gcl_cli_int` - Identificador interno do cliente
- `gcl_cli_ema` - Email do cliente
- `pev_eve_ide` - Identificador do evento
- `pev_eve_fpr` - Flag de processamento

### Prefixos Definidos

| Contexto | Prefixo | Significado |
|----------|---------|-------------|
| gestao_clientes | gcl | Gestão CLientes |
| processamento_eventos | pev | Processamento EVentos |
| cliente | cli | CLiente |
| evento | eve | EVento |

## Integração entre Contextos

### Fluxo: Criação de Cliente → Processamento de Evento

1. **Gestão de Clientes** recebe requisição HTTP
2. **Gestão de Clientes** valida dados via `validator.go`
3. **Gestão de Clientes** processa via `service.go`
4. **Integração Pipefy** cria card via `mutations.go`
5. **Gestão de Clientes** persiste via `repository.go`

### Fluxo: Webhook → Atualização de Cliente

1. **Processamento de Eventos** recebe webhook via `controller.go`
2. **Processamento de Eventos** valida dados via `validator.go`
3. **Processamento de Eventos** verifica idempotência via `repository.go`
4. **Processamento de Eventos** busca cliente via `Gestão de Clientes`
5. **Domínio** calcula prioridade via `prioridade_calculator.go`
6. **Processamento de Eventos** atualiza cliente via `Gestão de Clientes`
7. **Integração Pipefy** atualiza card via `mutations.go`
8. **Processamento de Eventos** persiste evento via `repository.go`

## Views do Banco de Dados

As views fornecem uma interface legível para acesso aos dados:

### `gestao_clientes.vw_cliente`
View que expõe os dados dos clientes com nomes legíveis:
- `cliente_nome` (gcl_cli_nom)
- `cliente_email` (gcl_cli_ema)
- `cliente_status` (gcl_cli_stc)
- `cliente_nivel_prioridade` (gcl_cli_npr)

## Sequências do Banco de Dados

Sequências para geração de identificadores internos:

### `gestao_clientes.seq_gcl_cli_int`
Gera identificadores internos para clientes.

### `processamento_eventos.seq_pev_eve_int`
Gera identificadores internos para eventos.

## Benefícios da Arquitetura

### 1. Separação de Responsabilidades
Cada contexto tem responsabilidade única e bem definida.

### 2. Linguagem Ubíqua
Cada contexto utiliza sua própria linguagem de domínio.

### 3. Manutenibilidade
Mudanças em um contexto não afetam outros contextos.

### 4. Testabilidade
Cada contexto pode ser testado independentemente.

### 5. Escalabilidade
Contextos podem ser escalados independentemente.

## Padrões de Comunicação

### Síncrona
- Controllers chamam Services diretamente
- Services chamam Repositories diretamente
- Integrações entre contexts via interfaces

### Assíncrona (Futuro)
- Webhooks para comunicação entre contextos
- Event-driven architecture para processamento

## Conclusão

A arquitetura de Feature Folders com Contextos Delimitados proporciona uma estrutura clara e manutenível, alinhada com os princípios de DDD e facilitando a evolução do sistema conforme as necessidades de negócio.
