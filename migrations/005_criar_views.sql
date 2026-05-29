-- Migration: Criar views para acesso legível
-- Contexto: Gestão de Clientes e Processamento de Eventos
-- Linguagem Ubíqua: Cliente, Evento, Views legíveis
-- Objetivo: Facilitar leitura e consultas sem trigramação

-- View legível para clientes
CREATE OR REPLACE VIEW vw_cliente AS
SELECT 
    gcl_cli_int AS identificador_interno,
    gcl_cli_ext AS identificador_externo,
    gcl_cli_nom AS nome,
    gcl_cli_ema AS email,
    gcl_cli_pat AS valor_patrimonio,
    gcl_cli_tso AS tipo_solicitacao,
    gcl_cli_stc AS status,
    gcl_cli_npr AS nivel_prioridade,
    gcl_cli_dcr AS data_criacao,
    gcl_cli_dat AS data_atualizacao
FROM gestao_clientes.cliente;

-- Comentário da view
COMMENT ON VIEW vw_cliente IS 'View legível para consulta de clientes sem trigramação';

-- View legível para eventos
CREATE OR REPLACE VIEW vw_evento AS
SELECT 
    pev_eve_int AS identificador_interno,
    pev_eve_ide AS identificador_evento,
    pev_eve_idc AS identificador_card,
    pev_eve_ema AS email_cliente,
    pev_eve_tms AS timestamp_evento,
    pev_eve_fpr AS foi_processado,
    pev_eve_dcr AS data_criacao,
    pev_eve_dat AS data_atualizacao
FROM processamento_eventos.evento;

-- Comentário da view
COMMENT ON VIEW vw_evento IS 'View legível para consulta de eventos sem trigramação';

-- View legível para clientes com eventos processados
CREATE OR REPLACE VIEW vw_cliente_eventos AS
SELECT 
    c.gcl_cli_int AS identificador_interno,
    c.gcl_cli_ext AS identificador_externo,
    c.gcl_cli_nom AS nome,
    c.gcl_cli_ema AS email,
    c.gcl_cli_pat AS valor_patrimonio,
    c.gcl_cli_tso AS tipo_solicitacao,
    c.gcl_cli_stc AS status,
    c.gcl_cli_npr AS nivel_prioridade,
    COUNT(e.pev_eve_int) AS total_eventos,
    COUNT(CASE WHEN e.pev_eve_fpr = TRUE THEN 1 END) AS eventos_processados
FROM gestao_clientes.cliente c
LEFT JOIN processamento_eventos.evento e ON c.gcl_cli_ema = e.pev_eve_ema
GROUP BY c.gcl_cli_int, c.gcl_cli_ext, c.gcl_cli_nom, c.gcl_cli_ema, c.gcl_cli_pat, c.gcl_cli_tso, c.gcl_cli_stc, c.gcl_cli_npr;

-- Comentário da view
COMMENT ON VIEW vw_cliente_eventos IS 'View legível para consulta de clientes com resumo de eventos processados';
