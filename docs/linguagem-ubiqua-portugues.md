# Linguagem Ubíqua - Mundo Invest

## Análise de Inconsistências de Idioma

### Problemas Identificados

O projeto atual apresenta uma **mistura inconsistente** de português e inglês em:

1. **Nomes de arquivos**: Maioria em português, alguns termos em inglês
2. **Nomes de pacotes**: 100% português (correto)
3. **Nomes de structs/interfaces**: Mix de português e inglês
4. **Nomes de funções/métodos**: Mix de português e inglês
5. **Nomes de variáveis**: Mix de português e inglês
6. **Mensagens de erro**: 100% português (correto)

### Mapeamento Atual

| Categoria | Português | Inglês | Status |
|-----------|----------|--------|--------|
| **Pacotes** | `gestao_clientes`, `processamento_eventos`, `dominio` | - | ✅ 100% PT |
| **Arquivos** | `service.go`, `controller.go`, `repository.go` | `client.go`, `mutations.go` | ⚠️ 85% PT |
| **Interfaces** | `ClienteService`, `EventoRepository` | `WebhookService`, `PipefyIntegrationService` | ⚠️ 70% PT |
| **Métodos** | `CriarCliente`, `ProcessarWebhook` | `EstruturarMutationCreateCard` | ⚠️ 60% PT |
| **Variáveis** | `nome`, `email`, `valorPatrimonio` | `pipeID`, `cardID` | ⚠️ 50% PT |
| **Mensagens** | "erro de validação", "Cliente criado com sucesso" | - | ✅ 100% PT |

## Linguagem Ubíqua Proposta (Português)

### Conceitos do Domínio

#### Entidades Principais
- **Cliente**: Pessoa física ou jurídica que solicita serviços
- **Evento**: Notificação do Pipefy sobre atualização de card
- **Solicitação**: Pedido de serviço do cliente (abertura conta, atualização cadastral, etc.)
- **Patrimônio**: Valor total dos bens do cliente
- **Prioridade**: Nível de urgência no atendimento (alta, normal)
- **Card**: Representação visual no Pipefy

#### Ações do Domínio
- **Criar**: Novo registro (cliente, evento)
- **Processar**: Tratar webhook/evento
- **Atualizar**: Modificar dados existentes
- **Calcular**: Determinar prioridade
- **Validar**: Verificar conformidade
- **Salvar**: Persistir no banco
- **Buscar**: Recuperar dados

#### Atributos
- **Identificador Interno**: ID gerado pelo sistema
- **Identificador Externo**: ID do sistema externo (Pipefy)
- **Nome**: Nome completo do cliente
- **Email**: Endereço eletrônico
- **Tipo de Solicitação**: Categoria do serviço
- **Valor do Patrimônio**: Montante financeiro
- **Nível de Prioridade**: Classificação de urgência
- **Status**: Estado atual do processo
- **Data de Criação**: Momento do registro
- **Data de Atualização**: Momento da última modificação

### Termos Técnicos (Mantidos em Inglês)

Alguns termos técnicos são mantidos em inglês por serem padrão da indústria:

- **Webhook**: Mecanismo de callback HTTP
- **Card**: Entidade específica do Pipefy
- **Pipe**: Pipeline do Pipefy
- **Mutation**: Operação GraphQL
- **Handler**: Manipulador de requisição HTTP
- **Repository**: Padrão de persistência
- **Service**: Camada de lógica de negócio
- **Controller**: Camada de apresentação
- **DTO**: Data Transfer Object
- **Strategy**: Padrão de projeto

## Proposta de Padronização

### Estratégia de Migração

#### Fase 1: Conceitos do Domínio (Alta Prioridade)
Migrar termos centrais do negócio para português:

**Termos Atuais → Termos Propostos:**
- `WebhookService` → `ServicoWebhook` ou `ServicoEvento`
- `PipefyIntegrationService` → `ServicoIntegracaoPipefy`
- `WebhookRequest` → `RequisicaoEvento` ou `RequisicaoWebhook`
- `WebhookController` → `ControladorEvento` ou `ControladorWebhook`
- `CardClienteData` → `DadosCardCliente`
- `CriarCardCliente` → `CriarCardCliente` (mantido)
- `AtualizarCardPrioridade` → `AtualizarPrioridadeCard`

#### Fase 2: Métodos Técnicos (Média Prioridade)
Migrar métodos que não são termos técnicos padrão:

**Termos Atuais → Termos Propostos:**
- `EstruturarMutationCreateCard` → `EstruturarMutacaoCriarCard`
- `EstruturarMutationUpdateCard` → `EstruturarMutacaoAtualizarCard`
- `ExecutarMutation` → `ExecutarMutacao`
- `FieldID` → `IdentificadorCampo` (mantido como FieldID por ser Pipefy)

#### Fase 3: Variáveis de Integração (Baixa Prioridade)
Variáveis relacionadas a sistemas externos podem manter inglês:

**Mantidos em Inglês:**
- `pipeID` (ID do Pipefy)
- `cardID` (ID do card)
- `fieldID` (ID do campo no Pipefy)
- `apiToken` (token de API)
- `apiURL` (URL de API)

### Regras de Nomenclatura

#### 1. Entidades do Domínio (100% Português)
```go
// ✅ CORRETO
type Cliente struct {
    IdentificadorInterno int64
    IdentificadorExterno string
    Nome                 string
    Email                string
    ValorPatrimonio      float64
    TipoSolicitacao      string
    NivelPrioridade      string
    Status               string
    DataCriacao          time.Time
}

// ❅ INCORRETO
type Client struct {
    InternalID int64
    ExternalID string
    Name string
    Email string
    PatrimonyValue float64
}
```

#### 2. Interfaces de Serviço (100% Português)
```go
// ✅ CORRETO
type ClienteService interface {
    CriarCliente(ctx context.Context, request CriarClienteRequest) (*Cliente, error)
    BuscarClientePorID(ctx context.Context, id int64) (*Cliente, error)
}

type ServicoEvento interface {
    ProcessarEvento(ctx context.Context, request RequisicaoEvento) error
}

// ❅ INCORRETO
type ClientService interface {
    CreateClient(ctx context.Context, request CreateClientRequest) (*Client, error)
}
```

#### 3. Interfaces de Repositório (100% Português)
```go
// ✅ CORRETO
type ClienteRepository interface {
    Salvar(ctx context.Context, cliente Cliente) (*Cliente, error)
    BuscarPorIdentificadorInterno(ctx context.Context, id int64) (*Cliente, error)
    Atualizar(ctx context.Context, cliente Cliente) error
}

type RepositorioEvento interface {
    Salvar(ctx context.Context, evento Evento) (*Evento, error)
    BuscarPorIdentificadorEvento(ctx context.Context, id string) (*Evento, error)
}

// ❅ INCORRETO
type ClientRepository interface {
    Save(ctx context.Context, client Client) (*Client, error)
    FindByID(ctx context.Context, id int64) (*Client, error)
}
```

#### 4. Métodos de Integração (Termos Técnicos Mantidos)
```go
// ✅ CORRETO - Termos técnicos mantidos
type ServicoIntegracaoPipefy interface {
    CriarCardCliente(ctx context.Context, pipeID string, dados DadosCardCliente) (string, error)
    AtualizarPrioridadeCard(ctx context.Context, cardID string, nivelPrioridade string) error
}

// ✅ CORRETO - Mutation mantido como termo técnico
func (m *Mutacoes) CriarMutacaoCard(pipeID string, nome, email, tipoSolicitacao string, valorPatrimonio float64) (string, error)

// ❅ INCORRETO - Mistura desnecessária
func (m *Mutations) CreateCardMutation(pipeID string, name, email, requestType string, patrimonyValue float64) (string, error)
```

#### 5. Variáveis de Sistemas Externos (Inglês Aceito)
```go
// ✅ CORRETO - Variáveis de sistemas externos
pipeID := "pipe_123"
cardID := "card_456"
fieldID := "field_789"
apiToken := "token_abc"
apiURL := "https://api.pipefy.com/graphql"

// ✅ CORRETO - Variáveis internas em português
identificadorInterno := cliente.IdentificadorInterno
nomeCliente := cliente.Nome
valorPatrimonio := cliente.ValorPatrimonio
```

### Dicionário de Tradução

| Inglês | Português | Contexto |
|--------|-----------|-----------|
| Client | Cliente | Entidade principal |
| Customer | Cliente | Entidade principal |
| User | Usuário | Sistema de autenticação |
| Event | Evento | Notificação do Pipefy |
| Webhook | Webhook | Termo técnico mantido |
| Card | Card | Termo técnico do Pipefy |
| Request | Requisição | Dados de entrada |
| Response | Resposta | Dados de saída |
| Service | Serviço | Camada de negócio |
| Repository | Repositório | Camada de persistência |
| Controller | Controlador | Camada de apresentação |
| Handler | Handler | Manipulador HTTP |
| Create | Criar | Ação de domínio |
| Update | Atualizar | Ação de domínio |
| Delete | Excluir | Ação de domínio |
| Save | Salvar | Ação de persistência |
| Find | Buscar | Ação de consulta |
| Validate | Validar | Ação de verificação |
| Calculate | Calcular | Ação de processamento |
| Process | Processar | Ação de tratamento |
| ID | Identificador | Atributo |
| Name | Nome | Atributo |
| Email | Email | Atributo (mantido) |
| Value | Valor | Atributo |
| Type | Tipo | Atributo |
| Status | Status | Atributo (mantido) |
| Priority | Prioridade | Atributo |
| Level | Nível | Atributo |
| Patrimony | Patrimônio | Atributo |
| Creation Date | Data de Criação | Atributo |
| Update Date | Data de Atualização | Atributo |

## Exemplos de Código Padronizado

### Estrutura de Arquivos
```
internal/
├── gestao_clientes/                    # ✅ Pacote em português
│   ├── controlador.go                 # ✅ Arquivo em português
│   ├── servico.go                     # ✅ Arquivo em português
│   ├── repositorio.go                 # ✅ Arquivo em português
│   ├── validador.go                  # ✅ Arquivo em português
│   ├── dto.go                        # ✅ DTO mantido (termo técnico)
│   └── modelo_cliente.go             # ✅ Modelo em português
├── processamento_eventos/             # ✅ Pacote em português
│   ├── controlador.go                 # ✅ Arquivo em português
│   ├── servico.go                     # ✅ Arquivo em português
│   ├── repositorio.go                 # ✅ Arquivo em português
│   └── modelo_evento.go              # ✅ Modelo em português
├── integracao_pipefy/                # ✅ Pacote em português
│   ├── cliente_pipefy.go             # ✅ Arquivo em português
│   ├── servico.go                     # ✅ Arquivo em português
│   └── mutacoes.go                   # ✅ Arquivo em português
└── dominio/                          # ✅ Pacote em português
    ├── calculadora_prioridade.go      # ✅ Arquivo em português
    ├── eventos.go                    # ✅ Arquivo em português
    └── constantes.go                 # ✅ Arquivo em português
```

### Código Go Padronizado
```go
package gestao_clientes

// ✅ CORRETO - Interface em português
type ServicoCliente interface {
    CriarCliente(ctx context.Context, requisicao RequisicaoCriarCliente) (*Cliente, error)
    BuscarClientePorIdentificador(ctx context.Context, identificador int64) (*Cliente, error)
    AtualizarCliente(ctx context.Context, cliente Cliente) error
}

// ✅ CORRETO - Struct em português
type Cliente struct {
    IdentificadorInterno int64
    IdentificadorExterno string
    Nome                 string
    Email                string
    ValorPatrimonio      float64
    TipoSolicitacao      string
    NivelPrioridade      string
    Status               string
    DataCriacao          time.Time
    DataAtualizacao      time.Time
}

// ✅ CORRETO - Método em português
func (s *servicoCliente) CriarCliente(ctx context.Context, requisicao RequisicaoCriarCliente) (*Cliente, error) {
    // Validação
    if err := ValidarRequisicaoCriarCliente(requisicao); err != nil {
        return nil, err
    }
    
    // Lógica de negócio
    cliente := &Cliente{
        Nome:            requisicao.Nome,
        Email:           requisicao.Email,
        ValorPatrimonio: requisicao.ValorPatrimonio,
        TipoSolicitacao: requisicao.TipoSolicitacao,
        NivelPrioridade: s.calculadoraPrioridade.CalcularNivelPrioridade(requisicao.ValorPatrimonio),
        Status:          dominio.StatusAguardandoAnalise,
        DataCriacao:      time.Now(),
    }
    
    // Persistência
    clienteSalvo, err := s.repositorio.Salvar(ctx, *cliente)
    if err != nil {
        return nil, err
    }
    
    // Integração (termos técnicos mantidos)
    cardID, err := s.servicoPipefy.CriarCardCliente(ctx, s.pipeID, DadosCardCliente{
        Nome:            cliente.Nome,
        Email:           cliente.Email,
        ValorPatrimonio: cliente.ValorPatrimonio,
        TipoSolicitacao: cliente.TipoSolicitacao,
    })
    if err != nil {
        return nil, err
    }
    
    clienteSalvo.IdentificadorExterno = cardID
    return clienteSalvo, nil
}
```

## Plano de Implementação

### Fase 1: Conceitos do Domínio (Semanas 1-2)
1. Renomear interfaces de serviço para português
2. Renomear structs de domínio para português
3. Atualizar métodos de serviço para português
4. Atualizar testes unitários

### Fase 2: Camada de Persistência (Semana 3)
1. Renomear interfaces de repositório para português
2. Renomear métodos de repositório para português
3. Atualizar implementações
4. Atualizar testes de integração

### Fase 3: Camada de Apresentação (Semana 4)
1. Renomear controllers para português
2. Renomear handlers para português
3. Atualizar rotas HTTP
4. Atualizar testes de controllers

### Fase 4: Integrações (Semana 5)
1. Renomear serviços de integração para português
2. Manter termos técnicos (webhook, card, mutation)
3. Atualizar clientes de integração
4. Atualizar testes de integração

### Fase 5: Documentação (Semana 6)
1. Atualizar README
2. Atualizar documentação de API
3. Atualizar exemplos de código
4. Criar guia de estilo

## Benefícios da Padronização

1. **Consistência**: Código mais uniforme e fácil de entender
2. **Alinhamento com Negócio**: Linguagem reflete o domínio do negócio
3. **Manutenibilidade**: Novos desenvolvedores entendem mais rápido
4. **Colaboração**: Equipe brasileira trabalha mais naturalmente
5. **Documentação**: Menos necessidade de tradução mental

## Conclusão

A adoção da linguagem ubíqua em português alinha o código com o domínio do negócio do Mundo Invest, mantendo termos técnicos padrão da indústria quando apropriado. Esta padronização melhorará a manutenibilidade e colaboração da equipe.