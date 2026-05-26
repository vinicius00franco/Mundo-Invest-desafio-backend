package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// ConfiguracaoBancoDados representa as configurações de conexão com o banco de dados
type ConfiguracaoBancoDados struct {
	Host     string
	Port     string
	Usuario  string
	Senha    string
	Banco    string
	SSLMode  string
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

	banco.SetMaxOpenConns(25)
	banco.SetMaxIdleConns(5)

	log.Println("Conexão com banco de dados estabelecida com sucesso")
	return banco, nil
}
