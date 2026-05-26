# Resumo Executivo - Plano de Implementação

## Modelagem de Dados (DDD + Normalização + Linguagem Ubíqua + Trigramação)

### Contextos Delimitados (Schemas)

O sistema está organizado em três schemas no banco de dados, cada um representando um contexto delimitado:

1. **gestao_clientes**: Cliente, Patrimonio, Solicitacao, StatusCliente, NivelPrioridade
2. **integracao_pipefy**: Card, IdentificadorCard, Mutation, Query
3. **processamento_eventos**: Webhook, Evento, IdentificadorEvento, Idempotencia

### Normalização (3FN)

A modelagem segue a Terceira Forma Normal:
- **1FN**: Atributos atômicos, sem grupos repetitivos
- **2FN**: Dependência total da chave primária
- **3FN**: Sem dependências transitivas

### Trigramação de Colunas

Todas as colunas usam trigramação (3 letras por palavra) para garantir unicidade e clareza:

**Exemplos:**
- `gcl_cli_int` = gestao_clientes_cliente_identificador_interno
- `pev_eve_ide` = processamento_eventos_evento_identificador_evento
- `gcl_cli_nom` = gestao_clientes_cliente_nome

### Sequências (Sequences)

Identificadores internos são gerados por sequências SQL:
- `gestao_clientes.seq_gcl_cli_int` (para cliente)
- `processamento_eventos.seq_pev_eve_int` (para evento)

### Views (Visões)

Views são criadas para fornecer acesso legível às tabelas com colunas trigramadas:
- `gestao_clientes.vw_cliente` (acesso legível)
- `processamento_eventos.vw_evento` (acesso legível)
- `gestao_clientes.vw_cliente_com_eventos` (visão consolidada)
- `processamento_eventos.vw_eventos_por_cliente` (visão consolidada)

### Controle de Acesso (RBAC)

Sistema de controle de acesso baseado em roles seguindo o princípio do menor privilégio:

**Roles Definidos:**
- `gestao_clientes_read` (SELECT apenas)
- `gestao_clientes_write` (SELECT, INSERT, UPDATE)
- `processamento_eventos_write` (SELECT, INSERT, UPDATE)
- `backup_operator` (SELECT, pg_dump)
- `dba_admin` (ALL PRIVILEGES)

**Usuários Definidos:**
- `app_cliente_read` (aplicação de leitura)
- `app_cliente_write` (aplicação principal)
- `app_webhook_write` (aplicação de webhooks)
- `backup_user` (scripts de backup)
- `dba_senior` (administração)

**Segurança Adicional:**
- Senhas com SCRAM-SHA-256
- Expiração de senhas a cada 90 dias
- Limites de conexão por aplicação
- Auditoria de acesso com triggers
- Logging de todas as modificações

### Trigramação de Nomes

**Tabelas:**
- `clientes` → `cliente` (singular, nome do domínio)
- `eventos` → `evento` (singular, nome do domínio)

**Colunas (Trigramadas):**
- `id` → `gcl_cli_int` (gestao_clientes_cliente_identificador_interno)
- `uuid` → `gcl_cli_ext` (gestao_clientes_cliente_identificador_externo)
- `nome` → `gcl_cli_nom` (gestao_clientes_cliente_nome)
- `email` → `gcl_cli_ema` (gestao_clientes_cliente_email)
- `status` → `gcl_cli_stc` (gestao_clientes_cliente_status_cliente)
- `card_id` → `gcl_cli_idc` (gestao_clientes_cliente_identificador_card)
- `prioridade` → `gcl_cli_npr` (gestao_clientes_cliente_nivel_prioridade)
- `event_id` → `pev_eve_ide` (processamento_eventos_evento_identificador_evento)
- `processado` → `pev_eve_fpr` (processamento_eventos_evento_foi_processado)
- `created_at` → `gcl_cli_dcr` (gestao_clientes_cliente_data_criacao)
- `updated_at` → `gcl_cli_dat` (gestao_clientes_cliente_data_atualizacao)
- `timestamp` → `pev_eve_dev` (processamento_eventos_evento_data_evento)

**Schemas por Domínio:**
- `gestao_clientes` (Contexto: Gestão de Clientes)
- `integracao_pipefy` (Contexto: Integração Pipefy)
- `processamento_eventos` (Contexto: Processamento de Eventos)

**Sequências:**
- `gestao_clientes.seq_gcl_cli_int` (para cliente)
- `processamento_eventos.seq_pev_eve_int` (para evento)

**Views (para legibilidade):**
- `gestao_clientes.vw_cliente` (acesso legível à tabela cliente)
- `processamento_eventos.vw_evento` (acesso legível à tabela evento)
- `gestao_clientes.vw_cliente_com_eventos` (visão consolidada)
- `processamento_eventos.vw_eventos_por_cliente` (visão consolidada)

**Arquivos e Estruturas:**
- `model/cliente.go` (não `model/customer.go`)
- `model/evento.go` (não `model/event.go`)
- `repository/cliente_repository.go` (não `repository/customer_repository.go`)
- `service/prioridade_calculator.go` (não `service/priority_calculator.go`)
- `migrations/001_criar_schemas.sql` (criar schemas por domínio)
- `migrations/002_criar_sequencia_cliente.sql` (criar sequência para cliente)
- `migrations/003_criar_tabela_cliente.sql` (criar tabela cliente com colunas trigramadas)

### Estratégia de Identificação

- **Identificador Interno**: INTEGER/SERIAL (banco, joins, performance)
- **Identificador Externo**: UUID v4 (APIs, JWT, segurança)

## Análise Completa dos Cenários

### Total de Cenários Identificados: 22

#### 📊 Distribuição por Tipo:
- **Cenários de Sucesso:** 5 (23%)
- **Cenários de Falha:** 4 (18%)
- **Cenários de Validação:** 13 (59%)

#### 📋 Por Funcionalidade:

**Criação de Cliente (POST /clientes) - 9 cenários:**
- ✅ Sucesso: 1 cenário
- ❌ Falha: 2 cenários (campos obrigatórios, erro de banco)
- ⚠️ Validação: 6 cenários (email, patrimônio negativo/zero, nome vazio)

**Webhook Card Updated (POST /webhooks/pipefy/card-updated) - 13 cenários:**
- ✅ Sucesso: 4 cenários (prioridade alta/normal/limite, idempotência)
- ❌ Falha: 2 cenários (cliente não encontrado, erro de banco)
- ⚠️ Validação: 7 cenários (campos obrigatórios, campos vazios, formatos inválidos)

## Atualizações Realizadas nos Documentos

### 📄 snapshot-banco-dados.md (NOVO)
- ✅ Snapshot completo da estrutura do banco de dados
- ✅ Documentação detalhada de schemas, tabelas, colunas
- ✅ Trigramação de todas as colunas documentada
- ✅ Sequências, índices, views, constraints e triggers
- ✅ Scripts de backup e restauração
- ✅ Validação de integridade
- ✅ Checklist de validação
- ✅ Dicionário de dados com glossário de trigramação

### 📄 modelagem-dados.md (NOVO)
- ✅ Contextos delimitados definidos (Gestão de Clientes, Integração Pipefy, Processamento de Eventos)
- ✅ Linguagem ubíqua para cada contexto
- ✅ Normalização 3FN aplicada
- ✅ Trigramação de nomes de tabelas e colunas
- ✅ Histórias de domínio para cada contexto
- ✅ Estratégia de identificação (interno vs externo)
- ✅ Autenticação JWT com identificador_externo
- ✅ Scripts PostgreSQL atualizados
- ✅ Glossário de linguagem ubíqua

### 📄 regras-negocio.md

### 📄 regras-negocio.md
- ✅ Adicionada RN-034: Validação de e-mail no webhook
- ✅ Adicionada RN-035: Validação de timestamp

### 📄 requisitos-funcionais.md
- ✅ Adicionados RF-036 a RF-042: Testes adicionais para cobrir todos os cenários
- ✅ Teste de validação de e-mail
- ✅ Teste de validação de patrimônio negativo
- ✅ Teste de validação de patrimônio zero
- ✅ Teste de validação de nome vazio
- ✅ Teste de prioridade no limite
- ✅ Teste de validação de webhook
- ✅ Teste de event_id vazio

### 📄 bdd/criacao-cliente/falha.md
- ✅ Adicionado cenário: Criar cliente sem e-mail
- ✅ Adicionado cenário: Criar cliente com erro de banco de dados

### 📄 bdd/webhook-card-updated/falha.md
- ✅ Adicionado cenário: Processar webhook com erro de banco de dados

### 📄 bdd/webhook-card-updated/validacao.md
- ✅ Adicionado cenário: Processar webhook com card_id vazio
- ✅ Adicionado cenário: Processar webhook com cliente_email vazio
- ✅ Adicionado cenário: Processar webhook com timestamp inválido

## Estrutura do Plano de Implementação

### 🏗️ Arquitetura
- Diagrama de blocos mostrando todas as camadas
- Separação clara de responsabilidades (Controller, Service, Repository, Integration)
- Integração simulada com Pipefy via GraphQL

### 🔄 Fluxos
- Diagrama de sequência para criação de cliente
- Diagrama de sequência para processamento de webhook
- Tratamento completo de erros e validações

### 📅 Cronograma (5 Sprints - 9 dias)

**Sprint 1: Fundamentos e Estrutura (2 dias)**
- Configuração do projeto
- Modelo de dados
- Camada de persistência

**Sprint 2: API de Criação de Cliente (2 dias)**
- Validação de dados
- Camada de serviço
- Camada de apresentação
- Integração Pipefy (simulação)

**Sprint 3: API de Webhook (2 dias)**
- Validação de dados
- Camada de serviço
- Cálculo de prioridade
- Integração Pipefy (simulação)

**Sprint 4: Testes Automatizados (2 dias)**
- 9 testes de criação de cliente
- 13 testes de webhook
- Testes de integração

**Sprint 5: Documentação e Finalização (1 dia)**
- README.md completo
- Documentação técnica
- Preparação para defesa

## Matriz de Riscos Identificados

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|------------|
| Documentação Pipefy incompleta | Média | Alta | Pesquisar múltiplas fontes |
| Idempotência não funcionar | Baixa | Alta | Verificação no início do processamento |
| Validação de e-mail falhar | Baixa | Média | Regex validado + testes extensivos |
| Concorrência em webhooks | Baixa | Média | Transações + índices únicos |
| Tempo insuficiente | Média | Alta | Priorizar funcionalidades críticas |

## Próximos Passos

1. ✅ Análise completa de cenários realizada
2. ✅ Documentos atualizados com regras e requisitos adicionais
3. ✅ Plano detalhado criado com diagramas
4. ✅ Modelagem de dados com DDD + 3FN + Linguagem Ubíqua
5. ✅ Trigramação de nomes aplicada
6. ✅ Contextos delimitados definidos
7. ✅ Sequências SQL para identificadores internos
8. ✅ Views para acesso legível
9. ✅ Estratégia de backup e snapshot do banco de dados
10. ✅ Scripts de automação para backup e restauração
11. ✅ Integração com pipeline CI/CD
12. ✅ Plano de recuperação de desastres
13. ⏭️ Revisar plano com stakeholders
14. ⏭️ Iniciar implementação (Sprint 1)

## Critérios de Sucesso

- [x] Todos os 22 cenários identificados e documentados
- [x] Regras de negócio atualizadas
- [x] Requisitos funcionais completados
- [x] Cenários BDD detalhados
- [x] Plano de implementação com diagramas
- [x] Modelagem de dados com DDD + 3FN + Linguagem Ubíqua
- [x] Trigramação de nomes de tabelas e colunas
- [x] Schemas por contexto delimitado definidos
- [x] Sequências SQL para identificadores internos
- [x] Views para acesso legível às tabelas
- [x] Contextos delimitados definidos
- [x] Estratégia de identificação (interno vs externo)
- [x] Estratégia de backup e snapshot do banco de dados
- [x] Scripts de automação para backup e restauração
- [x] Snapshot completo da estrutura documentado
- [x] Plano de recuperação de desastres definido
- [x] Integração com pipeline CI/CD para backups
- [x] Controle de acesso RBAC implementado
- [x] Roles e usuários com permissões mínimas
- [x] Auditoria de acesso configurada
- [x] Políticas de segurança aplicadas
- [ ] Implementação completa (9 dias)
- [ ] Todos os testes passando
- [ ] README funcional
- [ ] Projeto pronto para defesa