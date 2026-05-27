# Critérios de Aceite - Sprint 5: Documentação e Finalização

## Critérios de Aceite

- [ ] README.md completo e funcional
- [ ] Projeto pode ser executado seguindo README
- [ ] Testes podem ser executados seguindo README
- [ ] Exemplos de curl funcionam

## Validação

Para validar que a sprint foi concluída com sucesso:

### 1. Validação do README.md

Verificar que o README.md contém:

#### Seções Obrigatórias
- [ ] Título e descrição do projeto
- [ ] Pré-requisitos (Go/Python, Docker, PostgreSQL)
- [ ] Instruções de instalação
- [ ] Instruções de execução local
- [ ] Instruções de execução dos testes
- [ ] Exemplos de uso da API
- [ ] Estrutura do projeto
- [ ] Tecnologias utilizadas
- [ ] Autores e contribuidores

#### Instruções de Execução Local
```bash
# Deve funcionar seguindo estas instruções:
git clone <repositorio>
cd MundoInvest
docker-compose up -d
go run cmd/server/main.go
# ou
python cmd/server/main.py
```

#### Instruções de Execução dos Testes
```bash
# Deve funcionar seguindo estas instruções:
go test ./... -v
# ou
pytest tests/
```

### 2. Validação dos Exemplos de curl

#### Exemplo POST /clientes
```bash
curl -X POST http://localhost:8080/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "cliente_nome": "João Silva",
    "cliente_email": "joao.silva@example.com",
    "tipo_solicitacao": "abertura_conta",
    "valor_patrimonio": 150000.00
  }'

# Esperado: HTTP 201 com resposta JSON
```

#### Exemplo POST /webhooks/pipefy/card-updated
```bash
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
  -H "Content-Type: application/json" \
  -d '{
    "identificador_evento": "evt_12345",
    "identificador_card": "card_67890",
    "cliente_email": "joao.silva@example.com",
    "data_evento": "2026-05-27T10:00:00Z"
  }'

# Esperado: HTTP 200 com resposta JSON
```

### 3. Validação da Documentação Técnica

Verificar que a documentação técnica está atualizada:

- [ ] Regras de negócio documentadas em docs/regras-negocio.md
- [ ] Requisitos funcionais documentados em docs/requisitos-funcionais.md
- [ ] Modelagem de dados documentada em docs/modelagem-dados.md
- [ ] Snapshot do banco de dados documentado em docs/snapshot-banco-dados.md
- [ ] Estrutura de pastas documentada com Feature Folders
- [ ] Mapeamento de arquivos por contexto delimitado
- [ ] Mutations GraphQL documentadas com referências

### 4. Validação da Preparação para Defesa

Verificar que:

- [ ] Script de apresentação está preparado
- [ ] Pontos altos do código identificados (mutations GraphQL, estrutura DDD)
- [ ] Estrutura de pastas está organizada para explicação
- [ ] Demonstração do sistema está preparada
- [ ] Documentação é consistente em todos os arquivos

### 5. Validação Final do Projeto

#### Verificação de Compilação
```bash
# Go
go build ./...

# Python
python -m py_compile cmd/server/main.py
```

#### Verificação de Testes
```bash
# Executar todos os testes
go test ./... -v

# Esperado: Todos os 22 cenários passam
```

#### Verificação de Docker
```bash
# Iniciar PostgreSQL
docker-compose up -d

# Verificar status
docker-compose ps

# Esperado: PostgreSQL rodando corretamente
```

#### Verificação de Código
- [ ] Código está formatado corretamente
- [ ] Não há comentários desnecessários
- [ ] Nomenclatura segue linguagem ubíqua
- [ ] Não há código morto ou não utilizado
- [ ] Imports estão organizados

### 6. Checklist de Entrega

#### Funcionalidade
- [ ] POST /clientes funciona corretamente
- [ ] POST /webhooks/pipefy/card-updated funciona corretamente
- [ ] Idempotência funciona
- [ ] Cálculo de prioridade funciona
- [ ] Validações funcionam

#### Qualidade
- [ ] Todos os testes passam
- [ ] Cobertura de testes > 80%
- [ ] Código segue padrões da linguagem
- [ ] Não há warnings ou erros

#### Documentação
- [ ] README.md está completo
- [ ] Documentação técnica está atualizada
- [ ] Exemplos funcionam
- [ ] Código está comentado onde necessário

#### Estrutura
- [ ] Feature Folders implementados corretamente
- [ ] Contextos delimitados respeitados
- [ ] Linguagem ubíqua seguida
- [ ] Trigramação aplicada

## Critérios de Sucesso do Projeto

Verificar que todos os critérios de sucesso do projeto foram atendidos:

- [ ] Todos os requisitos funcionais implementados
- [ ] Todas as regras de negócio aplicadas
- [ ] Todos os testes obrigatórios passando
- [ ] Mutations GraphQL estruturadas corretamente
- [ ] README completo e funcional
- [ ] Projeto pronto para defesa em vídeo
- [ ] Modelagem de dados segue DDD + 3FN + Linguagem Ubíqua
- [ ] Nomenclatura de tabelas e colunas usa trigramação
- [ ] Schemas por contexto delimitado implementados
- [ ] Sequências SQL para identificadores internos funcionando
- [ ] Views criadas para acesso legível
- [ ] Contextos delimitados respeitados na implementação
- [ ] Estratégia de backup e snapshot implementada
- [ ] Scripts de automação de backup funcionando
- [ ] Snapshot do banco de dados documentado
- [ ] Plano de recuperação de desastres definido
- [ ] Integração com pipeline CI/CD para backups
- [ ] Controle de acesso RBAC implementado
- [ ] Roles e usuários com permissões mínimas configurados
- [ ] Auditoria de acesso funcionando
- [ ] Políticas de segurança aplicadas
