-- Migration: Criar tabela de auditoria de acesso
-- Contexto: Segurança do Banco de Dados
-- Linguagem Ubíqua: Auditoria, Acesso, Logs
-- Objetivo: Implementar auditoria de acesso com triggers

-- Criar schema para auditoria
CREATE SCHEMA IF NOT EXISTS auditoria;

-- Criar tabela de auditoria de acesso
CREATE TABLE IF NOT EXISTS auditoria.log_acesso (
    log_acs_int BIGSERIAL PRIMARY KEY,
    log_acs_usu VARCHAR(255) NOT NULL, -- Usuário que realizou a operação
    log_acs_ope VARCHAR(50) NOT NULL, -- Operação realizada (SELECT, INSERT, UPDATE, DELETE)
    log_acs_sch VARCHAR(255) NOT NULL, -- Schema acessado
    log_acs_tba VARCHAR(255) NOT NULL, -- Tabela acessada
    log_acs_reg TEXTO, -- Registro afetado (JSON)
    log_acs_dcr TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Data da operação
    log_acs_ip VARCHAR(45), -- Endereço IP (se disponível)
    log_acs_det TEXTE -- Detalhes adicionais
);

-- Criar índices para busca eficiente
CREATE INDEX IF NOT EXISTS idx_log_acs_usu ON auditoria.log_acesso(log_acs_usu);
CREATE INDEX IF NOT EXISTS idx_log_acs_ope ON auditoria.log_acesso(log_acs_ope);
CREATE INDEX IF NOT EXISTS idx_log_acs_sch ON auditoria.log_acesso(log_acs_sch);
CREATE INDEX IF NOT EXISTS idx_log_acs_tba ON auditoria.log_acesso(log_acs_tba);
CREATE INDEX IF NOT EXISTS idx_log_acs_dcr ON auditoria.log_acesso(log_acs_dcr);

-- Comentários da tabela e colunas
COMMENT ON SCHEMA auditoria IS 'Schema para auditoria de acesso ao banco de dados';
COMMENT ON TABLE auditoria.log_acesso IS 'Tabela de logs de auditoria de acesso';
COMMENT ON COLUMN auditoria.log_acesso.log_acs_int IS 'Identificador interno do log';
COMMENT ON COLUMN auditoria.log_acesso.log_acs_usu IS 'Usuário que realizou a operação';
COMMENT ON COLUMN auditoria.log_acesso.log_acs_ope IS 'Operação realizada (SELECT, INSERT, UPDATE, DELETE)';
COMMENT ON COLUMN auditoria.log_acesso.log_acs_sch IS 'Schema acessado';
COMMENT ON COLUMN auditoria.log_acesso.log_acs_tba IS 'Tabela acessada';
COMMENT ON COLUMN auditoria.log_acesso.log_acs_reg IS 'Registro afetado (JSON)';
COMMENT ON COLUMN auditoria.log_acesso.log_acs_dcr IS 'Data da operação';
COMMENT ON COLUMN auditoria.log_acesso.log_acs_ip IS 'Endereço IP (se disponível)';
COMMENT ON COLUMN auditoria.log_acesso.log_acs_det IS 'Detalhes adicionais';
