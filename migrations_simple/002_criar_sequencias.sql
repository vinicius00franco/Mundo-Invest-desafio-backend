-- Migration: Criar sequências para identificadores internos
-- Contexto: Configuração de sequências SQL
-- Linguagem Ubíqua: Identificador Interno (gcl_cli_int, pev_eve_int, ipf_car_int)

-- Criar sequência para identificador interno de cliente (gcl_cli_int)
CREATE SEQUENCE IF NOT EXISTS gestao_clientes.seq_gcl_cli_int
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- Criar sequência para identificador interno de evento (pev_eve_int)
CREATE SEQUENCE IF NOT EXISTS processamento_eventos.seq_pev_eve_int
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- Criar sequência para identificador interno de card (ipf_car_int)
CREATE SEQUENCE IF NOT EXISTS integracao_pipefy.seq_ipf_car_int
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- Comentários das sequências
COMMENT ON SEQUENCE gestao_clientes.seq_gcl_cli_int IS 'Sequência para gerar identificadores internos de clientes (gcl_cli_int)';
COMMENT ON SEQUENCE processamento_eventos.seq_pev_eve_int IS 'Sequência para gerar identificadores internos de eventos (pev_eve_int)';
COMMENT ON SEQUENCE integracao_pipefy.seq_ipf_car_int IS 'Sequência para gerar identificadores internos de cards (ipf_car_int)';
