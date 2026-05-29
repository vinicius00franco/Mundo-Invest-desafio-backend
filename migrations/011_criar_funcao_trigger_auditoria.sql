-- Migration: Criar função de trigger para auditoria
-- Contexto: Segurança do Banco de Dados
-- Linguagem Ubíqua: Auditoria, Trigger, Função
-- Objetivo: Implementar auditoria automática de operações

-- Criar função para registrar operações de INSERT
CREATE OR REPLACE FUNCTION auditoria.registrar_insert()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO auditoria.log_acesso (
        log_acs_usu, log_acs_ope, log_acs_sch, log_acs_tba, log_acs_reg
    ) VALUES (
        current_user,
        'INSERT',
        TG_TABLE_SCHEMA,
        TG_TABLE_NAME,
        row_to_json(NEW)::text
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Criar função para registrar operações de UPDATE
CREATE OR REPLACE FUNCTION auditoria.registrar_update()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO auditoria.log_acesso (
        log_acs_usu, log_acs_ope, log_acs_sch, log_acs_tba, log_acs_reg
    ) VALUES (
        current_user,
        'UPDATE',
        TG_TABLE_SCHEMA,
        TG_TABLE_NAME,
        json_build_object(
            'antigo', row_to_json(OLD),
            'novo', row_to_json(NEW)
        )::text
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Criar função para registrar operações de DELETE
CREATE OR REPLACE FUNCTION auditoria.registrar_delete()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO auditoria.log_acesso (
        log_acs_usu, log_acs_ope, log_acs_sch, log_acs_tba, log_acs_reg
    ) VALUES (
        current_user,
        'DELETE',
        TG_TABLE_SCHEMA,
        TG_TABLE_NAME,
        row_to_json(OLD)::text
    );
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Criar triggers para tabela cliente
DROP TRIGGER IF EXISTS trg_auditoria_cliente_insert ON gestao_clientes.cliente;
CREATE TRIGGER trg_auditoria_cliente_insert
    AFTER INSERT ON gestao_clientes.cliente
    FOR EACH ROW EXECUTE FUNCTION auditoria.registrar_insert();

DROP TRIGGER IF EXISTS trg_auditoria_cliente_update ON gestao_clientes.cliente;
CREATE TRIGGER trg_auditoria_cliente_update
    AFTER UPDATE ON gestao_clientes.cliente
    FOR EACH ROW EXECUTE FUNCTION auditoria.registrar_update();

DROP TRIGGER IF EXISTS trg_auditoria_cliente_delete ON gestao_clientes.cliente;
CREATE TRIGGER trg_auditoria_cliente_delete
    AFTER DELETE ON gestao_clientes.cliente
    FOR EACH ROW EXECUTE FUNCTION auditoria.registrar_delete();

-- Criar triggers para tabela evento
DROP TRIGGER IF EXISTS trg_auditoria_evento_insert ON processamento_eventos.evento;
CREATE TRIGGER trg_auditoria_evento_insert
    AFTER INSERT ON processamento_eventos.evento
    FOR EACH ROW EXECUTE FUNCTION auditoria.registrar_insert();

DROP TRIGGER IF EXISTS trg_auditoria_evento_update ON processamento_eventos.evento;
CREATE TRIGGER trg_auditoria_evento_update
    AFTER UPDATE ON processamento_eventos.evento
    FOR EACH ROW EXECUTE FUNCTION auditoria.registrar_update();

DROP TRIGGER IF EXISTS trg_auditoria_evento_delete ON processamento_eventos.evento;
CREATE TRIGGER trg_auditoria_evento_delete
    AFTER DELETE ON processamento_eventos.evento
    FOR EACH ROW EXECUTE FUNCTION auditoria.registrar_delete();

-- Criar triggers para tabela card
DROP TRIGGER IF EXISTS trg_auditoria_card_insert ON integracao_pipefy.card;
CREATE TRIGGER trg_auditoria_card_insert
    AFTER INSERT ON integracao_pipefy.card
    FOR EACH ROW EXECUTE FUNCTION auditoria.registrar_insert();

DROP TRIGGER IF EXISTS trg_auditoria_card_update ON integracao_pipefy.card;
CREATE TRIGGER trg_auditoria_card_update
    AFTER UPDATE ON integracao_pipefy.card
    FOR EACH ROW EXECUTE FUNCTION auditoria.registrar_update();

DROP TRIGGER IF EXISTS trg_auditoria_card_delete ON integracao_pipefy.card;
CREATE TRIGGER trg_auditoria_card_delete
    AFTER DELETE ON integracao_pipefy.card
    FOR EACH ROW EXECUTE FUNCTION auditoria.registrar_delete();

-- Comentários das funções
COMMENT ON FUNCTION auditoria.registrar_insert() IS 'Função para registrar operações de INSERT na auditoria';
COMMENT ON FUNCTION auditoria.registrar_update() IS 'Função para registrar operações de UPDATE na auditoria';
COMMENT ON FUNCTION auditoria.registrar_delete() IS 'Função para registrar operações de DELETE na auditoria';
