-- Migration: Criar usuários específicos com roles
-- Contexto: Segurança do Banco de Dados
-- Linguagem Ubíqua: Usuários, Roles, Autenticação
-- Objetivo: Criar usuários específicos com roles e políticas de senha

-- Criar usuário para leitura de clientes
CREATE USER IF NOT EXISTS app_cliente_read WITH PASSWORD 'senha_temp_cliente_read_2024';

-- Criar usuário para escrita de clientes
CREATE USER IF NOT EXISTS app_cliente_write WITH PASSWORD 'senha_temp_cliente_write_2024';

-- Criar usuário para escrita de webhooks
CREATE USER IF NOT EXISTS app_webhook_write WITH PASSWORD 'senha_temp_webhook_write_2024';

-- Criar usuário para backup
CREATE USER IF NOT EXISTS backup_user WITH PASSWORD 'senha_temp_backup_2024';

-- Criar usuário para DBA
CREATE USER IF NOT EXISTS dba_senior WITH PASSWORD 'senha_temp_dba_2024';

-- Atribuir roles aos usuários
GRANT gestao_clientes_read TO app_cliente_read;
GRANT gestao_clientes_write TO app_cliente_write;
GRANT processamento_eventos_write TO app_webhook_write;
GRANT backup_operator TO backup_user;
GRANT dba_admin TO dba_senior;

-- Configurar limites de conexão por aplicação
ALTER USER app_cliente_read CONNECTION LIMIT 10;
ALTER USER app_cliente_write CONNECTION LIMIT 20;
ALTER USER app_webhook_write CONNECTION LIMIT 15;
ALTER USER backup_user CONNECTION LIMIT 5;
ALTER USER dba_senior CONNECTION LIMIT 3;

-- Comentários dos usuários
COMMENT ON USER app_cliente_read IS 'Usuário de aplicação para leitura de dados de clientes';
COMMENT ON USER app_cliente_write IS 'Usuário de aplicação para escrita de dados de clientes';
COMMENT ON USER app_webhook_write IS 'Usuário de aplicação para processamento de webhooks';
COMMENT ON USER backup_user IS 'Usuário para operações de backup';
COMMENT ON USER dba_senior IS 'Usuário administrador do banco de dados';
