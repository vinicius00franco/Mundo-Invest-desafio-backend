package database_test

import (
	"testing"

	"github.com/MundoInvest/backend/internal/shared/config"
	"github.com/MundoInvest/backend/internal/shared/database"
	"github.com/MundoInvest/backend/internal/shared/logger"
	"github.com/MundoInvest/backend/internal/shared/mensagens"
)

// TestNovaConexaoSucesso testa a criação de conexão com sucesso
func TestNovaConexaoSucesso(t *testing.T) {
	// Inicializar logger antes do teste
	logger.Init()

	// Este teste requer variáveis de ambiente configuradas
	// Para executar: export DB_HOST=localhost DB_PORT=5434 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=mundo_invest

	configDB := database.ConfiguracaoBancoDados{
		Host:    "localhost",
		Port:    "5434",
		Usuario: "postgres",
		Senha:   "postgres",
		Banco:   "mundo_invest",
		SSLMode: "disable",
	}

	cfg := config.Load()
	db, err := database.NovoBancoDados(configDB, cfg)
	if err != nil {
		t.Skipf("Teste de integração pulado: não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Verificar se a conexão está ativa
	if err := db.Ping(); err != nil {
		t.Errorf("Erro ao verificar conexão: %v", err)
	}
}

// TestNovaConexaoFalha testa falha na conexão
func TestNovaConexaoFalha(t *testing.T) {
	configDB := database.ConfiguracaoBancoDados{
		Host:    "localhost",
		Port:    "9999", // Porta inválida
		Usuario: "postgres",
		Senha:   "postgres",
		Banco:   "mundo_invest",
		SSLMode: "disable",
	}

	cfg := config.Load()
	catalogo := mensagens.ObterCatalogo()
	_, err := database.NovoBancoDados(configDB, cfg)
	if err == nil {
		t.Error(catalogo.Texto(mensagens.ErrEsperadoErro))
	}
}

// TestNovaConexaoVariaveisAmbiente testa conexão usando variáveis de ambiente
func TestNovaConexaoVariaveisAmbiente(t *testing.T) {
	// Inicializar logger
	logger.Init()

	// Configurar variáveis de ambiente para teste
	// Em produção, estas variáveis viriam do ambiente
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5434")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "postgres")
	t.Setenv("DB_NAME", "mundo_invest")
	t.Setenv("DB_SSLMODE", "disable")

	cfg := config.Load()
	db, err := database.NovaConexao(cfg)
	if err != nil {
		t.Skipf("Teste de integração pulado: não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Verificar se a conexão está ativa
	if err := db.Ping(); err != nil {
		t.Errorf("Erro ao verificar conexão: %v", err)
	}
}

// TestTransacaoAtômica testa se transações são atômicas
func TestTransacaoAtômica(t *testing.T) {
	configDB := database.ConfiguracaoBancoDados{
		Host:    "localhost",
		Port:    "5434",
		Usuario: "postgres",
		Senha:   "postgres",
		Banco:   "mundo_invest",
		SSLMode: "disable",
	}

	cfg := config.Load()
	db, err := database.NovoBancoDados(configDB, cfg)
	if err != nil {
		t.Skipf("Teste de integração pulado: não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Iniciar transação
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Erro ao iniciar transação: %v", err)
	}

	// Tentar inserir um registro
	_, err = tx.Exec("INSERT INTO gestao_clientes.cliente (gcl_cli_nom, gcl_cli_ema, gcl_cli_tso, gcl_cli_pat, gcl_cli_stc) VALUES ($1, $2, $3, $4, $5)",
		"Teste Transação", "teste.transacao@example.com", "teste", 1000.00, "Teste")
	if err != nil {
		tx.Rollback()
		t.Fatalf("Erro ao inserir registro: %v", err)
	}

	// Rollback para não afetar o banco de dados
	if err := tx.Rollback(); err != nil {
		t.Fatalf("Erro ao fazer rollback: %v", err)
	}

	// Verificar se o registro não foi inserido (devido ao rollback)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM gestao_clientes.cliente WHERE gcl_cli_ema = $1", "teste.transacao@example.com").Scan(&count)
	if err != nil {
		t.Fatalf("Erro ao verificar registro: %v", err)
	}

	if count > 0 {
		// Limpar o registro se foi inserido
		db.Exec("DELETE FROM gestao_clientes.cliente WHERE gcl_cli_ema = $1", "teste.transacao@example.com")
		t.Error("Registro foi inserido mesmo após rollback, transação não é atômica")
	}
}

// TestRestricaoUnicidade testa restrição de unicidade de email
func TestRestricaoUnicidade(t *testing.T) {
	catalogo := mensagens.ObterCatalogo()
	configDB := database.ConfiguracaoBancoDados{
		Host:    "localhost",
		Port:    "5434",
		Usuario: "postgres",
		Senha:   "postgres",
		Banco:   "mundo_invest",
		SSLMode: "disable",
	}

	cfg := config.Load()
	db, err := database.NovoBancoDados(configDB, cfg)
	if err != nil {
		t.Skipf("Teste de integração pulado: não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Inserir um cliente de teste
	email := "teste.unicidade@example.com"
	_, err = db.Exec("INSERT INTO gestao_clientes.cliente (gcl_cli_nom, gcl_cli_ema, gcl_cli_tso, gcl_cli_pat, gcl_cli_stc) VALUES ($1, $2, $3, $4, $5)",
		"Teste Unicidade", email, "teste", 1000.00, "Teste")
	if err != nil {
		t.Fatalf("Erro ao inserir cliente de teste: %v", err)
	}
	// Limpar o registro após o teste
	defer db.Exec("DELETE FROM gestao_clientes.cliente WHERE gcl_cli_ema = $1", email)

	// Tentar inserir outro cliente com o mesmo email (deve falhar)
	_, err = db.Exec("INSERT INTO gestao_clientes.cliente (gcl_cli_nom, gcl_cli_ema, gcl_cli_tso, gcl_cli_pat, gcl_cli_stc) VALUES ($1, $2, $3, $4, $5)",
		"Teste Unicidade 2", email, "teste", 2000.00, "Teste")

	if err == nil {
		t.Error(catalogo.Texto(mensagens.ErrViolacaoUnicidade))
	}
}

// TestConfiguracaoPoolConexoes testa configuração de pool de conexões
func TestConfiguracaoPoolConexoes(t *testing.T) {
	configDB := database.ConfiguracaoBancoDados{
		Host:    "localhost",
		Port:    "5434",
		Usuario: "postgres",
		Senha:   "postgres",
		Banco:   "mundo_invest",
		SSLMode: "disable",
	}

	cfg := config.Load()
	db, err := database.NovoBancoDados(configDB, cfg)
	if err != nil {
		t.Skipf("Teste de integração pulado: não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Verificar configurações do pool
	stats := db.Stats()

	if stats.MaxOpenConnections != 25 {
		t.Errorf("MaxOpenConnections esperado 25, obtido %d", stats.MaxOpenConnections)
	}

	if stats.Idle != 0 && stats.Idle > 5 {
		t.Logf("Idle connections: %d (configurado para máximo 5)", stats.Idle)
	}
}

// TestIntegracaoSchemas testa se schemas estão criados corretamente
func TestIntegracaoSchemas(t *testing.T) {
	configDB := database.ConfiguracaoBancoDados{
		Host:    "localhost",
		Port:    "5434",
		Usuario: "postgres",
		Senha:   "postgres",
		Banco:   "mundo_invest",
		SSLMode: "disable",
	}

	cfg := config.Load()
	db, err := database.NovoBancoDados(configDB, cfg)
	if err != nil {
		t.Skipf("Teste de integração pulado: não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Verificar se schema gestao_clientes existe
	var schemaExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = 'gestao_clientes')").Scan(&schemaExists)
	if err != nil {
		t.Fatalf("Erro ao verificar schema gestao_clientes: %v", err)
	}

	if !schemaExists {
		t.Error("Schema gestao_clientes não existe")
	}

	// Verificar se schema processamento_eventos existe
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = 'processamento_eventos')").Scan(&schemaExists)
	if err != nil {
		t.Fatalf("Erro ao verificar schema processamento_eventos: %v", err)
	}

	if !schemaExists {
		t.Error("Schema processamento_eventos não existe")
	}
}

// TestIntegracaoTabelas testa se tabelas estão criadas corretamente
func TestIntegracaoTabelas(t *testing.T) {
	configDB := database.ConfiguracaoBancoDados{
		Host:    "localhost",
		Port:    "5434",
		Usuario: "postgres",
		Senha:   "postgres",
		Banco:   "mundo_invest",
		SSLMode: "disable",
	}

	cfg := config.Load()
	db, err := database.NovoBancoDados(configDB, cfg)
	if err != nil {
		t.Skipf("Teste de integração pulado: não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Verificar se tabela cliente existe
	var tableExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'gestao_clientes' AND table_name = 'cliente')").Scan(&tableExists)
	if err != nil {
		t.Fatalf("Erro ao verificar tabela cliente: %v", err)
	}

	if !tableExists {
		t.Error("Tabela cliente não existe no schema gestao_clientes")
	}

	// Verificar se tabela evento existe
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'processamento_eventos' AND table_name = 'evento')").Scan(&tableExists)
	if err != nil {
		t.Fatalf("Erro ao verificar tabela evento: %v", err)
	}

	if !tableExists {
		t.Error("Tabela evento não existe no schema processamento_eventos")
	}
}
