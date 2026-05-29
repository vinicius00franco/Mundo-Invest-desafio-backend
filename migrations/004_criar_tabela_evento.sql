-- Migration: Criar tabela de eventos com trigramação
-- Contexto: Processamento de Eventos
-- Linguagem Ubíqua: Evento, IdentificadorEvento, Idempotencia, Processamento
-- Trigramação: pev_eve_int, pev_eve_ide, pev_eve_idc, pev_eve_fpr, pev_eve_dcr

CREATE TABLE IF NOT EXISTS processamento_eventos.evento (
    -- Identificador interno (gerado pela sequência)
    pev_eve_int BIGINT PRIMARY KEY DEFAULT nextval('processamento_eventos.seq_pev_eve_int'),
    
    -- Identificador do evento (event_id do Pipefy)
    pev_eve_ide VARCHAR(255) NOT NULL, -- Identificador único do evento
    
    -- Identificador do card (card_id do Pipefy)
    pev_eve_idc VARCHAR(255) NOT NULL, -- Identificador do card relacionado
    
    -- Email do cliente (para correlação)
    pev_eve_ema VARCHAR(255) NOT NULL, -- Email do cliente
    
    -- Timestamp do evento
    pev_eve_tms TIMESTAMP NOT NULL, -- Timestamp do evento recebido do Pipefy
    
    -- Status de processamento
    pev_eve_fpr BOOLEAN DEFAULT FALSE, -- Flag indicando se o evento foi processado (idempotência)
    
    -- Timestamps
    pev_eve_dcr TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Data de criação do registro
    pev_eve_dat TIMESTAMP, -- Data de atualização do registro
    
    -- Índices
    CONSTRAINT uk_pev_eve_ide UNIQUE (pev_eve_ide) -- Identificador de evento único (idempotência)
);

-- Criar índice para busca por identificador de evento
CREATE INDEX IF NOT EXISTS idx_pev_eve_ide ON processamento_eventos.evento(pev_eve_ide);

-- Criar índice para busca por identificador de card
CREATE INDEX IF NOT EXISTS idx_pev_eve_idc ON processamento_eventos.evento(pev_eve_idc);

-- Criar índice para busca por email do cliente
CREATE INDEX IF NOT EXISTS idx_pev_eve_ema ON processamento_eventos.evento(pev_eve_ema);

-- Criar índice para busca por flag de processamento
CREATE INDEX IF NOT EXISTS idx_pev_eve_fpr ON processamento_eventos.evento(pev_eve_fpr);

-- Comentários da tabela e colunas
COMMENT ON TABLE processamento_eventos.evento IS 'Tabela de eventos do contexto de processamento de eventos';
COMMENT ON COLUMN processamento_eventos.evento.pev_eve_int IS 'Identificador interno do evento (gerado pela sequência)';
COMMENT ON COLUMN processamento_eventos.evento.pev_eve_ide IS 'Identificador do evento (event_id do Pipefy)';
COMMENT ON COLUMN processamento_eventos.evento.pev_eve_idc IS 'Identificador do card (card_id do Pipefy)';
COMMENT ON COLUMN processamento_eventos.evento.pev_eve_ema IS 'Email do cliente';
COMMENT ON COLUMN processamento_eventos.evento.pev_eve_tms IS 'Timestamp do evento recebido do Pipefy';
COMMENT ON COLUMN processamento_eventos.evento.pev_eve_fpr IS 'Flag indicando se o evento foi processado (idempotência)';
COMMENT ON COLUMN processamento_eventos.evento.pev_eve_dcr IS 'Data de criação do registro';
COMMENT ON COLUMN processamento_eventos.evento.pev_eve_dat IS 'Data de atualização do registro';
