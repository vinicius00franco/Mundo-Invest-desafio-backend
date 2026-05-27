# Plano de Centralização e Organização de Mensagens (Revisado)

## Objetivo
Centralizar e organizar todas as mensagens retornadas pelos endpoints da API para facilitar manutenção, consistência e internacionalização futura, seguindo os princípios de Domain-Driven Design com linguagem ubíqua em português e as melhores práticas de Go.

## Análise Atual

### Problemas Identificados
1. **Mensagens dispersas**: Mensagens de erro e sucesso estão distribuídas por múltiplos arquivos
2. **Inconsistência**: Formatos diferentes para mensagens similares
3. **Dificuldade de manutenção**: Alterações requerem busca em múltiplos arquivos
4. **Sem suporte a internacionalização**: Mensagens em português "hardcoded"
5. **Mistura de níveis**: Mensagens técnicas misturadas com mensagens de usuário
6. **Linguagem não ubíqua**: Nomes em inglês misturados com português
7. **Problemas de performance**: Instâncias criadas desnecessariamente
8. **Falta de type safety**: Uso de "magic strings"

### Inventário de Mensagens Atuais

#### Mensagens de Sucesso
- `internal/gestao_clientes/dto.go`: "Cliente criado com sucesso"
- `internal/processamento_eventos/controller.go`: "Webhook processado com sucesso"

#### Mensagens de Validação
- `internal/gestao_clientes/validator.go`: "nome é obrigatório", "email é obrigatório", "tipoSolicitacao é obrigatório"
- `internal/gestao_clientes/validation_strategy.go`: "nome é obrigatório", "email é obrigatório", "tipoSolicitacao é obrigatório"
- `internal/processamento_eventos/validator.go`: "identificadorEvento é obrigatório", "identificadorCard é obrigatório", "emailCliente é obrigatório", "dataEvento é obrigatório"
- `internal/integracao_pipefy/cliente_pipefy.go`: "pipe_id é obrigatório", "fields_attributes é obrigatório", "card_id é obrigatório"

#### Mensagens de Erro de Serviço
- `internal/gestao_clientes/service.go`: "erro ao salvar cliente", "erro ao criar card no Pipefy", "erro ao atualizar cliente com card ID"
- `internal/processamento_eventos/service.go`: "erro ao verificar idempotência", "erro ao atualizar cliente", "erro ao parsear data do evento", "erro ao salvar evento"
- `internal/integracao_pipefy/service.go`: "erro ao estruturar mutation createCard", "erro ao executar mutation createCard", "erro ao estruturar mutation updateCard", "erro ao executar mutation updateCard"

#### Mensagens de Erro HTTP
- `internal/gestao_clientes/controller.go`: "Método não permitido", "Erro ao parsear JSON", "Erro ao gerar resposta"
- `internal/processamento_eventos/controller.go`: "Método não permitido", "Cliente não encontrado", "Erro ao parsear JSON", "Erro ao processar webhook", "Erro ao gerar resposta"

#### Mensagens de Erro de Banco de Dados
- `internal/shared/database/connection.go`: "erro ao abrir conexão com banco de dados", "erro ao testar conexão com banco de dados"
- `internal/shared/database/transacao.go`: "erro ao iniciar transação", "erro ao fazer commit da transação", "erro ao fazer rollback da transação"

## Proposta de Solução (Revisada)

### Estrutura de Centralização

#### 1. Criar Pacote `internal/shared/mensagens`
```
internal/shared/mensagens/
├── catalogo.go           # Catálogo principal com singleton
├── catalogo_test.go      # Testes do catálogo
├── tipos.go              # Tipos e interfaces
├── chaves.go             # Chaves type-safe
└── mapeador_http.go      # Mapeamento separado de HTTP
```

#### 2. Tipos e Interfaces (Type Safety e Testabilidade)

```go
package mensagens

import (
    "fmt"
    "sync"
)

// Chave type-safe para evitar "magic strings"
type Chave string

// ProvedorMensagens interface para testabilidade
type ProvedorMensagens interface {
    Texto(chave Chave) string
    TextoFormatado(chave Chave, args ...interface{}) string
}

// TipoMensagem define o tipo da mensagem no contexto do domínio
type TipoMensagem string

const (
    TipoSucesso    TipoMensagem = "sucesso"
    TipoValidacao  TipoMensagem = "validacao"
    TipoServico    TipoMensagem = "servico"
    TipoIntegracao TipoMensagem = "integracao"
    TipoHTTP       TipoMensagem = "http"
    TipoBancoDados TipoMensagem = "banco_dados"
)

// ContextoMensagem define o contexto de negócio da mensagem
type ContextoMensagem string

const (
    ContextoGestaoClientes      ContextoMensagem = "gestao_clientes"
    ContextoProcessamentoEventos ContextoMensagem = "processamento_eventos"
    ContextoIntegracaoPipefy    ContextoMensagem = "integracao_pipefy"
    ContextoBancoDados          ContextoMensagem = "banco_dados"
    ContextoGenerico            ContextoMensagem = "generico"
)
```

#### 3. Estrutura Simplificada de Mensagens

```go
// Mensagem representa uma mensagem simplificada e focada
type Mensagem struct {
    Texto string
}

// Catalogo contém todas as mensagens centralizadas
type Catalogo struct {
    mensagens map[Chave]Mensagem
    mu        sync.RWMutex // Para thread-safety em extensões futuras
}

// Singleton thread-safe
var (
    catalogoSingleton *Catalogo
    catalogoOnce     sync.Once
)

// ObterCatalogo retorna o singleton do catálogo
func ObterCatalogo() *Catalogo {
    catalogoOnce.Do(func() {
        catalogoSingleton = NovoCatalogo()
    })
    return catalogoSingleton
}

// NovoCatalogo cria uma nova instância do catálogo (para testes)
func NovoCatalogo() *Catalogo {
    return &Catalogo{
        mensagens: make(map[Chave]Mensagem, 50), // Pre-alocar capacidade
    }
}
```

#### 4. Chaves Type-Safe com Nomes Curtos

```go
package mensagens

// Chaves de mensagens - nomes curtos mas descritivos
const (
    // Sucesso
    CliCriadoSucesso    Chave = "cliente.criado.sucesso"
    WebhookProcessado   Chave = "webhook.processado.sucesso"
    
    // Validação - Cliente
    ValNomeObrigatorio          Chave = "validacao.nome.obrigatorio"
    ValEmailObrigatorio         Chave = "validacao.email.obrigatorio"
    ValTipoSolicitacaoObrigatorio Chave = "validacao.tipo_solicitacao.obrigatorio"
    ValValorPatrimonioObrigatorio Chave = "validacao.valor_patrimonio.obrigatorio"
    ValNomeMinimoCaracteres     Chave = "validacao.nome.minimo_caracteres"
    ValEmailInvalido            Chave = "validacao.email.invalido"
    ValValorPatrimonioNegativo  Chave = "validacao.valor_patrimonio.negativo"
    ValValorPatrimonioZero      Chave = "validacao.valor_patrimonio.zero"
    
    // Validação - Webhook
    ValIdEventoObrigatorio      Chave = "validacao.id_evento.obrigatorio"
    ValIdCardObrigatorio        Chave = "validacao.id_card.obrigatorio"
    ValEmailClienteObrigatorio  Chave = "validacao.email_cliente.obrigatorio"
    ValDataEventoObrigatorio    Chave = "validacao.data_evento.obrigatorio"
    ValEmailClienteInvalido     Chave = "validacao.email_cliente.invalido"
    ValDataEventoInvalido       Chave = "validacao.data_evento.invalido"
    
    // Erros - Serviço
    ErrSalvarCliente            Chave = "erro.cliente.salvar"
    ErrAtualizarCliente         Chave = "erro.cliente.atualizar"
    ErrVerificarIdempotencia    Chave = "erro.webhook.verificar_idempotencia"
    ErrProcessarDataEvento      Chave = "erro.webhook.processar_data_evento"
    ErrSalvarEvento             Chave = "erro.webhook.salvar_evento"
    
    // Erros - Integração
    ErrCriarCardPipefy          Chave = "erro.pipefy.criar_card"
    ErrAtualizarCardPipefy      Chave = "erro.pipefy.atualizar_card"
    ErrEstruturarMutacaoCriar   Chave = "erro.pipefy.estruturar_mutacao_criar"
    ErrEstruturarMutacaoAtualizar Chave = "erro.pipefy.estruturar_mutacao_atualizar"
    ErrExecutarMutacaoCriar     Chave = "erro.pipefy.executar_mutacao_criar"
    ErrExecutarMutacaoAtualizar Chave = "erro.pipefy.executar_mutacao_atualizar"
    
    // Erros - HTTP
    ErrMetodoNaoPermitido       Chave = "erro.http.metodo_nao_permitido"
    ErrProcessarJSON            Chave = "erro.http.processar_json"
    ErrGerarResposta            Chave = "erro.http.gerar_resposta"
    ErrClienteNaoEncontrado     Chave = "erro.http.cliente_nao_encontrado"
    ErrProcessarWebhook         Chave = "erro.http.processar_webhook"
    
    // Erros - Banco de Dados
    ErrAbrirConexao             Chave = "erro.banco_dados.abrir_conexao"
    ErrTestarConexao            Chave = "erro.banco_dados.testar_conexao"
    ErrIniciarTransacao         Chave = "erro.banco_dados.iniciar_transacao"
    ErrConfirmarTransacao       Chave = "erro.banco_dados.confirmar_transacao"
    ErrReverterTransacao        Chave = "erro.banco_dados.reverter_transacao"
)
```

#### 5. Inicialização do Catálogo com Singleton

```go
// inicializarMensagens popula o catálogo com todas as mensagens
func (c *Catalogo) inicializarMensagens() {
    // Mensagens de Sucesso
    c.mensagens[CliCriadoSucesso] = Mensagem{
        Texto: "Cliente criado com sucesso",
    }
    
    c.mensagens[WebhookProcessado] = Mensagem{
        Texto: "Webhook processado com sucesso",
    }
    
    // Mensagens de Validação
    c.mensagens[ValNomeObrigatorio] = Mensagem{
        Texto: "nome é obrigatório",
    }
    
    c.mensagens[ValEmailObrigatorio] = Mensagem{
        Texto: "email é obrigatório",
    }
    
    c.mensagens[ValTipoSolicitacaoObrigatorio] = Mensagem{
        Texto: "tipoSolicitacao é obrigatório",
    }
    
    c.mensagens[ValValorPatrimonioObrigatorio] = Mensagem{
        Texto: "valorPatrimonio é obrigatório",
    }
    
    c.mensagens[ValNomeMinimoCaracteres] = Mensagem{
        Texto: "nome deve ter pelo menos 3 caracteres",
    }
    
    c.mensagens[ValEmailInvalido] = Mensagem{
        Texto: "email inválido",
    }
    
    c.mensagens[ValValorPatrimonioNegativo] = Mensagem{
        Texto: "valorPatrimonio deve ser positivo",
    }
    
    c.mensagens[ValValorPatrimonioZero] = Mensagem{
        Texto: "valorPatrimonio deve ser maior que zero",
    }
    
    // Validação Webhook
    c.mensagens[ValIdEventoObrigatorio] = Mensagem{
        Texto: "identificadorEvento é obrigatório",
    }
    
    c.mensagens[ValIdCardObrigatorio] = Mensagem{
        Texto: "identificadorCard é obrigatório",
    }
    
    c.mensagens[ValEmailClienteObrigatorio] = Mensagem{
        Texto: "emailCliente é obrigatório",
    }
    
    c.mensagens[ValDataEventoObrigatorio] = Mensagem{
        Texto: "dataEvento é obrigatório",
    }
    
    c.mensagens[ValEmailClienteInvalido] = Mensagem{
        Texto: "emailCliente inválido",
    }
    
    c.mensagens[ValDataEventoInvalido] = Mensagem{
        Texto: "dataEvento inválido",
    }
    
    // Erros de Serviço
    c.mensagens[ErrSalvarCliente] = Mensagem{
        Texto: "erro ao salvar cliente",
    }
    
    c.mensagens[ErrAtualizarCliente] = Mensagem{
        Texto: "erro ao atualizar cliente",
    }
    
    c.mensagens[ErrVerificarIdempotencia] = Mensagem{
        Texto: "erro ao verificar idempotência",
    }
    
    c.mensagens[ErrProcessarDataEvento] = Mensagem{
        Texto: "erro ao processar data do evento",
    }
    
    c.mensagens[ErrSalvarEvento] = Mensagem{
        Texto: "erro ao salvar evento",
    }
    
    // Erros de Integração
    c.mensagens[ErrCriarCardPipefy] = Mensagem{
        Texto: "erro ao criar card no Pipefy",
    }
    
    c.mensagens[ErrAtualizarCardPipefy] = Mensagem{
        Texto: "erro ao atualizar card no Pipefy",
    }
    
    c.mensagens[ErrEstruturarMutacaoCriar] = Mensagem{
        Texto: "erro ao estruturar mutation createCard",
    }
    
    c.mensagens[ErrEstruturarMutacaoAtualizar] = Mensagem{
        Texto: "erro ao estruturar mutation updateCard",
    }
    
    c.mensagens[ErrExecutarMutacaoCriar] = Mensagem{
        Texto: "erro ao executar mutation createCard",
    }
    
    c.mensagens[ErrExecutarMutacaoAtualizar] = Mensagem{
        Texto: "erro ao executar mutation updateCard",
    }
    
    // Erros HTTP
    c.mensagens[ErrMetodoNaoPermitido] = Mensagem{
        Texto: "Método não permitido",
    }
    
    c.mensagens[ErrProcessarJSON] = Mensagem{
        Texto: "Erro ao processar JSON",
    }
    
    c.mensagens[ErrGerarResposta] = Mensagem{
        Texto: "Erro ao gerar resposta",
    }
    
    c.mensagens[ErrClienteNaoEncontrado] = Mensagem{
        Texto: "Cliente não encontrado",
    }
    
    c.mensagens[ErrProcessarWebhook] = Mensagem{
        Texto: "Erro ao processar webhook",
    }
    
    // Erros de Banco de Dados
    c.mensagens[ErrAbrirConexao] = Mensagem{
        Texto: "erro ao abrir conexão com banco de dados",
    }
    
    c.mensagens[ErrTestarConexao] = Mensagem{
        Texto: "erro ao testar conexão com banco de dados",
    }
    
    c.mensagens[ErrIniciarTransacao] = Mensagem{
        Texto: "erro ao iniciar transação",
    }
    
    c.mensagens[ErrConfirmarTransacao] = Mensagem{
        Texto: "erro ao confirmar transação",
    }
    
    c.mensagens[ErrReverterTransacao] = Mensagem{
        Texto: "erro ao reverter transação",
    }
}

// Texto retorna o texto de uma mensagem pela chave
func (c *Catalogo) Texto(chave Chave) string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if msg, exists := c.mensagens[chave]; exists {
        return msg.Texto
    }
    return fmt.Sprintf("mensagem não encontrada: %s", chave)
}

// TextoFormatado retorna o texto de uma mensagem com parâmetros substituídos
func (c *Catalogo) TextoFormatado(chave Chave, args ...interface{}) string {
    texto := c.Texto(chave)
    return fmt.Sprintf(texto, args...)
}
```

#### 6. Mapeador HTTP (Separação de Responsabilidades)

```go
package mensagens

import "net/http"

// MapeadorHTTP separa responsabilidade de mapeamento de status HTTP
type MapeadorHTTP struct {
    statusPorTipo map[TipoMensagem]int
}

// NovoMapeadorHTTP cria um novo mapeador HTTP
func NovoMapeadorHTTP() *MapeadorHTTP {
    return &MapeadorHTTP{
        statusPorTipo: map[TipoMensagem]int{
            TipoSucesso:    http.StatusCreated, // ou http.StatusOK dependendo do contexto
            TipoValidacao:  http.StatusBadRequest,
            TipoServico:    http.StatusInternalServerError,
            TipoIntegracao: http.StatusBadGateway,
            TipoHTTP:       http.StatusBadRequest,
            TipoBancoDados: http.StatusInternalServerError,
        },
    }
}

// StatusPara retorna o status HTTP para um tipo de mensagem
func (m *MapeadorHTTP) StatusPara(tipo TipoMensagem) int {
    if status, exists := m.statusPorTipo[tipo]; exists {
        return status
    }
    return http.StatusInternalServerError // Fallback seguro
}
```

## Plano de Implementação (Revisado)

### Fase 1: Estrutura Básica com Singleton (Dia 1)
1. Criar pacote `internal/shared/mensagens`
2. Implementar tipos e interfaces com type safety
3. Criar catálogo com singleton pattern (`sync.Once`)
4. Implementar chaves type-safe com nomes curtos
5. Criar mapeador HTTP separado
6. Escrever testes unitários do catálogo
7. Testes de performance do singleton
8. Documentar estrutura do catálogo

### Fase 2: Migração - Gestão de Clientes (Dia 2)
1. Atualizar `internal/gestao_clientes/validator.go` para usar singleton
2. Atualizar `internal/gestao_clientes/validation_strategy.go` para usar singleton
3. Atualizar `internal/gestao_clientes/service.go` para usar singleton
4. Atualizar `internal/gestao_clientes/controller.go` para usar singleton
5. Atualizar `internal/gestao_clientes/dto.go` para usar singleton
6. Atualizar testes do pacote para usar interface ProvedorMensagens
7. Validar build e testes unitários
8. Executar testes de integração do módulo
9. Verificar performance com benchmarks

### Fase 3: Migração - Processamento de Eventos (Dia 3)
1. Atualizar `internal/processamento_eventos/validator.go` para usar singleton
2. Atualizar `internal/processamento_eventos/service.go` para usar singleton
3. Atualizar `internal/processamento_eventos/controller.go` para usar singleton
4. Atualizar testes do pacote para usar interface ProvedorMensagens
5. Validar build e testes unitários
6. Executar testes de integração do módulo
7. Verificar performance com benchmarks

### Fase 4: Migração - Integração Pipefy (Dia 4)
1. Atualizar `internal/integracao_pipefy/cliente_pipefy.go` para usar singleton
2. Atualizar `internal/integracao_pipefy/service.go` para usar singleton
3. Atualizar testes do pacote para usar interface ProvedorMensagens
4. Validar build e testes unitários
5. Executar testes de integração do módulo
6. Verificar performance com benchmarks

### Fase 5: Migração - Infraestrutura Compartilhada (Dia 5)
1. Atualizar `internal/shared/database/connection.go` para usar singleton
2. Atualizar `internal/shared/database/transacao.go` para usar singleton
3. Atualizar `internal/shared/errors/errors.go` para usar singleton
4. Atualizar testes do pacote para usar interface ProvedorMensagens
5. Validar build e testes unitários
6. Executar testes de integração do módulo
7. Verificar performance com benchmarks

### Fase 6: Validação Final e Otimizações (Dia 6)
1. Executar todos os testes unitários
2. Executar todos os testes de integração
3. Executar benchmarks de performance
4. Verificar consistência de mensagens em todo o sistema
5. Documentar uso do catálogo de mensagens
6. Criar guia para adicionar novas mensagens
7. Revisar linguagem ubíqua nas mensagens
8. Otimizações finais baseadas em métricas

## Benefícios Esperados (Revisados)

1. **Manutenção simplificada**: Alterações de mensagens em um único local centralizado
2. **Consistência**: Formato padronizado para todas as mensagens do sistema
3. **Performance**: Singleton pattern evita alocações desnecessárias
4. **Type safety**: Chaves type-safe eliminam "magic strings"
5. **Testabilidade**: Interface ProvedorMensagens facilita mocking
6. **Internacionalização**: Estrutura preparada para suportar múltiplos idiomas
7. **Separação de responsabilidades**: Mensagens separadas de configuração HTTP
8. **Linguagem ubíqua**: Nomes descritivos em português em todo o sistema
9. **Organização por contexto**: Mensagens agrupadas por contexto de negócio
10. **Idiomático em Go**: Segue melhores práticas da linguagem

## Riscos e Mitigação (Revisados)

### Risco 1: Problemas de Performance com Singleton
- **Descrição**: Singleton pode se tornar bottleneck em alta concorrência
- **Mitigação**: Implementar benchmarking e considerar `sync.Map` se necessário
- **Plano de contingência**: Usar dependency injection se singleton não for adequado
- **Validação**: Testes de performance em cada fase

### Risco 2: Quebra de Testes Existentes
- **Descrição**: Migração pode causar falhas em testes que verificam mensagens específicas
- **Mitigação**: Migração faseada com validação após cada fase
- **Plano de contingência**: Commit separado para cada fase para fácil rollback
- **Validação**: Executar testes completos após cada migração

### Risco 3: Regressão em Funcionalidades
- **Descrição**: Alterações podem introduzir bugs em funcionalidades existentes
- **Mitigação**: Execução completa de testes após cada migração
- **Validação**: Testes de integração para verificar comportamento
- **Monitoramento**: Verificar logs de erro em ambiente de desenvolvimento

### Risco 4: Complexidade Adicional
- **Descrição**: Nova camada de abstração pode aumentar complexidade do código
- **Mitigação**: Documentação clara e exemplos de uso
- **Simplicidade**: API simples e intuitiva para uso do catálogo
- **Treinamento**: Guia de uso para desenvolvedores

### Risco 5: Inconsistência na Linguagem Ubíqua
- **Descrição**: Novas mensagens podem não seguir padrão de linguagem ubíqua
- **Mitigação**: Revisão de código focada em linguagem ubíqua
- **Padrões**: Documentação clara de padrões de nomenclatura
- **Validação**: Code review específico para consistência linguística

### Risco 6: Dificuldade em Testes com Singleton
- **Descrição**: Singleton global pode dificultar testes unitários
- **Mitigação**: Implementar interface ProvedorMensagens para mocking
- **Validação**: Criar mocks específicos para testes
- **Padrões**: Documentar padrões de testabilidade

## Cronograma

- **Dia 1**: Estrutura básica do pacote mensagens + testes unitários
- **Dia 2**: Migração gestão de clientes + validação
- **Dia 3**: Migração processamento de eventos + validação
- **Dia 4**: Migração integração Pipefy + validação
- **Dia 5**: Migração infraestrutura compartilhada + validação
- **Dia 6**: Validação final, documentação e revisão de linguagem ubíqua

## Critérios de Sucesso (Revisados)

1. ✅ Todas as mensagens centralizadas no catálogo de mensagens
2. ✅ Nenhuma mensagem "hardcoded" nos arquivos de negócio
3. ✅ Todos os nomes de tipos, funções e variáveis em português
4. ✅ Linguagem ubíqua consistente em todo o sistema
5. ✅ Singleton pattern implementado com `sync.Once`
6. ✅ Type safety com chaves personalizadas (sem "magic strings")
7. ✅ Interface ProvedorMensagens para testabilidade
8. ✅ Separação de responsabilidades (mensagens vs HTTP)
9. ✅ Todos os testes unitários passando
10. ✅ Todos os testes de integração passando
11. ✅ Benchmarks de performance satisfatórios
12. ✅ Documentação completa do catálogo de mensagens
13. ✅ Guia para adicionar novas mensagens seguindo linguagem ubíqua
14. ✅ Revisão de código focada em consistência linguística

## Exemplos de Uso (Revisados)

### Exemplo 1: Uso em Validador de Cliente com Singleton

```go
package gestao_clientes

import (
    "github.com/MundoInvest/backend/internal/shared/mensagens"
)

func ValidarRequisicaoCriarCliente(requisicao RequisicaoCriarCliente) error {
    catalogo := mensagens.ObterCatalogo() // ✅ Singleton thread-safe
    var erros []string
    
    if strings.TrimSpace(requisicao.Nome) == "" {
        erros = append(erros, catalogo.Texto(mensagens.ValNomeObrigatorio))
    }
    
    if strings.TrimSpace(requisicao.Email) == "" {
        erros = append(erros, catalogo.Texto(mensagens.ValEmailObrigatorio))
    }
    
    if len(erros) > 0 {
        return errors.NewValidationError("", strings.Join(erros, "; "))
    }
    
    return nil
}
```

### Exemplo 2: Uso em Controller com Singleton

```go
package gestao_clientes

import (
    "github.com/MundoInvest/backend/internal/shared/mensagens"
)

func (c *ClienteController) CriarClienteHandler(w http.ResponseWriter, r *http.Request) {
    catalogo := mensagens.ObterCatalogo() // ✅ Singleton thread-safe
    mapeadorHTTP := mensagens.NovoMapeadorHTTP()
    
    // Parsear JSON
    var requisicao RequisicaoCriarCliente
    if err := json.NewDecoder(r.Body).Decode(&requisicao); err != nil {
        mensagem := catalogo.Texto(mensagens.ErrProcessarJSON)
        http.Error(w, mensagem, mapeadorHTTP.StatusPara(mensagens.TipoHTTP))
        return
    }
    
    // Validar
    if err := ValidarRequisicaoCriarCliente(requisicao); err != nil {
        http.Error(w, err.Error(), mapeadorHTTP.StatusPara(mensagens.TipoValidacao))
        return
    }
    
    // Criar cliente
    cliente, err := c.service.CriarCliente(r.Context(), requisicao)
    if err != nil {
        http.Error(w, err.Error(), mapeadorHTTP.StatusPara(mensagens.TipoServico))
        return
    }
    
    // Retornar sucesso
    mensagemSucesso := catalogo.Texto(mensagens.CliCriadoSucesso)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    
    response := NewCriarClienteResponse(cliente)
    response.Mensagem = mensagemSucesso
    
    json.NewEncoder(w).Encode(response)
}
```

### Exemplo 3: Uso com Parâmetros Dinâmicos

```go
package gestao_clientes

import (
    "github.com/MundoInvest/backend/internal/shared/mensagens"
)

func (s *clienteService) BuscarClientePorEmail(ctx context.Context, email string) (*Cliente, error) {
    catalogo := mensagens.ObterCatalogo()
    
    cliente, err := s.repository.BuscarPorEmail(ctx, email)
    if err != nil {
        // Usar mensagem formatada com parâmetros
        mensagem := catalogo.TextoFormatado(mensagens.ErrClienteNaoEncontrado, email)
        return nil, errors.NewNotFoundError("cliente", mensagem)
    }
    
    return cliente, nil
}
```

### Exemplo 4: Uso em Testes com Mock

```go
package gestao_clientes_test

import (
    "github.com/MundoInvest/backend/internal/shared/mensagens"
)

// MockProvedorMensagens para testes
type MockProvedorMensagens struct {
    mensagens map[mensagens.Chave]string
}

func (m *MockProvedorMensagens) Texto(chave mensagens.Chave) string {
    return m.mensagens[chave]
}

func (m *MockProvedorMensagens) TextoFormatado(chave mensagens.Chave, args ...interface{}) string {
    return fmt.Sprintf(m.Texto(chave), args...)
}

func TestCriarClienteHandler(t *testing.T) {
    // Criar mock
    mock := &MockProvedorMensagens{
        mensagens: map[mensagens.Chave]string{
            mensagens.ValNomeObrigatorio: "nome é obrigatório",
        },
    }
    
    // Injetar mock no service
    service := NovoClienteServiceComProvedorMensagens(mock)
    
    // Testar...
}
```

## Padrões de Nomenclatura (Linguagem Ubíqua Revisada)

### Regras para Chaves de Mensagens
1. **Hierarquia**: Contexto.Tipo.Entidade.Ação
2. **Português**: Todos os termos em português
3. **camelCase**: Separar palavras com underscore
4. **Descritivo**: Nomes que descrevem claramente o contexto
5. **Concisos**: Evitar nomes excessivamente longos (>40 caracteres)

### Exemplos de Chaves (Revisados)
- ✅ `CliCriadoSucesso` - Correto (curto e descritivo)
- ✅ `ValNomeObrigatorio` - Correto (prefixo Val para validação)
- ✅ `ErrSalvarCliente` - Correto (prefixo Err para erros)
- ❌ `ChaveMensagemValidacaoIdentificadorEventoObrigatorio` - Incorreto (muito longo)
- ❌ `MsgClienteCriadoSucesso` - Incorreto (prefixo desnecessário)
- ❌ `cliente.created.success` - Incorreto (inglês)

### Prefixos Padronizados
- `Cli` - Mensagens relacionadas a clientes
- `Val` - Mensagens de validação
- `Err` - Mensagens de erro
- `Web` - Mensagens relacionadas a webhooks

### Regras para Nomes de Tipos e Funções
1. **Português**: Todos os nomes em português
2. **camelCase**: Variáveis e funções
3. **PascalCase**: Tipos e structs
4. **Descritivo**: Nomes que indicam propósito
5. **Concisos**: Evitar nomes excessivamente longos

### Exemplos de Nomes (Revisados)
- ✅ `Catalogo` - Correto (simples e direto)
- ✅ `Texto` - Correto (função clara)
- ✅ `Chave` - Correto (type customizado)
- ✅ `ProvedorMensagens` - Correto (interface descritiva)
- ❌ `CatalogoMensagens` - Redundante (pode ser apenas `Catalogo`)
- ❌ `ObterMensagem` - Pode ser simplificado para `Texto`
- ❌ `MessageCatalog` - Incorreto (inglês)

### Campos de Estruturas (Simplificados)
- ✅ `Texto` - Campo essencial e claro
- ✅ `Chave` - Type customizado para chaves
- ❌ `Chave` na struct - Redundante (já é chave do map)
- ❌ `Tipo` - Raramente usado em runtime
- ❌ `Contexto` - Raramente usado em runtime
- ❌ `StatusHTTP` - Mistura de responsabilidades

## Integração com Sistema Existente

### Uso do Singleton em Produção

```go
// Em qualquer parte do código
import "github.com/MundoInvest/backend/internal/shared/mensagens"

func AlgumaFuncao() {
    catalogo := mensagens.ObterCatalogo() // ✅ Thread-safe singleton
    mensagem := catalogo.Texto(mensagens.CliCriadoSucesso)
    // usar mensagem...
}
```

### Uso em Testes com Mock

```go
// Em testes
func TestAlgumaCoisa(t *testing.T) {
    mock := &MockProvedorMensagens{
        mensagens: map[mensagens.Chave]string{
            mensagens.CliCriadoSucesso: "Teste sucesso",
        },
    }
    
    // Injetar mock no sistema em teste
    service := NovoService(mock)
    
    // Testar...
}
```

## Impacto das Mudanças Técnicas

### Melhorias de Performance
- **Singleton Pattern**: Elimina alocação de memória repetitiva
- **Pre-alocação**: Map com capacidade conhecida reduz reallocation
- **Thread-safety**: `sync.RWMutex` para acesso concorrente seguro
- **Benchmarking**: Métricas de performance em cada fase

### Melhorias de Qualidade de Código
- **Type Safety**: Elimina "magic strings" com tipos personalizados
- **Testabilidade**: Interface `ProvedorMensagens` facilita mocking
- **Separação de Responsabilidades**: Mensagens separadas de configuração HTTP
- **Idiomático Go**: Segue melhores práticas da comunidade Go

### Manutenibilidade
- **Centralização**: Todas as mensagens em um único local
- **Consistência**: Formato padronizado em todo o sistema
- **Documentação**: Catálogo serve como documentação viva
- **Extensibilidade**: Fácil adicionar novas mensagens

## Resumo Executivo

Este plano revisado incorpora as recomendações de um especialista senior em Go, mantendo a excelência da linguagem ubíqua em português enquanto implementa as melhores práticas da linguagem:

1. **Performance**: Singleton pattern com `sync.Once` elimina overhead
2. **Type Safety**: Chaves personalizadas eliminam erros de compilação
3. **Testabilidade**: Interface para dependency injection facilita testes
4. **Simplicidade**: Estrutura simplificada seguindo filosofia Go
5. **Separação de Responsabilidades**: Mensagens separadas de HTTP
6. **Linguagem Ubíqua**: Nomes descritivos em português mantidos

O plano está pronto para implementação com uma abordagem faseada de 6 dias, validação contínua e critérios de sucesso claros.

## Próximos Passos

1. Aprovar plano revisado com melhorias técnicas
2. Iniciar implementação Fase 1 (estrutura básica com singleton)
3. Validar cada fase antes de prosseguir para a próxima
4. Documentar aprendizados durante implementação
5. Revisar continuamente a consistência da linguagem ubíqua
6. Atualizar documentação de arquitetura com novos padrões
7. Executar benchmarks de performance em cada fase
8. Coletar feedback da equipe e ajustar conforme necessário

1. Aprovar plano ajustado com linguagem ubíqua
2. Iniciar implementação Fase 1 (estrutura básica)
3. Validar cada fase antes de prosseguir para a próxima
4. Documentar aprendizados durante implementação
5. Revisar continuamente a consistência da linguagem ubíqua
6. Atualizar documentação de arquitetura com novos padrões

## Resumo dos Ajustes de Linguagem Ubíqua e Melhorias Técnicas

### Ajustes Baseados na Avaliação de Especialista Go

#### 1. Performance e Singleton Pattern
- **Antes**: `NovoCatalogoMensagens()` chamado a cada requisição
- **Depois**: `ObterCatalogo()` com singleton `sync.Once`
- **Benefício**: Elimina alocações desnecessárias e overhead

#### 2. Simplificação de Estrutura
- **Antes**: Struct com 6 campos (muitos redundantes)
- **Depois**: Struct simplificada com apenas `Texto`
- **Benefício**: Menos complexidade, mais idiomático em Go

#### 3. Type Safety
- **Antes**: Strings "magic" como chaves
- **Depois**: Tipo personalizado `Chave` com constantes type-safe
- **Benefício**: Compile-time checks, elimina erros de digitação

#### 4. Nomes de Constantes Otimizados
- **Antes**: `ChaveMensagemValidacaoIdentificadorEventoObrigatorio` (60+ caracteres)
- **Depois**: `ValIdEventoObrigatorio` (prefixos curtos)
- **Benefício**: Melhor legibilidade, linhas mais curtas

#### 5. Separação de Responsabilidades
- **Antes**: `StatusHTTP` dentro da struct de mensagem
- **Depois**: `MapeadorHTTP` separado
- **Benefício**: Separation of concerns, mais flexível

#### 6. Testabilidade
- **Antes**: Singleton global dificulta mocking
- **Depois**: Interface `ProvedorMensagens` para DI
- **Benefício**: Fácil criar mocks para testes

#### 7. Suporte a Parâmetros
- **Antes**: Não abordado no plano original
- **Depois**: `TextoFormatado()` com `fmt.Sprintf`
- **Benefício**: Suporte a mensagens dinâmicas

### Nomes de Pacotes e Arquivos (Mantidos)
- `internal/shared/mensagens` (já ajustado)
- `catalogo.go` (simplificado de `catalogo_mensagens.go`)
- `tipos.go` (simplificado de `tipos_mensagens.go`)
- `chaves.go` (novo arquivo para chaves type-safe)
- `mapeador_http.go` (novo arquivo para separação de responsabilidades)

### Nomes de Tipos (Simplificados)
- `CatalogoMensagens` → `Catalogo` (mais conciso)
- `ChaveMensagem` → `Chave` (mais direto)
- `ProvedorMensagens` (nova interface para testabilidade)

### Nomes de Constantes (Otimizados)
- `ChaveMensagemClienteCriadoComSucesso` → `CliCriadoSucesso`
- `ChaveMensagemValidacaoNomeObrigatorio` → `ValNomeObrigatorio`
- `ChaveMensagemErroSalvarCliente` → `ErrSalvarCliente`

### Nomes de Funções (Simplificados)
- `NovoCatalogoMensagens()` → `NovoCatalogo()`
- `ObterMensagem()` → `Texto()`
- `ObterMensagemComParametros()` → `TextoFormatado()`

### Valores de Constantes (Mantidos)
- `TipoMensagemSucesso` → `TipoSucesso` (mais curto)
- `ContextoGestaoClientes` (mantido por clareza)

### Campos de Estruturas (Simplificados)
- Removidos: `Chave`, `Tipo`, `Contexto`, `TextoPadrao`, `StatusHTTP`
- Mantido: `Texto` (essencial)
- Adicionado: Suporte a parâmetros via `TextoFormatado()`

### Terminologia Técnica (Mantida)
- `processar` em vez de `parsear` (já ajustado)
- `confirmar` em vez de `commit` (já ajustado)
- `reverter` em vez de `rollback` (já ajustado)
- `mutacao` (já padronizado)

Esses ajustes incorporam as recomendações da avaliação especializada em Go, mantendo a excelência da linguagem ubíqua em português enquanto segue as melhores práticas da linguagem para performance, type safety e testabilidade.
