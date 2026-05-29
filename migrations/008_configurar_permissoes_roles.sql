-- Migration: Configurar permissões mínimas por role (RBAC)
-- Contexto: Segurança do Banco de Dados
-- Linguagem Ubíqua: Permissões, Princípio do Menor Privilégio
-- Objetivo: Configurar permissões mínimas por role seguindo o princípio do menor privilégio

-- Permissões para gestao_clientes_read (apenas leitura)
GRANT USAGE ON SCHEMA gestao_clientes TO gestao_clientes_read;
GRANT SELECT ON ALL TABLES IN SCHEMA gestao_clientes TO gestao_clientes_read;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA gestao_clientes TO gestao_clientes_read;

-- Permissões para gestao_clientes_write (leitura e escrita)
GRANT USAGE ON SCHEMA gestao_clientes TO gestao_clientes_write;
GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA gestao_clientes TO gestao_clientes_write;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA gestao_clientes TO gestao_clientes_write;

-- Permissões para processamento_eventos_write (leitura e escrita)
GRANT USAGE ON SCHEMA processamento_eventos TO processamento_eventos_write;
GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA processamento_eventos TO processamento_eventos_write;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA processamento_eventos TO processamento_eventos_write;

-- Permissões para backup_operator (backup e restauração)
GRANT USAGE ON SCHEMA gestao_clientes TO backup_operator;
GRANT USAGE ON SCHEMA processamento_eventos TO backup_operator;
GRANT USAGE ON SCHEMA integracao_pipefy TO backup_operator;
GRANT SELECT ON ALL TABLES IN SCHEMA gestao_clientes TO backup_operator;
GRANT SELECT ON ALL TABLES IN SCHEMA processamento_eventos TO backup_operator;
GRANT SELECT ON ALL TABLES IN SCHEMA integracao_pipefy TO backup_operator;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA gestao_clientes TO backup_operator;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA processamento_eventos TO backup_operator;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA integracao_pipefy TO backup_operator;

-- Permissões para dba_admin (privilégios completos)
GRANT ALL PRIVILEGES ON SCHEMA gestao_clientes TO dba_admin;
GRANT ALL PRIVILEGES ON SCHEMA processamento_eventos TO dba_admin;
GRANT ALL PRIVILEGES ON SCHEMA integracao_pipefy TO dba_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA gestao_clientes TO dba_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA processamento_eventos TO dba_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA integracao_pipefy TO dba_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA gestao_clientes TO dba_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA processamento_eventos TO dba_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA integracao_pipefy TO dba_admin;

-- Configurar permissões futuras (ALTER DEFAULT)
ALTER DEFAULT PRIVILEGES IN SCHEMA gestao_clientes GRANT SELECT ON TABLES TO gestao_clientes_read;
ALTER DEFAULT PRIVILEGES IN SCHEMA gestao_clientes GRANT SELECT, INSERT, UPDATE ON TABLES TO gestao_clientes_write;
ALTER DEFAULT PRIVILEGES IN SCHEMA processamento_eventos GRANT SELECT, INSERT, UPDATE ON TABLES TO processamento_eventos_write;
ALTER DEFAULT PRIVILEGES IN SCHEMA gestao_clientes GRANT USAGE, SELECT ON SEQUENCES TO gestao_clientes_write;
ALTER DEFAULT PRIVILEGES IN SCHEMA processamento_eventos GRANT USAGE, SELECT ON SEQUENCES TO processamento_eventos_write;
