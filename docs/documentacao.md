# Documentação Técnica - Mundo Invest

## Arquitetura do Sistema

### Visão Geral
```
┌─────────────────┐         ┌─────────────────┐
│   Cliente       │         │   API Backend   │
│   (HTTP Client) │◄────────┤   (Golang)       │
│                 │  HTTP    │                 │
│  - POST /clientes│         │  - Controller   │
│  - POST /webhooks│         │  - Service      │
└─────────────────┘         │  - Repository   │
                            │  - GraphQL Client│
                            └─────────────────┘
                                      │
                                      ▼
                              ┌─────────────────┐
                              │  Banco de Dados │
                              │  (SQLite/PostgreSQL)│
                              └─────────────────┘
```

## Estrutura de Pastas Sugerida

```
MundoInvest/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── controllers/
│   │   ├── clienteController.go
│   │   └── webhookController.go
│   ├── services/
│   │   ├── clienteService.go
│   │   ├── cardService.go
│   │   └── pipefyService.go
│   ├── repositories/
│   │   ├── clienteRepository.go
│   │   └── eventoRepository.go
│   ├── models/
│   │   ├── cliente.go
│   │   ├── evento.go
│   │   └── card.go
│   ├── graphql/
│   │   ├── mutations.go
│   │   └── queries.go
│   └── middleware/
│       └── logging.go
├── pkg/
│   ├── validation/
│   │   └── validator.go
│   └── database/
│       └── connection.go
├── bdd/
│   ├── criacao-cliente.feature
│   └── webhook-card-updated.feature
├── docs/
│   ├── documentacao.md
│   ├── regras-negocio.md
│   └── requisitos-funcionais.md
├── go.mod
├── go.sum
└── README.md
```

## Modelos de Dados

### Cliente
```go
type Cliente struct {
    ID              string    `json:"id" db:"id"`
    Nome            string    `json:"cliente_nome" db:"nome"`
    Email           string    `json:"cliente_email" db:"email"`
    TipoSolicitacao string    `json:"tipo_solicitacao" db:"tipo_solicitacao"`
    ValorPatrimonio float64  `json:"valor_patrimonio" db:"valor_patrimonio"`
    Status          string    `json:"status" db:"status"`
    Prioridade      string    `json:"prioridade" db:"prioridade"`
    CardID          string    `json:"card_id" db:"card_id"`
    CriadoEm        time.Time `json:"criado_em" db:"criado_em"`
    AtualizadoEm    time.Time `json:"atualizado_em" db:"atualizado_em"`
}
```

### Evento (Webhook)
```go
type Evento struct {
    ID           string    `json:"id" db:"id"`
    EventID      string    `json:"event_id" db:"event_id"`
    CardID       string    `json:"card_id" db:"card_id"`
    ClienteEmail string    `json:"cliente_email" db:"cliente_email"`
    Timestamp    time.Time `json:"timestamp" db:"timestamp"`
    Processado   bool      `json:"processado" db:"processado"`
    CriadoEm     time.Time `json:"criado_em" db:"criado_em"`
}
```

### Requisição de Criação de Cliente
```go
type CriarClienteRequest struct {
    Nome            string  `json:"cliente_nome"`
    Email           string  `json:"cliente_email"`
    TipoSolicitacao string  `json:"tipo_solicitacao"`
    ValorPatrimonio float64 `json:"valor_patrimonio"`
}
```

### Requisição de Webhook
```go
type WebhookRequest struct {
    EventID      string    `json:"event_id"`
    CardID       string    `json:"card_id"`
    ClienteEmail string    `json:"cliente_email"`
    Timestamp    time.Time `json:"timestamp"`
}
```

### Resposta de Sucesso
```go
type SucessoResponse struct {
    Sucesso  bool   `json:"sucesso"`
    Mensagem string `json:"mensagem"`
    Dados    any    `json:"dados,omitempty"`
}
```

### Resposta de Erro
```go
type ErroResponse struct {
    Sucesso  bool   `json:"sucesso"`
    Mensagem string `json:"mensagem"`
    Erro     string `json:"erro"`
}
```

## Fluxo de Processamento

### Fluxo 1: Criação de Cliente (Sucesso)
```
1. Cliente envia POST /clientes
2. Controller valida payload (campos obrigatórios, e-mail válido)
3. Service chama Repository para salvar no banco
4. Status inicial: "Aguardando Análise"
5. Service estrutura mutation GraphQL createCard
6. Service simula envio para Pipefy (não conecta realmente)
7. Repository retorna cliente criado
8. Controller retorna HTTP 201 com dados do cliente
```

### Fluxo 2: Processamento de Webhook (Sucesso)
```
1. Pipefy envia POST /webhooks/pipefy/card-updated
2. Controller valida payload
3. Service verifica idempotência pelo event_id
4. Se já processado, retorna HTTP 200 sem processar novamente
5. Service busca cliente por e-mail no banco
6. Service aplica regra de prioridade:
   - Se patrimonio >= 200.000 → prioridade_alta
   - Se patrimonio < 200.000 → prioridade_normal
7. Service estrutura mutation GraphQL updateCard
8. Service simula envio para Pipefy
9. Service atualiza status para "Processado" e prioridade no banco
10. Controller retorna HTTP 200
```

## Endpoints da API

### POST /clientes
**Descrição**: Cria um novo cliente e mapeia para card no Pipefy

**Request Body**:
```json
{
  "cliente_nome": "João Silva",
  "cliente_email": "joao.silva@example.com",
  "tipo_solicitacao": "Atualização cadastral",
  "valor_patrimonio": 250000
}
```

**Respostas**:

- **201 Created** (Sucesso)
```json
{
  "sucesso": true,
  "mensagem": "Cliente criado com sucesso",
  "dados": {
    "id": "uuid-do-cliente",
    "cliente_nome": "João Silva",
    "cliente_email": "joao.silva@example.com",
    "tipo_solicitacao": "Atualização cadastral",
    "valor_patrimonio": 250000,
    "status": "Aguardando Análise",
    "card_id": "card-gerado"
  }
}
```

- **400 Bad Request** (Erro de validação)
```json
{
  "sucesso": false,
  "mensagem": "Dados inválidos",
  "erro": "VALIDACAO_CAMPOS_OBRIGATORIOS"
}
```

- **400 Bad Request** (E-mail inválido)
```json
{
  "sucesso": false,
  "mensagem": "E-mail inválido",
  "erro": "VALIDACAO_EMAIL"
}
```

### POST /webhooks/pipefy/card-updated
**Descrição**: Processa webhook de atualização de card do Pipefy

**Request Body**:
```json
{
  "event_id": "evt_123",
  "card_id": "card_456",
  "cliente_email": "joao.silva@example.com",
  "timestamp": "2026-05-18T12:00:00Z"
}
```

**Respostas**:

- **200 OK** (Sucesso)
```json
{
  "sucesso": true,
  "mensagem": "Webhook processado com sucesso",
  "dados": {
    "cliente_email": "joao.silva@example.com",
    "prioridade": "prioridade_alta",
    "status": "Processado"
  }
}
```

- **200 OK** (Idempotência - já processado)
```json
{
  "sucesso": true,
  "mensagem": "Evento já processado anteriormente",
  "dados": {
    "event_id": "evt_123"
  }
}
```

- **404 Not Found** (Cliente não encontrado)
```json
{
  "sucesso": false,
  "mensagem": "Cliente não encontrado",
  "erro": "CLIENTE_NAO_ENCONTRADO"
}
```

## GraphQL Mutations (Pipefy)

### createCard
```graphql
mutation createCard($input: CreateCardInput!) {
  createCard(input: $input) {
    success
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

### updateCardField
```graphql
mutation updateCardField($input: UpdateCardFieldInput!) {
  updateCardField(input: $input) {
    success
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

## Cenários de Teste

### Cenários de Sucesso
| Cenário | Descrição | Resultado Esperado |
|---------|-----------|-------------------|
| TC-001 | Criar cliente com dados válidos | HTTP 201, cliente salvo no banco |
| TC-002 | Processar webhook com patrimonio >= 200.000 | HTTP 200, prioridade_alta definida |
| TC-003 | Processar webhook com patrimonio < 200.000 | HTTP 200, prioridade_normal definida |
| TC-004 | Processar webhook duplicado (mesmo event_id) | HTTP 200, não processa novamente |

### Cenários de Erro
| Cenário | Descrição | Resultado Esperado |
|---------|-----------|-------------------|
| TC-005 | Criar cliente sem campos obrigatórios | HTTP 400, mensagem de erro |
| TC-006 | Criar cliente com e-mail inválido | HTTP 400, mensagem de erro |
| TC-007 | Processar webhook com cliente inexistente | HTTP 404, mensagem de erro |

## Decisões e Trade-offs

| Decisão | Justificativa | Trade-off |
|---------|---------------|-----------|
| SQLite para desenvolvimento | Simplicidade para o desafio | Não escalável para produção |
| Simulação de Pipefy | Não depende de serviço externo | Não testa integração real |
| Idempotência por event_id | Evita processamento duplicado | Requer armazenamento de eventos |
| Prioridade baseada em patrimônio | Regra de negócio simples | Limite fixo (200.000) |

## Dependências Sugeridas

### Backend Golang
- `github.com/gorilla/mux`: Router HTTP
- `github.com/mattn/go-sqlite3`: Driver SQLite
- `github.com/lib/pq`: Driver PostgreSQL (opcional)
- `github.com/stretchr/testify`: Assertions e mocks para testes
- `github.com/jmoiron/sqlx`: Helper para SQL

## Como Executar

### Desenvolvimento (SQLite)
```bash
# Instalar dependências
go mod download

# Rodar servidor
go run cmd/server/main.go
# Server rodando em http://localhost:8080
```

### Produção com Docker (PostgreSQL)

#### Dockerfile
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

#### docker-compose.yml
```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: mundoinvest
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  server:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: mundoinvest
      PORT: 8080
    restart: unless-stopped

volumes:
  postgres_data:
```

#### Comandos Docker
```bash
# Subir todos os serviços
docker-compose up -d

# Ver logs
docker-compose logs -f

# Parar serviços
docker-compose down

# Parar serviços e remover volumes
docker-compose down -v

# Reconstruir imagem
docker-compose up -d --build

# Executar comandos no container
docker-compose exec server sh
```

### Testes
```bash
# Rodar todos os testes
go test ./...

# Rodar testes com coverage
go test -cover ./...

# Rodar testes de integração
go test -tags=integration ./...
```

## Próximos Passos
1. Implementar estrutura de pastas
2. Implementar modelos de dados
3. Implementar repository e conexão com banco
4. Implementar services e regras de negócio
5. Implementar controllers e endpoints
6. Implementar mutations GraphQL
7. Implementar testes de integração
8. Criar README com instruções
