-- Migration: Criar schemas por contexto delimitado
-- Contexto: Configuração inicial do banco de dados
-- Linguagem Ubíqua: Contextos Delimitados (Gestão de Clientes, Integração Pipefy, Processamento de Eventos)

-- Criar schema para contexto de Gestão de Clientes
CREATE SCHEMA IF NOT EXISTS gestao_clientes;

-- Criar schema para contexto de Integração Pipefy
CREATE SCHEMA IF NOT EXISTS integracao_pipefy;

-- Criar schema para contexto de Processamento de Eventos
CREATE SCHEMA IF NOT EXISTS processamento_eventos;

-- Comentários dos schemas
COMMENT ON SCHEMA gestao_clientes IS 'Contexto delimitado para gestão de clientes e seus patrimônios';
COMMENT ON SCHEMA integracao_pipefy IS 'Contexto delimitado para integração com Pipefy';
COMMENT ON SCHEMA processamento_eventos IS 'Contexto delimitado para processamento de eventos do Pipefy';
