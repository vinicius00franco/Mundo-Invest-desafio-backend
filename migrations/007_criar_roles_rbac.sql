-- Migration: Criar roles de acesso (RBAC)
-- Contexto: Segurança do Banco de Dados
-- Linguagem Ubíqua: Roles, Permissões, Princípio do Menor Privilégio
-- Objetivo: Implementar controle de acesso baseado em roles

-- Criar role para leitura de dados de clientes
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'gestao_clientes_read') THEN
        CREATE ROLE gestao_clientes_read WITH NOLOGIN;
    END IF;
END $$;

-- Criar role para escrita de dados de clientes
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'gestao_clientes_write') THEN
        CREATE ROLE gestao_clientes_write WITH NOLOGIN;
    END IF;
END $$;

-- Criar role para escrita de eventos de processamento
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'processamento_eventos_write') THEN
        CREATE ROLE processamento_eventos_write WITH NOLOGIN;
    END IF;
END $$;

-- Criar role para operações de backup
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'backup_operator') THEN
        CREATE ROLE backup_operator WITH NOLOGIN;
    END IF;
END $$;

-- Criar role para administradores do banco
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'dba_admin') THEN
        CREATE ROLE dba_admin WITH NOLOGIN;
    END IF;
END $$;

-- Comentários das roles
COMMENT ON ROLE gestao_clientes_read IS 'Role para leitura de dados do contexto de gestão de clientes';
COMMENT ON ROLE gestao_clientes_write IS 'Role para escrita de dados do contexto de gestão de clientes';
COMMENT ON ROLE processamento_eventos_write IS 'Role para escrita de eventos do contexto de processamento de eventos';
COMMENT ON ROLE backup_operator IS 'Role para operações de backup do banco de dados';
COMMENT ON ROLE dba_admin IS 'Role para administradores do banco de dados com privilégios completos';
