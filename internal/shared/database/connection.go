package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/MundoInvest/backend/internal/shared/config"
	"github.com/MundoInvest/backend/internal/shared/logger"
	"github.com/MundoInvest/backend/internal/shared/mensagens"
	_ "github.com/lib/pq"
)

// ConfiguracaoBancoDados representa as configurações de conexão com o banco de dados
type ConfiguracaoBancoDados struct {
	Host    string
	Port    string
	Usuario string
	Senha   string
	Banco   string
	SSLMode string
}

// NovoBancoDados cria uma nova conexão com o banco de dados PostgreSQL
func NovoBancoDados(cfg ConfiguracaoBancoDados, appConfig *config.Config) (*sql.DB, error) {
	catalogo := mensagens.ObterCatalogo()
	stringConexao := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Usuario,
		cfg.Senha,
		cfg.Banco,
		cfg.SSLMode,
	)

	banco, err := sql.Open("postgres", stringConexao)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", catalogo.Texto(mensagens.ErrAbrirConexao), err)
	}

	if err = banco.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", catalogo.Texto(mensagens.ErrTestarConexao), err)
	}

	banco.SetMaxOpenConns(appConfig.Database.MaxOpenConns)
	banco.SetMaxIdleConns(appConfig.Database.MaxIdleConns)
	banco.SetConnMaxLifetime(appConfig.Database.ConnMaxLifetime)
	banco.SetConnMaxIdleTime(appConfig.Database.ConnMaxIdleTime)

	logger.Info("Conexão com banco de dados estabelecida com sucesso",
		"host", cfg.Host,
		"database", cfg.Banco,
	)
	return banco, nil
}

// NovaConexao cria uma nova conexão com o banco de dados usando variáveis de ambiente
func NovaConexao(appConfig *config.Config) (*sql.DB, error) {
	config := ConfiguracaoBancoDados{
		Host:    getEnv("DB_HOST", "localhost"),
		Port:    getEnv("DB_PORT", "5432"),
		Usuario: getEnv("DB_USER", "postgres"),
		Senha:   getEnv("DB_PASSWORD", "postgres"),
		Banco:   getEnv("DB_NAME", "mundo_invest"),
		SSLMode: getEnv("DB_SSLMODE", "disable"),
	}

	return NovoBancoDados(config, appConfig)
}

// getEnv obtém valor de variável de ambiente ou retorna valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
