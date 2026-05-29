-- Migration: Criar tabela de clientes com trigramação
-- Contexto: Gestão de Clientes
-- Linguagem Ubíqua: Cliente, Patrimonio, StatusCliente, NivelPrioridade
-- Trigramação: gcl_cli_int, gcl_cli_ext, gcl_cli_nom, gcl_cli_ema, gcl_cli_pat, gcl_cli_tso, gcl_cli_stc, gcl_cli_npr

CREATE TABLE IF NOT EXISTS gestao_clientes.cliente (
    -- Identificador interno (gerado pela sequência)
    gcl_cli_int BIGINT PRIMARY KEY DEFAULT nextval('gestao_clientes.seq_gcl_cli_int'),
    
    -- Identificador externo (card_id do Pipefy)
    gcl_cli_ext VARCHAR(255),
    
    -- Dados do cliente
    gcl_cli_nom VARCHAR(255) NOT NULL, -- Nome do cliente
    gcl_cli_ema VARCHAR(255) NOT NULL, -- Email do cliente
    gcl_cli_pat DECIMAL(15, 2) NOT NULL, -- Valor do patrimônio
    gcl_cli_tso VARCHAR(50) NOT NULL, -- Tipo de solicitação
    
    -- Status e prioridade
    gcl_cli_stc VARCHAR(50) DEFAULT 'Aguardando Análise', -- Status do cliente
    gcl_cli_npr VARCHAR(50), -- Nível de prioridade (alta/normal)
    
    -- Timestamps
    gcl_cli_dcr TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Data de criação
    gcl_cli_dat TIMESTAMP, -- Data de atualização
    
    -- Índices
    CONSTRAINT uk_gcl_cli_ema UNIQUE (gcl_cli_ema), -- Email único
    CONSTRAINT uk_gcl_cli_ext UNIQUE (gcl_cli_ext) -- Identificador externo único
);

-- Criar índice para busca por email
CREATE INDEX IF NOT EXISTS idx_gcl_cli_ema ON gestao_clientes.cliente(gcl_cli_ema);

-- Criar índice para busca por identificador externo
CREATE INDEX IF NOT EXISTS idx_gcl_cli_ext ON gestao_clientes.cliente(gcl_cli_ext);

-- Criar índice para busca por status
CREATE INDEX IF NOT EXISTS idx_gcl_cli_stc ON gestao_clientes.cliente(gcl_cli_stc);

-- Comentários da tabela e colunas
COMMENT ON TABLE gestao_clientes.cliente IS 'Tabela de clientes do contexto de gestão de clientes';
COMMENT ON COLUMN gestao_clientes.cliente.gcl_cli_int IS 'Identificador interno do cliente (gerado pela sequência)';
COMMENT ON COLUMN gestao_clientes.cliente.gcl_cli_ext IS 'Identificador externo do cliente (card_id do Pipefy)';
COMMENT ON COLUMN gestao_clientes.cliente.gcl_cli_nom IS 'Nome do cliente';
COMMENT ON COLUMN gestao_clientes.cliente.gcl_cli_ema IS 'Email do cliente';
COMMENT ON COLUMN gestao_clientes.cliente.gcl_cli_pat IS 'Valor do patrimônio do cliente';
COMMENT ON COLUMN gestao_clientes.cliente.gcl_cli_tso IS 'Tipo de solicitação';
COMMENT ON COLUMN gestao_clientes.cliente.gcl_cli_stc IS 'Status do cliente (Aguardando Análise, Processado, etc.)';
COMMENT ON COLUMN gestao_clientes.cliente.gcl_cli_npr IS 'Nível de prioridade (alta, normal)';
COMMENT ON COLUMN gestao_clientes.cliente.gcl_cli_dcr IS 'Data de criação do cliente';
COMMENT ON COLUMN gestao_clientes.cliente.gcl_cli_dat IS 'Data de atualização do cliente';
