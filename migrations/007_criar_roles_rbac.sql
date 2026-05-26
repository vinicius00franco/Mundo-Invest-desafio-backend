-- Migration: Criar roles de acesso (RBAC)
-- Contexto: Segurança do Banco de Dados
-- Linguagem Ubíqua: Roles, Permissões, Princípio do Menor Privilégio
-- Objetivo: Implementar controle de acesso baseado em roles

-- Criar role para leitura de dados de clientes
CREATE ROLE IF NOT EXISTS gestao_clientes_read WITH NOLOGIN;

-- Criar role para escrita de dados de clientes
CREATE ROLE IF NOT EXISTS gestao_clientes_write WITH NOLOGIN;

-- Criar role para escrita de eventos de processamento
CREATE ROLE IF NOT EXISTS processamento_eventos_write WITH NOLOGIN;

-- Criar role para operações de backup
CREATE ROLE IF NOT EXISTS backup_operator WITH NOLOGIN;

-- Criar role para administradores do banco
CREATE ROLE IF NOT EXISTS dba_admin WITH NOLOGIN;

-- Comentários das roles
COMMENT ON ROLE gestao_clientes_read IS 'Role para leitura de dados do contexto de gestão de clientes';
COMMENT ON ROLE gestao_clientes_write IS 'Role para escrita de dados do contexto de gestão de clientes';
COMMENT ON ROLE processamento_eventos_write IS 'Role para escrita de eventos do contexto de processamento de eventos';
COMMENT ON ROLE backup_operator IS 'Role para operações de backup do banco de dados';
COMMENT ON ROLE dba_admin IS 'Role para administradores do banco de dados com privilégios completos';
