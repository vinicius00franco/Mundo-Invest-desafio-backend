# Diagrama de rotas, classes e sequência

Documentação visual do fluxo HTTP: arquivos, tipos, propriedades e métodos envolvidos em cada rota.

> **Legenda:** `+` método público · `-` dependência/campo privado · tipos em *itálico* são DTOs/entradas HTTP.

---

## 1. Composição na subida do servidor

Montagem em `internal/server/server.go` (`NewServer`).

```mermaid
flowchart LR
    subgraph server["server/server.go"]
        S["Server<br/>-httpServer<br/>-config<br/>-db"]
    end

    subgraph infra["shared"]
        DB["database.NovaConexao()"]
        CFG["config.Config"]
    end

    subgraph gcl["gestao_clientes"]
        CR["ClienteRepository"]
        CS["ClienteService"]
        CC["ClienteController"]
    end

    subgraph pev["processamento_eventos"]
        ER["EventoRepository"]
        WS["WebhookService"]
        EC["EventoController"]
    end

    subgraph dom["dominio"]
        CP["CalculadoraPrioridade"]
        ED["EventDispatcher"]
    end

    subgraph pipefy["integracao_pipefy"]
        PC["PipefyGraphQLClient"]
        PS["PipefyIntegrationService"]
    end

    S --> DB
    S --> CR & ER
    S --> CP & ED & PC
    PC --> PS
    CR & PS & ED --> CS
    CS --> CC
    ER & CR & CP --> WS
    WS --> EC
    CC & EC --> MUX["http.ServeMux<br/>POST /clientes<br/>POST /webhooks/pipefy/card-updated"]
```

---

## 2. Diagrama de classes por contexto

### 2.1 Gestão de Clientes — `POST /clientes`

```mermaid
classDiagram
    direction TB

    class RequisicaoCriarCliente {
        +string Nome
        +string Email
        +float64 ValorPatrimonio
        +string TipoSolicitacao
    }

    class ClienteController {
        -ClienteService service
        +CriarClienteHandler(w, r)
        +RegistrarRotas(mux)
    }

    class ClienteService {
        <<interface>>
        +CriarCliente(ctx, RequisicaoCriarCliente) Cliente
    }

    class clienteService {
        -ClienteRepository repository
        -PipefyIntegrationService pipefyService
        -string pipeID
        -EventDispatcher eventDispatcher
        -Config config
        +CriarCliente(ctx, requisicao)
    }

    class Cliente {
        +int64 IdentificadorInterno
        +string IdentificadorExterno
        +string Nome
        +string Email
        +float64 ValorPatrimonio
        +string TipoSolicitacao
        +string Status
        +string NivelPrioridade
        +time Time DataCriacao
        +time Time DataAtualizacao
    }

    class ClienteRepository {
        <<interface>>
        +Salvar(ctx, Cliente) Cliente
        +Atualizar(ctx, Cliente)
        +BuscarPorEmail(ctx, email) Cliente
    }

    class CriarClienteResponse {
        +string Mensagem
        +int64 IdentificadorInterno
        +string IdentificadorExterno
        +string Nome
        +string Email
        +float64 ValorPatrimonio
        +string TipoSolicitacao
        +string Status
        +time Time DataCriacao
    }

    class ValidarRequisicaoCriarCliente {
        <<função>>
        +Validate via ClienteValidationStrategy
    }

    ClienteController --> ClienteService
    clienteService ..|> ClienteService
    clienteService --> ClienteRepository
    clienteService --> RequisicaoCriarCliente
    clienteService --> Cliente
    ClienteController --> RequisicaoCriarCliente
    ClienteController --> ValidarRequisicaoCriarCliente
    ClienteController --> CriarClienteResponse
```

**Arquivos**

| Tipo | Arquivo |
|------|---------|
| Controller | `internal/gestao_clientes/controller.go` |
| Service | `internal/gestao_clientes/service.go` |
| Entidade / Repository | `internal/gestao_clientes/repository.go` |
| DTO entrada | `internal/gestao_clientes/validator.go` (`RequisicaoCriarCliente`) |
| DTO saída | `internal/gestao_clientes/dto.go` |
| Validação | `internal/gestao_clientes/validator.go`, `validation_strategy.go` |

### 2.2 Integração Pipefy (fluxo criação de card)

```mermaid
classDiagram
    direction TB

    class PipefyIntegrationService {
        <<interface>>
        +CriarCardCliente(ctx, pipeID, CardClienteData) string
        +AtualizarCardPrioridade(ctx, cardID, nivelPrioridade)
    }

    class pipefyIntegrationService {
        -PipefyGraphQLClient client
        -Config config
        +CriarCardCliente()
        +AtualizarCardPrioridade()
    }

    class CardClienteData {
        +string Nome
        +string Email
        +float64 ValorPatrimonio
        +string TipoSolicitacao
    }

    class PipefyGraphQLClient {
        <<interface>>
        +EstruturarMutationCreateCard(pipeID, fields) string
        +EstruturarMutationUpdateCard(cardID, fields) string
        +ExecutarMutation(mutation) string
    }

    class pipefyGraphQLClient {
        -string apiToken
        -string apiURL
        +EstruturarMutationCreateCard()
        +EstruturarMutationUpdateCard()
        +ExecutarMutation()
    }

    class FieldAttribute {
        +string FieldID
        +[]string Values
    }

    class Mutacoes {
        -PipefyGraphQLClient cliente
        +CreateCardMutation()
        +UpdateCardMutation()
    }

    pipefyIntegrationService ..|> PipefyIntegrationService
    pipefyIntegrationService --> PipefyGraphQLClient
    pipefyGraphQLClient ..|> PipefyGraphQLClient
    pipefyIntegrationService --> CardClienteData
    pipefyGraphQLClient --> FieldAttribute
    Mutacoes --> PipefyGraphQLClient
    Mutacoes ..> "não usado em server.go" : legado
```

**Arquivos**

| Camada | Arquivo |
|--------|---------|
| Serviço aplicação | `internal/integracao_pipefy/service.go` |
| Cliente GraphQL | `internal/integracao_pipefy/cliente_pipefy.go` |
| Facade opcional | `internal/integracao_pipefy/mutacoes.go` *(não ligado ao `server`)* |

### 2.3 Processamento de Eventos — `POST /webhooks/pipefy/card-updated`

```mermaid
classDiagram
    direction TB

    class RequisicaoWebhook {
        +string IdentificadorEvento
        +string IdentificadorCard
        +string EmailCliente
        +string DataEvento
    }

    class EventoController {
        -WebhookService service
        +ProcessarWebhookHandler(w, r)
        +RegistrarRotas(mux)
    }

    class WebhookService {
        <<interface>>
        +ProcessarWebhook(ctx, RequisicaoWebhook)
    }

    class webhookService {
        -EventoRepository eventoRepository
        -ClienteRepository clienteRepository
        -CalculadoraPrioridade calculadoraPrioridade
        -PipefyGraphQLClient pipefyClient
        -Config config
        +ProcessarWebhook()
    }

    class Evento {
        +int64 IdentificadorInterno
        +string IdentificadorEvento
        +string IdentificadorCard
        +string EmailCliente
        +time Time TimestampEvento
        +bool FoiProcessado
        +time Time DataCriacao
        +*time Time DataAtualizacao
    }

    class EventoRepository {
        <<interface>>
        +Salvar(ctx, Evento) Evento
        +VerificarFoiProcessado(ctx, eventID) bool
    }

    class CalculadoraPrioridade {
        <<interface>>
        +CalcularNivelPrioridade(valorPatrimonio) string
    }

    class calculadoraPrioridade {
        -Config config
        +CalcularNivelPrioridade()
    }

    EventoController --> WebhookService
    webhookService ..|> WebhookService
    webhookService --> EventoRepository
    webhookService --> ClienteRepository
    webhookService --> CalculadoraPrioridade
    webhookService --> PipefyGraphQLClient
    webhookService --> RequisicaoWebhook
    webhookService --> Evento
    calculadoraPrioridade ..|> CalculadoraPrioridade
```

**Arquivos**

| Tipo | Arquivo |
|------|---------|
| Controller | `internal/processamento_eventos/controller.go` |
| Service | `internal/processamento_eventos/service.go` |
| Entidade / Repository | `internal/processamento_eventos/repository.go` |
| DTO entrada | `internal/processamento_eventos/validator.go` |
| Domínio | `internal/dominio/calculadora_prioridade.go` |
| Pipefy (update) | `internal/integracao_pipefy/cliente_pipefy.go` |

---

## 3. Sequência — `POST /clientes`

```mermaid
sequenceDiagram
    autonumber
    actor Cliente as Cliente HTTP
    participant MUX as http.ServeMux
    participant CTRL as controller.go<br/>ClienteController
    participant VAL as validator.go<br/>ValidarRequisicaoCriarCliente
    participant SVC as service.go<br/>clienteService
    participant REPO as repository.go<br/>clienteRepository
    participant DOM as dominio<br/>StatusAguardandoAnalise
    participant PSV as pipefy/service.go<br/>pipefyIntegrationService
    participant PCL as pipefy/cliente_pipefy.go<br/>pipefyGraphQLClient
    participant DB as PostgreSQL<br/>gestao_clientes.cliente

    Cliente->>MUX: POST /clientes<br/>RequisicaoCriarCliente JSON
    MUX->>CTRL: CriarClienteHandler()

    CTRL->>VAL: ValidarRequisicaoCriarCliente(requisicao)
    VAL-->>CTRL: ok | erro validação

    CTRL->>SVC: CriarCliente(ctx, requisicao)

    SVC->>VAL: ValidarRequisicaoCriarCliente (novamente)
    SVC->>SVC: Monta Cliente<br/>Status = Aguardando Análise
    Note over SVC,DOM: propriedades: Nome, Email,<br/>ValorPatrimonio, TipoSolicitacao

    SVC->>REPO: Salvar(ctx, cliente)
    REPO->>DB: INSERT cliente
    DB-->>REPO: IdentificadorInterno
    REPO-->>SVC: *Cliente

    SVC->>PSV: CriarCardCliente(ctx, pipeID, CardClienteData)
    PSV->>PSV: fieldsAttributes[]<br/>nome, email, patrimônio, tipo
    PSV->>PCL: EstruturarMutationCreateCard(pipeID, fields)
    PCL-->>PSV: mutation createCard (string GraphQL)
    PSV->>PCL: ExecutarMutation(mutation)
    PCL-->>PSV: cardID simulado

    SVC->>REPO: Atualizar(ctx, cliente)<br/>IdentificadorExterno = cardID
    REPO->>DB: UPDATE cliente

    SVC->>SVC: eventDispatcher.Dispatch(ClienteCriadoEvent)
    SVC-->>CTRL: *Cliente
    CTRL->>CTRL: NewCriarClienteResponse(cliente)
    CTRL-->>Cliente: 201 CriarClienteResponse JSON
```

---

## 4. Sequência — `POST /webhooks/pipefy/card-updated`

```mermaid
sequenceDiagram
    autonumber
    actor Pipefy as Pipefy / simulador
    participant MUX as http.ServeMux
    participant CTRL as controller.go<br/>EventoController
    participant VAL as validator.go<br/>ValidarRequisicaoWebhook
    participant SVC as service.go<br/>webhookService
    participant EVR as repository.go<br/>eventoRepository
    participant CLR as gestao_clientes/<br/>clienteRepository
    participant CALC as dominio/<br/>calculadoraPrioridade
    participant PCL as pipefy/cliente_pipefy.go<br/>pipefyGraphQLClient
    participant DBE as PostgreSQL<br/>processamento_eventos.evento
    participant DBC as PostgreSQL<br/>gestao_clientes.cliente

    Pipefy->>MUX: POST /webhooks/pipefy/card-updated<br/>RequisicaoWebhook JSON
    MUX->>CTRL: ProcessarWebhookHandler()

    CTRL->>VAL: ValidarRequisicaoWebhook(requisicao)
    VAL-->>CTRL: ok | erro

    CTRL->>SVC: ProcessarWebhook(ctx, requisicao)

    SVC->>EVR: VerificarFoiProcessado(ctx, IdentificadorEvento)
    EVR->>DBE: SELECT por event_id
  alt já processado
        EVR-->>SVC: foiProcessado = true
        SVC-->>CTRL: nil (idempotente)
        CTRL-->>Pipefy: 200 sucesso
    else novo evento
        EVR-->>SVC: foiProcessado = false

        SVC->>CLR: BuscarPorEmail(ctx, EmailCliente)
        CLR->>DBC: SELECT cliente
        DBC-->>CLR: *Cliente (ValorPatrimonio, Status, …)
        CLR-->>SVC: *Cliente

        SVC->>CALC: CalcularNivelPrioridade(ValorPatrimonio)
        Note over CALC: ≥ 200.000 → prioridade_alta<br/>< 200.000 → prioridade_normal
        CALC-->>SVC: nivelPrioridade

        SVC->>SVC: cliente.Status = Processado<br/>cliente.NivelPrioridade = nivel
        SVC->>CLR: Atualizar(ctx, cliente)
        CLR->>DBC: UPDATE cliente

        SVC->>PCL: EstruturarMutationUpdateCard(IdentificadorCard, fields)
        PCL-->>SVC: mutation updateCard
        SVC->>PCL: ExecutarMutation(mutation)

        SVC->>SVC: Monta Evento<br/>FoiProcessado = true
        SVC->>EVR: Salvar(ctx, evento)
        EVR->>DBE: INSERT evento
        EVR-->>SVC: *Evento

        SVC-->>CTRL: nil
        CTRL-->>Pipefy: 200 JSON sucesso
    end
```

---

## 5. Mapa arquivo → responsabilidade

| Ordem | Arquivo | Papel na rota |
|------|---------|----------------|
| 1 | `cmd/server/main.go` | `godotenv`, `config.Load`, `server.NewServer` |
| 2 | `internal/server/server.go` | DI, rotas, `ListenAndServe` |
| 3a | `gestao_clientes/controller.go` | HTTP `POST /clientes` |
| 3b | `processamento_eventos/controller.go` | HTTP `POST /webhooks/...` |
| 4a | `gestao_clientes/service.go` | Regra de criação + orquestração Pipefy |
| 4b | `processamento_eventos/service.go` | Idempotência, prioridade, update Pipefy |
| 5 | `gestao_clientes/repository.go` | Persistência `Cliente` |
| 6 | `processamento_eventos/repository.go` | Persistência `Evento` |
| 7 | `dominio/calculadora_prioridade.go` | Regra patrimônio → prioridade |
| 8 | `integracao_pipefy/service.go` | Caso de uso createCard (rota clientes) |
| 9 | `integracao_pipefy/cliente_pipefy.go` | Strings GraphQL createCard / updateCard |
| — | `integracao_pipefy/mutacoes.go` | Alternativa não usada pelo `server` |

---

## 6. Referência rápida das rotas

| Método | Rota | Handler | Service | Persistência | Pipefy |
|--------|------|---------|---------|--------------|--------|
| POST | `/clientes` | `ClienteController.CriarClienteHandler` | `clienteService.CriarCliente` | `ClienteRepository` | `PipefyIntegrationService` → `PipefyGraphQLClient` |
| POST | `/webhooks/pipefy/card-updated` | `EventoController.ProcessarWebhookHandler` | `webhookService.ProcessarWebhook` | `EventoRepository` + `ClienteRepository` | `PipefyGraphQLClient` direto (`updateCard`) |

---

*Gerado a partir do código em `internal/` — alinhar após refactors (ex.: unificar webhook com `PipefyIntegrationService`).*
