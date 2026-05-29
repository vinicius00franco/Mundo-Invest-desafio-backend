# Plano de Renomeação para Linguagem Ubíqua em Português

## Resumo Executivo

Este documento detalha o plano de migração dos nomes de arquivos, variáveis e funções para português, seguindo a linguagem ubíqua do domínio do Mundo Invest.

## Diagnóstico Atual

### Inconsistências Identificadas

**Arquivos com nomes em inglês:**
- `client.go` → deveria ser `cliente_pipefy.go`
- `mutations.go` → deveria ser `mutacoes.go`
- `dto.go` → manter como DTO (termo técnico)

**Interfaces com nomes em inglês:**
- `WebhookService` → `ServicoWebhook`
- `PipefyIntegrationService` → `ServicoIntegracaoPipefy`
- `PipefyGraphQLClient` → `ClienteGraphQLPipefy`

**Métodos com nomes em inglês:**
- `EstruturarMutationCreateCard` → `EstruturarMutacaoCriarCard`
- `EstruturarMutationUpdateCard` → `EstruturarMutacaoAtualizarCard`
- `ExecutarMutation` → `ExecutarMutacao`

**Variáveis com nomes em inglês:**
- `pipeID` → manter (sistema externo)
- `cardID` → manter (sistema externo)
- `fieldID` → manter (sistema externo)

## Plano de Renomeação

### Fase 1: Arquivos de Integração Pipefy

#### Arquivo: `internal/integracao_pipefy/client.go`
**Renomear para:** `internal/integracao_pipefy/cliente_pipefy.go`

**Mudanças no código:**
```go
// Antes
package integracao_pipefy

type PipefyGraphQLClient interface {
    EstruturarMutationCreateCard(pipeID string, fieldsAttributes []FieldAttribute) (string, error)
    EstruturarMutationUpdateCard(cardID string, fieldsAttributes []FieldAttribute) (string, error)
    ExecutarMutation(mutation string) (string, error)
}

type pipefyGraphQLClient struct {
    apiToken string
    apiURL   string
}

// Depois
package integracao_pipefy

type ClienteGraphQLPipefy interface {
    EstruturarMutacaoCriarCard(pipeID string, atributosCampo []AtributoCampo) (string, error)
    EstruturarMutacaoAtualizarCard(cardID string, atributosCampo []AtributoCampo) (string, error)
    ExecutarMutacao(mutacao string) (string, error)
}

type clienteGraphQLPipefy struct {
    tokenAPI string
    urlAPI   string
}
```

#### Arquivo: `internal/integracao_pipefy/mutations.go`
**Renomear para:** `internal/integracao_pipefy/mutacoes.go`

**Mudanças no código:**
```go
// Antes
package integracao_pipefy

type Mutations struct {
    client PipefyGraphQLClient
}

func (m *Mutations) CreateCardMutation(...) (string, error) {
    // ...
}

// Depois
package integracao_pipefy

type Mutacoes struct {
    cliente ClienteGraphQLPipefy
}

func (m *Mutacoes) CriarMutacaoCard(...) (string, error) {
    // ...
}
```

#### Arquivo: `internal/integracao_pipefy/service.go`
**Mudanças no código:**
```go
// Antes
type PipefyIntegrationService interface {
    CriarCardCliente(ctx context.Context, pipeID string, dados CardClienteData) (string, error)
    AtualizarCardPrioridade(ctx context.Context, cardID string, nivelPrioridade string) error
}

type pipefyIntegrationService struct {
    client PipefyGraphQLClient
    config *config.Config
}

// Depois
type ServicoIntegracaoPipefy interface {
    CriarCardCliente(ctx context.Context, pipeID string, dados DadosCardCliente) (string, error)
    AtualizarPrioridadeCard(ctx context.Context, cardID string, nivelPrioridade string) error
}

type servicoIntegracaoPipefy struct {
    cliente ClienteGraphQLPipefy
    config  *config.Config
}
```

### Fase 2: Processamento de Eventos

#### Arquivo: `internal/processamento_eventos/service.go`
**Mudanças no código:**
```go
// Antes
type WebhookService interface {
    ProcessarWebhook(ctx context.Context, request WebhookRequest) error
}

type webhookService struct {
    eventoRepository     EventoRepository
    clienteRepository    gestao_clientes.ClienteRepository
    prioridadeCalculator dominio.PrioridadeCalculator
    pipefyClient         integracao_pipefy.PipefyGraphQLClient
    config               *config.Config
}

// Depois
type ServicoWebhook interface {
    ProcessarWebhook(ctx context.Context, requisicao RequisicaoWebhook) error
}

type servicoWebhook struct {
    repositorioEvento          gestao_clientes.RepositorioEvento
    repositorioCliente         gestao_clientes.RepositorioCliente
    calculadoraPrioridade      dominio.CalculadoraPrioridade
    clientePipefy             integracao_pipefy.ClienteGraphQLPipefy
    config                    *config.Config
}
```

#### Arquivo: `internal/processamento_eventos/validator.go`
**Mudanças no código:**
```go
// Antes
type WebhookRequest struct {
    IdentificadorEvento string `json:"identificadorEvento"`
    IdentificadorCard   string `json:"identificadorCard"`
    ClienteEmail        string `json:"clienteEmail"`
    DataEvento          string `json:"dataEvento"`
}

func ValidarWebhookRequest(request WebhookRequest) error {
    // ...
}

// Depois
type RequisicaoWebhook struct {
    IdentificadorEvento string `json:"identificadorEvento"`
    IdentificadorCard   string `json:"identificadorCard"`
    EmailCliente        string `json:"emailCliente"`
    DataEvento          string `json:"dataEvento"`
}

func ValidarRequisicaoWebhook(requisicao RequisicaoWebhook) error {
    // ...
}
```

#### Arquivo: `internal/processamento_eventos/controller.go`
**Mudanças no código:**
```go
// Antes
type EventoController struct {
    service WebhookService
}

func NovoEventoController(service WebhookService) *EventoController {
    // ...
}

func (c *EventoController) ProcessarWebhookHandler(w http.ResponseWriter, r *http.Request) {
    // ...
}

// Depois
type ControladorEvento struct {
    servico ServicoWebhook
}

func NovoControladorEvento(servico ServicoWebhook) *ControladorEvento {
    // ...
}

func (c *ControladorEvento) ProcessarWebhookHandler(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

### Fase 3: Domínio

#### Arquivo: `internal/dominio/prioridade_calculator.go`
**Renomear para:** `internal/dominio/calculadora_prioridade.go`

**Mudanças no código:**
```go
// Antes
type PrioridadeCalculator interface {
    CalcularNivelPrioridade(valorPatrimonio float64) string
}

type prioridadeCalculator struct {
    config *config.Config
}

func NovoPrioridadeCalculator(cfg *config.Config) PrioridadeCalculator {
    // ...
}

// Depois
type CalculadoraPrioridade interface {
    CalcularNivelPrioridade(valorPatrimonio float64) string
}

type calculadoraPrioridade struct {
    config *config.Config
}

func NovaCalculadoraPrioridade(cfg *config.Config) CalculadoraPrioridade {
    // ...
}
```

### Fase 4: Gestão de Clientes

#### Arquivo: `internal/gestao_clientes/service.go`
**Mudanças no código:**
```go
// Antes
type ClienteService interface {
    CriarCliente(ctx context.Context, request CriarClienteRequest) (*Cliente, error)
}

type clienteService struct {
    repository      ClienteRepository
    pipefyService   integracao_pipefy.PipefyIntegrationService
    pipeID          string
    eventDispatcher dominio.EventDispatcher
    config          *config.Config
}

// Depois
type ServicoCliente interface {
    CriarCliente(ctx context.Context, requisicao RequisicaoCriarCliente) (*Cliente, error)
}

type servicoCliente struct {
    repositorio          RepositorioCliente
    servicoPipefy        integracao_pipefy.ServicoIntegracaoPipefy
    identificadorPipe    string
    despachanteEventos    dominio.DespachanteEventos
    config               *config.Config
}
```

#### Arquivo: `internal/gestao_clientes/controller.go`
**Mudanças no código:**
```go
// Antes
type ClienteController struct {
    service ClienteService
}

func NovoClienteController(service ClienteService) *ClienteController {
    // ...
}

func (c *ClienteController) CriarClienteHandler(w http.ResponseWriter, r *http.Request) {
    // ...
}

// Depois
type ControladorCliente struct {
    servico ServicoCliente
}

func NovoControladorCliente(servico ServicoCliente) *ControladorCliente {
    // ...
}

func (c *ControladorCliente) CriarClienteHandler(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

#### Arquivo: `internal/gestao_clientes/validator.go`
**Mudanças no código:**
```go
// Antes
type CriarClienteRequest struct {
    Nome            string  `json:"nome"`
    Email           string  `json:"email"`
    ValorPatrimonio float64 `json:"valorPatrimonio"`
    TipoSolicitacao string  `json:"tipoSolicitacao"`
}

func ValidarCriarClienteRequest(request CriarClienteRequest) error {
    // ...
}

// Depois
type RequisicaoCriarCliente struct {
    Nome            string  `json:"nome"`
    Email           string  `json:"email"`
    ValorPatrimonio float64 `json:"valorPatrimonio"`
    TipoSolicitacao string  `json:"tipoSolicitacao"`
}

func ValidarRequisicaoCriarCliente(requisicao RequisicaoCriarCliente) error {
    // ...
}
```

### Fase 5: Infraestrutura Compartilhada

#### Arquivo: `internal/shared/database/connection.go`
**Mudanças no código:**
```go
// Antes
type ConfiguracaoBancoDados struct {
    Host    string
    Port    string
    Usuario string
    Senha   string
    Banco   string
    SSLMode string
}

func NovoBancoDados(cfg ConfiguracaoBancoDados, appConfig *config.Config) (*sql.DB, error) {
    // ...
}

// Depois
type ConfiguracaoBancoDados struct {
    Host    string
    Porta    string
    Usuario  string
    Senha    string
    Banco    string
    ModoSSL  string
}

func NovoBancoDados(cfg ConfiguracaoBancoDados, configApp *config.Config) (*sql.DB, error) {
    // ...
}
```

#### Arquivo: `internal/shared/errors/errors.go`
**Mantido como está** - Os nomes das estruturas de erro são termos técnicos padrão.

## Mapeamento Completo de Renomeação

### Pacotes (Sem alteração - já em português)
- `gestao_clientes` ✅
- `processamento_eventos` ✅
- `dominio` ✅
- `integracao_pipefy` ✅
- `shared` ✅

### Arquivos
| Arquivo Atual | Arquivo Proposto | Motivo |
|---------------|------------------|---------|
| `integracao_pipefy/client.go` | `integracao_pipefy/cliente_pipefy.go` | Descritivo em português |
| `integracao_pipefy/mutations.go` | `integracao_pipefy/mutacoes.go` | Tradução direta |
| `dominio/prioridade_calculator.go` | `dominio/calculadora_prioridade.go` | Descritivo em português |
| `shared/database/connection.go` | `shared/database/conexao.go` | Tradução direta |
| `shared/database/transacao.go` | `shared/database/transacao.go` | Já está correto |

### Interfaces
| Interface Atual | Interface Proposta | Motivo |
|------------------|---------------------|---------|
| `WebhookService` | `ServicoWebhook` | Termo técnico mantido |
| `PipefyIntegrationService` | `ServicoIntegracaoPipefy` | Tradução |
| `PipefyGraphQLClient` | `ClienteGraphQLPipefy` | Tradução |
| `PrioridadeCalculator` | `CalculadoraPrioridade` | Tradução |
| `EventDispatcher` | `DespachanteEventos` | Tradução |

### Structs
| Struct Atual | Struct Proposto | Motivo |
|--------------|----------------|---------|
| `WebhookRequest` | `RequisicaoWebhook` | Tradução |
| `CardClienteData` | `DadosCardCliente` | Tradução |
| `FieldAttribute` | `AtributoCampo` | Tradução |
| `ValidationError` | `ErroValidacao` | Tradução |
| `RepositoryError` | `ErroRepositorio` | Tradução |

### Métodos
| Método Atual | Método Proposto | Motivo |
|--------------|-----------------|---------|
| `EstruturarMutationCreateCard` | `EstruturarMutacaoCriarCard` | Tradução |
| `EstruturarMutationUpdateCard` | `EstruturarMutacaoAtualizarCard` | Tradução |
| `ExecutarMutation` | `ExecutarMutacao` | Tradução |
| `ValidarWebhookRequest` | `ValidarRequisicaoWebhook` | Tradução |
| `NovoPrioridadeCalculator` | `NovaCalculadoraPrioridade` | Tradução |

### Variáveis (Mantidas em Inglês)
| Variável | Motivo |
|----------|---------|
| `pipeID` | Sistema externo Pipefy |
| `cardID` | Sistema externo Pipefy |
| `fieldID` | Sistema externo Pipefy |
| `apiToken` | Termo técnico padrão |
| `apiURL` | Termo técnico padrão |

## Ordem de Implementação

### 1. Preparação (Dia 1)
- Criar branch `feature/linguagem-ubiqua-portugues`
- Atualizar documentação
- Comunicar equipe sobre mudanças

### 2. Domínio (Dia 2)
- Renomear `prioridade_calculator.go` → `calculadora_prioridade.go`
- Atualizar interfaces e structs
- Atualizar testes

### 3. Integração Pipefy (Dias 3-4)
- Renomear `client.go` → `cliente_pipefy.go`
- Renomear `mutations.go` → `mutacoes.go`
- Atualizar interfaces e métodos
- Atualizar testes

### 4. Processamento Eventos (Dia 5)
- Atualizar interfaces e structs
- Atualizar métodos
- Atualizar testes

### 5. Gestão Clientes (Dia 6)
- Atualizar interfaces e structs
- Atualizar métodos
- Atualizar testes

### 6. Infraestrutura (Dia 7)
- Atualizar structs de configuração
- Atualizar testes

### 7. Validação Final (Dia 8)
- Executar todos os testes
- Verificar build
- Atualizar documentação

## Testes de Validação

### Checklist de Validação
- [ ] Todos os testes unitários passam
- [ ] Todos os testes de integração passam
- [ ] Build funciona sem erros
- [ ] API responde corretamente
- [ ] Documentação atualizada
- [ ] Code review realizado
- [ ] Merge request aprovada

## Riscos e Mitigações

### Risco 1: Quebra de Integrações Externas
**Mitigação:** Manter IDs de sistemas externos em inglês (pipeID, cardID)

### Risco 2: Curva de Aprendizado da Equipe
**Mitigação:** Documentação detalhada e sessões de treinamento

### Risco 3: Conflitos em Branches Paralelos
**Mitigação:** Coordenação com equipe e merge strategy adequado

## Conclusão

Este plano proporciona uma migração sistemática para linguagem ubíqua em português, mantendo termos técnicos quando apropriado e garantindo consistência em todo o código.