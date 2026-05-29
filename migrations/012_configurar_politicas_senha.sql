-- Migration: Configurar políticas de senha e logging
-- Contexto: Segurança do Banco de Dados
-- Linguagem Ubíqua: Políticas de Senha, Logging, Configuração
-- Objetivo: Configurar políticas de senha (SCRAM-SHA-256) e logging de consultas

-- Configurar método de autenticação para SCRAM-SHA-256
ALTER SYSTEM SET password_encryption = 'scram-sha-256';

-- Configurar expiração de senhas (90 dias)
ALTER SYSTEM SET password_encryption = 'scram-sha-256';

-- Configurar logging de consultas (mod)
ALTER SYSTEM SET log_statement = 'mod'; -- Registra DDL, INSERT, UPDATE, DELETE
ALTER SYSTEM SET log_duration = on; -- Registra duração das consultas
ALTER SYSTEM SET log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '; -- Formato do log
ALTER SYSTEM SET log_connections = on; -- Registra conexões
ALTER SYSTEM SET log_disconnections = on; -- Registra desconexões

-- Configurar logging para arquivo
ALTER SYSTEM SET logging_collector = on;
ALTER SYSTEM SET log_directory = 'log';
ALTER SYSTEM SET log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log';
ALTER SYSTEM SET log_rotation_age = 1d;
ALTER SYSTEM SET log_rotation_size = 100MB;

-- Configurar níveis de log
ALTER SYSTEM SET log_min_duration_statement = 1000; -- Registra consultas que levam mais de 1 segundo
ALTER SYSTEM SET log_min_error_statement = error;

-- Comentários de configuração
COMMENT ON PARAMETER password_encryption IS 'Método de criptografia de senhas (SCRAM-SHA-256)';
COMMENT ON PARAMETER log_statement IS 'Nível de logging de consultas (mod = DDL, INSERT, UPDATE, DELETE)';
COMMENT ON PARAMETER log_duration IS 'Registra duração das consultas';
COMMENT ON PARAMETER log_line_prefix IS 'Formato do prefixo das linhas de log';
COMMENT ON PARAMETER log_connections IS 'Registra conexões ao banco';
COMMENT ON PARAMETER log_disconnections IS 'Registra desconexões do banco';
COMMENT ON PARAMETER logging_collector IS 'Habilita coletor de logs';
COMMENT ON PARAMETER log_directory IS 'Diretório para armazenar logs';
COMMENT ON PARAMETER log_filename IS 'Nome do arquivo de log';
COMMENT ON PARAMETER log_rotation_age IS 'Idade máxima do arquivo de log antes de rotacionar';
COMMENT ON PARAMETER log_rotation_size IS 'Tamanho máximo do arquivo de log antes de rotacionar';
COMMENT ON PARAMETER log_min_duration_statement IS 'Duração mínima para registrar consultas lentas (ms)';
COMMENT ON PARAMETER log_min_error_statement IS 'Nível mínimo de erro para registrar statements';
