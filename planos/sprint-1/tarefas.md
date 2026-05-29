# Tarefas - Sprint 1: Fundamentos e Estrutura

## 1. Configuração do Projeto

- [ ] Escolher linguagem (Golang ou Python)
- [ ] Inicializar projeto com dependências básicas
- [ ] Configurar estrutura de pastas seguindo Feature Folders (DDD + contextos delimitados)
- [ ] Criar estrutura: cmd/server/, internal/gestao_clientes/, internal/integracao_pipefy/, internal/processamento_eventos/, internal/dominio/, internal/shared/, pkg/validation/
- [ ] Configurar PostgreSQL via Docker (desenvolvimento e testes)
- [ ] Criar docker-compose.yml com serviço PostgreSQL
- [ ] Criar pasta migrations com versionamento

## 2. Modelo de Dados (3FN + Linguagem Ubíqua + Trigramação)

- [ ] Criar schemas por contexto delimitado (gestao_clientes, integracao_pipefy, processamento_eventos)
- [ ] Criar sequências SQL para identificadores internos (seq_gcl_cli_int, seq_pev_eve_int)
- [ ] Criar entidade Cliente com colunas trigramadas (gcl_cli_int, gcl_cli_ext, gcl_cli_nom, etc.)
- [ ] Criar entidade Evento com colunas trigramadas (pev_eve_int, pev_eve_ide, pev_eve_idc, etc.)
- [ ] Criar migrations com nomenclatura descritiva (001_criar_schemas.sql, etc.)
- [ ] Criar views para acesso legível (vw_cliente, vw_evento, etc.)
- [ ] Configurar índices (gcl_cli_ema único, pev_eve_ide único)

## 3. Camada de Persistência (Feature Folders)

- [ ] Implementar internal/gestao_clientes/repository.go com ClienteRepository: salvar, buscarPorIdentificadorInterno, buscarPorClienteEmail, buscarPorIdentificadorExterno, atualizar
- [ ] Implementar internal/processamento_eventos/repository.go com EventoRepository: salvar, buscarPorIdentificadorEvento, verificarFoiProcessado
- [ ] Implementar internal/shared/database/connection.go para gerenciar conexões
- [ ] Implementar transações atômicas

## 4. Segurança do Banco de Dados (RBAC)

- [ ] Criar roles de acesso (gestao_clientes_read, gestao_clientes_write, processamento_eventos_write, backup_operator, dba_admin)
- [ ] Criar usuários específicos (app_cliente_read, app_cliente_write, app_webhook_write, backup_user, dba_senior)
- [ ] Configurar permissões mínimas por role (princípio do menor privilégio)
- [ ] Configurar limites de conexão por aplicação
- [ ] Implementar auditoria de acesso com triggers
- [ ] Configurar políticas de senha (SCRAM-SHA-256, expiração 90 dias)
- [ ] Configurar logging de consultas (mod)
