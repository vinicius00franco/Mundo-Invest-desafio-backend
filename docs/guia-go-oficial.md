# Guia Oficial Go - Boas Práticas para Backend

> Documento compilado a partir da documentação oficial do Go (go.dev)

## Índice
- [Boas Práticas de Código](#boas-práticas-de-código)
- [Testes](#testes)
- [CI/CD](#cicd)
- [Libs Essenciais para Backend](#libs-essenciais-para-backend)
- [Segurança](#segurança)
- [Erros Comuns](#erros-comuns)

---

## Boas Práticas de Código

### Formatação
- **Use gofmt**: A ferramenta `gofmt` formata automaticamente o código Go no estilo padrão
- **go fmt**: Execute `go fmt` no nível de pacote para formatar todos os arquivos
- **Padrão universal**: Todo código Go nos pacotes padrão usa gofmt
- **Não perca tempo com alinhamento manual**: gofmt alinha comentários e campos de structs automaticamente

```bash
gofmt -w .
go fmt ./...
```

### Comentários
- **Comentários de linha são a norma**: Use `//` em vez de `/* */`
- **Comentários de documentação**: Comentários antes de declarações top-level documentam a declaração
- **Frases completas**: Comentários devem ser frases completas que começam com o nome do elemento descrito e terminam com ponto

```go
// Request representa uma requisição para executar um comando.
type Request struct {
    // ...
}

// Encode escreve a codificação JSON de req em w.
func Encode(w io.Writer, req *Request) {
    // ...
}
```

### Nomes
- **Visibilidade**: A primeira letra maiúscula torna o nome exportado (público)
- **Nomes de pacotes**: Devem ser curtos, minúsculos e descritivos
- **Getters**: Se você tem um campo chamado `owner`, o getter deve ser `Owner()`, não `GetOwner()`
- **Nomes de interfaces**: Interfaces com um método geralmente terminam em `-er` (ex: `Reader`, `Writer`)
- **MixedCaps**: Use MixedCaps para nomes exportados e mixedCaps para não exportados

### Context
- **Primeiro parâmetro**: Funções que usam Context devem aceitá-lo como primeiro parâmetro
- **Propagação explícita**: Context deve ser passado explicitamente por toda a cadeia de chamadas
- **Context.Background()**: Use apenas quando não há um contexto específico
- **Não adicione Context a structs**: Adicione parâmetro `ctx` aos métodos que precisam
- **Imutabilidade**: Context é imutável, pode ser passado para múltiplas chamadas

```go
func ProcessarCliente(ctx context.Context, clienteID string) error {
    // ...
}
```

### Tratamento de Erros
- **Retornar detalhes**: Use múltiplos valores de retorno para fornecer informações detalhadas de erro
- **Tipo error**: Por convenção, erros têm tipo `error`, uma interface simples
- **Strings de erro**: Devem identificar sua origem com prefixo (ex: "image: unknown format")
- **Type assertions**: Use type switch ou type assertion para extrair detalhes específicos

```go
type PathError struct {
    Op   string
    Path string
    Err  error
}

func (e *PathError) Error() string {
    return e.Op + " " + e.Path + ": " + e.Err.Error()
}
```

### Copying
- **Cuidado com structs de outros pacotes**: Copiar structs pode causar aliasing inesperado
- **Métodos de ponteiro**: Não copie valores do tipo T se seus métodos são associados a *T
- **Exemplo**: `bytes.Buffer` contém um `[]byte` slice; copiar pode causar aliasing

### Criptografia
- **NUNCA use math/rand para chaves**: Use `crypto/rand.Reader` para gerar chaves
- **Entropia insuficiente**: `math/rand` seeded com `Time.Nanoseconds()` tem poucos bits de entropia
- **Para texto**: Use `crypto/rand.Text` ou encode bytes com `encoding/hex` ou `encoding/base64`

```go
import (
    "crypto/rand"
    "encoding/hex"
)

func GerarChave() (string, error) {
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}
```

---

## Testes

### Estrutura de Testes
- **Arquivos de teste**: Terminam com `_test.go`
- **Funções de teste**: Nome começa com `Test` (ex: `TestHelloName`)
- **Parâmetro**: Funções de teste recebem `*testing.T`
- **Mesmo pacote**: Implemente testes no mesmo pacote do código testado

```go
package cliente

import (
    "testing"
)

func TestObterCliente(t *testing.T) {
    cliente, err := ObterCliente("123")
    if err != nil {
        t.Errorf("ObterCliente() = %v, want nil", err)
    }
    if cliente.Nome != "João" {
        t.Errorf("Nome = %v, want João", cliente.Nome)
    }
}
```

### Executando Testes
- **go test**: Executa funções de teste no pacote atual
- **-v**: Output verbose lista todos os testes e resultados
- **-race**: Detecta race conditions durante os testes

```bash
go test
go test -v
go test -race
go test ./...
```

### Boas Práticas de Testes
- **Testar erros**: Crie testes para casos de erro (ex: entrada vazia)
- **Mensagens úteis**: Use `t.Errorf` com mensagens descritivas
- **Subtestes**: Use `t.Run()` para organizar testes relacionados
- **Table-driven tests**: Use tabelas para testar múltiplos casos

### Fuzzing
- **Teste automatizado**: Usa coverage guidance para manipular inputs aleatórios
- **Edge cases**: Encontra casos extremos que programadores podem perder
- **Vulnerabilidades**: Pode descobrir SQL injection, buffer overflows, DoS, XSS
- **Tutorial**: [Go Fuzzing Tutorial](https://go.dev/doc/tutorial/fuzz)

```go
func FuzzProcessarWebhook(f *testing.F) {
    f.Add([]byte(`{"id": "123"}`))
    f.Fuzz(func(t *testing.T, data []byte) {
        ProcessarWebhook(data)
    })
}
```

---

## CI/CD

### Integração Contínua
- **Plataformas suportadas**: Go é bem suportado pela maioria de frameworks CI/CD
- **Govulncheck**: Ferramenta oficial para escanear vulnerabilidades
- **GitHub Action**: Action oficial disponível no GitHub Marketplace
- **Flag -json**: Govulncheck suporta -json para integração com outros sistemas

### Workflow Básico
```yaml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go test -race ./...
      - run: go vet ./...
      - uses: golang/govulncheck-action@v1
```

### Comandos Essenciais
- **go test**: Executa testes
- **go vet**: Examina construtos suspeitos
- **go build**: Compila o projeto
- **go mod tidy**: Limpa dependências

---

## Libs Essenciais para Backend

### Web Frameworks
- **net/http**: Pacote HTTP da biblioteca padrão
- **Gin**: Framework web com API estilo martini
- **Echo**: Framework de alta performance, extensível e minimalista
- **Gorilla**: Toolkit web para Go
- **Chi**: Router leve, idiomático e composicional

### Routers
- **net/http**: Standard library HTTP package
- **julienschmidt/httprouter**: Router HTTP leve e de alta performance
- **gorilla/mux**: Router HTTP poderoso e matcher de URL
- **go-chi/chi**: Router leve, idiomático e composicional

### Bancos de Dados
- **database/sql**: Interface padrão com suporte para MySQL, Postgres, Oracle, MS SQL, BigQuery
- **mongo-driver/mongo**: Driver oficial MongoDB para Go
- **elastic/go-elasticsearch**: Cliente Elasticsearch para Go
- **GORM**: ORM library para Go
- **Bleve**: Full-text search e indexing para Go
- **CockroachDB**: Banco de dados distribuído SQL

### Web Libraries
- **markbates/goth**: Autenticação para web apps
- **dgrijalva/jwt-go**: Implementação Go de JSON Web Tokens
- **html/template**: Engine de template HTML da biblioteca padrão

### Template Engines
- **html/template**: Engine de template HTML da biblioteca padrão
- **flosch/pongo2**: Linguagem de template com sintaxe Django

### Operações com Banco de Dados
- **Query**: Para queries que retornam múltiplas linhas
- **QueryRow**: Para queries que retornam uma única linha (mais eficiente)
- **Exec**: Para INSERT, UPDATE, DELETE (não retorna dados)
- **sql.Tx**: Para transações (commit/rollback)
- **context.Context**: Para cancelamento de queries

```go
// Query única
row := db.QueryRowContext(ctx, "SELECT nome FROM clientes WHERE id = ?", id)
var nome string
err := row.Scan(&nome)

// Transação
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback()

_, err = tx.Exec("INSERT INTO clientes (nome) VALUES (?)", nome)
if err != nil {
    return err
}

err = tx.Commit()
```

---

## Segurança

### Scan de Vulnerabilidades
- **govulncheck**: Ferramenta oficial para escanhar código e binários
- **Database**: Backed by Go vulnerability database
- **CI/CD**: Pode ser integrado em workflows CI/CD
- **IDE**: Disponível via Go extension para VS Code

### Manter Go e Dependências Atualizadas
- **Versão Go**: Mantenha atualizada para acessar features, performance e patches de segurança
- **Dependências**: Atualize dependências de terceiros regularmente
- **Revisão**: Cada atualização deve ser cuidadosamente revisada e testada
- **Risco**: Atualizar sem revisão pode introduzir bugs ou código malicioso

### Race Detector
- **Race conditions**: Ocorrem quando goroutines acessam mesmo recurso concorrentemente
- **-race flag**: Adicione ao rodar testes ou build
- **Runtime**: Detecta races que ocorrem em runtime
- **Limitação**: Não encontra races em code paths não executados

```bash
go test -race ./...
go build -race
```

### Go Vet
- **Análise estática**: Examina código para construtos suspeitos
- **Issues**: Código inalcançável, variáveis não usadas, erros comuns em goroutines
- **Qualidade**: Mantém qualidade de código, reduz tempo de debug
- **Execução**: `go vet ./...`

### Boas Práticas de Segurança
- **Validação de input**: Sempre valide e sanitize inputs
- **SQL Injection**: Use prepared statements ou ORM
- **XSS**: Escape de conteúdo HTML adequado
- **CSRF**: Implemente tokens CSRF para state-changing operations
- **HTTPS**: Sempre use HTTPS em produção
- **Secrets**: Nunca hardcode secrets; use environment variables ou secret managers

---

## Erros Comuns

### Referência a Variável de Loop
- **Problema**: Variável de loop é única, reutilizada em cada iteração
- **Go < 1.22**: Causa comportamento inesperado ao usar referências
- **Solução**: Copie variável para nova variável no loop
- **Go >= 1.22**: Variáveis são scoped à iteração (fixado)

```go
// ERRADO (Go < 1.22)
var out []*int
for i := 0; i < 3; i++ {
    out = append(out, &i) // Mesmo endereço
}

// CORRETO
var out []*int
for i := 0; i < 3; i++ {
    i := i // Copia para nova variável
    out = append(out, &i)
}
```

### Goroutines em Loop
- **Problema**: Mesmo problema de referência a variável de loop
- **Solução**: Passar valor como parâmetro para goroutine
- **Go >= 1.22**: Fixado automaticamente

```go
// ERRADO (Go < 1.22)
for _, cliente := range clientes {
    go func() {
        processarCliente(cliente) // cliente pode mudar
    }()
}

// CORRETO
for _, cliente := range clientes {
    go func(c Cliente) {
        processarCliente(c)
    }(cliente)
}
```

### Outros Erros Comuns
- **Não tratar erros**: Sempre check erros, não ignore
- **Panic em libraries**: Não use panic em código de biblioteca
- **Context não propagado**: Esquecer de passar Context pela cadeia
- **Resource leaks**: Não fechar connections, files, goroutines
- **Nil pointer dereference**: Sempre check nil antes de usar
- **Interface nil**: Interface com valor nil mas tipo não-nil

### Boas Práticas para Evitar Erros
- **Sempre trate erros**: Nunca ignore `_` sem razão
- **Use defer**: Para cleanup (close files, unlock mutexes)
- **Race detector**: Execute testes com -race regularmente
- **Go vet**: Execute go vet antes de commits
- **Code review**: Peer review ajuda a pegar erros comuns

---

## Referências Oficiais

- [Effective Go](https://go.dev/doc/effective_go)
- [Security Best Practices](https://go.dev/doc/security/best-practices)
- [Common Mistakes](https://go.dev/wiki/CommonMistakes)
- [Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Testing Tutorial](https://go.dev/doc/tutorial/add-a-test)
- [Database Access](https://go.dev/doc/database/)
- [Web Development](https://go.dev/solutions/webdev)
- [Go Packages](https://pkg.go.dev/)
