-- Migration: Criar tabela de cards do Pipefy com trigramação
-- Contexto: Integração Pipefy
-- Linguagem Ubíqua: Card, IdentificadorCard, Mutation, Query
-- Trigramação: ipf_car_int, ipf_car_ext, ipf_car_ema, ipf_car_mut, ipf_car_que

CREATE TABLE IF NOT EXISTS integracao_pipefy.card (
    -- Identificador interno (gerado pela sequência)
    ipf_car_int BIGINT PRIMARY KEY DEFAULT nextval('integracao_pipefy.seq_ipf_car_int'),
    
    -- Identificador externo (card_id do Pipefy)
    ipf_car_ext VARCHAR(255) NOT NULL UNIQUE, -- Identificador do card no Pipefy
    
    -- Email do cliente (para correlação)
    ipf_car_ema VARCHAR(255) NOT NULL, -- Email do cliente
    
    -- Mutations e Queries
    ipf_car_mut TEXT, -- Mutation GraphQL estruturada
    ipf_car_que TEXT, -- Query GraphQL estruturada
    
    -- Status do card
    ipf_car_stc VARCHAR(50) DEFAULT 'Pendente', -- Status do card
    
    -- Timestamps
    ipf_car_dcr TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Data de criação
    ipf_car_dat TIMESTAMP, -- Data de atualização
    
    -- Índices
    CONSTRAINT uk_ipf_car_ext UNIQUE (ipf_car_ext), -- Identificador externo único
    CONSTRAINT fk_ipf_car_ema FOREIGN KEY (ipf_car_ema) REFERENCES gestao_clientes.cliente(gcl_cli_ema) ON DELETE CASCADE
);

-- Criar índice para busca por identificador externo
CREATE INDEX IF NOT EXISTS idx_ipf_car_ext ON integracao_pipefy.card(ipf_car_ext);

-- Criar índice para busca por email do cliente
CREATE INDEX IF NOT EXISTS idx_ipf_car_ema ON integracao_pipefy.card(ipf_car_ema);

-- Criar índice para busca por status
CREATE INDEX IF NOT EXISTS idx_ipf_car_stc ON integracao_pipefy.card(ipf_car_stc);

-- Comentários da tabela e colunas
COMMENT ON TABLE integracao_pipefy.card IS 'Tabela de cards do contexto de integração Pipefy';
COMMENT ON COLUMN integracao_pipefy.card.ipf_car_int IS 'Identificador interno do card';
COMMENT ON COLUMN integracao_pipefy.card.ipf_car_ext IS 'Identificador do card no Pipefy';
COMMENT ON COLUMN integracao_pipefy.card.ipf_car_ema IS 'Email do cliente';
COMMENT ON COLUMN integracao_pipefy.card.ipf_car_mut IS 'Mutation GraphQL estruturada';
COMMENT ON COLUMN integracao_pipefy.card.ipf_car_que IS 'Query GraphQL estruturada';
COMMENT ON COLUMN integracao_pipefy.card.ipf_car_stc IS 'Status do card (Pendente, Processado, etc.)';
COMMENT ON COLUMN integracao_pipefy.card.ipf_car_dcr IS 'Data de criação do card';
COMMENT ON COLUMN integracao_pipefy.card.ipf_car_dat IS 'Data de atualização do card';
