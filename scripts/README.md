# Scripts de Teste API

Este diretório contém scripts bash para testar o backend usando curl, baseados nos cenários BDD definidos em `bdd/criacao-cliente/`.

## Estrutura

```
scripts/
├── exports/                    # Diretório onde os resultados são salvos
├── testar_cenarios_sucesso.sh  # Testa cenários de sucesso
├── testar_cenarios_falha.sh    # Testa cenários de falha
├── testar_cenarios_validacao.sh # Testa cenários de validação
└── README.md                   # Este arquivo
```

## Pré-requisitos

1. **Servidor rodando**: O backend deve estar rodando em `http://localhost:8080`
2. **Banco de dados**: PostgreSQL deve estar acessível (via docker-compose)
3. **Dependências**: `curl` e `jq` (para formatação JSON)

### Instalar jq (se necessário)

```bash
# Ubuntu/Debian
sudo apt-get install jq

# macOS
brew install jq
```

## Como Usar

### 1. Iniciar o servidor

```bash
# Iniciar o banco de dados
docker-compose up -d

# Configurar variáveis de ambiente
export DB_HOST=localhost
export DB_PORT=5434
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=mundo_invest
export DB_SSLMODE=disable
export PORT=8080

# Executar o servidor
go run cmd/server/main.go
```

### 2. Executar os testes

Execute cada script em um terminal separado:

```bash
# Testar cenários de sucesso
./scripts/testar_cenarios_sucesso.sh

# Testar cenários de falha
./scripts/testar_cenarios_falha.sh

# Testar cenários de validação
./scripts/testar_cenarios_validacao.sh
```

### 3. Verificar os resultados

Os resultados são salvos automaticamente em `scripts/exports/` com o formato:
- `cenarios_sucesso_YYYYMMDD_HHMMSS.txt`
- `cenarios_falha_YYYYMMDD_HHMMSS.txt`
- `cenarios_validacao_YYYYMMDD_HHMMSS.txt`

## Cenários Testados

### Cenários de Sucesso (`testar_cenarios_sucesso.sh`)

- **Criar cliente com dados válidos**: Criação bem-sucedida com todos os campos obrigatórios

### Cenários de Falha (`testar_cenarios_falha.sh`)

- **Sem campo cliente_nome**: Falha ao omitir nome
- **Sem campo tipo_solicitacao**: Falha ao omitir tipo de solicitação
- **Sem campo valor_patrimonio**: Falha ao omitir valor patrimonial
- **Sem campo cliente_email**: Falha ao omitir e-mail
- **Erro de banco de dados**: Falha quando banco está indisponível (requer parar o PostgreSQL)

### Cenários de Validação (`testar_cenarios_validacao.sh`)

- **E-mail inválido**: Falha com formato de e-mail incorreto
- **Patrimônio negativo**: Falha com valor patrimonial negativo
- **Patrimônio zero**: Falha com valor patrimonial igual a zero
- **Nome vazio**: Falha com nome em branco
- **E-mail vazio**: Falha com e-mail em branco
- **Tipo_solicitacao vazio**: Falha com tipo de solicitação em branco
- **Nome muito curto**: Falha com nome com menos de 3 caracteres

## Configuração

Para alterar a URL base do servidor, edite a variável `BASE_URL` no início de cada script:

```bash
BASE_URL="http://localhost:8080"
```

## Exemplo de Saída

```
========================================
TESTES DE CENÁRIOS DE SUCESSO
Data: Mon May 27 11:30:00 UTC 2026
========================================

=== CENÁRIO 1: Criar cliente com dados válidos ===
Descrição: Criar um novo cliente com todos os campos obrigatórios preenchidos corretamente

Payload enviado:
{
  "cliente_nome": "João Silva",
  "cliente_email": "joao.silva@example.com",
  "tipo_solicitacao": "Atualização cadastral",
  "valor_patrimonio": 250000
}

Requisição: POST http://localhost:8080/clientes

Código HTTP: 201
Resposta:
{
  "mensagem": "Cliente criado com sucesso",
  "identificador_interno": "uuid-aqui",
  "identificador_externo": "card-id-aqui",
  ...
}

Verificações:
✓ HTTP 201 Created - PASSOU
✓ Campo identificador_interno presente - PASSOU
✓ Campo identificador_externo (card_id) presente - PASSOU
✓ Campo status presente: Aguardando Análise
✓ Status é 'Aguardando Análise' - PASSOU
```

## Troubleshooting

### Erro: "curl: command not found"

Instale o curl:
```bash
# Ubuntu/Debian
sudo apt-get install curl

# macOS
brew install curl
```

### Erro: "jq: command not found"

Instale o jq (veja pré-requisitos acima)

### Servidor não responde

Verifique se o servidor está rodando:
```bash
curl http://localhost:8080/clientes
```

### Erro de conexão com banco

Verifique se o PostgreSQL está rodando:
```bash
docker-compose ps
```

## Notas

- Os scripts são idempotentes - podem ser executados múltiplas vezes
- Cada execução cria um novo arquivo de resultado com timestamp
- Os scripts não modificam dados existentes no banco (apenas criam novos registros)
- Para o cenário de erro de banco de dados, pare o PostgreSQL antes de executar:
  ```bash
  docker-compose stop postgres
  ```
  E reinicie após o teste:
  ```bash
  docker-compose start postgres
  ```
