package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MundoInvest/backend/internal/dominio"
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
func NovoBancoDados(config ConfiguracaoBancoDados) (*sql.DB, error) {
	stringConexao := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.Usuario,
		config.Senha,
		config.Banco,
		config.SSLMode,
	)

	banco, err := sql.Open("postgres", stringConexao)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com banco de dados: %w", err)
	}

	if err = banco.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao testar conexão com banco de dados: %w", err)
	}

	banco.SetMaxOpenConns(dominio.MaxOpenConns)
	banco.SetMaxIdleConns(dominio.MaxIdleConns)
	banco.SetConnMaxLifetime(time.Duration(dominio.ConnMaxLifetime) * time.Second)
	banco.SetConnMaxIdleTime(time.Duration(dominio.ConnMaxIdleTime) * time.Second)

	log.Println("Conexão com banco de dados estabelecida com sucesso")
	return banco, nil
}

// NovaConexao cria uma nova conexão com o banco de dados usando variáveis de ambiente
func NovaConexao() (*sql.DB, error) {
	config := ConfiguracaoBancoDados{
		Host:    getEnv("DB_HOST", "localhost"),
		Port:    getEnv("DB_PORT", "5432"),
		Usuario: getEnv("DB_USER", "postgres"),
		Senha:   getEnv("DB_PASSWORD", "postgres"),
		Banco:   getEnv("DB_NAME", "mundo_invest"),
		SSLMode: getEnv("DB_SSLMODE", "disable"),
	}

	return NovoBancoDados(config)
}

// getEnv obtém valor de variável de ambiente ou retorna valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
